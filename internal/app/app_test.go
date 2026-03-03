package app_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"desent-api-quest/internal/app"
)

func TestAPIFlow(t *testing.T) {
	application := app.New()
	handler := application.Handler()
	token := issueToken(t, handler, "admin", "secret")

	t.Run("ping returns success true", func(t *testing.T) {
		resp := mustRequest(t, handler, "GET", "/ping", "", "")
		assertStatus(t, resp, http.StatusOK)

		body := decodeBody[map[string]bool](t, resp)
		if body["success"] != true {
			t.Fatalf("expected true, got %#v", body)
		}
	})

	t.Run("echo returns same JSON", func(t *testing.T) {
		requestBody := `{"z":1,"a":{"second":2,"first":1},"items":[{"b":2,"a":1}]}`
		resp := mustRequest(t, handler, "POST", "/echo", requestBody, "")
		assertStatus(t, resp, http.StatusOK)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read response body: %v", err)
		}
		if string(body) != requestBody {
			t.Fatalf("expected exact echo body %q, got %q", requestBody, string(body))
		}
	})

	t.Run("echo invalid JSON returns bad request", func(t *testing.T) {
		resp := mustRequest(t, handler, "POST", "/echo", `{"hello":`, "")
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("create book returns created book", func(t *testing.T) {
		resp := mustRequest(t, handler, "POST", "/books", `{"title":"Clean Code","author":"Robert C. Martin","year":2008}`, "")
		assertStatus(t, resp, http.StatusCreated)

		book := decodeBody[map[string]any](t, resp)
		if book["id"] == "" {
			t.Fatalf("expected generated id, got %#v", book)
		}
	})

	t.Run("create invalid book returns bad request", func(t *testing.T) {
		resp := mustRequest(t, handler, "POST", "/books", `{"title":"","author":"Robert C. Martin","year":2008}`, "")
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("protected books listing requires token", func(t *testing.T) {
		resp := mustRequest(t, handler, "GET", "/books", "", "")
		assertStatus(t, resp, http.StatusUnauthorized)
	})

	t.Run("protected books listing rejects malformed bearer", func(t *testing.T) {
		resp := mustRequest(t, handler, "GET", "/books", "", "Token nope")
		assertStatus(t, resp, http.StatusUnauthorized)
	})

	t.Run("get books list supports pagination", func(t *testing.T) {
		_ = mustRequest(t, handler, "POST", "/books", `{"title":"Refactoring","author":"Martin Fowler","year":1999}`, "")
		_ = mustRequest(t, handler, "POST", "/books", `{"title":"Patterns","author":"Erich Gamma","year":1994}`, "")

		resp := mustRequest(t, handler, "GET", "/books?page=1&limit=2", "", "Bearer "+token)
		assertStatus(t, resp, http.StatusOK)

		body := decodeBody[[]map[string]any](t, resp)
		if len(body) != 2 {
			t.Fatalf("expected 2 items, got %#v", body)
		}
	})

	t.Run("get books list supports author filter", func(t *testing.T) {
		resp := mustRequest(t, handler, "GET", "/books?author=martin%20fowler", "", "Bearer "+token)
		assertStatus(t, resp, http.StatusOK)

		body := decodeBody[[]map[string]any](t, resp)
		if len(body) != 1 {
			t.Fatalf("expected one filtered item, got %#v", body)
		}
	})

	t.Run("get books list rejects invalid page", func(t *testing.T) {
		resp := mustRequest(t, handler, "GET", "/books?page=0", "", "Bearer "+token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("get books list rejects invalid limit", func(t *testing.T) {
		resp := mustRequest(t, handler, "GET", "/books?limit=0", "", "Bearer "+token)
		assertStatus(t, resp, http.StatusBadRequest)
	})

	t.Run("get book by id returns existing book", func(t *testing.T) {
		created := mustRequest(t, handler, "POST", "/books", `{"title":"DDD","author":"Eric Evans","year":2003}`, "")
		assertStatus(t, created, http.StatusCreated)
		book := decodeBody[map[string]any](t, created)

		resp := mustRequest(t, handler, "GET", "/books/"+book["id"].(string), "", "")
		assertStatus(t, resp, http.StatusOK)
	})

	t.Run("get book by id returns not found", func(t *testing.T) {
		resp := mustRequest(t, handler, "GET", "/books/99999", "", "")
		assertStatus(t, resp, http.StatusNotFound)
	})

	t.Run("update book returns updated book", func(t *testing.T) {
		created := mustRequest(t, handler, "POST", "/books", `{"title":"Go","author":"Anon","year":2020}`, "")
		assertStatus(t, created, http.StatusCreated)
		book := decodeBody[map[string]any](t, created)

		resp := mustRequest(t, handler, "PUT", "/books/"+book["id"].(string), `{"title":"Go in Action","author":"Anon","year":2021}`, "")
		assertStatus(t, resp, http.StatusOK)

		updated := decodeBody[map[string]any](t, resp)
		if updated["title"] != "Go in Action" {
			t.Fatalf("unexpected updated payload: %#v", updated)
		}
	})

	t.Run("update missing book returns not found", func(t *testing.T) {
		resp := mustRequest(t, handler, "PUT", "/books/99999", `{"title":"Missing","author":"Anon","year":2021}`, "")
		assertStatus(t, resp, http.StatusNotFound)
	})

	t.Run("delete existing book returns no content", func(t *testing.T) {
		created := mustRequest(t, handler, "POST", "/books", `{"title":"Delete Me","author":"Anon","year":2021}`, "")
		assertStatus(t, created, http.StatusCreated)
		book := decodeBody[map[string]any](t, created)

		resp := mustRequest(t, handler, "DELETE", "/books/"+book["id"].(string), "", "")
		assertStatus(t, resp, http.StatusNoContent)
	})

	t.Run("delete missing book returns not found", func(t *testing.T) {
		resp := mustRequest(t, handler, "DELETE", "/books/99999", "", "")
		assertStatus(t, resp, http.StatusNotFound)
	})

	t.Run("invalid credentials return unauthorized", func(t *testing.T) {
		resp := mustRequest(t, handler, "POST", "/auth/token", `{"username":"admin","password":"wrong"}`, "")
		assertStatus(t, resp, http.StatusUnauthorized)
	})

	t.Run("alternate valid credentials return token", func(t *testing.T) {
		resp := mustRequest(t, handler, "POST", "/auth/token", `{"username":"user","password":"password"}`, "")
		assertStatus(t, resp, http.StatusOK)

		body := decodeBody[map[string]string](t, resp)
		if body["token"] == "" {
			t.Fatalf("expected token in response, got %#v", body)
		}
	})
}

func issueToken(t *testing.T, handler http.Handler, username, password string) string {
	t.Helper()

	resp := mustRequest(t, handler, "POST", "/auth/token", `{"username":"`+username+`","password":"`+password+`"}`, "")
	assertStatus(t, resp, http.StatusOK)

	body := decodeBody[map[string]string](t, resp)
	if body["token"] == "" {
		t.Fatal("expected token in response")
	}

	return body["token"]
}

func mustRequest(t *testing.T, handler http.Handler, method, path, rawBody, authHeader string) *http.Response {
	t.Helper()

	var body io.Reader
	if rawBody != "" {
		body = bytes.NewBufferString(rawBody)
	}

	req, err := http.NewRequest(method, "http://example.com"+path, body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if rawBody != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	resp := recorder.Result()
	t.Cleanup(func() {
		resp.Body.Close()
	})
	return resp
}

func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status %d, got %d: %s", want, resp.StatusCode, string(data))
	}
}

func decodeBody[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	var value T
	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return value
}

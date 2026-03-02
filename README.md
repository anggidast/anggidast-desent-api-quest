# desent-api-quest

Simple REST API in Go using the standard library only.

## Run

```bash
go run ./cmd/api
```

Server address:

```text
http://localhost:8080
```

## Endpoints

- `GET /ping`
- `POST /echo`
- `POST /auth/token`
- `POST /books`
- `GET /books` (requires Bearer token)
- `GET /books/{id}`
- `PUT /books/{id}`
- `DELETE /books/{id}`

## Example

```bash
curl -X POST http://localhost:8080/auth/token \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"secret"}'
```

```bash
curl http://localhost:8080/books \
  -H 'Authorization: Bearer YOUR_TOKEN'
```

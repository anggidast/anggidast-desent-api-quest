package domain

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type BookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type BookListResult struct {
	Items      []Book `json:"items"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalItems int    `json:"total_items"`
	TotalPages int    `json:"total_pages"`
}

type ListBooksParams struct {
	Author string
	Page   int
	Limit  int
}

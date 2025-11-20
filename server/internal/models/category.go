package models

type Category struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type CategoryWithBooks struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Books       []Book `json:"books"`
}
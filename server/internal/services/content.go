package services

import (
    "sync"
    "picturebook/server/internal/db"
    "picturebook/server/internal/models"
    "picturebook/server/internal/repository"
)

type ContentService struct {
    mu   sync.RWMutex
    repo *repository.Repo
}

func NewContentService() *ContentService {
    d, err := db.Connect()
    if err != nil { panic(err) }
    if err := db.EnsureSchema(d); err != nil { panic(err) }
    r := repository.New(d)
    s := &ContentService{repo: r}
    // seed if empty
    if err := r.SeedIfEmpty(sampleBooks()); err != nil { panic(err) }
    return s
}

func (s *ContentService) ListEditorPicks() []models.Book {
    res, err := s.repo.ListEditorPicks()
    if err != nil { return []models.Book{} }
    return res
}

func (s *ContentService) ListBooks(age int, sort string, page, pageSize int) ([]models.Book, bool) {
    items, hasMore, err := s.repo.ListBooks(age, sort, page, pageSize)
    if err != nil { return []models.Book{}, false }
    return items, hasMore
}

func (s *ContentService) GetBook(id string) (models.Book, bool) {
    b, ok, _ := s.repo.GetBook(id)
    return b, ok
}

func (s *ContentService) AddBook(b models.Book) error {
    return s.repo.CreateBook(b)
}

func (s *ContentService) AddPages(bookID string, fileURLs []string) []models.Page {
    // create pages with generated IDs
    var pages []models.Page
    start := 0
    // fetch existing to determine start
    if book, ok, _ := s.repo.GetBook(bookID); ok {
        start = len(book.Pages)
    }
    for idx, url := range fileURLs {
        pages = append(pages, models.Page{ID: bookID + "-p-" + itoa(start+idx), Index: start + idx, ImageURL: url, Duration: nil})
    }
    _ = s.repo.AddPages(bookID, pages)
    return pages
}

func (s *ContentService) ReorderPages(bookID string, indexByID map[string]int) []models.Page {
    _ = s.repo.ReorderPages(bookID, indexByID)
    b, ok, _ := s.repo.GetBook(bookID)
    if !ok { return nil }
    return b.Pages
}

func (s *ContentService) SetEditorPick(bookID string, v bool) error {
    if err := s.repo.SetEditorPick(bookID, v); err != nil { return err }
    return nil
}

func (s *ContentService) DeleteBook(id string) error {
    return s.repo.DeleteBook(id)
}

func (s *ContentService) UpdateBook(b models.Book) error {
    return s.repo.UpdateBook(b)
}

func (s *ContentService) ListBooksAdmin(sort string, page, pageSize int, q string) ([]models.Book, bool) {
    items, more, err := s.repo.ListBooksAdmin(sort, page, pageSize, q)
    if err != nil { return []models.Book{}, false }
    return items, more
}

// Categories
func (s *ContentService) CreateCategory(c models.Category) error { return s.repo.CreateCategory(c) }
func (s *ContentService) DeleteCategory(id string) error { return s.repo.DeleteCategory(id) }
func (s *ContentService) ListCategories() ([]models.Category, error) { return s.repo.ListCategories() }
func (s *ContentService) ListCategoriesWithBooks() ([]models.CategoryWithBooks, error) { return s.repo.ListCategoriesWithBooks() }
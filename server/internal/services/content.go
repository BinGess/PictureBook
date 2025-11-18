package services

import (
    "sync"
    "picturebook/server/internal/models"
)

type ContentService struct {
    mu    sync.RWMutex
    books []models.Book
}

func NewContentService() *ContentService {
    s := &ContentService{}
    s.books = sampleBooks()
    return s
}

func (s *ContentService) ListEditorPicks() []models.Book {
    s.mu.RLock()
    defer s.mu.RUnlock()
    res := make([]models.Book, 0, 5)
    for _, b := range s.books {
        if b.Status == "published" && b.IsEditorPick {
            res = append(res, b)
        }
        if len(res) >= 5 {
            break
        }
    }
    return res
}

func (s *ContentService) ListBooks(age int, sort string, page, pageSize int) ([]models.Book, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    items := make([]models.Book, 0, len(s.books))
    for _, b := range s.books {
        if b.Status == "published" {
            items = append(items, b)
        }
    }
    sortBooks(items, age, sort)
    start := (page - 1) * pageSize
    if start < 0 {
        start = 0
    }
    if start >= len(items) {
        return []models.Book{}, false
    }
    end := start + pageSize
    if end > len(items) {
        end = len(items)
    }
    hasMore := end < len(items)
    return items[start:end], hasMore
}

func (s *ContentService) GetBook(id string) (models.Book, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, b := range s.books {
        if b.ID == id && b.Status == "published" {
            return b, true
        }
    }
    return models.Book{}, false
}

func (s *ContentService) AddBook(b models.Book) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.books = append(s.books, b)
}

func (s *ContentService) AddPages(bookID string, fileURLs []string) []models.Page {
    s.mu.Lock()
    defer s.mu.Unlock()
    var created []models.Page
    for i := range s.books {
        if s.books[i].ID == bookID {
            start := len(s.books[i].Pages)
            for idx, url := range fileURLs {
                p := models.Page{ID: bookID + "-p-" + itoa(start+idx), Index: start + idx, ImageURL: url, Duration: nil}
                s.books[i].Pages = append(s.books[i].Pages, p)
                created = append(created, p)
            }
            return created
        }
    }
    return created
}

func (s *ContentService) ReorderPages(bookID string, indexByID map[string]int) []models.Page {
    s.mu.Lock()
    defer s.mu.Unlock()
    for i := range s.books {
        if s.books[i].ID == bookID {
            for j := range s.books[i].Pages {
                id := s.books[i].Pages[j].ID
                if newIdx, ok := indexByID[id]; ok {
                    s.books[i].Pages[j].Index = newIdx
                }
            }
            pages := s.books[i].Pages
            sortPages(pages)
            s.books[i].Pages = pages
            return pages
        }
    }
    return nil
}
package repository

import (
    "database/sql"
    "encoding/json"
    "errors"
    "picturebook/server/internal/models"
)

type Repo struct { DB *sql.DB }

func New(db *sql.DB) *Repo { return &Repo{DB: db} }

func (r *Repo) SeedIfEmpty(sample []models.Book) error {
    var cnt int
    if err := r.DB.QueryRow("SELECT COUNT(1) FROM books").Scan(&cnt); err != nil { return err }
    if cnt > 0 { return nil }
    for _, b := range sample {
        if err := r.CreateBook(b); err != nil { return err }
        if len(b.Pages) > 0 { _ = r.AddPages(b.ID, b.Pages) }
    }
    return nil
}

func (r *Repo) CreateBook(b models.Book) error {
    tags, _ := json.Marshal(b.Tags)
    themes, _ := json.Marshal(b.ThemeKeywords)
    if b.CategoryID == "" {
        return errors.New("missing_category")
    }
    var exists int
    _ = r.DB.QueryRow("SELECT COUNT(1) FROM categories WHERE id=?", b.CategoryID).Scan(&exists)
    if exists == 0 { return errors.New("missing_category") }
    _, err := r.DB.Exec(`INSERT INTO books(id,title,cover_url,age_min,age_max,tags,popularity_score,theme_keywords,is_editor_pick,status,category_id) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
        b.ID, b.Title, b.CoverURL, b.AgeMin, b.AgeMax, string(tags), b.PopularityScore, string(themes), boolToInt(b.IsEditorPick), b.Status, b.CategoryID,
    )
    return err
}

func (r *Repo) SetEditorPick(id string, v bool) error {
    if v {
        var cnt int
        if err := r.DB.QueryRow("SELECT COUNT(1) FROM books WHERE is_editor_pick=1 AND status='published'").Scan(&cnt); err != nil { return err }
        if cnt >= 5 { return errors.New("editor_picks_limit") }
    }
    _, err := r.DB.Exec("UPDATE books SET is_editor_pick=? WHERE id=?", boolToInt(v), id)
    return err
}

func (r *Repo) ListEditorPicks() ([]models.Book, error) {
    rows, err := r.DB.Query("SELECT id,title,cover_url,age_min,age_max,tags,popularity_score,theme_keywords,is_editor_pick,status,category_id FROM books WHERE status='published' AND is_editor_pick=1 LIMIT 5")
    if err != nil { return nil, err }
    defer rows.Close()
    var res []models.Book
    for rows.Next() {
        var b models.Book
        var tags, themes string
        var pick int
        if err := rows.Scan(&b.ID, &b.Title, &b.CoverURL, &b.AgeMin, &b.AgeMax, &tags, &b.PopularityScore, &themes, &pick, &b.Status, &b.CategoryID); err != nil { return nil, err }
        _ = json.Unmarshal([]byte(tags), &b.Tags)
        _ = json.Unmarshal([]byte(themes), &b.ThemeKeywords)
        b.IsEditorPick = pick == 1
        res = append(res, b)
    }
    return res, nil
}

func (r *Repo) ListBooks(age int, sortBy string, page, pageSize int) ([]models.Book, bool, error) {
    rows, err := r.DB.Query("SELECT id,title,cover_url,age_min,age_max,tags,popularity_score,theme_keywords,is_editor_pick,status,category_id FROM books WHERE status='published' ORDER BY popularity_score DESC LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize)
    if err != nil { return nil, false, err }
    defer rows.Close()
    var res []models.Book
    for rows.Next() {
        var b models.Book
        var tags, themes string
        var pick int
        if err := rows.Scan(&b.ID, &b.Title, &b.CoverURL, &b.AgeMin, &b.AgeMax, &tags, &b.PopularityScore, &themes, &pick, &b.Status, &b.CategoryID); err != nil { return nil, false, err }
        _ = json.Unmarshal([]byte(tags), &b.Tags)
        _ = json.Unmarshal([]byte(themes), &b.ThemeKeywords)
        b.IsEditorPick = pick == 1
        res = append(res, b)
    }
    var total int
    _ = r.DB.QueryRow("SELECT COUNT(1) FROM books WHERE status='published'").Scan(&total)
    hasMore := page*pageSize < total
    return res, hasMore, nil
}

func (r *Repo) GetBook(id string) (models.Book, bool, error) {
    var b models.Book
    var tags, themes string
    var pick int
    err := r.DB.QueryRow("SELECT id,title,cover_url,age_min,age_max,tags,popularity_score,theme_keywords,is_editor_pick,status,category_id FROM books WHERE id=? AND status='published'", id).
        Scan(&b.ID, &b.Title, &b.CoverURL, &b.AgeMin, &b.AgeMax, &tags, &b.PopularityScore, &themes, &pick, &b.Status, &b.CategoryID)
    if err != nil { return models.Book{}, false, nil }
    _ = json.Unmarshal([]byte(tags), &b.Tags)
    _ = json.Unmarshal([]byte(themes), &b.ThemeKeywords)
    b.IsEditorPick = pick == 1
    // pages
    rows, err := r.DB.Query("SELECT id,idx,image_url,duration_hint FROM pages WHERE book_id=? ORDER BY idx ASC", id)
    if err == nil {
        for rows.Next() {
            var p models.Page
            var dur sql.NullInt64
            if err := rows.Scan(&p.ID, &p.Index, &p.ImageURL, &dur); err == nil {
                if dur.Valid { d := int(dur.Int64); p.Duration = &d }
                b.Pages = append(b.Pages, p)
            }
        }
        rows.Close()
    }
    return b, true, nil
}

func (r *Repo) AddPages(bookID string, pages []models.Page) error {
    for _, p := range pages {
        var dur *int = p.Duration
        var d any
        if dur != nil { d = *dur } else { d = nil }
        if _, err := r.DB.Exec("INSERT INTO pages(id,book_id,idx,image_url,duration_hint) VALUES(?,?,?,?,?)", p.ID, bookID, p.Index, p.ImageURL, d); err != nil { return err }
    }
    return nil
}

func (r *Repo) ReorderPages(bookID string, indexByID map[string]int) error {
    for id, idx := range indexByID {
        if _, err := r.DB.Exec("UPDATE pages SET idx=? WHERE id=? AND book_id=?", idx, id, bookID); err != nil { return err }
    }
    return nil
}

func (r *Repo) DeleteBook(id string) error {
    _, err := r.DB.Exec("DELETE FROM books WHERE id=?", id)
    return err
}

func (r *Repo) UpdateBook(b models.Book) error {
    tags, _ := json.Marshal(b.Tags)
    themes, _ := json.Marshal(b.ThemeKeywords)
    _, err := r.DB.Exec("UPDATE books SET title=?, cover_url=?, age_min=?, age_max=?, tags=?, popularity_score=?, theme_keywords=?, status=?, category_id=? WHERE id=?",
        b.Title, b.CoverURL, b.AgeMin, b.AgeMax, string(tags), b.PopularityScore, string(themes), b.Status, b.CategoryID, b.ID,
    )
    return err
}

func (r *Repo) ListBooksAdmin(sortBy string, page, pageSize int, q string) ([]models.Book, bool, error) {
    base := "SELECT id,title,cover_url,age_min,age_max,tags,popularity_score,theme_keywords,is_editor_pick,status,category_id FROM books"
    where := ""
    args := []any{}
    if q != "" {
        where = " WHERE title LIKE ? OR id LIKE ?"
        like := "%" + q + "%"
        args = append(args, like, like)
    }
    order := " ORDER BY popularity_score DESC"
    if sortBy == "title" { order = " ORDER BY title ASC" }
    limit := " LIMIT ? OFFSET ?"
    args = append(args, pageSize, (page-1)*pageSize)
    rows, err := r.DB.Query(base+where+order+limit, args...)
    if err != nil { return nil, false, err }
    defer rows.Close()
    var res []models.Book
    for rows.Next() {
        var b models.Book
        var tags, themes string
        var pick int
        if err := rows.Scan(&b.ID, &b.Title, &b.CoverURL, &b.AgeMin, &b.AgeMax, &tags, &b.PopularityScore, &themes, &pick, &b.Status, &b.CategoryID); err != nil { return nil, false, err }
        _ = json.Unmarshal([]byte(tags), &b.Tags)
        _ = json.Unmarshal([]byte(themes), &b.ThemeKeywords)
        b.IsEditorPick = pick == 1
        res = append(res, b)
    }
    var total int
    if q != "" {
        _ = r.DB.QueryRow("SELECT COUNT(1) FROM books WHERE title LIKE ? OR id LIKE ?", "%"+q+"%", "%"+q+"%").Scan(&total)
    } else {
        _ = r.DB.QueryRow("SELECT COUNT(1) FROM books").Scan(&total)
    }
    hasMore := page*pageSize < total
    return res, hasMore, nil
}

// Categories CRUD and aggregations
func (r *Repo) CreateCategory(c models.Category) error {
    _, err := r.DB.Exec("INSERT INTO categories(id,name,description) VALUES(?,?,?)", c.ID, c.Name, c.Description)
    return err
}

func (r *Repo) DeleteCategory(id string) error {
    var cnt int
    if err := r.DB.QueryRow("SELECT COUNT(1) FROM books WHERE category_id=?", id).Scan(&cnt); err != nil { return err }
    if cnt > 0 { return errors.New("category_has_books") }
    _, err := r.DB.Exec("DELETE FROM categories WHERE id=?", id)
    return err
}

func (r *Repo) ListCategories() ([]models.Category, error) {
    rows, err := r.DB.Query("SELECT id,name,description FROM categories ORDER BY name ASC")
    if err != nil { return nil, err }
    defer rows.Close()
    var res []models.Category
    for rows.Next() {
        var c models.Category
        if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil { return nil, err }
        res = append(res, c)
    }
    return res, nil
}

func (r *Repo) ListCategoriesWithBooks() ([]models.CategoryWithBooks, error) {
    cats, err := r.ListCategories()
    if err != nil { return nil, err }
    var out []models.CategoryWithBooks
    for _, c := range cats {
        rows, err := r.DB.Query("SELECT id,title,cover_url,age_min,age_max,tags,popularity_score,theme_keywords,is_editor_pick,status,category_id FROM books WHERE status='published' AND category_id=? ORDER BY popularity_score DESC", c.ID)
        if err != nil { return nil, err }
        var books []models.Book
        for rows.Next() {
            var b models.Book
            var tags, themes string
            var pick int
            if err := rows.Scan(&b.ID, &b.Title, &b.CoverURL, &b.AgeMin, &b.AgeMax, &tags, &b.PopularityScore, &themes, &pick, &b.Status, &b.CategoryID); err != nil { rows.Close(); return nil, err }
            _ = json.Unmarshal([]byte(tags), &b.Tags)
            _ = json.Unmarshal([]byte(themes), &b.ThemeKeywords)
            b.IsEditorPick = pick == 1
            b.CategoryName = c.Name
            books = append(books, b)
        }
        rows.Close()
        out = append(out, models.CategoryWithBooks{ID: c.ID, Name: c.Name, Description: c.Description, Books: books})
    }
    return out, nil
}

func boolToInt(b bool) int { if b { return 1 } ; return 0 }
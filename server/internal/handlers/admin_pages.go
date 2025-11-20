package handlers

import (
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "github.com/gin-gonic/gin"
    "picturebook/server/internal/models"
    "picturebook/server/internal/services"
)

type AdminPages struct {
    content *services.ContentService
}

func NewAdminPages(content *services.ContentService) *AdminPages { return &AdminPages{content: content} }

func (h *AdminPages) LoginPage(c *gin.Context) { c.HTML(http.StatusOK, "login.html", gin.H{}) }

func (h *AdminPages) Login(c *gin.Context) {
    email := c.PostForm("email")
    password := c.PostForm("password")
    if email != "admin@example.com" || password == "" {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "invalid"})
        return
    }
    c.SetCookie("admin_session", "ok", 3600, "/", "", false, true)
    c.Redirect(http.StatusFound, "/admin/books")
}

func (h *AdminPages) BooksList(c *gin.Context) {
    page := atoiDefault(c.Query("page"), 1)
    size := atoiDefault(c.Query("page_size"), 20)
    sort := c.DefaultQuery("sort", "popular")
    q := c.Query("q")
    items, more := h.content.ListBooksAdmin(sort, page, size, q)
    c.HTML(http.StatusOK, "books_list.html", gin.H{"items": items, "page": page, "page_size": size, "has_more": more, "sort": sort, "q": q})
}

func (h *AdminPages) NewBookPage(c *gin.Context) {
    cats, _ := h.content.ListCategories()
    c.HTML(http.StatusOK, "book_new.html", gin.H{"categories": cats})
}

func (h *AdminPages) CreateBook(c *gin.Context) {
    ageMin := atoiDefault(c.PostForm("ageMin"), 3)
    ageMax := atoiDefault(c.PostForm("ageMax"), 10)
    pop := atofDefault(c.PostForm("popularityScore"), 0.5)
    tags := splitCSV(c.PostForm("tags"))
    themes := splitCSV(c.PostForm("themeKeywords"))
    b := models.Book{
        ID: c.PostForm("id"), Title: c.PostForm("title"), CoverURL: c.PostForm("coverURL"), AgeMin: ageMin, AgeMax: ageMax,
        Tags: tags, PopularityScore: pop, ThemeKeywords: themes, IsEditorPick: c.PostForm("isEditorPick") == "on", Status: "published", Pages: []models.Page{},
    }
    b.CategoryID = c.PostForm("category_id")
    if b.CategoryID == "" {
        cats, _ := h.content.ListCategories()
        c.HTML(http.StatusBadRequest, "book_new.html", gin.H{"error": "missing_category", "categories": cats})
        return
    }
    picks := h.content.ListEditorPicks()
    if b.IsEditorPick && len(picks) >= 5 {
        cats, _ := h.content.ListCategories()
        c.HTML(http.StatusBadRequest, "book_new.html", gin.H{"error": "limit", "categories": cats})
        return
    }
    if err := h.content.AddBook(b); err != nil {
        cats, _ := h.content.ListCategories()
        c.HTML(http.StatusBadRequest, "book_new.html", gin.H{"error": "db_error", "categories": cats, "errMsg": err.Error()})
        return
    }
    c.Redirect(http.StatusFound, "/admin/books")
}

func (h *AdminPages) PagesOfBook(c *gin.Context) {
    id := c.Param("id")
    b, ok := h.content.GetBook(id)
    if !ok { c.String(http.StatusNotFound, "not_found"); return }
    c.HTML(http.StatusOK, "pages.html", gin.H{"book": b})
}

func (h *AdminPages) ToggleEditorPick(c *gin.Context) {
    id := c.Param("id")
    v := c.PostForm("isEditorPick") == "on"
    if err := h.content.SetEditorPick(id, v); err != nil {
        c.HTML(http.StatusBadRequest, "pages.html", gin.H{"error": "editor_picks_limit"})
        return
    }
    c.Redirect(http.StatusFound, "/admin/books/"+id+"/pages")
}

func (h *AdminPages) DeleteBook(c *gin.Context) {
    id := c.Param("id")
    b, ok := h.content.GetBook(id)
    if ok {
        dir := uploadsDir()
        if b.CoverURL != "" && strings.HasPrefix(b.CoverURL, "/assets/") {
            name := strings.TrimPrefix(b.CoverURL, "/assets/")
            _ = os.Remove(filepath.Join(dir, name))
        }
        for _, p := range b.Pages {
            if p.ImageURL != "" && strings.HasPrefix(p.ImageURL, "/assets/") {
                name := strings.TrimPrefix(p.ImageURL, "/assets/")
                _ = os.Remove(filepath.Join(dir, name))
            }
        }
    }
    _ = h.content.DeleteBook(id)
    c.Redirect(http.StatusFound, "/admin/books")
}

// Category pages
func (h *AdminPages) CategoriesList(c *gin.Context) {
    cats, _ := h.content.ListCategories()
    c.HTML(http.StatusOK, "categories_list.html", gin.H{"categories": cats})
}

func (h *AdminPages) NewCategoryPage(c *gin.Context) { c.HTML(http.StatusOK, "category_new.html", gin.H{}) }

func (h *AdminPages) CreateCategory(c *gin.Context) {
    id := c.PostForm("id")
    name := c.PostForm("name")
    desc := c.PostForm("description")
    if id == "" || name == "" { c.HTML(http.StatusBadRequest, "category_new.html", gin.H{"error":"invalid"}); return }
    _ = h.content.CreateCategory(models.Category{ID: id, Name: name, Description: desc})
    c.Redirect(http.StatusFound, "/admin/categories")
}

func (h *AdminPages) DeleteCategory(c *gin.Context) {
    id := c.Param("id")
    if err := h.content.DeleteCategory(id); err != nil {
        cats, _ := h.content.ListCategories()
        c.HTML(http.StatusBadRequest, "categories_list.html", gin.H{"categories": cats, "error": "category_has_books"})
        return
    }
    c.Redirect(http.StatusFound, "/admin/categories")
}

func (h *AdminPages) EditBookPage(c *gin.Context) {
    id := c.Param("id")
    b, ok := h.content.GetBook(id)
    if !ok { c.String(http.StatusNotFound, "not_found"); return }
    cats, _ := h.content.ListCategories()
    c.HTML(http.StatusOK, "book_edit.html", gin.H{"book": b, "categories": cats})
}

func (h *AdminPages) EditBook(c *gin.Context) {
    id := c.Param("id")
    b, ok := h.content.GetBook(id)
    if !ok { c.String(http.StatusNotFound, "not_found"); return }
    title := c.PostForm("title")
    cover := c.PostForm("coverURL")
    ageMin := atoiDefault(c.PostForm("ageMin"), b.AgeMin)
    ageMax := atoiDefault(c.PostForm("ageMax"), b.AgeMax)
    pop := atofDefault(c.PostForm("popularityScore"), b.PopularityScore)
    tags := splitCSV(c.PostForm("tags"))
    themes := splitCSV(c.PostForm("themeKeywords"))
    status := c.PostForm("status")
    if status == "" { status = b.Status }
    cat := c.PostForm("category_id")
    nb := models.Book{ID: b.ID, Title: title, CoverURL: cover, AgeMin: ageMin, AgeMax: ageMax, Tags: tags, PopularityScore: pop, ThemeKeywords: themes, IsEditorPick: b.IsEditorPick, Pages: b.Pages, Status: status, CategoryID: cat}
    _ = h.content.UpdateBook(nb)
    c.Redirect(http.StatusFound, "/admin/books/"+id+"/pages")
}
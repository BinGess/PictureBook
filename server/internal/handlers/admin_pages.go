package handlers

import (
    "net/http"
    "strconv"
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
    items, _ := h.content.ListBooks(6, "popular", 1, 1000)
    c.HTML(http.StatusOK, "books_list.html", gin.H{"items": items})
}

func (h *AdminPages) NewBookPage(c *gin.Context) { c.HTML(http.StatusOK, "book_new.html", gin.H{}) }

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
    picks := h.content.ListEditorPicks()
    if b.IsEditorPick && len(picks) >= 5 {
        c.HTML(http.StatusBadRequest, "book_new.html", gin.H{"error": "limit"})
        return
    }
    h.content.AddBook(b)
    c.Redirect(http.StatusFound, "/admin/books")
}

func (h *AdminPages) PagesOfBook(c *gin.Context) {
    id := c.Param("id")
    b, ok := h.content.GetBook(id)
    if !ok { c.String(http.StatusNotFound, "not_found"); return }
    c.HTML(http.StatusOK, "pages.html", gin.H{"book": b})
}
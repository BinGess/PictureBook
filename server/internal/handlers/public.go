package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "picturebook/server/internal/services"
)

type PublicHandler struct {
    content *services.ContentService
}

func NewPublicHandler(content *services.ContentService) *PublicHandler {
    return &PublicHandler{content: content}
}

func (h *PublicHandler) EditorPicks(c *gin.Context) {
    res := h.content.ListEditorPicks()
    c.JSON(http.StatusOK, res)
}

func (h *PublicHandler) Books(c *gin.Context) {
    age := atoi(c.Query("age"))
    sortBy := c.DefaultQuery("sort", "popular")
    page := atoiDefault(c.Query("page"), 1)
    pageSize := atoiDefault(c.Query("page_size"), 24)
    items, hasMore := h.content.ListBooks(age, sortBy, page, pageSize)
    c.JSON(http.StatusOK, gin.H{
        "items": items,
        "paging": gin.H{
            "page": page,
            "page_size": pageSize,
            "has_more": hasMore,
        },
    })
}

func (h *PublicHandler) Book(c *gin.Context) {
    id := c.Param("id")
    b, ok := h.content.GetBook(id)
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
        return
    }
    c.JSON(http.StatusOK, b)
}

func (h *PublicHandler) Recommend(c *gin.Context) {
    id := c.Param("id")
    b, ok := h.content.GetBook(id)
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
        return
    }
    age := atoiDefault(c.Query("age"), (b.AgeMin+b.AgeMax)/2)
    limit := atoiDefault(c.Query("limit"), 5)
    all, _ := h.content.ListBooks(age, "combined", 1, 1000)
    rec := services.Recommend(b, all, age, limit)
    c.JSON(http.StatusOK, rec)
}
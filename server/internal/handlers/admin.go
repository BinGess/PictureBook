package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "picturebook/server/internal/models"
    "picturebook/server/internal/services"
)

type AdminHandler struct {
    content *services.ContentService
}

func NewAdminHandler(content *services.ContentService) *AdminHandler {
    return &AdminHandler{content: content}
}

func (h *AdminHandler) CreateBook(c *gin.Context) {
    var b models.Book
    if err := c.ShouldBindJSON(&b); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
        return
    }
    if b.IsEditorPick {
        picks := h.content.ListEditorPicks()
        if len(picks) >= 5 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "editor_picks_limit"})
            return
        }
    }
    if b.Status == "" { b.Status = "draft" }
    h.content.AddBook(b)
    c.JSON(http.StatusOK, b)
}
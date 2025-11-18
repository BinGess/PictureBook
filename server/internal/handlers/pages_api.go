package handlers

import (
    "net/http"
    "path/filepath"
    "github.com/gin-gonic/gin"
    "picturebook/server/internal/services"
)

type PagesAPI struct {
    content *services.ContentService
}

func NewPagesAPI(content *services.ContentService) *PagesAPI { return &PagesAPI{content: content} }

func (h *PagesAPI) UploadToBook(c *gin.Context) {
    bookID := c.Param("id")
    form, err := c.MultipartForm()
    if err != nil || form == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing_files"})
        return
    }
    files := form.File["files"]
    if len(files) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing_files"})
        return
    }
    dir := "uploads"
    urls := make([]string, 0, len(files))
    for _, f := range files {
        name := filepath.Base(f.Filename)
        dst := filepath.Join(dir, name)
        if err := c.SaveUploadedFile(f, dst); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed"})
            return
        }
        urls = append(urls, "/assets/"+name)
    }
    created := h.content.AddPages(bookID, urls)
    c.JSON(http.StatusOK, gin.H{"pages": created})
}

func (h *PagesAPI) Reorder(c *gin.Context) {
    bookID := c.Param("id")
    m := map[string]int{}
    for k, v := range c.Request.PostForm {
        if len(k) > 7 && k[:7] == "index[" {
            id := k[7 : len(k)-1]
            if len(v) > 0 {
                if n, err := strconv.Atoi(v[0]); err == nil { m[id] = n }
            }
        }
    }
    pages := h.content.ReorderPages(bookID, m)
    if pages == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
        return
    }
    c.Redirect(http.StatusFound, "/admin/books/"+bookID+"/pages")
}
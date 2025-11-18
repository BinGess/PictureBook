package handlers

import (
    "net/http"
    "path/filepath"
    "github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing_file"})
        return
    }
    dir := "uploads"
    name := filepath.Base(file.Filename)
    dst := filepath.Join(dir, name)
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"file_url": "/assets/" + name})
}
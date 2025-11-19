package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "picturebook/server/internal/handlers"
    "picturebook/server/internal/middleware"
    "picturebook/server/internal/services"
)

func main() {
    allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
    if allowed == "" { allowed = "*" }
    content := NewContent()
    r := gin.Default()
    r.Use(middleware.CORS(allowed))
    r.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })
    r.Static("/assets", "uploads")
    r.LoadHTMLGlob("web/admin/*.html")
    pub := handlers.NewPublicHandler(content)
    v1 := r.Group("/v1")
    v1.GET("/editor-picks", pub.EditorPicks)
    v1.GET("/books", pub.Books)
    v1.GET("/books/:id", pub.Book)
    v1.GET("/books/:id/recommendations", pub.Recommend)
    admin := handlers.NewAdminHandler(content)
    v1.POST("/admin/books", admin.CreateBook)
    v1.POST("/admin/upload", handlers.Upload)
    pa := handlers.NewPagesAPI(content)
    v1.POST("/admin/books/:id/pages/upload", pa.UploadToBook)
    ap := handlers.NewAdminPages(content)
    r.GET("/admin/login", ap.LoginPage)
    r.POST("/admin/login", ap.Login)
    adminGroup := r.Group("/admin")
    adminGroup.Use(middleware.AdminAuth())
    adminGroup.GET("/books", ap.BooksList)
    adminGroup.GET("/books/new", ap.NewBookPage)
    adminGroup.POST("/books/new", ap.CreateBook)
    adminGroup.GET("/books/:id/pages", ap.PagesOfBook)
    adminGroup.POST("/books/:id/pages/reorder", pa.Reorder)
    adminGroup.POST("/books/:id/editor-pick", ap.ToggleEditorPick)
    r.Run()
}

func NewContent() *services.ContentService { return services.NewContentService() }
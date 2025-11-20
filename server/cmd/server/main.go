package main

import (
    "os"
    "path/filepath"
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
    up := os.Getenv("UPLOADS_DIR")
    if up == "" { up = "/app/uploads" }
    _ = os.MkdirAll(filepath.Clean(up), 0755)
    r.Static("/assets", up)
    r.Static("/static", "web/static")
    dbp := os.Getenv("SQLITE_PATH")
    if dbp == "" { dbp = "/app/data/picturebook.db" }
    _ = os.MkdirAll(filepath.Dir(filepath.Clean(dbp)), 0755)
    r.LoadHTMLGlob("web/admin/*.html")
    pub := handlers.NewPublicHandler(content)
    v1 := r.Group("/v1")
    v1.GET("/editor-picks", pub.EditorPicks)
    v1.GET("/books", pub.Books)
    v1.GET("/books/:id", pub.Book)
    v1.GET("/books/:id/recommendations", pub.Recommend)
    v1.GET("/categories", pub.Categories)
    v1.GET("/categories-with-books", pub.CategoriesWithBooks)
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
    adminGroup.GET("/books/:id/edit", ap.EditBookPage)
    adminGroup.POST("/books/:id/edit", ap.EditBook)
    adminGroup.GET("/books/:id/pages", ap.PagesOfBook)
    adminGroup.POST("/books/:id/pages/reorder", pa.Reorder)
    adminGroup.POST("/books/:id/editor-pick", ap.ToggleEditorPick)
    adminGroup.POST("/books/:id/delete", ap.DeleteBook)
    adminGroup.GET("/categories", ap.CategoriesList)
    adminGroup.GET("/categories/new", ap.NewCategoryPage)
    adminGroup.POST("/categories/new", ap.CreateCategory)
    adminGroup.POST("/categories/:id/delete", ap.DeleteCategory)
    r.Run()
}

func NewContent() *services.ContentService { return services.NewContentService() }
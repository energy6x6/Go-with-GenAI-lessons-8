package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "libratrack/docs"
	"libratrack/internal/handler"
	"libratrack/internal/middleware"
	"libratrack/internal/repository"
	"net/http"
)

func NewRouter(repos repository.Repositories) http.Handler {
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true
	r.Use(middleware.RecoveryMiddleware(), middleware.LoggingMiddleware(), middleware.CORSMiddleware())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	b := handler.NewBookHandler(repos.Books)
	a := handler.NewAuthorHandler(repos.Authors)
	cat := handler.NewCategoryHandler(repos.Categories)
	copies := handler.NewCopyHandler(repos.Copies)
	api := r.Group("/api/v1")
	api.POST("/books", b.Create)
	api.GET("/books", b.List)
	api.GET("/books/:id", b.Get)
	api.PUT("/books/:id", b.Update)
	api.DELETE("/books/:id", b.Delete)
	api.POST("/authors", a.Create)
	api.GET("/authors", a.List)
	api.GET("/authors/:id", a.Get)
	api.POST("/categories", cat.Create)
	api.GET("/categories", cat.List)
	api.GET("/categories/:id", cat.Get)
	api.POST("/books/:id/copies", copies.Create)
	api.GET("/books/:id/copies", copies.List)
	api.GET("/books/:id/availability", copies.Availability)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(func(c *gin.Context) { handler.Error(c, 404, "not_found", "Route not found") })
	r.NoMethod(func(c *gin.Context) { handler.Error(c, 405, "method_not_allowed", "Method not allowed") })
	return r
}

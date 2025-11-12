package router

import (
	"net/http"
	"time"

	"gogogo/config"
	"gogogo/controllers"
	"gogogo/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	server := gin.Default()

	corsConfig := cors.Config{
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if len(config.AppConfig.CORS.AllowOrigins) == 0 {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = config.AppConfig.CORS.AllowOrigins
	}

	server.Use(cors.New(corsConfig))

	api := server.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/login", controllers.Login)
	auth.POST("/register", controllers.Register)

	api.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", controllers.GetProfile)
		protected.GET("/me/posts", controllers.ListMyPosts)

		protected.POST("/posts", controllers.CreatePost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)

		protected.POST("/categories", controllers.CreateCategory)
		protected.PUT("/categories/:id", controllers.UpdateCategory)
		protected.DELETE("/categories/:id", controllers.DeleteCategory)

		protected.POST("/tags", controllers.CreateTag)
		protected.PUT("/tags/:id", controllers.UpdateTag)
		protected.DELETE("/tags/:id", controllers.DeleteTag)
	}

	api.GET("/posts", controllers.ListPosts)
	api.GET("/posts/slug/:slug", controllers.GetPostBySlug)
	api.GET("/posts/:id", controllers.GetPostByID)
	api.GET("/posts/:id/comments", controllers.ListComments)
	api.POST("/posts/:id/comments", controllers.CreateComment)

	api.GET("/categories", controllers.ListCategories)
	api.GET("/categories/:id/posts", controllers.ListPostsByCategory)
	api.GET("/tags", controllers.ListTags)
	api.GET("/tags/:slug/posts", controllers.ListPostsByTag)

	return server
}

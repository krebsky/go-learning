package main

import (
	"blog/database"
	"blog/handlers"
	"blog/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDB()

	r := gin.Default()

	// 公开接口 - 不需要认证
	api := r.Group("/api")
	{
		// 用户认证
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		// 文章公开接口
		api.GET("/posts", handlers.GetPosts)
		api.GET("/posts/:id", handlers.GetPost)

		// 评论公开接口（使用 :id 作为 postId）
		api.GET("/posts/:id/comments", handlers.GetComments)
	}

	// 需要认证的接口
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		// 文章管理
		auth.POST("/posts", handlers.CreatePost)
		auth.PUT("/posts/:id", handlers.UpdatePost)
		auth.DELETE("/posts/:id", handlers.DeletePost)

		// 评论管理（使用 :id 作为 postId）
		auth.POST("/posts/:id/comments", handlers.CreateComment)
	}

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

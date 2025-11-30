package handlers

import (
	"blog/database"
	"blog/middleware"
	"blog/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePostRequest 创建文章请求结构
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UpdatePostRequest 更新文章请求结构
type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("CreatePost validation error: %v", err)
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		log.Printf("Post creation error: %v", err)
		return
	}

	// 加载用户信息
	database.DB.Preload("User").First(&post, post.ID)

	log.Printf("Post created successfully: ID=%d, UserID=%d", post.ID, userID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

// GetPosts 获取所有文章列表
func GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := database.DB.Preload("User").Preload("Comments.User").Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		log.Printf("GetPosts error: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"count": len(posts),
	})
}

// GetPost 获取单个文章详情
func GetPost(c *gin.Context) {
	postID := c.Param("id")

	var post models.Post
	if err := database.DB.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		log.Printf("GetPost error: post ID %s not found", postID)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID := c.Param("id")

	// 查找文章
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		log.Printf("UpdatePost error: post ID %s not found", postID)
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		log.Printf("UpdatePost error: user %d tried to update post %d owned by user %d", userID, post.ID, post.UserID)
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("UpdatePost validation error: %v", err)
		return
	}

	// 更新文章
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		log.Printf("Post update error: %v", err)
		return
	}

	// 重新加载用户信息
	database.DB.Preload("User").First(&post, post.ID)

	log.Printf("Post updated successfully: ID=%d", post.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    post,
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID := c.Param("id")

	// 查找文章
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		log.Printf("DeletePost error: post ID %s not found", postID)
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		log.Printf("DeletePost error: user %d tried to delete post %d owned by user %d", userID, post.ID, post.UserID)
		return
	}

	// 删除文章（级联删除评论）
	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		log.Printf("Post deletion error: %v", err)
		return
	}

	log.Printf("Post deleted successfully: ID=%d", post.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}

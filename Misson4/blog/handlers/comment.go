package handlers

import (
	"blog/database"
	"blog/middleware"
	"blog/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID := c.Param("id")

	// 检查文章是否存在
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		log.Printf("CreateComment error: post ID %s not found", postID)
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("CreateComment validation error: %v", err)
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  post.ID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		log.Printf("Comment creation error: %v", err)
		return
	}

	// 加载用户信息
	database.DB.Preload("User").First(&comment, comment.ID)

	log.Printf("Comment created successfully: ID=%d, PostID=%d, UserID=%d", comment.ID, post.ID, userID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetComments 获取文章的所有评论
func GetComments(c *gin.Context) {
	postID := c.Param("id")

	// 检查文章是否存在
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		log.Printf("GetComments error: post ID %s not found", postID)
		return
	}

	var comments []models.Comment
	if err := database.DB.Preload("User").Where("post_id = ?", postID).Order("created_at desc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		log.Printf("GetComments error: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"count":    len(comments),
	})
}

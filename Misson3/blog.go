package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"type:varchar(50)"`
	Email     string `gorm:"type:varchar(100)"`
	PostCount int    `gorm:"default:0"` // 文章数量统计字段
	Posts     []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Post struct {
	ID            uint      `gorm:"primaryKey;autoIncrement`
	Title         string    `gorm:"type:varchar(200)"`
	Content       string    `gorm:"type:text"`
	UserID        uint      `gorm:"not null;index"`
	CommentCount  int       `gorm:"default:0"`
	CommentStatus string    `gorm:"type:varchar(20);default:'有评论'"`
	User          User      `gorm:"foreignKey:UserID"`
	Comments      []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Content string `gorm:"type:text;not null"`
	PostID  uint   `gorm:"not null;index"`
	Post    Post   `gorm:"foreignKey:PostID"`
}

func connectDB() (*gorm.DB, error) {
	// 数据库连接字符串
	dsn := "root:Zhaoyang@100297@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %v", err)
	}

	return db, nil
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 更新用户的文章数量
	result := tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1"))

	if result.Error != nil {
		return fmt.Errorf("更新用户文章数量失败: %v", result.Error)
	}

	fmt.Printf("已更新用户 %d 的文章数量\n", p.UserID)
	return nil
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 查询该文章的评论数量
	var count int64
	tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count)

	// 如果评论数量为 0，更新文章的评论状态
	if count == 0 {
		result := tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Updates(map[string]interface{}{
				"comment_count":  0,
				"comment_status": "无评论",
			})

		if result.Error != nil {
			return fmt.Errorf("更新文章评论状态失败: %v", result.Error)
		}

		fmt.Printf("文章 %d 的评论已全部删除，已更新评论状态为 '无评论'\n", c.PostID)
	} else {
		// 更新评论数量
		result := tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			UpdateColumn("comment_count", count)

		if result.Error != nil {
			return fmt.Errorf("更新文章评论数量失败: %v", result.Error)
		}

		fmt.Printf("已更新文章 %d 的评论数量为 %d\n", c.PostID, count)
	}

	return nil
}

func (c *Comment) AfterCreate(tx *gorm.DB) error {
	// 更新文章的评论数量
	result := tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	if result.Error != nil {
		return fmt.Errorf("更新文章评论数量失败: %v", result.Error)
	}

	// 更新评论状态
	result = tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_status", "有评论")

	if result.Error != nil {
		return fmt.Errorf("更新文章评论状态失败: %v", result.Error)
	}

	fmt.Printf("已更新文章 %d 的评论数量\n", c.PostID)
	return nil
}

// 题目2-1：查询某个用户发布的所有文章及其对应的评论信息
func getUserPostsWithComments(db *gorm.DB, userID uint) (User, error) {
	var user User

	result := db.Preload("Posts").
		Preload("Posts.Comments").
		First(&user, userID)

	if result.Error != nil {
		return user, fmt.Errorf("查询失败: %v", result.Error)
	}

	return user, nil
}

func getPostWithMostComments(db *gorm.DB) (*Post, error) {
	var post Post

	result := db.Order("comment_count DESC").
		First(&post)

	if result.Error != nil {
		return nil, fmt.Errorf("查询失败: %v", result.Error)
	}

	return &post, nil
}

func main() {

	db, err := connectDB()
	if err != nil {
		fmt.Printf("连接数据库失败: %v", err)
		return
	}

	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	fmt.Println("\n创建测试数据...")

	// 创建用户
	user1 := User{
		Username: "张三",
		Email:    "zhangsan@example.com",
	}
	user2 := User{
		Username: "李四",
		Email:    "lisi@example.com",
	}
	db.Create(&user1)
	db.Create(&user2)
	fmt.Printf("创建用户: %s (ID: %d), %s (ID: %d)\n",
		user1.Username, user1.ID, user2.Username, user2.ID)

	// 创建文章（钩子函数会自动更新用户的文章数量）
	post1 := Post{
		Title:   "Go语言入门教程",
		Content: "这是一篇关于Go语言入门的详细教程，涵盖了基础语法、数据类型、函数等内容。",
		UserID:  user1.ID,
	}
	post2 := Post{
		Title:   "GORM使用指南",
		Content: "本文介绍了如何使用GORM进行数据库操作，包括模型定义、关联查询、钩子函数等。",
		UserID:  user1.ID,
	}
	post3 := Post{
		Title:   "数据库设计原则",
		Content: "良好的数据库设计是系统成功的关键，本文介绍了一些重要的设计原则。",
		UserID:  user2.ID,
	}
	db.Create(&post1)
	db.Create(&post2)
	db.Create(&post3)
	fmt.Println("创建文章完成（钩子函数已自动更新用户文章数量）")

	// 验证用户文章数量是否已更新
	var updatedUser1 User
	db.First(&updatedUser1, user1.ID)
	fmt.Printf("用户 %s 的文章数量已更新为: %d\n", updatedUser1.Username, updatedUser1.PostCount)

	// 创建评论（钩子函数会自动更新文章的评论数量和状态）
	comments := []Comment{
		{Content: "很好的教程，学到了很多！", PostID: post1.ID},
		{Content: "期待更多内容", PostID: post1.ID},
		{Content: "GORM确实很方便", PostID: post2.ID},
		{Content: "设计原则很重要", PostID: post3.ID},
		{Content: "赞同你的观点", PostID: post3.ID},
		{Content: "补充一点：还要考虑性能", PostID: post3.ID},
	}

	for _, comment := range comments {
		db.Create(&comment)
	}
	fmt.Println("创建评论完成（钩子函数已自动更新文章评论数量和状态）")

	// 题目2-1：查询某个用户发布的所有文章及其对应的评论信息
	fmt.Println("\n【题目2-1】查询用户发布的所有文章及其评论")
	userWithPosts, err := getUserPostsWithComments(db, user1.ID)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
	}
	fmt.Println(userWithPosts)

	// 题目2-2：查询评论数量最多的文章信息
	fmt.Println("\n【题目2-2】查询评论数量最多的文章")
	mostCommentedPost, err := getPostWithMostComments(db)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
	}
	fmt.Println(mostCommentedPost)

	// 测试删除评论的钩子函数
	fmt.Println("\n【测试】删除文章的评论，验证钩子函数")

	// 删除 post1 的所有评论
	var post1Comments []Comment
	db.Where("post_id = ?", post1.ID).Find(&post1Comments)

	for _, comment := range post1Comments {
		fmt.Printf("删除评论 ID: %d\n", comment.ID)
		db.Delete(&comment) // 触发 AfterDelete 钩子
	}

	// 验证 post1 的评论状态
	var updatedPost1 Post
	db.First(&updatedPost1, post1.ID)
	fmt.Printf("\n文章 '%s' 的评论数量: %d, 评论状态: %s\n",
		updatedPost1.Title, updatedPost1.CommentCount, updatedPost1.CommentStatus)

}

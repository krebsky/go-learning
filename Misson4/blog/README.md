# 个人博客系统后端

基于 Go 语言、Gin 框架和 GORM 开发的个人博客系统后端，支持用户认证、文章管理和评论功能。

## 功能特性

- 用户注册和登录（JWT 认证）
- 文章 CRUD 操作（创建、读取、更新、删除）
- 评论功能（创建、读取）
- 权限控制（只有作者可以修改/删除自己的文章）
- 错误处理和日志记录

## 技术栈

- **Go 1.21+**
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **MySQL** - 数据库
- **JWT** - 用户认证
- **bcrypt** - 密码加密

## 项目结构

```
blog/
├── main.go              # 程序入口
├── models/              # 数据模型
│   └── models.go
├── database/            # 数据库连接
│   └── database.go
├── handlers/            # 请求处理
│   ├── auth.go         # 用户认证
│   ├── post.go         # 文章管理
│   └── comment.go      # 评论管理
├── middleware/          # 中间件
│   └── auth.go         # JWT 认证中间件
├── go.mod              # 依赖管理
├── go.sum              # 依赖校验
└── README.md           # 项目说明
```

## 环境要求

- Go 1.21 或更高版本
- MySQL 5.7+ 或 MySQL 8.0+
- Git（可选）

## 安装步骤

1. **克隆或下载项目**



2. **确保 MySQL 服务运行**

确保 MySQL 服务已启动并运行在 `localhost:3306`。

3. **安装依赖**

```bash
go mod download
```

4. **运行项目**

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

首次运行时会自动创建所需的数据表。


## 测试用例

### 使用 Postman 测试

1. **注册新用户**
   - Method: POST
   - URL: `http://localhost:8080/api/register`
   - Body (JSON):
     ```json
     {
       "username": "testuser",
       "password": "password123",
       "email": "test@example.com"
     }
     ```

2. **用户登录**
   - Method: POST
   - URL: `http://localhost:8080/api/login`
   - Body (JSON):
     ```json
     {
       "username": "testuser",
       "password": "password123"
     }
     ```
   - 保存返回的 `token` 用于后续请求
   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ1ODU1ODUsImlkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIn0.XANUEEQSxQ2WTWom6EFOaR3wbBm74NVKc0xg2X-yXSI"

3. **创建文章**
   - Method: POST
   - URL: `http://localhost:8080/api/posts`
   - Headers:
     - `Authorization: Bearer <your_token>`
   - Body (JSON):
     ```json
     {
       "title": "测试文章",
       "content": "这是测试文章的内容"
     }
     ```

4. **获取所有文章**
   - Method: GET
   - URL: `http://localhost:8080/api/posts`

5. **获取单个文章**
   - Method: GET
   - URL: `http://localhost:8080/api/posts/1`

6. **更新文章**
   - Method: PUT
   - URL: `http://localhost:8080/api/posts/1`
   - Headers:
     - `Authorization: Bearer <your_token>`
   - Body (JSON):
     ```json
     {
       "title": "更新后的标题",
       "content": "更新后的内容"
     }
     ```

7. **创建评论**
   - Method: POST
   - URL: `http://localhost:8080/api/posts/1/comments`
   - Headers:
     - `Authorization: Bearer <your_token>`
   - Body (JSON):
     ```json
     {
       "content": "这ß是一条评论"
     }
     ```

8. **获取文章评论**
   - Method: GET
   - URL: `http://localhost:8080/api/posts/1/comments`

9. **删除文章**
   - Method: DELETE
   - URL: `http://localhost:8080/api/posts/1`
   - Headers:
     - `Authorization: Bearer <your_token>`

## 数据库

项目使用 MySQL 数据库。请确保 MySQL 服务已启动并运行在 `localhost:3306`。

### 数据库配置

当前配置信息（在 `database/database.go` 中）：
- **主机**: localhost
- **端口**: 3306
- **用户名**: root
- **密码**: 
- **数据库名**: mysql


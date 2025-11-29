package main

import (
	"fmt"

	//
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Student 学生结构体
type Student struct {
	ID    uint `gorm:"primaryKey;autoIncrement"`
	Name  string
	Age   int
	Grade string
}

// Account 账户表
type Account struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	Balance float64
}

// Transaction 交易表
type Transaction struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	FromAccountID uint
	ToAccountID   uint
	Amount        float64
}

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// func connectDB() (*gorm.DB, error) {
// 	// 数据库连接字符串
// 	dsn := "root:Zhaoyang@100297@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"

// 	// 打开数据库连接
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("打开数据库连接失败: %v", err)
// 	}

// 	return db, nil
// }

func main() {
	//GORM连接数据库
	// db, err := connectDB()
	// if err != nil {
	// 	fmt.Printf("连接数据库失败: %v", err)
	// 	return
	// }

	dsn := "root:Zhaoyang@100297@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Println("打开数据库连接失败: %v", err)
		return
	}

	// db.AutoMigrate(&Account{}, &Transaction{})

	// db.Create(&Account{Balance: 1000.0})
	// db.Create(&Account{Balance: 2000.0})

	// db.Transaction(func(tx *gorm.DB) error {
	// 	var fromAccount Account
	// 	tx.First(&fromAccount, 1)

	// 	if fromAccount.Balance < 100 {
	// 		return fmt.Errorf("账户 %d 余额不足，当前余额: %.2f，需要: %.2f",
	// 			fromAccount.ID, fromAccount.Balance, 100)
	// 	}

	// 	var toAccount Account
	// 	tx.First(&toAccount, 2)

	// 	if err := tx.Model(&fromAccount).
	// 		Update("balance", gorm.Expr("balance - ?", 100)).Error; err != nil {
	// 		return fmt.Errorf("扣除转出账户余额失败: %v", err)
	// 	}

	// 	if err := tx.Model(&toAccount).
	// 		Update("balance", gorm.Expr("balance + ?", 100)).Error; err != nil {
	// 		return fmt.Errorf("增加转入账户余额失败: %v", err)
	// 	}

	// 	transaction := Transaction{
	// 		FromAccountID: fromAccount.ID,
	// 		ToAccountID:   toAccount.ID,
	// 		Amount:        100,
	// 	}
	// 	if err := tx.Create(&transaction).Error; err != nil {
	// 		return fmt.Errorf("记录交易失败: %v", err)
	// 	}

	// 	return nil
	// })

	query1 := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		department VARCHAR(50) NOT NULL,
		salary DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err = db.Exec(query1)
	if err != nil {
		fmt.Println("创建表失败: %v", err)
		return
	}

	query2 := `
	CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		author VARCHAR(100) NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err = db.Exec(query2)
	if err != nil {
		fmt.Println("创建表失败: %v", err)
		return
	}

	// 插入测试数据
	fmt.Println("\n插入测试数据...")
	testEmployees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 15000.00},
		{Name: "李四", Department: "技术部", Salary: 18000.00},
		{Name: "王五", Department: "销售部", Salary: 12000.00},
		{Name: "赵六", Department: "技术部", Salary: 20000.00},
		{Name: "孙七", Department: "人事部", Salary: 10000.00},
	}

	// 使用 sqlx.NamedExec 批量插入
	for _, emp := range testEmployees {
		query := `INSERT INTO employees (name, department, salary) 
		          VALUES (:name, :department, :salary)`
		_, err := db.NamedExec(query, emp)
		if err != nil {
			fmt.Printf("插入失败: %v\n", err)
		}
	}

	testBooks := []Book{
		{Title: "Go语言程序设计", Author: "张三", Price: 89.00},
		{Title: "数据库原理", Author: "李四", Price: 45.00},
		{Title: "算法导论", Author: "王五", Price: 128.00},
		{Title: "设计模式", Author: "赵六", Price: 68.00},
		{Title: "计算机网络", Author: "孙七", Price: 35.00},
	}

	// 使用 sqlx.NamedExec 批量插入
	for _, book := range testBooks {
		query := `INSERT INTO books (title, author, price) 
		          VALUES (:title, :author, :price)`
		_, err := db.NamedExec(query, book)
		if err != nil {
			fmt.Printf("插入失败: %v\n", err)
		}
	}

	var employees []Employee

	query3 := `SELECT id, name, department, salary 
	          FROM employees 
	          WHERE department = ? 
	          ORDER BY id`

	err = db.Select(&employees, query3, "技术部")
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	fmt.Println(employees)

	var employee Employee

	query4 := `SELECT id, name, department, salary 
	          FROM employees 
	          ORDER BY salary DESC 
	          LIMIT 1`

	err = db.Get(&employee, query4)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	fmt.Println(employee)

	var books []Book

	query5 := `SELECT id, title, author, price 
	          FROM books 
	          WHERE price > ? 
	          ORDER BY price DESC`

	err = db.Select(&books, query5, 50)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	fmt.Println(books)

}

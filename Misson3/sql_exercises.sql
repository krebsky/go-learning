-- ============================================
-- SQL 语句练习 - 基本 CRUD 操作
-- ============================================

-- 1. 创建 students 表
CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    grade VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- 题目1：插入数据
-- ============================================
-- 要求：向 students 表中插入一条新记录
--       学生姓名为 "张三"，年龄为 20，年级为 "三年级"

INSERT INTO students (name, age, grade) VALUES ('张三', 20, '三年级');

-- ============================================
-- 题目2：查询数据
-- ============================================
-- 要求：查询 students 表中所有年龄大于 18 岁的学生信息

SELECT id, name, age, grade FROM students WHERE age > 18;

-- ============================================
-- 题目3：更新数据
-- ============================================
-- 要求：将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"

UPDATE students SET grade = '四年级' WHERE name = '张三';

-- ============================================
-- 题目4：删除数据
-- ============================================
-- 要求：删除 students 表中年龄小于 15 岁的学生记录

DELETE FROM students WHERE age < 15;

-- ============================================
-- 额外练习：查询所有学生（用于验证）
-- ============================================
SELECT id, name, age, grade FROM students ORDER BY id;


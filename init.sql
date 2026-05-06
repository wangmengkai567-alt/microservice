-- 初始化数据库脚本

-- 创建用户服务数据库
CREATE DATABASE IF NOT EXISTS user_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建商品服务数据库
CREATE DATABASE IF NOT EXISTS product_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建订单服务数据库
CREATE DATABASE IF NOT EXISTS order_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用商品服务数据库
USE product_service;

-- 创建商品表
CREATE TABLE IF NOT EXISTS products (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试商品数据（仅在不存在时插入，避免重启后覆盖库存）
INSERT IGNORE INTO products (id, name, description, price, stock, created_at, updated_at) VALUES
(1, 'iPhone 15 Pro', '最新款苹果手机，A17 Pro 芯片', 7999.00, 50, NOW(), NOW()),
(2, 'MacBook Pro 16', 'M3 Max 芯片，64GB 内存', 25999.00, 20, NOW(), NOW()),
(3, 'AirPods Pro 2', '主动降噪无线耳机', 1899.00, 100, NOW(), NOW()),
(4, 'iPad Air', '10.9 英寸平板电脑', 4799.00, 50, NOW(), NOW()),
(5, 'Apple Watch Series 9', '智能手表', 3199.00, 80, NOW(), NOW());

-- 使用用户服务数据库
USE user_service;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 使用订单服务数据库
USE order_service;

-- 创建订单表
CREATE TABLE IF NOT EXISTS orders (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    user_id INT UNSIGNED NOT NULL,
    product_id INT UNSIGNED NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

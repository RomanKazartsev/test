CREATE DATABASE IF NOT EXISTS order_db;

-- Использование созданной базы данных
USE example_db;

-- Создание таблицы
CREATE TABLE IF NOT EXISTS current_order (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    order_uid VARCHAR(255),
    rack_number,
    entry
    name
    phone
    zip
);

-- Вставка данных в таблицу
INSERT INTO users (username) VALUES ('user1'), ('user2'), ('user3');
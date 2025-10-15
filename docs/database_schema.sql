-- MySQL Database Schema for Contact Management System
-- Run this script to create the database and tables with proper relationships and indexes

-- Create database
CREATE DATABASE IF NOT EXISTS getcontact CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE getcontact;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Indexes
    INDEX idx_users_full_name (full_name),
    INDEX idx_users_email (email),
    INDEX idx_users_phone (phone),
    INDEX idx_users_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create contacts table
CREATE TABLE IF NOT EXISTS contacts (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NULL,
    favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Foreign key constraint
    CONSTRAINT fk_contacts_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    -- Indexes
    INDEX idx_contacts_user_id (user_id),
    INDEX idx_contacts_full_name (full_name),
    INDEX idx_contacts_phone (phone),
    INDEX idx_contacts_email (email),
    INDEX idx_contacts_favorite (favorite),
    INDEX idx_contacts_created_at (created_at),

    -- Composite index for common queries
    INDEX idx_contacts_user_favorite (user_id, favorite),
    INDEX idx_contacts_user_created (user_id, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create a sample user for testing (optional)
-- Password is hashed version of 'password123'
INSERT IGNORE INTO users (id, full_name, email, phone, password) VALUES
(1, 'Admin User', 'admin@example.com', '+1234567890', '$2a$10$example.hash.here');

-- Create sample contacts for testing (optional)
INSERT IGNORE INTO contacts (user_id, full_name, phone, email, favorite) VALUES
(1, 'John Doe', '+1234567891', 'john@example.com', true),
(1, 'Jane Smith', '+1234567892', 'jane@example.com', false),
(1, 'Bob Johnson', '+1234567893', 'bob@example.com', true);
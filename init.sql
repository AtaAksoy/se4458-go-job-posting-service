-- Create jobs database if not exists
CREATE DATABASE IF NOT EXISTS jobsdb;
USE jobsdb;

-- Create jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    company VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    status BOOLEAN DEFAULT TRUE,
    created_at BIGINT NOT NULL,
    INDEX idx_title (title),
    INDEX idx_company (company),
    INDEX idx_city (city),
    INDEX idx_state (state),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FULLTEXT idx_search (title, description, company, city, state)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample data
INSERT INTO jobs (title, description, company, city, state, status, created_at) VALUES
('Senior Go Developer', 'We are looking for an experienced Go developer with 5+ years of experience in building scalable microservices.', 'TechCorp', 'Istanbul', 'TR', TRUE, UNIX_TIMESTAMP()),
('Frontend Developer', 'Join our team as a Frontend Developer specializing in React and TypeScript.', 'WebSolutions', 'Ankara', 'TR', TRUE, UNIX_TIMESTAMP()),
('DevOps Engineer', 'Experienced DevOps engineer needed for CI/CD pipeline management and cloud infrastructure.', 'CloudTech', 'Izmir', 'TR', TRUE, UNIX_TIMESTAMP()),
('Data Scientist', 'Looking for a Data Scientist with expertise in machine learning and big data processing.', 'DataAnalytics', 'Bursa', 'TR', TRUE, UNIX_TIMESTAMP()),
('Mobile Developer', 'iOS/Android developer with experience in Flutter or React Native.', 'MobileApps', 'Antalya', 'TR', TRUE, UNIX_TIMESTAMP()); 
-- Add SuperAdmin role
INSERT IGNORE INTO user_roles (id, name) VALUES (4, 'SUPERADMIN');

-- System logs table for audit, system, and security logs
CREATE TABLE IF NOT EXISTS system_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    category ENUM('AUDIT', 'SYSTEM', 'SECURITY') NOT NULL,
    action VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    user_email VARCHAR(100) DEFAULT NULL,
    target_email VARCHAR(100) DEFAULT NULL,
    ip_address VARCHAR(45) DEFAULT NULL,
    user_agent VARCHAR(255) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_system_logs_category (category),
    INDEX idx_system_logs_action (action),
    INDEX idx_system_logs_user_email (user_email),
    INDEX idx_system_logs_created_at (created_at),
    INDEX idx_system_logs_category_created (category, created_at)
);

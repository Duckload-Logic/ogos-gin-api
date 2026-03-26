-- ============================================================================
-- SYSTEM DOMAIN
-- ============================================================================

CREATE TABLE api_keys (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(64) NOT NULL UNIQUE,
    key_prefix VARCHAR(8) NOT NULL,
    scopes JSON DEFAULT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    last_used_at TIMESTAMP NULL DEFAULT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_api_keys_key_hash ON api_keys(key_hash ASC);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash ASC);
CREATE INDEX idx_api_keys_is_active ON api_keys(is_active ASC);

CREATE TABLE notifications (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type ENUM('Appointment', 'Slip', 'Guidance', 'System', 'General') DEFAULT 'System',
    is_read TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_notifications_user_id ON notifications(user_id ASC);

CREATE TABLE system_logs (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    user_id CHAR(36) NULL DEFAULT NULL,
    user_email VARCHAR(100) NULL DEFAULT NULL,
    ip_address VARCHAR(45) NULL DEFAULT NULL,
    user_agent VARCHAR(255) NULL DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category ENUM('SECURITY', 'SYSTEM', 'AUDIT', 'CONSENT') NOT NULL,
    CONSTRAINT fk_system_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_system_logs_action ON system_logs(action ASC);
CREATE INDEX idx_system_logs_user_email ON system_logs(user_email ASC);
CREATE INDEX idx_system_logs_created_at ON system_logs(created_at ASC);
CREATE INDEX idx_system_logs_category_created ON system_logs(created_at ASC);
CREATE INDEX idx_system_logs_user_id ON system_logs(user_id ASC);

CREATE TABLE counselor_profiles (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    license_number VARCHAR(50) DEFAULT NULL,
    specialization VARCHAR(100) DEFAULT NULL,
    is_available TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_counselor_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX idx_counselor_profiles_user_unique ON counselor_profiles(user_id ASC);
CREATE INDEX idx_counselor_profiles_user_id ON counselor_profiles(user_id ASC);

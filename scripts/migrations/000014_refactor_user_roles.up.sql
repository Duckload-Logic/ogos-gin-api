-- ============================================================================
-- REFACTOR USER ROLES TO MANY-TO-MANY
-- ============================================================================

-- 1. Rename existing lookup table to 'roles' for clarity
RENAME TABLE user_roles TO roles;

-- 2. Create the new junction table for multi-role support
CREATE TABLE user_roles (
    user_id CHAR(36) NOT NULL,
    role_id INT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by CHAR(36) DEFAULT NULL,
    reason TEXT DEFAULT NULL,
    reference_id VARCHAR(100) DEFAULT NULL,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_admin FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- 3. Migrate existing role data from users table to junction table
INSERT INTO user_roles (user_id, role_id)
SELECT id, role_id FROM users;

-- 4. Remove the old single role constraint and column
ALTER TABLE users DROP FOREIGN KEY users_ibfk_1;
DROP INDEX idx_users_role_id ON users;
ALTER TABLE users DROP COLUMN role_id;

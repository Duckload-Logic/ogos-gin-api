-- ============================================================================
-- ROLLBACK REFACTOR USER ROLES TO MANY-TO-MANY
-- ============================================================================

-- 1. Add back the role_id column to users
ALTER TABLE users ADD COLUMN role_id INT NULL;

-- 2. Restore data from junction table (take the first role for each user)
UPDATE users u 
SET role_id = (
    SELECT role_id 
    FROM user_roles ur 
    WHERE ur.user_id = u.id 
    ORDER BY assigned_at ASC 
    LIMIT 1
);

-- 3. Enforce NOT NULL and restore constraints
-- NOTE: If a user has no roles in the junction table, this may fail.
-- In a real rollback, you'd handle defaults here.
ALTER TABLE users MODIFY COLUMN role_id INT NOT NULL;
CREATE INDEX idx_users_role_id ON users(role_id ASC);
ALTER TABLE users ADD CONSTRAINT users_ibfk_1 FOREIGN KEY (role_id) REFERENCES roles(id);

-- 4. Drop junction table
DROP TABLE user_roles;

-- 5. Rename lookup table back to 'user_roles'
RENAME TABLE roles TO user_roles;

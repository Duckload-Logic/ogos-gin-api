-- ============================================================================
-- Add user_id FK to system_logs, drop target_email
-- ============================================================================

START TRANSACTION;

SET FOREIGN_KEY_CHECKS = 0;

-- 1. Add user_id column after message
ALTER TABLE system_logs
  ADD COLUMN user_id INT NULL AFTER message;

-- 2. (Optional) Populate user_id from existing user_email by joining with users table.
--    This ensures that existing log entries have a user_id where possible.
UPDATE system_logs sl
  LEFT JOIN users u ON sl.user_email = u.email
  SET sl.user_id = u.id;

-- 3. Add foreign key constraint
ALTER TABLE system_logs
  ADD CONSTRAINT fk_system_logs_user
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

-- 4. Drop the target_email column (no longer needed)
ALTER TABLE system_logs
  DROP COLUMN target_email;

-- 5. Recreate indexes (if any) to include user_id for performance.
--    (The table already has indexes on category, action, created_at.)
--    Optionally add an index on user_id for faster lookups.
CREATE INDEX idx_system_logs_user_id ON system_logs(user_id);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
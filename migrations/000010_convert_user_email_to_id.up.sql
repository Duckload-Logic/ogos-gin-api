-- ============================================================================
-- Add users.id as primary key and migrate all foreign keys from email to id
-- ============================================================================

START TRANSACTION;

SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------------------------------------------------------
-- 1. Modify users table: add auto_increment id as primary key
-- ----------------------------------------------------------------------------
ALTER TABLE users
  DROP PRIMARY KEY,
  ADD COLUMN id INT AUTO_INCREMENT PRIMARY KEY FIRST,
  ADD UNIQUE INDEX idx_users_email (email);

-- ----------------------------------------------------------------------------
-- 2. counselor_profiles
-- ----------------------------------------------------------------------------
-- Drop old foreign key and unique constraint (if any)
ALTER TABLE counselor_profiles DROP FOREIGN KEY counselor_profiles_ibfk_1;
-- Add new column and migrate data
ALTER TABLE counselor_profiles ADD COLUMN user_id INT AFTER user_email;
UPDATE counselor_profiles cp JOIN users u ON cp.user_email = u.email SET cp.user_id = u.id;
ALTER TABLE counselor_profiles MODIFY user_id INT NOT NULL;
ALTER TABLE counselor_profiles ADD UNIQUE INDEX idx_counselor_profiles_user_unique (user_id);
-- Drop old email column (automatically removes any indexes on it)
ALTER TABLE counselor_profiles DROP COLUMN user_email;
-- Add new foreign key
ALTER TABLE counselor_profiles ADD CONSTRAINT fk_counselor_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 3. appointments
-- ----------------------------------------------------------------------------
ALTER TABLE appointments DROP FOREIGN KEY appointments_ibfk_1;
ALTER TABLE appointments ADD COLUMN user_id INT AFTER user_email;
UPDATE appointments a JOIN users u ON a.user_email = u.email SET a.user_id = u.id;
ALTER TABLE appointments DROP COLUMN user_email;
ALTER TABLE appointments ADD CONSTRAINT fk_appointments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

-- ----------------------------------------------------------------------------
-- 4. iir_records
-- ----------------------------------------------------------------------------
ALTER TABLE iir_records DROP FOREIGN KEY iir_records_ibfk_1;
ALTER TABLE iir_records ADD COLUMN user_id INT AFTER user_email;
UPDATE iir_records ir JOIN users u ON ir.user_email = u.email SET ir.user_id = u.id;
ALTER TABLE iir_records MODIFY user_id INT NOT NULL;
ALTER TABLE iir_records ADD UNIQUE INDEX idx_iir_records_user_unique (user_id);
ALTER TABLE iir_records DROP COLUMN user_email;
ALTER TABLE iir_records ADD CONSTRAINT fk_iir_records_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 5. iir_drafts
-- ----------------------------------------------------------------------------
ALTER TABLE iir_drafts DROP FOREIGN KEY iir_drafts_ibfk_1;
ALTER TABLE iir_drafts ADD COLUMN user_id INT AFTER user_email;
UPDATE iir_drafts idr JOIN users u ON idr.user_email = u.email SET idr.user_id = u.id;
ALTER TABLE iir_drafts MODIFY user_id INT NOT NULL;
ALTER TABLE iir_drafts ADD UNIQUE INDEX idx_iir_drafts_user_unique (user_id);
ALTER TABLE iir_drafts DROP COLUMN user_email;
ALTER TABLE iir_drafts ADD CONSTRAINT fk_iir_drafts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 6. admission_slips
-- ----------------------------------------------------------------------------
ALTER TABLE admission_slips DROP FOREIGN KEY admission_slips_ibfk_1;
ALTER TABLE admission_slips ADD COLUMN user_id INT AFTER user_email;
UPDATE admission_slips aslip JOIN users u ON aslip.user_email = u.email SET aslip.user_id = u.id;
ALTER TABLE admission_slips MODIFY user_id INT NOT NULL;
ALTER TABLE admission_slips DROP COLUMN user_email;
ALTER TABLE admission_slips ADD CONSTRAINT fk_admission_slips_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 7. notifications
--    (the column is currently named user_id but references users.email)
-- ----------------------------------------------------------------------------
-- Rename the old column to user_email to avoid confusion
ALTER TABLE notifications CHANGE user_id user_email VARCHAR(100) NOT NULL;
-- Drop the old named foreign key
ALTER TABLE notifications DROP FOREIGN KEY fk_notification_user;
-- Add new user_id column
ALTER TABLE notifications ADD COLUMN user_id INT AFTER user_email;
UPDATE notifications n JOIN users u ON n.user_email = u.email SET n.user_id = u.id;
ALTER TABLE notifications MODIFY user_id INT NOT NULL;
ALTER TABLE notifications DROP COLUMN user_email;
ALTER TABLE notifications ADD CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 8. user_consents (no FK originally, but we add one now)
-- ----------------------------------------------------------------------------
ALTER TABLE user_consents ADD COLUMN user_id INT AFTER user_email;
UPDATE user_consents uc JOIN users u ON uc.user_email = u.email SET uc.user_id = u.id;
ALTER TABLE user_consents MODIFY user_id INT NOT NULL;
ALTER TABLE user_consents DROP COLUMN user_email;
ALTER TABLE user_consents ADD CONSTRAINT fk_user_consents_user FOREIGN KEY (user_id) REFERENCES users(id);

-- ----------------------------------------------------------------------------
-- 9. Re‑create performance indexes on the new user_id columns
-- ----------------------------------------------------------------------------
CREATE INDEX idx_counselor_profiles_user_id ON counselor_profiles(user_id);
CREATE INDEX idx_appointments_user_id ON appointments(user_id);
CREATE INDEX idx_iir_records_user_id ON iir_records(user_id);
CREATE INDEX idx_iir_drafts_user_id ON iir_drafts(user_id);
CREATE INDEX idx_admission_slips_user_id ON admission_slips(user_id);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_user_consents_user_id ON user_consents(user_id);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
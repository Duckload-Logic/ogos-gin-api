-- ============================================================================
-- Rollback: restore users.email as primary key and revert all FK changes
-- ============================================================================

START TRANSACTION;

SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------------------------------------------------------
-- 1. counselor_profiles
-- ----------------------------------------------------------------------------
-- Drop the new foreign key and index
ALTER TABLE counselor_profiles DROP FOREIGN KEY fk_counselor_profiles_user;
DROP INDEX idx_counselor_profiles_user_id ON counselor_profiles;
-- Add back user_email column
ALTER TABLE counselor_profiles ADD COLUMN user_email VARCHAR(100) AFTER user_id;
UPDATE counselor_profiles cp JOIN users u ON cp.user_id = u.id SET cp.user_email = u.email;
ALTER TABLE counselor_profiles MODIFY user_email VARCHAR(100) NOT NULL;
ALTER TABLE counselor_profiles ADD UNIQUE INDEX (user_email);
-- Drop user_id column
ALTER TABLE counselor_profiles DROP COLUMN user_id;
-- Restore old foreign key (name will be auto-generated)
ALTER TABLE counselor_profiles ADD FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 2. appointments
-- ----------------------------------------------------------------------------
ALTER TABLE appointments DROP FOREIGN KEY fk_appointments_user;
DROP INDEX idx_appointments_user_id ON appointments;
ALTER TABLE appointments ADD COLUMN user_email VARCHAR(100) AFTER user_id;
UPDATE appointments a JOIN users u ON a.user_id = u.id SET a.user_email = u.email;
ALTER TABLE appointments DROP COLUMN user_id;
ALTER TABLE appointments ADD FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE SET NULL;

-- ----------------------------------------------------------------------------
-- 3. iir_records
-- ----------------------------------------------------------------------------
ALTER TABLE iir_records DROP FOREIGN KEY fk_iir_records_user;
DROP INDEX idx_iir_records_user_id ON iir_records;
ALTER TABLE iir_records ADD COLUMN user_email VARCHAR(100) AFTER user_id;
UPDATE iir_records ir JOIN users u ON ir.user_id = u.id SET ir.user_email = u.email;
ALTER TABLE iir_records MODIFY user_email VARCHAR(100) NOT NULL;
ALTER TABLE iir_records ADD UNIQUE INDEX (user_email);
ALTER TABLE iir_records DROP COLUMN user_id;
ALTER TABLE iir_records ADD FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 4. iir_drafts
-- ----------------------------------------------------------------------------
ALTER TABLE iir_drafts DROP FOREIGN KEY fk_iir_drafts_user;
DROP INDEX idx_iir_drafts_user_id ON iir_drafts;
ALTER TABLE iir_drafts ADD COLUMN user_email VARCHAR(100) AFTER user_id;
UPDATE iir_drafts idr JOIN users u ON idr.user_id = u.id SET idr.user_email = u.email;
ALTER TABLE iir_drafts MODIFY user_email VARCHAR(100) NOT NULL;
ALTER TABLE iir_drafts ADD UNIQUE INDEX (user_email);
ALTER TABLE iir_drafts DROP COLUMN user_id;
ALTER TABLE iir_drafts ADD FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 5. admission_slips
-- ----------------------------------------------------------------------------
ALTER TABLE admission_slips DROP FOREIGN KEY fk_admission_slips_user;
DROP INDEX idx_admission_slips_user_id ON admission_slips;
ALTER TABLE admission_slips ADD COLUMN user_email VARCHAR(100) AFTER user_id;
UPDATE admission_slips aslip JOIN users u ON aslip.user_id = u.id SET aslip.user_email = u.email;
ALTER TABLE admission_slips MODIFY user_email VARCHAR(100) NOT NULL;
ALTER TABLE admission_slips DROP COLUMN user_id;
ALTER TABLE admission_slips ADD FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 6. notifications
-- ----------------------------------------------------------------------------
ALTER TABLE notifications DROP FOREIGN KEY fk_notifications_user;
DROP INDEX idx_notifications_user_id ON notifications;
ALTER TABLE notifications ADD COLUMN user_email VARCHAR(100) AFTER user_id;
UPDATE notifications n JOIN users u ON n.user_id = u.id SET n.user_email = u.email;
ALTER TABLE notifications MODIFY user_email VARCHAR(100) NOT NULL;
ALTER TABLE notifications DROP COLUMN user_id;
ALTER TABLE notifications ADD CONSTRAINT fk_notification_user FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE;

-- ----------------------------------------------------------------------------
-- 7. user_consents
-- ----------------------------------------------------------------------------
ALTER TABLE user_consents DROP FOREIGN KEY fk_user_consents_user;
DROP INDEX idx_user_consents_user_id ON user_consents;
ALTER TABLE user_consents ADD COLUMN user_email VARCHAR(255) AFTER user_id;
UPDATE user_consents uc JOIN users u ON uc.user_id = u.id SET uc.user_email = u.email;
ALTER TABLE user_consents MODIFY user_email VARCHAR(255) NOT NULL;
ALTER TABLE user_consents DROP COLUMN user_id;

-- ----------------------------------------------------------------------------
-- 8. Restore users table: remove id, make email primary key again
-- ----------------------------------------------------------------------------
ALTER TABLE users DROP INDEX idx_users_email;
ALTER TABLE users DROP PRIMARY KEY, DROP COLUMN id;
ALTER TABLE users ADD PRIMARY KEY (email);

-- ----------------------------------------------------------------------------
-- 9. Recreate the original indexes on user_email columns (from migration 000003)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_counselor_profiles_user_email ON counselor_profiles(user_email);
CREATE INDEX idx_appointments_user_email ON appointments(user_email);
CREATE INDEX idx_iir_records_user_email ON iir_records(user_email);
CREATE INDEX idx_iir_drafts_user_email ON iir_drafts(user_email);
CREATE INDEX idx_admission_slips_user_email ON admission_slips(user_email);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
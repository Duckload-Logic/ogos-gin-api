-- 1. Drop the junction/child table first (due to Foreign Key constraints)
DROP TABLE IF EXISTS user_consents;

-- 2. Drop the parent table
DROP TABLE IF EXISTS legal_documents;

-- 3. Revert system_logs changes
-- Note: Replace 'VARCHAR(255)' with whatever the original type was
-- if it wasn't a standard string.
ALTER TABLE system_logs
DROP COLUMN category;

ALTER TABLE system_logs
ADD COLUMN category ENUM('SECURITY', 'SYSTEM') NOT NULL; -- Revert to original ENUM values
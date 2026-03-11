-- 1. Add an index on document_id to satisfy the foreign key requirement
ALTER TABLE user_consents ADD INDEX idx_document_id (document_id);

-- 2. Drop the old unique constraint (on document_id alone)
ALTER TABLE user_consents DROP INDEX unique_user_consent;

-- 3. Add the composite unique key on (user_id, document_id)
ALTER TABLE user_consents ADD UNIQUE KEY unique_user_consent (user_id, document_id);
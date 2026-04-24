-- ============================================================================
-- ENABLE MULTIPLE ROLES PER WHITELISTED EMAIL
-- ============================================================================

-- 1. Drop existing primary key on 'email'
ALTER TABLE whitelists DROP PRIMARY KEY;

-- 2. Add composite primary key on '(email, role_id)'
ALTER TABLE whitelists ADD PRIMARY KEY (email, role_id);

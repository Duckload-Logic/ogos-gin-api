-- ============================================================================
-- REVERT MULTIPLE ROLES PER WHITELISTED EMAIL
-- ============================================================================

-- NOTE: This may fail if there are duplicate emails with different roles.
-- Manual cleanup of duplicates would be required before running this.

ALTER TABLE whitelists DROP PRIMARY KEY;
ALTER TABLE whitelists ADD PRIMARY KEY (email);

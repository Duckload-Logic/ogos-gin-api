-- ============================================================================
-- ROLLBACK: AUDIT TRAILS
-- ============================================================================

DROP INDEX idx_audit_trails_user_id ON audit_trails;
DROP INDEX idx_audit_trails_entity ON audit_trails;
DROP INDEX idx_audit_trails_action ON audit_trails;
DROP INDEX idx_audit_trails_created_at ON audit_trails;

DROP TABLE IF EXISTS audit_trails;

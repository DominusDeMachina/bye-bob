-- Migration: audit_logs (down)
-- Created at: 2025-05-21T11:00:00Z

BEGIN;

-- Drop triggers for audit logging
DROP TRIGGER IF EXISTS employees_audit ON employees;
DROP TRIGGER IF EXISTS positions_audit ON positions;
DROP TRIGGER IF EXISTS departments_audit ON departments;
DROP TRIGGER IF EXISTS sites_audit ON sites;
DROP TRIGGER IF EXISTS assessment_templates_audit ON assessment_templates;
DROP TRIGGER IF EXISTS assessments_audit ON assessments;
DROP TRIGGER IF EXISTS goals_audit ON goals;
DROP TRIGGER IF EXISTS goal_checkins_audit ON goal_checkins;

-- Drop function for audit logging
DROP FUNCTION IF EXISTS audit_log_func();

-- Drop indexes on audit_logs
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_table_name;
DROP INDEX IF EXISTS idx_audit_logs_record_id;
DROP INDEX IF EXISTS idx_audit_logs_created_at;

-- Drop audit_logs table
DROP TABLE IF EXISTS audit_logs;

COMMIT; 
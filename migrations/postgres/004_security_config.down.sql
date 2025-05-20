-- Migration: security_config (down)
-- Created at: 2025-05-21T12:00:00Z

BEGIN;

-- Drop trigger for sensitive field protection
DROP TRIGGER IF EXISTS employees_sensitive_fields_trigger ON employees;

-- Drop trigger function
DROP FUNCTION IF EXISTS restrict_sensitive_fields_update();

-- Disable Row-Level Security on tables
ALTER TABLE employees DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS employees_app_policy ON employees;
DROP POLICY IF EXISTS employees_admin_policy ON employees;
DROP POLICY IF EXISTS employees_readonly_policy ON employees;

-- Revoke roles from users
REVOKE byebob_app_role FROM byebob_app;
REVOKE byebob_admin_role FROM byebob_admin;
REVOKE byebob_readonly_role FROM byebob_readonly;

-- Drop users
DROP USER IF EXISTS byebob_app;
DROP USER IF EXISTS byebob_admin;
DROP USER IF EXISTS byebob_readonly;

-- Drop roles
DROP ROLE IF EXISTS byebob_app_role;
DROP ROLE IF EXISTS byebob_admin_role;
DROP ROLE IF EXISTS byebob_readonly_role;

COMMIT; 
-- Migration: initial_schema (down)
-- Created at: 2025-05-21T10:00:00Z

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_employees_email;
DROP INDEX IF EXISTS idx_employees_position_id;
DROP INDEX IF EXISTS idx_employees_department_id;
DROP INDEX IF EXISTS idx_employees_site_id;
DROP INDEX IF EXISTS idx_employees_manager_id;
DROP INDEX IF EXISTS idx_employees_status;

-- Drop foreign key constraints
ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employee_position;
ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employee_department;
ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employee_site;
ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employee_manager;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS fk_department_lead;

-- Drop tables
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS departments;
DROP TABLE IF EXISTS sites;

COMMIT; 
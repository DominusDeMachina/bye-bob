-- Migration: security_config (up)
-- Created at: 2025-05-21T12:00:00Z

BEGIN;

-- Create application roles
CREATE ROLE byebob_app_role;
CREATE ROLE byebob_admin_role;
CREATE ROLE byebob_readonly_role;

-- Grant privileges to app role (service account)
GRANT USAGE ON SCHEMA public TO byebob_app_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO byebob_app_role;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO byebob_app_role;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO byebob_app_role;

-- Restrict app role from modifying schema
REVOKE CREATE ON SCHEMA public FROM byebob_app_role;

-- Grant privileges to admin role (administrative access)
GRANT USAGE ON SCHEMA public TO byebob_admin_role;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO byebob_admin_role;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO byebob_admin_role;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO byebob_admin_role;
GRANT CREATE ON SCHEMA public TO byebob_admin_role;

-- Grant privileges to readonly role (reporting, analytics)
GRANT USAGE ON SCHEMA public TO byebob_readonly_role;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO byebob_readonly_role;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO byebob_readonly_role;

-- Create application users
CREATE USER byebob_app WITH PASSWORD 'app_password_placeholder';
CREATE USER byebob_admin WITH PASSWORD 'admin_password_placeholder';
CREATE USER byebob_readonly WITH PASSWORD 'readonly_password_placeholder';

-- Assign roles to users
GRANT byebob_app_role TO byebob_app;
GRANT byebob_admin_role TO byebob_admin;
GRANT byebob_readonly_role TO byebob_readonly;

-- Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO byebob_app_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON TABLES TO byebob_admin_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT ON TABLES TO byebob_readonly_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO byebob_app_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON SEQUENCES TO byebob_admin_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO byebob_readonly_role;

-- Implement Row-Level Security (RLS) for sensitive tables

-- Enable RLS on employees table
ALTER TABLE employees ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
-- App role can see all employees but with limited permissions
CREATE POLICY employees_app_policy ON employees
    FOR ALL
    TO byebob_app_role
    USING (true);

-- Create update trigger function to restrict sensitive field updates for app role
CREATE OR REPLACE FUNCTION restrict_sensitive_fields_update()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the current user has the app role but not admin role
    IF (SELECT pg_has_role(CURRENT_USER, 'byebob_app_role', 'MEMBER') AND 
        NOT pg_has_role(CURRENT_USER, 'byebob_admin_role', 'MEMBER')) THEN
        
        -- Check if sensitive fields were modified
        IF (OLD.employment_type IS DISTINCT FROM NEW.employment_type) OR
           (OLD.start_date IS DISTINCT FROM NEW.start_date) OR
           (OLD.end_date IS DISTINCT FROM NEW.end_date) OR
           (OLD.status IS DISTINCT FROM NEW.status) THEN
            RAISE EXCEPTION 'Not authorized to modify sensitive employee fields';
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for sensitive field protection
CREATE TRIGGER employees_sensitive_fields_trigger
BEFORE UPDATE ON employees
FOR EACH ROW
EXECUTE FUNCTION restrict_sensitive_fields_update();

-- Admin role has full access to all rows
CREATE POLICY employees_admin_policy ON employees
    FOR ALL
    TO byebob_admin_role
    USING (true);

-- Readonly role can only view employee data
CREATE POLICY employees_readonly_policy ON employees
    FOR SELECT
    TO byebob_readonly_role
    USING (true);

-- Force RLS for all users including superusers
ALTER TABLE employees FORCE ROW LEVEL SECURITY;

-- Set up secure connection requirements
-- Uncomment in production environment
-- ALTER ROLE byebob_app SET client_min_messages TO warning;
-- ALTER ROLE byebob_app SET ssl TO on;
-- ALTER ROLE byebob_admin SET ssl TO on;
-- ALTER ROLE byebob_readonly SET ssl TO on;

COMMIT; 
-- Migration: audit_logs (up)
-- Created at: 2025-05-21T11:00:00Z

BEGIN;

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    action VARCHAR(50) NOT NULL,
    table_name VARCHAR(100) NOT NULL,
    record_id UUID,
    changes JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create indexes on frequently queried columns
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_table_name ON audit_logs(table_name);
CREATE INDEX idx_audit_logs_record_id ON audit_logs(record_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Create function to record audit logs
CREATE OR REPLACE FUNCTION audit_log_func() RETURNS TRIGGER AS $$
DECLARE
    changes_json JSONB;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        changes_json = to_jsonb(OLD);
        INSERT INTO audit_logs (user_id, action, table_name, record_id, changes)
        VALUES (current_setting('app.user_id', TRUE)::UUID, 'DELETE', TG_TABLE_NAME, OLD.id, changes_json);
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        changes_json = jsonb_object_agg(key, value)
        FROM (
            SELECT key, value
            FROM jsonb_each(to_jsonb(NEW)) AS new_fields(key, value)
            JOIN jsonb_each(to_jsonb(OLD)) AS old_fields(key, value) USING (key)
            WHERE new_fields.value IS DISTINCT FROM old_fields.value
        ) AS changed_fields;

        IF changes_json IS NOT NULL AND changes_json <> '{}'::JSONB THEN
            INSERT INTO audit_logs (user_id, action, table_name, record_id, changes)
            VALUES (current_setting('app.user_id', TRUE)::UUID, 'UPDATE', TG_TABLE_NAME, NEW.id, changes_json);
        END IF;
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        changes_json = to_jsonb(NEW);
        INSERT INTO audit_logs (user_id, action, table_name, record_id, changes)
        VALUES (current_setting('app.user_id', TRUE)::UUID, 'INSERT', TG_TABLE_NAME, NEW.id, changes_json);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for audit logging
CREATE TRIGGER employees_audit
AFTER INSERT OR UPDATE OR DELETE ON employees
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

CREATE TRIGGER positions_audit
AFTER INSERT OR UPDATE OR DELETE ON positions
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

CREATE TRIGGER departments_audit
AFTER INSERT OR UPDATE OR DELETE ON departments
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

CREATE TRIGGER sites_audit
AFTER INSERT OR UPDATE OR DELETE ON sites
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

CREATE TRIGGER assessment_templates_audit
AFTER INSERT OR UPDATE OR DELETE ON assessment_templates
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

CREATE TRIGGER assessments_audit
AFTER INSERT OR UPDATE OR DELETE ON assessments
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

CREATE TRIGGER goals_audit
AFTER INSERT OR UPDATE OR DELETE ON goals
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

CREATE TRIGGER goal_checkins_audit
AFTER INSERT OR UPDATE OR DELETE ON goal_checkins
FOR EACH ROW EXECUTE FUNCTION audit_log_func();

COMMIT; 
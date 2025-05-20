-- Migration: assessment_and_goals (down)
-- Created at: 2025-05-21T10:30:00Z

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_assessments_template_id;
DROP INDEX IF EXISTS idx_assessments_employee_id;
DROP INDEX IF EXISTS idx_assessments_reviewer_id;
DROP INDEX IF EXISTS idx_assessments_status;
DROP INDEX IF EXISTS idx_goals_employee_id;
DROP INDEX IF EXISTS idx_goals_status;
DROP INDEX IF EXISTS idx_goal_checkins_goal_id;

-- Drop tables
DROP TABLE IF EXISTS goal_checkins;
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS assessments;
DROP TABLE IF EXISTS assessment_templates;

COMMIT; 
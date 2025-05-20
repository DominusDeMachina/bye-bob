-- Migration: assessment_and_goals (up)
-- Created at: 2025-05-21T10:30:00Z

BEGIN;

-- Create assessment_templates table
CREATE TABLE IF NOT EXISTS assessment_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create assessments table
CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    reviewer_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_assessment_template FOREIGN KEY (template_id) REFERENCES assessment_templates(id),
    CONSTRAINT fk_assessment_employee FOREIGN KEY (employee_id) REFERENCES employees(id),
    CONSTRAINT fk_assessment_reviewer FOREIGN KEY (reviewer_id) REFERENCES employees(id)
);

-- Create goals table
CREATE TABLE IF NOT EXISTS goals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    time_frame VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_goal_employee FOREIGN KEY (employee_id) REFERENCES employees(id)
);

-- Create goal_checkins table
CREATE TABLE IF NOT EXISTS goal_checkins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    goal_id UUID NOT NULL,
    note TEXT,
    progress INTEGER NOT NULL CHECK (progress >= 0 AND progress <= 100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_checkin_goal FOREIGN KEY (goal_id) REFERENCES goals(id)
);

-- Create indexes on frequently queried columns
CREATE INDEX idx_assessments_template_id ON assessments(template_id);
CREATE INDEX idx_assessments_employee_id ON assessments(employee_id);
CREATE INDEX idx_assessments_reviewer_id ON assessments(reviewer_id);
CREATE INDEX idx_assessments_status ON assessments(status);
CREATE INDEX idx_goals_employee_id ON goals(employee_id);
CREATE INDEX idx_goals_status ON goals(status);
CREATE INDEX idx_goal_checkins_goal_id ON goal_checkins(goal_id);

COMMIT; 
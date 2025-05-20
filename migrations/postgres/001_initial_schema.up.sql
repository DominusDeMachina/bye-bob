-- Migration: initial_schema (up)
-- Created at: 2025-05-21T10:00:00Z

BEGIN;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create employees table
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    address TEXT,
    position_id UUID,
    department_id UUID,
    site_id UUID,
    manager_id UUID,
    employment_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    profile_picture_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create positions table
CREATE TABLE IF NOT EXISTS positions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(100) NOT NULL,
    description TEXT,
    requirements TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create departments table
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    lead_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create sites table
CREATE TABLE IF NOT EXISTS sites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Add foreign key constraints
ALTER TABLE employees
    ADD CONSTRAINT fk_employee_position
    FOREIGN KEY (position_id)
    REFERENCES positions(id);

ALTER TABLE employees
    ADD CONSTRAINT fk_employee_department
    FOREIGN KEY (department_id)
    REFERENCES departments(id);

ALTER TABLE employees
    ADD CONSTRAINT fk_employee_site
    FOREIGN KEY (site_id)
    REFERENCES sites(id);

ALTER TABLE employees
    ADD CONSTRAINT fk_employee_manager
    FOREIGN KEY (manager_id)
    REFERENCES employees(id);

ALTER TABLE departments
    ADD CONSTRAINT fk_department_lead
    FOREIGN KEY (lead_id)
    REFERENCES employees(id);

-- Create indexes on frequently queried columns
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_position_id ON employees(position_id);
CREATE INDEX idx_employees_department_id ON employees(department_id);
CREATE INDEX idx_employees_site_id ON employees(site_id);
CREATE INDEX idx_employees_manager_id ON employees(manager_id);
CREATE INDEX idx_employees_status ON employees(status);

COMMIT; 
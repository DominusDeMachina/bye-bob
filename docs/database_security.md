# PostgreSQL Security Configuration

This document outlines the security configurations and access control mechanisms implemented in the ByeBob application's PostgreSQL database.

## Role-Based Access Control

The application uses a role-based access control system with three primary roles:

### 1. Application Role (`byebob_app_role`)

- **Purpose**: Used by the application service account for normal operations
- **Permissions**:
  - `SELECT`, `INSERT`, `UPDATE`, `DELETE` on all tables
  - `USAGE` and `SELECT` on all sequences
  - `EXECUTE` on all functions
- **Restrictions**:
  - Cannot modify database schema (no `CREATE` permission)
  - Row-level security restricts modification of sensitive employee fields

### 2. Admin Role (`byebob_admin_role`)

- **Purpose**: Used for administrative tasks and database management
- **Permissions**:
  - Full privileges on all tables and sequences
  - Can modify database schema (`CREATE` permission)
  - Can bypass row-level security restrictions

### 3. Read-Only Role (`byebob_readonly_role`)

- **Purpose**: Used for reporting, analytics, and read-only access
- **Permissions**:
  - `SELECT` on all tables
  - `USAGE` and `SELECT` on all sequences
  - No write access to any tables

## Database Users

The system includes three pre-configured database users:

1. `byebob_app` - Service account used by the application
2. `byebob_admin` - Administrative user for management tasks
3. `byebob_readonly` - Read-only user for reporting and analytics

## Row-Level Security (RLS)

Row-level security is implemented on sensitive tables to enforce additional access controls:

### Employees Table

- **App Role**: Can view all records but cannot modify sensitive fields (employment_type, start_date, end_date, status)
- **Admin Role**: Has full access to all records
- **Read-Only Role**: Can only view records, cannot modify them

## Default Privileges

Default privileges are configured to automatically apply appropriate permissions to any new database objects:

- App role gets `SELECT`, `INSERT`, `UPDATE`, `DELETE` on new tables
- Admin role gets all privileges on new tables
- Read-only role gets `SELECT` on new tables

## Connection Security

In production, the following security settings are enabled:

- SSL connections are required for all users
- Client minimum message level is set to "warning" to reduce information leakage

## Password Management

**IMPORTANT**: The migration files contain placeholder passwords. In a production environment:

1. Use environment variables or a secure configuration management system to set real passwords
2. Never store actual passwords in version control
3. Regularly rotate database passwords according to your security policy

## Railway.com Configuration

### Setting Up Database Connection

When connecting to Railway.com PostgreSQL, use the appropriate credentials for each context:

```go
// config/database.go example
package config

import (
	"fmt"
	"os"
)

// ApplicationDBURL returns the database URL for the application service
func ApplicationDBURL() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	
	// Use application-specific credentials
	user := os.Getenv("DB_APP_USER") // should be set to byebob_app
	password := os.Getenv("DB_APP_PASSWORD")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", 
		user, password, host, port, dbname)
}

// AdminDBURL returns the database URL for administrative tasks
func AdminDBURL() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	
	// Use admin credentials
	user := os.Getenv("DB_ADMIN_USER") // should be set to byebob_admin
	password := os.Getenv("DB_ADMIN_PASSWORD")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", 
		user, password, host, port, dbname)
}

// ReadOnlyDBURL returns the database URL for reporting/analytics
func ReadOnlyDBURL() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	
	// Use read-only credentials
	user := os.Getenv("DB_READONLY_USER") // should be set to byebob_readonly
	password := os.Getenv("DB_READONLY_PASSWORD")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", 
		user, password, host, port, dbname)
}
```

### Environment Variables Example

```
# Railway.com environment variables
DB_HOST=your-railway-postgres-host.railway.app
DB_PORT=5432
DB_NAME=railway
DB_APP_USER=byebob_app
DB_APP_PASSWORD=your-secure-app-password
DB_ADMIN_USER=byebob_admin
DB_ADMIN_PASSWORD=your-secure-admin-password
DB_READONLY_USER=byebob_readonly
DB_READONLY_PASSWORD=your-secure-readonly-password
```

## Applying Security Changes

Security configuration is applied through migration 004_security_config.up.sql and can be reverted with 004_security_config.down.sql if needed.

## Best Practices for Developers

1. Always connect using the least-privileged role necessary for your task
2. Use parameterized queries to prevent SQL injection
3. Do not expose database credentials in application code
4. Use the application service account (byebob_app) for normal application operations
5. Reserve admin access for authorized database administrators only 
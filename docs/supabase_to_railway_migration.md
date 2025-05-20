# Migrating from Supabase to Railway

This document outlines the process for migrating the ByeBob application's database from Supabase to Railway.com PostgreSQL.

## Why Migrate?

Railway.com offers several advantages for our PostgreSQL database needs:

1. **Simplified Deployment**: Railway provides a streamlined deployment experience with automatic database provisioning
2. **Integrated CI/CD**: Better integration with our continuous deployment pipeline
3. **Cost-Effective Scaling**: More flexible pricing model for our current stage
4. **Developer Experience**: Improved developer tooling and dashboard interface
5. **Performance**: Optimized PostgreSQL instances with better regional availability

## Migration Steps

### 1. Create Railway PostgreSQL Instance

1. Create a new Railway account or log in at [railway.app](https://railway.app)
2. Create a new project and provision a PostgreSQL database
3. Note your database connection details from the Railway dashboard

### 2. Export Data from Supabase

1. Access your Supabase dashboard
2. Navigate to SQL Editor
3. Run a pg_dump command or use the export function:
   ```sql
   -- In Supabase SQL Editor, generate CSV or SQL dumps of key tables
   COPY (SELECT * FROM employees) TO STDOUT WITH CSV HEADER;
   COPY (SELECT * FROM departments) TO STDOUT WITH CSV HEADER;
   -- Repeat for other tables
   ```
4. Alternatively, use `pg_dump` directly if you have access:
   ```bash
   pg_dump "postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres" > supabase_backup.sql
   ```

### 3. Update Environment Variables

Update your `.env` file to use Railway instead of Supabase:

```env
# Railway Database (primary)
RAILWAY_DB_URL=postgresql://postgres:password@containers-us-west-xxx.railway.app:5432/railway

# Supabase Database (kept for reference, can be removed once migration is confirmed)
# DB_HOST=db.xxx.supabase.co
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=your-supabase-password
# DB_NAME=postgres
# DB_SSLMODE=require
```

### 4. Import Data to Railway

1. If you have a SQL dump, import it into Railway:
   ```bash
   psql "postgresql://postgres:password@containers-us-west-xxx.railway.app:5432/railway" < supabase_backup.sql
   ```
   
2. Alternatively, import via CSV files or SQL scripts through your preferred PostgreSQL client connected to Railway

### 5. Run Migrations

Apply all schema migrations to ensure the Railway database is up to date:

```bash
export POSTGRESQL_URL="postgresql://postgres:password@containers-us-west-xxx.railway.app:5432/railway"
make migrate-up
```

### 6. Test the Connection

Verify that your application can connect to Railway:

```bash
make test-railway
```

### 7. Update CI/CD Pipelines

If you're using CI/CD pipelines, update the database connection variables in your pipeline configuration.

### 8. Monitoring and Verification

1. Monitor application performance with the new database
2. Verify that all database-dependent functionality works correctly
3. Check Railway's metrics dashboard for database performance

### 9. Cleanup

Once the migration is confirmed successful:

1. Update documentation to reflect the new Railway setup
2. Remove old Supabase environment variables
3. Archive or backup the Supabase project for reference

## Rollback Plan

If issues arise during migration:

1. Revert environment variables to Supabase connection details
2. Ensure Supabase database remains accessible during the migration period
3. Document any data changes made during the Railway testing period that might need to be replicated back to Supabase

## References

- [Railway PostgreSQL Documentation](https://docs.railway.app/database/postgresql)
- [PostgreSQL Import/Export Guide](https://www.postgresql.org/docs/current/backup.html) 
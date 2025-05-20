# Railway Setup Guide

This document outlines the process for setting up Railway.com for the ByeBob project's PostgreSQL database.

## Create a Railway Project

1. Go to [Railway](https://railway.app/) and log in or create an account
2. Click "New Project"
3. Select "Provision PostgreSQL"
4. Fill in the following details:
   - **Name**: byebob-db (or your preferred name)
   - **Region**: Choose the region closest to your target users
5. Click "Deploy"
6. Wait for the database to be provisioned (this should be quick)

## Get Connection Information

After your PostgreSQL database is created, you'll need to collect the following information:

1. From the Railway dashboard, select your newly created PostgreSQL database
2. Go to the "Connect" tab
3. Note down the following values:
   - **PostgreSQL Connection URL**: This contains all connection details in one string
   - **Database Name**: The name of your database
   - **Database User**: The username for connecting
   - **Database Password**: The password for connecting
   - **Host**: The database host address
   - **Port**: The port number (usually 5432)

## Update Environment Variables

Update your `.env` file with the following values from Railway:

```
# Railway Database
RAILWAY_DB_URL=postgresql://username:password@host.railway.app:port/database
```

Alternatively, you can use the individual connection components:

```
# Database
DB_HOST=host.railway.app
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=railway
DB_SSLMODE=require
```

## Test Connection

To verify that your connection is working:

1. Run the database test command:
   ```bash
   make test-db
   ```
   
2. Or specifically test the Railway connection:
   ```bash
   make test-railway
   ```

## Manage Your Database

Railway provides several ways to manage your PostgreSQL database:

1. **Railway Dashboard**:
   - View database metrics and logs
   - Scale your database as needed
   - Manage environment variables
   - Create database backups

2. **Connect with psql**:
   ```bash
   psql "postgresql://username:password@host.railway.app:port/database"
   ```

3. **Connect with pgAdmin** or other PostgreSQL management tools using the connection details

## Migrations

For running migrations against your Railway PostgreSQL database:

1. Set the PostgreSQL connection URL environment variable:
   ```bash
   export POSTGRESQL_URL="postgresql://username:password@host.railway.app:port/database"
   ```

2. Run the migration commands:
   ```bash
   make migrate-up
   ```

## Security Considerations

1. Never commit your `.env` file with real credentials to version control
2. Use Railway's environment variable functionality for production deployments
3. Regularly back up your database
4. Consider using connection pooling for production environments
5. Implement appropriate security measures in your application code
6. Use strong, unique passwords

## Railway Scaling and Features

Railway PostgreSQL offers several features that may be useful for your application:

1. **Automatic Backups**: Railway automatically creates backups of your database
2. **Metrics and Logging**: Monitor database performance and logs
3. **Scaling**: Easily scale your database as your application grows
4. **High Availability**: Options for higher availability configurations
5. **Networking**: Private networking and IP allowlisting

## Troubleshooting

If you encounter connection issues:

1. Verify your connection credentials in the `.env` file
2. Check if your IP is allowed (Railway has network controls)
3. Ensure your database is not paused (if using a plan with usage limits)
4. Check the Railway dashboard for database status and logs
5. Try connecting with `psql` to isolate application vs. connection issues 
# Supabase Setup Guide

This document outlines the process for setting up Supabase for the ByeBob project.

## Create a Supabase Project

1. Go to [Supabase](https://supabase.com) and log in or create an account
2. Click "New Project"
3. Fill in the following details:
   - **Name**: ByeBob (or your preferred project name)
   - **Database Password**: Create a secure password (save this for later)
   - **Region**: Choose the region closest to your target users
   - **Pricing Plan**: Select the appropriate plan (Free tier is suitable for development)
4. Click "Create New Project"
5. Wait for the project to be created (this may take a few minutes)

## Get Connection Information

After your project is created, you'll need to collect the following information:

1. From the Supabase dashboard, go to "Project Settings" > "API"
2. Note down the following values:
   - **Project URL**: `https://[your-project-id].supabase.co`
   - **API Key**: The `anon` and `service_role` keys are listed here
   - **Database Connection String**: Go to "Project Settings" > "Database" to find this

## Update Environment Variables

Update your `.env` file with the following values from Supabase:

```
# Supabase
SUPABASE_URL=https://[your-project-id].supabase.co
SUPABASE_API_KEY=[your-api-key]
SUPABASE_ANON_KEY=[your-anon-key]
SUPABASE_SERVICE_KEY=[your-service-role-key]
```

## Test Connection

To verify that your connection is working:

1. Run the application with the updated environment variables
2. Check logs for successful database connection messages
3. Alternatively, you can use Supabase's SQL editor to run queries directly

## Use Supabase Studio

Supabase provides a web-based interface called Supabase Studio for managing your database:

1. Go to the "Table Editor" to create and manage tables
2. Use the "SQL Editor" to run custom SQL queries
3. Set up authentication providers in the "Authentication" section
4. Configure storage buckets in the "Storage" section

## Security Considerations

1. Never commit your `.env` file with real credentials to version control
2. Use Row Level Security (RLS) policies to control data access
3. Keep your API keys secure and rotate them if they're compromised
4. Use the appropriate key for the context:
   - `anon` key for client-side code
   - `service_role` key for server-side operations (bypasses RLS)

## Troubleshooting

If you encounter connection issues:

1. Verify your credentials in the `.env` file
2. Check if your IP is allowed in Supabase's network restrictions
3. Ensure your database is not paused (free tier databases pause after inactivity)
4. Check the database logs in Supabase Studio 
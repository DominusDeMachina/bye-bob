#!/bin/bash

# Script to update database user passwords in PostgreSQL
# Usage: ./update_db_passwords.sh <environment>
# Environments: dev, staging, prod

set -e

if [ -z "$1" ]; then
  echo "Error: Environment not specified."
  echo "Usage: ./update_db_passwords.sh <environment>"
  echo "Environments: dev, staging, prod"
  exit 1
fi

ENV=$1

# Load environment variables
if [ -f ".env.$ENV" ]; then
  source ".env.$ENV"
else
  echo "Error: Environment file .env.$ENV not found"
  exit 1
fi

# Check required environment variables
if [ -z "$DB_CONNECTION_STRING" ] || [ -z "$DB_APP_PASSWORD" ] || [ -z "$DB_ADMIN_PASSWORD" ] || [ -z "$DB_READONLY_PASSWORD" ]; then
  echo "Error: Required environment variables not set."
  echo "Please ensure the following variables are set in your .env.$ENV file:"
  echo "  - DB_CONNECTION_STRING"
  echo "  - DB_APP_PASSWORD"
  echo "  - DB_ADMIN_PASSWORD"
  echo "  - DB_READONLY_PASSWORD"
  exit 1
fi

echo "Updating database passwords for $ENV environment..."

# Generate SQL commands
SQL_COMMANDS=$(cat <<EOF
ALTER USER byebob_app WITH PASSWORD '$DB_APP_PASSWORD';
ALTER USER byebob_admin WITH PASSWORD '$DB_ADMIN_PASSWORD';
ALTER USER byebob_readonly WITH PASSWORD '$DB_READONLY_PASSWORD';
SELECT 'Passwords updated successfully' as result;
EOF
)

# Execute SQL commands
echo "$SQL_COMMANDS" | psql "$DB_CONNECTION_STRING"

echo "Database passwords updated successfully for $ENV environment!" 
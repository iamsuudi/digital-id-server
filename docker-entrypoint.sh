#!/bin/sh

# Exit on error
set -e

# Wait for the database to be ready
# This is an extra precaution, as docker-compose's depends_on should handle this.
echo "Waiting for database..."
while ! nc -z db 5432; do
  sleep 1
done
echo "Database is up!"

# Run migrations and seed the database
make reload_schema
make seed_db

# Start the server
exec "/app/main"

#!/bin/sh
set -e

# Load environment variables
if [ -f /app/app.env ]; then
    source /app/app.env
else
    exit 1
fi

# Run DB migration
echo "Running database migration..."
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Start the app
echo "Starting the server..."
exec "$@"

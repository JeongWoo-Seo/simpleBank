#!/bin/sh
set -e

# Run database migration
echo "run db migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Start the app
echo "start the app"
exec "$@"
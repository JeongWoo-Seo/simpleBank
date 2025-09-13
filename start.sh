#!/bin/sh

set -e

# Wait for postgres to start
echo "waiting for postgres to start"
while ! nc -z postgres 5432; do
  sleep 1
done

# Run database migration
echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Start the app
echo "start the app"
exec "$@"

#!/bin/sh
set -e

echo "🟢 Waiting for PostgreSQL to be available..."

until nc -z postgres 5432; do
  echo "⏳ PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "✅ PostgreSQL is up - starting service..."

exec /app/authservice -config config/prod -migrations migrations

#!/bin/bash

# Setup script for Neon.db with Fly.io
# Run this script after getting your Neon.db connection string

echo "Setting up Neon.db with Fly.io..."
echo ""
echo "Please provide your Neon.db connection string."
echo "You can find it in your Neon dashboard under 'Connection Details'"
echo "It should look like: postgresql://username:password@ep-xxx-xxx.us-east-1.aws.neon.tech/dbname?sslmode=require"
echo ""

read -p "Enter your Neon.db connection string: " NEON_URL

if [ -z "$NEON_URL" ]; then
    echo "Error: No connection string provided"
    exit 1
fi

echo "Setting DATABASE_URL secret in Fly.io..."
flyctl secrets set DATABASE_URL="$NEON_URL" --app myapp-1757744589

echo ""
echo "✅ Neon.db connection configured!"
echo "You can now deploy your app with: flyctl deploy"
echo ""
echo "To test the connection, you can run:"
echo "flyctl ssh console --app myapp-1757744589"
echo "Then inside the container, run: psql \$DATABASE_URL"

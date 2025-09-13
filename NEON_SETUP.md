# Neon.db Setup with Fly.io

## Steps to configure Neon.db:

### 1. Get your Neon.db connection string
1. Go to your [Neon Dashboard](https://console.neon.tech/)
2. Select your project
3. Go to "Connection Details" 
4. Copy the connection string (it should look like):
   ```
   postgresql://username:password@ep-xxx-xxx.us-east-1.aws.neon.tech/dbname?sslmode=require
   ```

### 2. Set the DATABASE_URL secret
Run the setup script:
```bash
./setup-neon.sh
```

Or manually set it:
```bash
flyctl secrets set DATABASE_URL="your-neon-connection-string" --app myapp-1757744589
```

### 3. Deploy your application
```bash
flyctl deploy
```

### 4. Test the connection
```bash
# Check logs
flyctl logs --app myapp-1757744589

# SSH into the container to test database connection
flyctl ssh console --app myapp-1757744589
# Then inside the container:
psql $DATABASE_URL
```

## Environment Variables

Your app will automatically use the `DATABASE_URL` environment variable when deployed to Fly.io. For local development, you can still use individual database variables in your `.env` file.

## Database Migrations

The app will automatically run migrations on startup using GORM's AutoMigrate feature. Make sure your Neon.db database is accessible and the connection string is correct.

## Troubleshooting

- If you see connection errors, verify your Neon.db connection string
- Check that your Neon.db project is active and not paused
- Ensure the database user has the necessary permissions

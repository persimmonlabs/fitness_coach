# Railway Deployment Guide

## Quick Deploy Steps

### 1. Create Railway Project

1. Go to https://railway.app
2. Sign in with GitHub
3. Click "New Project"
4. Select "Deploy from GitHub repo"
5. Choose `persimmonlabs/fitness_coach`
6. Railway will detect it's a Go project automatically

### 2. Add PostgreSQL Database

1. In your Railway project, click "+ New"
2. Select "Database" → "PostgreSQL"
3. Railway will create a PostgreSQL instance
4. It will automatically add a `DATABASE_URL` environment variable

### 3. Configure Environment Variables

Click on your Go service, go to "Variables" tab, and add:

```
DATABASE_URL (automatically set by Railway when you add PostgreSQL)
JWT_SECRET=your-production-secret-key-change-this
PORT=8080
ENV=production
```

Optional (for AI features):
```
OPENROUTER_API_KEY=your-key
SUPABASE_URL=your-url
SUPABASE_ANON_KEY=your-key
```

### 4. Deploy

Railway will automatically:
- Build your Go app
- Run migrations (via GORM AutoMigrate)
- Start the server
- Provide a public URL

## Verification

Once deployed, test your API:

```bash
# Health check
curl https://your-app.railway.app/health

# Register a user
curl -X POST https://your-app.railway.app/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

## Database Connection

Railway automatically sets `DATABASE_URL` in this format:
```
postgresql://user:password@hostname:port/database
```

Your app reads this from the environment and connects automatically.

## Viewing Logs

1. Click on your service in Railway
2. Go to "Deployments" tab
3. Click on the latest deployment
4. View real-time logs

## Environment Variables Set by Railway

Railway automatically provides:
- `DATABASE_URL` - PostgreSQL connection string (when you add PostgreSQL)
- `PORT` - Port to listen on (optional, defaults to 8080)
- `RAILWAY_ENVIRONMENT` - Current environment (production/staging)

## Troubleshooting

### Build Fails

Check that:
- `go.mod` and `go.sum` are committed
- `cmd/api/main.go` exists
- Code compiles locally: `go build ./...`

### App Crashes on Start

Check logs for:
- Missing `DATABASE_URL`
- Missing `JWT_SECRET`
- Migration errors

### Database Connection Fails

Verify:
- PostgreSQL service is running in Railway
- `DATABASE_URL` is set correctly
- Your app can reach the database (Railway handles networking)

### Can't Access API

Check:
- Deployment status (should be "Active")
- Port is correct (Railway expects your app to listen on `$PORT`)
- Public domain is generated (under "Settings" → "Networking")

## Connecting to Database

### Via Railway Dashboard

1. Click on PostgreSQL service
2. Go to "Data" tab
3. Browse tables directly

### Via psql

1. Click on PostgreSQL service
2. Go to "Connect" tab
3. Copy the connection command:
```bash
psql postgresql://user:password@hostname:port/database
```

4. Connect and run queries:
```sql
\dt               -- List tables
SELECT * FROM users;
\q                -- Quit
```

## Continuous Deployment

Every push to `main` branch will automatically:
1. Trigger a new build
2. Run tests (if configured)
3. Deploy the new version
4. Zero-downtime deployment

## Custom Domain (Optional)

1. Go to your service settings
2. Click "Networking"
3. Add custom domain
4. Update DNS records as instructed

## Monitoring

Railway provides:
- CPU and Memory usage graphs
- Request logs
- Deployment history
- Health checks

## Scaling

To scale your app:
1. Go to service settings
2. Adjust replicas (horizontal scaling)
3. Adjust resources (vertical scaling)

## Cost

Railway free tier includes:
- $5 credit per month
- Shared PostgreSQL
- 512MB RAM per service
- 1GB disk storage

## Production Checklist

✅ Set strong `JWT_SECRET`
✅ Set `ENV=production`
✅ Add PostgreSQL database
✅ Verify `DATABASE_URL` is set
✅ Test health endpoint
✅ Test authentication endpoints
✅ Check logs for errors
✅ Set up custom domain (optional)
✅ Configure monitoring/alerts

## Next Steps

After deployment:
1. Test all endpoints
2. Register a test user
3. Verify database tables exist
4. Check application logs
5. Set up frontend to connect to Railway URL

Your API will be available at:
```
https://your-app-name.railway.app
```

Use this URL in your frontend's API configuration!

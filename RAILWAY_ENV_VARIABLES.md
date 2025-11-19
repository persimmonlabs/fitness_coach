# Railway Environment Variables Setup

## ✅ REQUIRED Variables (Only 2!)

Add these in Railway Dashboard → Your Service → Variables:

### 1. DATABASE_URL
**Automatically set by Railway when you add PostgreSQL ✅**

If not set automatically:
```
postgresql://postgres:password@hostname:5432/railway
```
(Get this from your PostgreSQL service in Railway)

### 2. JWT_SECRET
**You MUST add this manually:**
```
JWT_SECRET=your-super-secret-production-key-min-32-characters-long
```

Generate a secure secret:
```bash
# Option 1: OpenSSL
openssl rand -base64 32

# Option 2: Python
python -c "import secrets; print(secrets.token_urlsafe(32))"

# Option 3: Just use this one for now:
JWT_SECRET=8x9mK2pL4nQ7rT1vW3yZ5aB6cD0eF8gH9iJ2kL4mN6oP8qR
```

## Optional Variables (Defaults are fine)

These have sensible defaults in the code:

```
PORT=8080                    # Railway sets this automatically
ENV=production               # Defaults to "development"
```

## For AI Features (Add Later)

Only add these when you want to enable AI meal parsing:

```
OPENROUTER_API_KEY=your-openrouter-key
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-supabase-anon-key
```

## ❌ REMOVE These Variables

All those other 32 variables (AI_MAX_TOKENS, SERVER_HOST, etc.) are **NOT needed** for this app. Those look like they're from a different project.

You can safely delete:
- AI_MAX_TOKENS
- AI_MODEL
- AI_TEMPERATURE
- CORS_* (app handles this)
- DB_HOST, DB_NAME, etc. (use DATABASE_URL instead)
- JWT_EXPIRATION (app defaults to 24h)
- LOG_* (app handles this)
- MIGRATION_* (not used)
- OPENAI_API_KEY (we use OpenRouter, not OpenAI)
- RATE_LIMIT_* (app handles this)
- SERVER_* (not needed)

## Quick Setup

1. Go to Railway Dashboard
2. Click your `fitness_coach` service
3. Click "Variables" tab
4. Add just these two:

```
DATABASE_URL  (should already be there)
JWT_SECRET    (add this one manually)
```

5. Click "Deploy" or push a new commit to trigger redeploy

## How to Add Variables in Railway

1. Click "New Variable"
2. Enter name: `JWT_SECRET`
3. Enter value: `8x9mK2pL4nQ7rT1vW3yZ5aB6cD0eF8gH9iJ2kL4mN6oP8qR`
4. Save

That's it! 🎉

## Verify DATABASE_URL is Set

Click your PostgreSQL service → "Variables" tab

You should see `DATABASE_URL` is automatically shared with your Go service.

If not:
1. Click PostgreSQL service
2. Click "Connect" tab
3. Copy the "Postgres Connection URL"
4. Add it manually to your Go service as `DATABASE_URL`

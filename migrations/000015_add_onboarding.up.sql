-- Add onboarding fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS onboarding_completed BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS age INT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS sex VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS height_cm DECIMAL(10,2);
ALTER TABLE users ADD COLUMN IF NOT EXISTS activity_level VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS dietary_preferences JSONB DEFAULT '{}'::jsonb;

CREATE INDEX IF NOT EXISTS idx_users_onboarding ON users(onboarding_completed);

COMMENT ON COLUMN users.sex IS 'male, female, other';
COMMENT ON COLUMN users.activity_level IS 'sedentary, lightly_active, moderately_active, very_active, extremely_active';


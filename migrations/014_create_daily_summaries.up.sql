-- Create daily_summaries table
CREATE TABLE daily_summaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    summary_date DATE NOT NULL,
    total_calories_consumed DECIMAL(8,2),
    total_protein_g DECIMAL(7,2),
    total_carbs_g DECIMAL(7,2),
    total_fat_g DECIMAL(7,2),
    total_fiber_g DECIMAL(7,2),
    total_calories_burned DECIMAL(8,2),
    total_exercise_minutes INTEGER,
    steps_count INTEGER,
    water_ml INTEGER,
    sleep_hours DECIMAL(4,2),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_summary_date UNIQUE (user_id, summary_date)
);

CREATE INDEX idx_daily_summaries_user ON daily_summaries(user_id);
CREATE INDEX idx_daily_summaries_date ON daily_summaries(summary_date);
CREATE INDEX idx_daily_summaries_user_date ON daily_summaries(user_id, summary_date);

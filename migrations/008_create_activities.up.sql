-- Create activities table
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type VARCHAR(100) NOT NULL, -- 'cardio', 'strength', 'sports', 'other'
    name VARCHAR(255) NOT NULL,
    duration_minutes INTEGER,
    calories_burned DECIMAL(8,2),
    distance_km DECIMAL(8,2),
    activity_time TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activities_user ON activities(user_id);
CREATE INDEX idx_activities_time ON activities(activity_time);
CREATE INDEX idx_activities_type ON activities(activity_type);
CREATE INDEX idx_activities_user_time ON activities(user_id, activity_time);

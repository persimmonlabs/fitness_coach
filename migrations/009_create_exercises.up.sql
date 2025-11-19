-- Create exercises table
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    category VARCHAR(100) NOT NULL, -- 'strength', 'cardio', 'flexibility', 'balance'
    primary_muscle_group VARCHAR(100),
    secondary_muscle_groups TEXT[], -- Array of muscle groups
    equipment VARCHAR(100),
    difficulty_level VARCHAR(50), -- 'beginner', 'intermediate', 'advanced'
    description TEXT,
    instructions TEXT,
    video_url VARCHAR(500),
    is_public BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_exercises_name ON exercises(name);
CREATE INDEX idx_exercises_category ON exercises(category);
CREATE INDEX idx_exercises_primary_muscle ON exercises(primary_muscle_group);
CREATE INDEX idx_exercises_equipment ON exercises(equipment);
CREATE INDEX idx_exercises_is_public ON exercises(is_public);

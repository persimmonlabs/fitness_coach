-- Create workout_sets table
CREATE TABLE workout_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_exercise_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    set_number INTEGER NOT NULL,
    reps INTEGER,
    weight_kg DECIMAL(6,2),
    duration_seconds INTEGER,
    distance_meters DECIMAL(8,2),
    estimated_1rm DECIMAL(6,2) GENERATED ALWAYS AS (
        CASE
            WHEN reps > 0 AND weight_kg > 0 THEN weight_kg * (1 + reps / 30.0)
            ELSE NULL
        END
    ) STORED,
    rpe DECIMAL(3,1), -- Rate of Perceived Exertion (1-10)
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_workout_sets_exercise ON workout_sets(workout_exercise_id);
CREATE INDEX idx_workout_sets_set_number ON workout_sets(workout_exercise_id, set_number);

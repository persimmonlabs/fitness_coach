-- Create metrics table
CREATE TABLE metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    metric_date DATE NOT NULL,
    weight_kg DECIMAL(5,2),
    body_fat_percentage DECIMAL(4,2),
    muscle_mass_kg DECIMAL(5,2),
    waist_cm DECIMAL(5,2),
    chest_cm DECIMAL(5,2),
    hips_cm DECIMAL(5,2),
    thigh_cm DECIMAL(5,2),
    arm_cm DECIMAL(5,2),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_metric_date UNIQUE (user_id, metric_date)
);

CREATE INDEX idx_metrics_user ON metrics(user_id);
CREATE INDEX idx_metrics_date ON metrics(metric_date);
CREATE INDEX idx_metrics_user_date ON metrics(user_id, metric_date);

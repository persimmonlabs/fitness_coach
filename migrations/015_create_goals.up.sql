-- Create goals table
CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    goal_type VARCHAR(100) NOT NULL, -- 'weight', 'body_fat', 'strength', 'endurance', 'nutrition', 'custom'
    target_value DECIMAL(10,2),
    target_unit VARCHAR(50),
    current_value DECIMAL(10,2),
    start_date DATE NOT NULL,
    target_date DATE,
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'completed', 'abandoned'
    description TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_goals_user ON goals(user_id);
CREATE INDEX idx_goals_type ON goals(goal_type);
CREATE INDEX idx_goals_status ON goals(status);
CREATE INDEX idx_goals_user_status ON goals(user_id, status);
CREATE INDEX idx_goals_target_date ON goals(target_date);

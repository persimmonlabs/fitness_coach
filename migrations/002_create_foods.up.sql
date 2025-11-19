-- Create foods table
CREATE TABLE foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255),
    calories_per_100g DECIMAL(7,2) GENERATED ALWAYS AS (
        (protein_per_100g * 4) + (carbs_per_100g * 4) + (fat_per_100g * 9)
    ) STORED,
    protein_per_100g DECIMAL(6,2) NOT NULL,
    carbs_per_100g DECIMAL(6,2) NOT NULL,
    fat_per_100g DECIMAL(6,2) NOT NULL,
    fiber_per_100g DECIMAL(6,2) DEFAULT 0,
    sugar_per_100g DECIMAL(6,2) DEFAULT 0,
    category VARCHAR(100),
    is_public BOOLEAN DEFAULT false,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_fiber_carbs CHECK (fiber_per_100g <= carbs_per_100g),
    CONSTRAINT check_sugar_carbs CHECK (sugar_per_100g <= carbs_per_100g)
);

CREATE INDEX idx_foods_name ON foods(name);
CREATE INDEX idx_foods_category ON foods(category);
CREATE INDEX idx_foods_is_public ON foods(is_public);
CREATE INDEX idx_foods_created_by ON foods(created_by);

-- Create serving_units table
CREATE TABLE serving_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    abbreviation VARCHAR(20),
    unit_type VARCHAR(50) NOT NULL, -- 'weight', 'volume', 'count'
    grams_equivalent DECIMAL(10,3),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_serving_units_name ON serving_units(name);
CREATE INDEX idx_serving_units_unit_type ON serving_units(unit_type);

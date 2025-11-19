-- Create food_serving_conversions table
CREATE TABLE food_serving_conversions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_id UUID NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
    serving_unit_id UUID NOT NULL REFERENCES serving_units(id) ON DELETE CASCADE,
    grams_per_serving DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_food_serving UNIQUE (food_id, serving_unit_id)
);

CREATE INDEX idx_food_serving_conversions_food ON food_serving_conversions(food_id);
CREATE INDEX idx_food_serving_conversions_unit ON food_serving_conversions(serving_unit_id);

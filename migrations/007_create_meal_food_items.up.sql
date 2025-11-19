-- Create meal_food_items table
CREATE TABLE meal_food_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meal_id UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    food_id UUID NOT NULL REFERENCES foods(id) ON DELETE RESTRICT,
    quantity DECIMAL(10,2) NOT NULL,
    serving_unit_id UUID REFERENCES serving_units(id) ON DELETE RESTRICT,
    grams DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_meal_food_items_meal ON meal_food_items(meal_id);
CREATE INDEX idx_meal_food_items_food ON meal_food_items(food_id);

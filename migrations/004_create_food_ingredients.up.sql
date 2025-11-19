-- Create food_ingredients table (for composite foods)
CREATE TABLE food_ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    composite_food_id UUID NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
    ingredient_food_id UUID NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL,
    serving_unit_id UUID REFERENCES serving_units(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_food_ingredient UNIQUE (composite_food_id, ingredient_food_id)
);

CREATE INDEX idx_food_ingredients_composite ON food_ingredients(composite_food_id);
CREATE INDEX idx_food_ingredients_ingredient ON food_ingredients(ingredient_food_id);

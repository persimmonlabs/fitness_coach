-- Seed serving units
INSERT INTO serving_units (name, abbreviation, unit_type, grams_equivalent) VALUES
-- Weight units
('gram', 'g', 'weight', 1.000),
('kilogram', 'kg', 'weight', 1000.000),
('ounce', 'oz', 'weight', 28.350),
('pound', 'lb', 'weight', 453.592),

-- Volume units
('milliliter', 'ml', 'volume', 1.000),
('liter', 'l', 'volume', 1000.000),
('cup', 'cup', 'volume', 240.000),
('tablespoon', 'tbsp', 'volume', 15.000),
('teaspoon', 'tsp', 'volume', 5.000),
('fluid ounce', 'fl oz', 'volume', 30.000),

-- Count-based serving sizes
('small', 'sm', 'count', NULL),
('medium', 'md', 'count', NULL),
('large', 'lg', 'count', NULL),
('extra large', 'xl', 'count', NULL),
('piece', 'pc', 'count', NULL),
('slice', 'slice', 'count', NULL),
('serving', 'serving', 'count', NULL),
('scoop', 'scoop', 'count', NULL),
('handful', 'handful', 'count', NULL),
('packet', 'packet', 'count', NULL),
('can', 'can', 'count', NULL),
('bottle', 'bottle', 'count', NULL),
('container', 'container', 'count', NULL);

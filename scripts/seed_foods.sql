-- Seed common foods with accurate USDA nutrition data per 100g
INSERT INTO foods (name, brand, protein_per_100g, carbs_per_100g, fat_per_100g, fiber_per_100g, sugar_per_100g, category, is_public) VALUES
-- Proteins
('Chicken Breast (skinless, cooked)', NULL, 31.00, 0.00, 3.60, 0.00, 0.00, 'protein', true),
('Ground Beef (85% lean, cooked)', NULL, 25.00, 0.00, 15.00, 0.00, 0.00, 'protein', true),
('Salmon (cooked)', NULL, 25.40, 0.00, 13.40, 0.00, 0.00, 'protein', true),
('Eggs (whole, cooked)', NULL, 13.00, 1.10, 11.00, 0.00, 1.10, 'protein', true),
('Tuna (canned in water)', NULL, 29.00, 0.00, 1.00, 0.00, 0.00, 'protein', true),
('Greek Yogurt (plain, non-fat)', NULL, 10.00, 3.60, 0.40, 0.00, 3.20, 'dairy', true),
('Cottage Cheese (low-fat)', NULL, 11.00, 3.40, 1.00, 0.00, 2.70, 'dairy', true),
('Tofu (firm)', NULL, 8.00, 1.90, 4.80, 0.30, 0.70, 'protein', true),
('Pork Chop (lean, cooked)', NULL, 28.00, 0.00, 7.00, 0.00, 0.00, 'protein', true),
('Turkey Breast (skinless, cooked)', NULL, 29.00, 0.00, 1.00, 0.00, 0.00, 'protein', true),

-- Carbohydrates
('White Rice (cooked)', NULL, 2.70, 28.00, 0.30, 0.40, 0.05, 'grain', true),
('Brown Rice (cooked)', NULL, 2.60, 23.00, 0.90, 1.80, 0.35, 'grain', true),
('Quinoa (cooked)', NULL, 4.40, 21.30, 1.90, 2.80, 0.87, 'grain', true),
('Oatmeal (cooked)', NULL, 2.40, 12.00, 1.40, 1.70, 0.30, 'grain', true),
('Whole Wheat Bread', NULL, 13.00, 43.00, 3.50, 7.00, 5.00, 'grain', true),
('Sweet Potato (baked)', NULL, 2.00, 20.00, 0.20, 3.00, 4.20, 'vegetable', true),
('Pasta (whole wheat, cooked)', NULL, 5.30, 25.00, 0.60, 3.50, 0.90, 'grain', true),
('White Potato (baked)', NULL, 2.50, 21.00, 0.10, 2.10, 1.20, 'vegetable', true),
('Banana', NULL, 1.10, 22.80, 0.30, 2.60, 12.20, 'fruit', true),
('Apple', NULL, 0.30, 13.80, 0.20, 2.40, 10.40, 'fruit', true),

-- Vegetables
('Broccoli (cooked)', NULL, 2.40, 6.00, 0.40, 3.30, 1.40, 'vegetable', true),
('Spinach (raw)', NULL, 2.90, 3.60, 0.40, 2.20, 0.40, 'vegetable', true),
('Carrots (raw)', NULL, 0.90, 9.60, 0.20, 2.80, 4.70, 'vegetable', true),
('Bell Pepper (raw)', NULL, 1.00, 6.00, 0.30, 2.10, 4.20, 'vegetable', true),
('Tomato (raw)', NULL, 0.90, 3.90, 0.20, 1.20, 2.60, 'vegetable', true),
('Cucumber (raw)', NULL, 0.70, 3.60, 0.10, 0.50, 1.70, 'vegetable', true),
('Lettuce (romaine)', NULL, 1.20, 3.30, 0.30, 2.10, 1.20, 'vegetable', true),
('Cauliflower (cooked)', NULL, 1.80, 4.00, 0.50, 2.30, 2.00, 'vegetable', true),
('Green Beans (cooked)', NULL, 1.80, 7.00, 0.20, 3.40, 3.30, 'vegetable', true),
('Asparagus (cooked)', NULL, 2.40, 3.90, 0.20, 2.10, 1.30, 'vegetable', true),

-- Fats & Nuts
('Almonds (raw)', NULL, 21.00, 22.00, 49.00, 12.50, 4.35, 'nuts', true),
('Peanut Butter', NULL, 25.00, 20.00, 50.00, 6.00, 9.00, 'nuts', true),
('Avocado', NULL, 2.00, 8.50, 15.00, 6.70, 0.70, 'fruit', true),
('Olive Oil', NULL, 0.00, 0.00, 100.00, 0.00, 0.00, 'fat', true),
('Walnuts (raw)', NULL, 15.00, 14.00, 65.00, 6.70, 2.60, 'nuts', true),
('Cashews (raw)', NULL, 18.00, 30.00, 44.00, 3.30, 5.90, 'nuts', true),
('Chia Seeds', NULL, 17.00, 42.00, 31.00, 34.00, 0.00, 'seeds', true),
('Flaxseeds (ground)', NULL, 18.00, 29.00, 42.00, 27.00, 1.55, 'seeds', true),

-- Dairy
('Milk (whole)', NULL, 3.20, 4.80, 3.30, 0.00, 5.00, 'dairy', true),
('Milk (skim)', NULL, 3.40, 5.00, 0.20, 0.00, 5.00, 'dairy', true),
('Cheddar Cheese', NULL, 25.00, 1.30, 33.00, 0.00, 0.50, 'dairy', true),
('Mozzarella Cheese (part skim)', NULL, 24.00, 2.20, 16.00, 0.00, 1.00, 'dairy', true),
('Whey Protein Powder', NULL, 80.00, 8.00, 3.00, 0.00, 3.00, 'supplement', true),

-- Fruits
('Orange', NULL, 0.90, 11.80, 0.10, 2.40, 9.40, 'fruit', true),
('Blueberries', NULL, 0.70, 14.50, 0.30, 2.40, 10.00, 'fruit', true),
('Strawberries', NULL, 0.70, 7.70, 0.30, 2.00, 4.90, 'fruit', true),
('Grapes', NULL, 0.70, 18.00, 0.20, 0.90, 15.50, 'fruit', true),
('Watermelon', NULL, 0.60, 7.60, 0.20, 0.40, 6.20, 'fruit', true),

-- Legumes
('Black Beans (cooked)', NULL, 8.90, 23.70, 0.50, 8.70, 0.30, 'legume', true),
('Lentils (cooked)', NULL, 9.00, 20.00, 0.40, 7.90, 1.80, 'legume', true),
('Chickpeas (cooked)', NULL, 8.90, 27.40, 2.60, 7.60, 4.80, 'legume', true);

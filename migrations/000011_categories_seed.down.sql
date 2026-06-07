ALTER TABLE payment_categories ALTER COLUMN group_id DROP NOT NULL;
UPDATE payments SET category_id = NULL;
DELETE FROM payment_categories;
DELETE FROM category_groups;

ALTER TABLE payment_categories
  DROP COLUMN group_id,
  ALTER COLUMN icon TYPE VARCHAR(10),
  ADD  COLUMN is_income BOOLEAN NOT NULL DEFAULT false;

DROP TABLE category_groups;

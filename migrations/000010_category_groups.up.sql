CREATE TABLE category_groups (
  id        SERIAL       PRIMARY KEY,
  name      VARCHAR(100) NOT NULL,
  is_income BOOLEAN      NOT NULL DEFAULT false
);

ALTER TABLE payment_categories
  ADD  COLUMN group_id INTEGER REFERENCES category_groups(id),
  ALTER COLUMN icon TYPE TEXT,
  DROP COLUMN is_income;

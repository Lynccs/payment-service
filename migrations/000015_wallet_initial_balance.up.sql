-- Замінюємо balance на initial_balance у wallets
ALTER TABLE wallets
  ADD COLUMN initial_balance NUMERIC(15, 2) NOT NULL DEFAULT 0;

UPDATE wallets SET initial_balance = balance;

ALTER TABLE wallets DROP COLUMN balance;

-- Системні категорії для коригування балансу
INSERT INTO category_groups (name, is_income) VALUES
  ('Система', false),
  ('Система', true);

INSERT INTO payment_categories (name, icon, group_id) VALUES
  ('Коригування балансу', '⚖️', (SELECT id FROM category_groups WHERE name = 'Система' AND is_income = false)),
  ('Коригування балансу', '⚖️', (SELECT id FROM category_groups WHERE name = 'Система' AND is_income = true));

DELETE FROM payment_categories
  WHERE name = 'Коригування балансу';

DELETE FROM category_groups
  WHERE name = 'Система';

ALTER TABLE wallets
  ADD COLUMN balance NUMERIC(15, 2) NOT NULL DEFAULT 0;

UPDATE wallets SET balance = initial_balance;

ALTER TABLE wallets DROP COLUMN initial_balance;

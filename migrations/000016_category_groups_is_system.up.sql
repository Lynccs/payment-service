ALTER TABLE category_groups ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT false;

UPDATE category_groups SET is_system = true WHERE name = 'Система';

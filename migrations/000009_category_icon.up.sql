ALTER TABLE payment_categories ADD COLUMN icon varchar(10);

UPDATE payment_categories SET icon = '🛒'  WHERE name = 'Продукти';
UPDATE payment_categories SET icon = '🍽️' WHERE name = 'Кафе та ресторани';
UPDATE payment_categories SET icon = '🚗'  WHERE name = 'Транспорт';
UPDATE payment_categories SET icon = '🏠'  WHERE name = 'Комунальні послуги';
UPDATE payment_categories SET icon = '👕'  WHERE name = 'Одяг та взуття';
UPDATE payment_categories SET icon = '💊'  WHERE name = 'Здоровя та аптека';
UPDATE payment_categories SET icon = '🎬'  WHERE name = 'Розваги';
UPDATE payment_categories SET icon = '📚'  WHERE name = 'Освіта';
UPDATE payment_categories SET icon = '💻'  WHERE name = 'Техніка та електроніка';
UPDATE payment_categories SET icon = '🏃'  WHERE name = 'Спорт';
UPDATE payment_categories SET icon = '💰'  WHERE name = 'Зарплата';
UPDATE payment_categories SET icon = '💼'  WHERE name = 'Фріланс';
UPDATE payment_categories SET icon = '🎁'  WHERE name = 'Подарунок';
UPDATE payment_categories SET icon = '💳'  WHERE name = 'Кешбек';

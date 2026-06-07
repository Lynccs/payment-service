-- Видаляємо платежі, бо вони мають FK на payment_methods
DELETE FROM payments;

-- Очищаємо старі методи
DELETE FROM payment_methods;

-- Скидаємо лічильник послідовності щоб ID починались з 1
ALTER SEQUENCE payment_methods_id_seq RESTART WITH 1;

-- Вносимо нові методи
INSERT INTO payment_methods (name) VALUES
  ('Готівка'),
  ('Дебетова карта'),
  ('Кредитна карта'),
  ('Банківський переказ'),
  ('Мобільний платіж');

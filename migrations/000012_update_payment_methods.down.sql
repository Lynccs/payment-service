DELETE FROM payments;

DELETE FROM payment_methods;

ALTER SEQUENCE payment_methods_id_seq RESTART WITH 1;

INSERT INTO payment_methods (name) VALUES
  ('Monobank'),
  ('ПриватБанк'),
  ('Готівка'),
  ('Apple Pay');

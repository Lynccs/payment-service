CREATE TABLE IF NOT EXISTS payment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30)
);

INSERT INTO payment_types (name) VALUES
    ('Переказ'),
    ('Поповнення'),
    ('Зняття'),
    ('Оплата');
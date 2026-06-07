DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_methods;

CREATE TABLE IF NOT EXISTS payment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30)
);

INSERT INTO payment_types (name) VALUES
    ('Переказ'),
    ('Поповнення'),
    ('Зняття'),
    ('Оплата');

CREATE TABLE IF NOT EXISTS payment_statuses (
    id SERIAL PRIMARY KEY,
    status VARCHAR(30)
);

INSERT INTO payment_statuses (status) VALUES
    ('В обробці'),
    ('Завершено'),
    ('Скасовано'),
    ('Помилка');

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    from_wallet_id INTEGER REFERENCES wallets(id),
    to_wallet_id INTEGER REFERENCES wallets(id),
    payment_type_id INTEGER REFERENCES payment_types(id),
    payment_status_id INTEGER REFERENCES payment_statuses(id),
    amount NUMERIC(15, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

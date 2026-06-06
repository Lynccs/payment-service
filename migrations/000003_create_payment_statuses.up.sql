CREATE TABLE IF NOT EXISTS payment_statuses (
    id SERIAL PRIMARY KEY,
    status VARCHAR(30)
);

INSERT INTO payment_statuses (status) VALUES
    ('В обробці'),
    ('Завершено'),
    ('Скасовано'),
    ('Помилка');
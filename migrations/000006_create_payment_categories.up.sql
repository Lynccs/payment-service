CREATE TABLE IF NOT EXISTS payment_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    is_income BOOLEAN NOT NULL
);

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_statuses;
DROP TABLE IF EXISTS payment_types;

CREATE TABLE IF NOT EXISTS payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    wallet_id INTEGER NOT NULL REFERENCES wallets(id),
    is_income BOOLEAN NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    category_id INTEGER REFERENCES payment_categories(id),
    payment_method_id INTEGER REFERENCES payment_methods(id),
    description TEXT,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

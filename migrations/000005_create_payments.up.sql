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
DROP TABLE IF EXISTS budget_limits;

CREATE TABLE budget_limits (
    id          SERIAL PRIMARY KEY,
    wallet_id   INTEGER NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES payment_categories(id) ON DELETE CASCADE,
    year        INTEGER NOT NULL,
    month       INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    amount      NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    UNIQUE(wallet_id, category_id, year, month)
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bonus_calculation_id UUID UNIQUE,
    status TEXT NOT NULL CHECK (status IN ('NEW', 'INVALID', 'PROCESSING', 'PROCESSED'))
);

CREATE TABLE bonus_calculations (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    loyalty_status TEXT NOT NULL CHECK (loyalty_status IN ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED')),
    accrual INTEGER NOT NULL
);

CREATE TABLE loyalty_account_transactions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    transaction_type TEXT NOT NULL CHECK (transaction_type IN ('WITHDRAW', 'ACCRUAL')),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE
    amount INTEGER NOT NULL DEFAULT 0
);
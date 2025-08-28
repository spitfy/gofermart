CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number TEXT NOT NULL,          -- номер заказа, который передаёт пользователь
    status TEXT NOT NULL,                -- например: 'NEW', 'PROCESSING', 'PROCESSED'
    accrual NUMERIC(12,2),               -- начисленные лояльные баллы (может быть NULL, если ещё не начислены)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(order_number)                 -- каждый заказ уникален
);

CREATE TABLE balance_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(12,2) NOT NULL,       -- >0 = начисление, <0 = списание
    type TEXT NOT NULL,                  -- 'ACCRUAL', 'WITHDRAW'
    order_id BIGINT REFERENCES orders(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Таблица выводов баллов пользователем (история запросов на вывод)
CREATE TABLE withdrawals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number TEXT NOT NULL,          -- номер заказа, в счёт которого списываются баллы
    amount NUMERIC(12,2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
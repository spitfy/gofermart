DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_withdrawals_user_id;
DROP INDEX IF EXISTS idx_accruals_user_id;

DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS accruals;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
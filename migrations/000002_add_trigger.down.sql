DROP TRIGGER IF EXISTS accrual_update_balance ON accruals;
DROP TRIGGER IF EXISTS withdrawal_update_balance ON withdrawals;

DROP FUNCTION IF EXISTS update_user_balance();

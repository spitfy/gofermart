CREATE OR REPLACE FUNCTION update_user_balance() RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'accruals' THEN
        UPDATE users SET balance = balance + NEW.amount WHERE id = NEW.user_id;
    ELSIF TG_TABLE_NAME = 'withdrawals' THEN
        UPDATE users SET balance = balance - NEW.amount WHERE id = NEW.user_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER accrual_update_balance
    AFTER INSERT ON accruals
    FOR EACH ROW EXECUTE FUNCTION update_user_balance();

CREATE TRIGGER withdrawal_update_balance
    AFTER INSERT ON withdrawals
    FOR EACH ROW EXECUTE FUNCTION update_user_balance();

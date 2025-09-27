BEGIN;
ALTER TABLE withdrawals DROP CONSTRAINT IF EXISTS withdrawals_order_id_fkey;

ALTER TABLE withdrawals RENAME COLUMN order_id TO "order";

ALTER TABLE withdrawals ALTER COLUMN "order" TYPE TEXT;
COMMIT;
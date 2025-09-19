BEGIN;
ALTER TABLE withdrawals RENAME COLUMN "order" TO order_id;

ALTER TABLE withdrawals ALTER COLUMN order_id TYPE INTEGER USING order_id::integer;

ALTER TABLE withdrawals ADD CONSTRAINT withdrawal_order_id_fkey
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
COMMIT;
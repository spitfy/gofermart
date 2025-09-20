BEGIN;
ALTER TABLE accruals
    ADD CONSTRAINT accruals_order_id_fkey FOREIGN KEY (order_id) REFERENCES orders(id);
COMMIT;
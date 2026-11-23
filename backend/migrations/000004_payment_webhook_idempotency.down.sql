DROP INDEX IF EXISTS uq_payment_logs_event_key;
ALTER TABLE payment_logs DROP COLUMN IF EXISTS event_key;

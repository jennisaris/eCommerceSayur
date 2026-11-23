-- ============================================================================
-- 000004_payment_webhook_idempotency.up.sql
-- Additive storage for deduplicating payment gateway callbacks.
-- ============================================================================

ALTER TABLE payment_logs
    ADD COLUMN IF NOT EXISTS event_key VARCHAR(128);

CREATE UNIQUE INDEX IF NOT EXISTS uq_payment_logs_event_key
    ON payment_logs(event_key)
    WHERE event_key IS NOT NULL AND btrim(event_key) <> '';

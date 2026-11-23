-- ============================================================================
-- 000003_security_indexes.up.sql
-- Additive, data-preserving indexes for FK integrity and critical lookups.
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_user_roles_role_id
    ON user_roles(role_id);

CREATE INDEX IF NOT EXISTS idx_verification_tokens_user_id
    ON verification_tokens(user_id);

CREATE INDEX IF NOT EXISTS idx_categories_parent_id
    ON categories(parent_id);

CREATE INDEX IF NOT EXISTS idx_products_parent_id
    ON products(parent_id);

CREATE INDEX IF NOT EXISTS idx_cart_items_product_id
    ON cart_items(product_id);

CREATE INDEX IF NOT EXISTS idx_order_items_product_id
    ON order_items(product_id);

CREATE INDEX IF NOT EXISTS idx_payments_user_id
    ON payments(user_id);

CREATE INDEX IF NOT EXISTS idx_products_status_category_created_at
    ON products(status, category_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_orders_buyer_created_at
    ON orders(buyer_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_orders_status_created_at
    ON orders(status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_payments_order_created_at
    ON payments(order_id, created_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS uq_payments_gateway_id_nonempty
    ON payments(payment_gateway_id)
    WHERE payment_gateway_id IS NOT NULL AND btrim(payment_gateway_id) <> '';

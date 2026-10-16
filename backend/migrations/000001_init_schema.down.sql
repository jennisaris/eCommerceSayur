-- Drop indexes
DROP INDEX IF EXISTS idx_notifications_receiver_id;
DROP INDEX IF EXISTS idx_payment_logs_payment_id;
DROP INDEX IF EXISTS idx_payments_gateway_id;
DROP INDEX IF EXISTS idx_payments_order_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_order_code;
DROP INDEX IF EXISTS idx_orders_buyer_id;
DROP INDEX IF EXISTS idx_products_status;
DROP INDEX IF EXISTS idx_products_slug;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables in reverse topological order
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS payment_logs CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS cart_items CASCADE;
DROP TABLE IF EXISTS carts CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS verification_tokens CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;

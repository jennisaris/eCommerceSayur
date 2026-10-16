-- ==============================================================================
-- 000002_seed_data.up.sql
-- Initial Data Seeder: Roles, Super Admin, Demo Customer, Categories & Sayur FreshCart Products
-- ==============================================================================

-- 1. Roles
INSERT INTO roles (id, name) VALUES
(1, 'Super Admin'),
(2, 'Customer')
ON CONFLICT (id) DO NOTHING;

-- 2. Default Users (Passwords hashed with bcrypt: 'Admin123!')
INSERT INTO users (id, name, email, password, phone, address, photo, is_verified, created_at, updated_at) VALUES
(1, 'Super Admin Sayur', 'admin@ecommerce.com', '$2a$10$iBCD9/KDT/xYY5a5RhEaOuTD9gnwnXnzR4fUehKaFdwwHchaV1OE.', '081234567890', 'Sentra Distribusi Sayur No. 1, Jakarta Selatan', '/assets/images/avatar/avatar-1.jpg', TRUE, NOW(), NOW()),
(2, 'Budi Santoso', 'budi@customer.com', '$2a$10$iBCD9/KDT/xYY5a5RhEaOuTD9gnwnXnzR4fUehKaFdwwHchaV1OE.', '081987654321', 'Jl. Kemang Raya No. 45B, Mampang Prapatan, Jakarta Selatan', '/assets/images/avatar/avatar-2.jpg', TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Assign Roles
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1), -- User 1 is Super Admin
(2, 2)  -- User 2 is Customer
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Reset sequence for users
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));

-- 3. Categories (Sayur Project FreshCart Master Categories)
INSERT INTO categories (id, parent_id, name, slug, icon, status, description, created_at, updated_at) VALUES
(1, NULL, 'Fruits & Vegetables', 'fruits-vegetables', '/assets/images/category/category-fruits-vegetables.jpg', TRUE, 'Sayuran hijau organik segar, buah-buahan lokal dan impor higienis siap masak.', NOW(), NOW()),
(2, NULL, 'Dairy, Bread & Eggs', 'dairy-bread-eggs', '/assets/images/category/category-dairy-bread-eggs.jpg', TRUE, 'Susu segar murni, telur ayam kampung omega 3, dan roti tawar artisan.', NOW(), NOW()),
(3, NULL, 'Snacks & Munchies', 'snack-munchies', '/assets/images/category/category-snack-munchies.jpg', TRUE, 'Camilan sehat keripik sayur, kacang panggang, dan biskuit gurih.', NOW(), NOW()),
(4, NULL, 'Bakery & Biscuits', 'bakery-biscuits', '/assets/images/category/category-bakery-biscuits.jpg', TRUE, 'Roti gandum utuh, pastry lezat, dan kue basah harian segar dari oven.', NOW(), NOW()),
(5, NULL, 'Chicken, Meat & Fish', 'chicken-meat-fish', '/assets/images/category/category-chicken-meat-fish.jpg', TRUE, 'Daging sapi segar potongan rapi, ayam karkas higienis, dan ikan laut segar.', NOW(), NOW()),
(6, NULL, 'Tea, Coffee & Drinks', 'tea-coffee-drinks', '/assets/images/category/category-tea-coffee-drinks.jpg', TRUE, 'Kopi arabika nusantara, teh hijau melati, dan jus peras murni tanpa pengawet.', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));

-- 4. Sample Products (Sayur FreshCart Product Catalog)
INSERT INTO products (id, category_id, parent_id, name, slug, image, description, regular_price, sale_price, unit, weight, stock, variant_name, status, created_at, updated_at) VALUES
(1, 1, NULL, 'Sayur Bayam Hijau Hidroponik Organik 250g', 'sayur-bayam-hijau-hidroponik-250g', '/assets/images/products/product-img-1.jpg', 'Bayam hijau segar hasil panen hidroponik tanpa pestisida kimia. Kaya zat besi dan serat tinggi, dipetik langsung subuh hari.', 15000, 12500, 'ikat', 300, 150, 'Ikat Segar 250g', 'ACTIVE', NOW(), NOW()),

(2, 1, NULL, 'Tomat Merah Segar Curah Pilihan 1 Kg', 'tomat-merah-segar-pilihan-1kg', '/assets/images/products/product-img-2.jpg', 'Tomat merah segar berdaging tebal dan manis alami. Sangat cocok untuk jus segar, bumbu tumis, maupun salad harian.', 24000, 19500, 'kg', 1050, 120, 'Kemasan 1 Kg', 'ACTIVE', NOW(), NOW()),

(3, 2, NULL, 'Susu Sapi Murni Pasteurisasi Fresh 1 Liter', 'susu-sapi-murni-pasteurisasi-1l', '/assets/images/products/product-img-3.jpg', 'Susu sapi perah murni tanpa pemanis buatan yang dipasteurisasi higienis. Kaya kalsium dan vitamin D untuk seluruh keluarga.', 32000, 28000, 'botol', 1100, 80, 'Botol 1 Liter', 'ACTIVE', NOW(), NOW()),

(4, 2, NULL, 'Telur Ayam Kampung Asli Omega 3 (Isi 10 Butir)', 'telur-ayam-kampung-omega3-10pcs', '/assets/images/products/product-img-4.jpg', 'Telur ayam kampung organik kaya omega 3 dan rendah kolesterol dengan kuning telur oranye pekat berkualitas premium.', 38000, 34000, 'pack', 700, 95, 'Pack Isi 10 Butir', 'ACTIVE', NOW(), NOW()),

(5, 5, NULL, 'Daging Dada Ayam Fillet Tanpa Kulit 500g', 'dada-ayam-fillet-tanpa-kulit-500g', '/assets/images/products/product-img-5.jpg', 'Dada ayam tanpa tulang dan tanpa lemak berlebih. Dihasilkan dari ayam potong higienis standar sertifikasi halal, siap olah.', 45000, 39900, 'pack', 550, 60, 'Tray Pack 500g', 'ACTIVE', NOW(), NOW()),

(6, 4, NULL, 'Roti Gandum Utuh Whole Wheat Toast Bread', 'roti-gandum-utuh-whole-wheat-bread', '/assets/images/products/product-img-6.jpg', 'Roti tawar gandum asli bertekstur empuk, tinggi serat dan mengenyangkan lebih lama. Tanpa bahan pengawet sintesis.', 26000, 22000, 'pack', 400, 75, 'Kemasan 10 Slices', 'ACTIVE', NOW(), NOW()),

(7, 6, NULL, 'Kopi Arabika Gayo Specialty Drip Bag (Isi 5 Sachet)', 'kopi-arabika-gayo-specialty-drip-bag', '/assets/images/products/product-img-7.jpg', 'Kopi gayo arabika single origin dalam kemasan drip bag praktis. Aroma spicy floral dengan rasa cokelat karamel manis.', 48000, 42000, 'box', 200, 100, 'Box 5 Drip Bags', 'ACTIVE', NOW(), NOW()),

(8, 3, NULL, 'Keripik Pisang Madu Oven Crispy 200g', 'keripik-pisang-madu-oven-crispy-200g', '/assets/images/products/product-img-8.jpg', 'Keripik pisang pilihan yang dipanggang dengan balutan madu murni. Renyah gurih tanpa minyak berlebih, cemilan sehat.', 22000, 18500, 'pack', 250, 110, 'Standing Pouch 200g', 'ACTIVE', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));

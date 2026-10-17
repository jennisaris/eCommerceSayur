-- ==============================================================================
-- 000002_seed_data.up.sql
-- Initial Data Seeder: Roles, Super Admin, Demo Customer, Book Categories & BookVerse Bestsellers
-- ==============================================================================

-- 1. Roles
INSERT INTO roles (id, name) VALUES
(1, 'Super Admin'),
(2, 'Customer')
ON CONFLICT (id) DO NOTHING;

-- 2. Default Users (Passwords hashed with bcrypt: 'Admin123!')
INSERT INTO users (id, name, email, password, phone, address, photo, is_verified, created_at, updated_at) VALUES
(1, 'Super Admin BookVerse', 'admin@ecommerce.com', '$2a$10$iBCD9/KDT/xYY5a5RhEaOuTD9gnwnXnzR4fUehKaFdwwHchaV1OE.', '081234567890', 'Gedung Literasi BookVerse Lantai 5, Jakarta Selatan', '/assets/images/avatar/avatar-1.jpg', TRUE, NOW(), NOW()),
(2, 'Budi Pembaca', 'budi@customer.com', '$2a$10$iBCD9/KDT/xYY5a5RhEaOuTD9gnwnXnzR4fUehKaFdwwHchaV1OE.', '081987654321', 'Jl. Kemang Raya No. 45B, Mampang Prapatan, Jakarta Selatan', '/assets/images/avatar/avatar-2.jpg', TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Assign Roles
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1), -- User 1 is Super Admin
(2, 2)  -- User 2 is Customer
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Reset sequence for users
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));

-- 3. Categories (BookVerse Master Book Categories)
INSERT INTO categories (id, parent_id, name, slug, icon, status, description, created_at, updated_at) VALUES
(1, NULL, 'Fiksi & Sastra', 'fiksi-sastra', '/assets/images/category/category-fruits-vegetables.jpg', TRUE, 'Koleksi novel fiksi populer, sastra klasik, cerpen, fiksi sejarah, dan karya fiksi pemenang penghargaan bergengsi.', NOW(), NOW()),
(2, NULL, 'Pengembangan Diri & Motivasi', 'self-improvement', '/assets/images/category/category-dairy-bread-eggs.jpg', TRUE, 'Buku psikologi populer, pembentukan habit positif, produktivitas kerja, ketenangan mental, dan kepemimpinan modern.', NOW(), NOW()),
(3, NULL, 'Bisnis, Finansial & Investasi', 'bisnis-investasi', '/assets/images/category/category-snack-munchies.jpg', TRUE, 'Panduan mengelola finansial pribadi, investasi pasar modal, saham, analisis bisnis, dan strategi kewirausahaan.', NOW(), NOW()),
(4, NULL, 'Teknologi, Coding & Sains', 'teknologi-sains', '/assets/images/category/category-bakery-biscuits.jpg', TRUE, 'Buku arsitektur software engineering, pemrograman praktis, kecerdasan buatan (AI), data science, dan sains modern.', NOW(), NOW()),
(5, NULL, 'Komik, Manga & Grafis', 'komik-manga', '/assets/images/category/category-chicken-meat-fish.jpg', TRUE, 'Manga Jepang terjemahan resmi berlisensi, komik grafis nusantara, manhwa petualangan, dan serial komik best seller.', NOW(), NOW()),
(6, NULL, 'Anak & Edukasi Pelajar', 'anak-edukasi', '/assets/images/category/category-tea-coffee-drinks.jpg', TRUE, 'Buku cerita bergambar edukatif untuk anak, ensiklopedia sains cilik, serta materi referensi belajar sekolah lengkap.', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  parent_id = EXCLUDED.parent_id,
  name = EXCLUDED.name,
  slug = EXCLUDED.slug,
  icon = EXCLUDED.icon,
  status = EXCLUDED.status,
  description = EXCLUDED.description,
  updated_at = NOW();

SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));

-- 4. Sample Products (BookVerse Bestselling Books Catalog)
INSERT INTO products (id, category_id, parent_id, name, slug, image, description, regular_price, sale_price, unit, weight, stock, variant_name, status, created_at, updated_at) VALUES
(1, 2, NULL, 'Atomic Habits: Perubahan Kecil yang Memberikan Hasil Luar Biasa', 'atomic-habits-james-clear', '/images/products/book-atomic-habits.jpg', 'Penulis: James Clear | Penerbit: Gramedia Pustaka Utama | Halaman: 356 | Bahasa: Indonesia. Buku terlaris #1 New York Times tentang cara mudah membangun kebiasaan baik dan membuang kebiasaan buruk dengan kaidah perubahan 1% setiap hari.', 108000, 95000, 'buku', 350, 150, 'Softcover Edisi Bahasa Indonesia', 'ACTIVE', NOW(), NOW()),

(2, 4, NULL, 'Clean Code: A Handbook of Agile Software Craftsmanship', 'clean-code-robert-c-martin', '/images/products/book-clean-code.jpg', 'Penulis: Robert C. Martin (Uncle Bob) | Penerbit: Prentice Hall | Halaman: 464 | Bahasa: Inggris. Panduan legendaris wajib bagi programmer dan software engineer dalam menulis kode berkualitas, refactoring, dan penerapan clean architecture.', 385000, 320000, 'buku', 680, 80, 'Paperback Original International Edition', 'ACTIVE', NOW(), NOW()),

(3, 2, NULL, 'Filosofi Teras: Panduan Menemukan Ketenangan Hidup Modern', 'filosofi-teras-henry-manampiring', '/images/products/product-img-3.jpg', 'Penulis: Henry Manampiring | Penerbit: Penerbit Buku Kompas | Halaman: 320 | Bahasa: Indonesia. Menyelami kebijaksanaan filsafat Stoisisme kuno untuk mengendalikan emosi negatif, cemas berlebih, dan meraih ketenangan batin era modern.', 98000, 85000, 'buku', 320, 120, 'Softcover Bookpaper Original', 'ACTIVE', NOW(), NOW()),

(4, 3, NULL, 'The Psychology of Money: Pelajaran Abadi Mengenai Kekayaan & Kebahagiaan', 'the-psychology-of-money-morgan-housel', '/images/products/product-img-4.jpg', 'Penulis: Morgan Housel | Penerbit: Penerbit Baca | Halaman: 268 | Bahasa: Indonesia. 19 kisah inspiratif tentang bagaimana perilaku, psikologi emosional, dan bias kognitif manusia memengaruhi keputusan finansial terbesar dalam hidup.', 85000, 75000, 'buku', 290, 140, 'Softcover Cetakan Terbaru', 'ACTIVE', NOW(), NOW()),

(5, 1, NULL, 'Laut Bercerita (Edisi Khusus Sampul Keras Kolektor)', 'laut-bercerita-leila-s-chudori', '/images/products/product-img-5.jpg', 'Penulis: Leila S. Chudori | Penerbit: Kepustakaan Populer Gramedia | Halaman: 394 | Bahasa: Indonesia. Karya sastra fiksi sejarah yang menggugah emosi tentang persahabatan, cinta, dan kehilangan mendalam aktivis mahasiswa era 1998.', 135000, 115000, 'buku', 480, 90, 'Hardcover Special Collector Edition', 'ACTIVE', NOW(), NOW()),

(6, 1, NULL, 'Bumi Manusia: Tetralogi Buru Karya Agung Sastra Nusantara', 'bumi-manusia-pramoedya-ananta-toer', '/images/products/product-img-6.jpg', 'Penulis: Pramoedya Ananta Toer | Penerbit: Lentera Dipantara | Halaman: 535 | Bahasa: Indonesia. Karya masterpiece sastra dunia mengisahkan Minke dan Annelies Mellema dengan latar pergulatan identitas Hindia Belanda awal abad ke-20.', 140000, 120000, 'buku', 550, 85, 'Paperback Edisi Resital', 'ACTIVE', NOW(), NOW()),

(7, 5, NULL, 'One Piece Vol. 100: Haoshoku (Manga Resmi Terjemahan Indonesia)', 'one-piece-vol-100-eiichiro-oda', '/images/products/product-img-7.jpg', 'Karya: Eiichiro Oda | Penerbit: Elex Media Komputindo | Halaman: 208 | Bahasa: Indonesia. Edisi spesial volume ke-100 pertempuran puncak Onigashima antara Luffy & Aliansi Topi Jerami menghadapi Yonkou Kaido.', 45000, 38000, 'buku', 190, 200, 'Tankobon Manga Bersegel Resmi', 'ACTIVE', NOW(), NOW()),

(8, 6, NULL, 'Dongeng Edukasi Sains Seru untuk Anak Pintar', 'dongeng-edukasi-sains-seru-anak-pintar', '/images/products/product-img-8.jpg', 'Penulis: Tim Edukasi Cilik | Penerbit: Bhuana Ilmu Populer | Halaman: 128 | Bahasa: Indonesia Full Color. Cerita dongeng bergambar mengenalkan rahasia tata surya, hewan, tumbuhan, dan sains sehari-hari secara menyenangkan untuk anak usia 4-9 tahun.', 68000, 55000, 'buku', 260, 100, 'Buku Cerita Bergambar Full Color', 'ACTIVE', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  category_id = EXCLUDED.category_id,
  parent_id = EXCLUDED.parent_id,
  name = EXCLUDED.name,
  slug = EXCLUDED.slug,
  image = EXCLUDED.image,
  description = EXCLUDED.description,
  regular_price = EXCLUDED.regular_price,
  sale_price = EXCLUDED.sale_price,
  unit = EXCLUDED.unit,
  weight = EXCLUDED.weight,
  stock = EXCLUDED.stock,
  variant_name = EXCLUDED.variant_name,
  status = EXCLUDED.status,
  updated_at = NOW();

SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));


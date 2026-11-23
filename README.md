# 🥬 Sayur FreshCart — Enterprise Full-Stack eCommerce Web Platform

Platform web e-commerce *production-ready* berkinerja tinggi yang dibangun dari awal (*from scratch*) dengan arsitektur **Clean Architecture (DDD)** di sisi backend dan **Nuxt 3 (Vue 3, Pinia, Tailwind CSS)** di sisi frontend, diintegrasikan secara presisi dengan template UI **FreshCart Sayur Project** (`ed56duyt-sayur-template.zip`). Dirancang khusus untuk ketahanan beban tinggi (1.000+ VU), integritas data finansial ACID, model pengiriman penuh (**100% Home Delivery**), serta integrasi aman payment gateway **Midtrans Snap**.

---

## 🏗️ 1. TECH STACK & ARSITEKTUR

| Komponen | Teknologi | Keterangan |
|---|---|---|
| **Backend API** | **Golang 1.22 + Fiber v2** | Clean Architecture / DDD (Domain, UseCase, Repository, Delivery HTTP) |
| **Frontend Web** | **Nuxt 3 + Vue 3 + TypeScript** | Composition API, Pinia State Management, Tailwind CSS |
| **UI Design System** | **FreshCart Sayur Template** | Tabler Icons, FreshCart Theme CSS, Category & Product Card Systems |
| **Primary Database** | **PostgreSQL 16** | ACID Transactions, JSONB webhook logs, Foreign Keys & Soft Deletes |
| **In-Memory Cache** | **Redis 7** | Cache-aside untuk katalog produk & sinkronisasi sesi keranjang (*cross-device*) |
| **Payment Gateway** | **Midtrans Snap & Webhook** | Snap popup modal + verifikasi SHA-512 Signature Key & restock otomatis |
| **Load Testing** | **Grafana k6** | Skenario simulasi 1.000 VUs dengan target latensi $\le 200\text{ ms}$ |
| **Containerization**| **Docker & Docker Compose** | Multi-stage Docker builds untuk rilis image yang ringan dan aman |

---

## 📋 2. ATURAN BISNIS UTAMA (CRITICAL RULES)

1. **Model Pengiriman: 100% Delivery Model**  
   - Tidak ada opsi *pickup* (ambil sendiri).  
   - Setiap checkout mewajibkan input: `recipient_name`, `recipient_phone`, `shipping_address`, dan penghitungan dinamis `shipping_fee`.
2. **Jam Operasional: 24/7 Non-Stop**  
   - Sistem transaksi, keranjang, dan pemesanan aktif 24 jam sehari tanpa batas jam operasional.
3. **Integritas Nilai Transaksi (Price Snapshotting)**  
   - Di tabel `order_items`, kolom `price` dan `subtotal` **wajib di-snapshot** saat checkout dibuat. Tidak ada kalkulasi ulang harga order lama dari tabel master `products`.
4. **Idempotensi Pembayaran Midtrans**  
   - Webhook callback memverifikasi **SHA-512 Signature Key**:  
     $$\text{SHA512}(\text{order\_id} + \text{status\_code} + \text{gross\_amount} + \text{ServerKey})$$  
   - Seluruh payload mentah dicatat ke tabel `payment_logs.raw_response` (tipe `JSONB`) untuk audit trail.
5. **Manajemen Stok Otomatis**  
   - Mengunci dan memotong stok dengan *pessimistic row locking* (`SELECT ... FOR UPDATE`) saat order dibuat (`status: PENDING_PAYMENT`).  
   - Jika pembayaran gagal/kedaluwarsa/dibatalkan (`expire`/`cancel`/`deny`), stok secara otomatis dikembalikan (*restocked*).

---

## 🗄️ 3. STRUKTUR 11 TABEL POSTGRESQL & INDEXES

1. `users`: Data akun pelanggan dan pengelola, verifikasi status, profil.
2. `roles`: Master peran akses (`Super Admin`, `Customer`).
3. `user_roles`: Pivot role-based access control (RBAC).
4. `verification_tokens`: Token aktivasi email dan reset password.
5. `categories`: Kategori produk bertingkat (*self-referencing parent_id*).
6. `products`: Katalog produk, varian, satuan, berat (gram), dan stok.
7. `carts`: Keranjang belanja pengguna (*cross-device synchronization*).
8. `cart_items`: Rincian item dalam keranjang belanja.
9. `orders`: Header transaksi pengiriman 100% delivery dengan snapshot alamat & ongkir.
10. `order_items`: Detail item transaksi dengan snapshot harga mutlak.
11. `payments`: Status pembayaran Midtrans (`settlement`, `pending`, `expire`, dll).
12. `payment_logs`: Log audit idempotent webhook (kolom `raw_response JSONB`).
13. `notifications`: Notifikasi real-time pelanggan.

### Indeks Kinerja Tinggi (Tuning k6 1.000+ VU):
- `idx_users_email` pada `users(email)`
- `idx_products_category_id` pada `products(category_id)`
- `idx_products_slug` pada `products(slug)`
- `idx_products_status` pada `products(status)`
- `idx_orders_buyer_id` pada `orders(buyer_id)`
- `idx_orders_order_code` pada `orders(order_code)`
- `idx_orders_status` pada `orders(status)`
- `idx_order_items_order_id` pada `order_items(order_id)`
- `idx_payments_order_id` pada `payments(order_id)`
- `idx_payments_gateway_id` pada `payments(payment_gateway_id)`
- `idx_payment_logs_payment_id` pada `payment_logs(payment_id)`
- `idx_notifications_receiver_id` pada `notifications(receiver_id)`

---

## 🚀 4. CARA MENJALANKAN DENGAN SATU PERINTAH

Seluruh stack (PostgreSQL, Redis, Golang API Backend, dan Nuxt 3 Frontend) dapat dijalankan langsung menggunakan Docker Compose:

```bash
docker compose up --build -d
```

Atau menggunakan `Makefile`:
```bash
make up
```

### Akses Layanan:
- **Frontend Customer & Admin:** [http://localhost:3000](http://localhost:3000)
- **Backend API:** [http://localhost:8080/api/v1](http://localhost:8080/api/v1)
- **Health Check:** [http://localhost:8080/health](http://localhost:8080/health)

Untuk development lokal, buat akun melalui halaman registrasi atau jalankan seeder secara eksplisit dengan `RUN_SEEDER=true`. Seeder development tidak boleh digunakan pada production.

---

## 📡 5. DAFTAR API ENDPOINTS

### Autentikasi & Akun
- `POST /api/v1/auth/register` — Pendaftaran pelanggan baru (otomatis membuat keranjang & role Customer).
- `POST /api/v1/auth/login` — Login dan penerbitan token JWT.
- `GET /api/v1/auth/csrf` — Inisialisasi token CSRF untuk request mutating.
- `GET /api/v1/auth/me` — Ambil profil pengguna terautentikasi melalui HttpOnly session cookie.
- `POST /api/v1/auth/refresh` — Rotasi refresh session dan access cookie.
- `POST /api/v1/auth/logout` — Revoke refresh session dan hapus cookie.
- `GET /api/v1/auth/verify-email?token=...` — Verifikasi email pengguna.

### Katalog & Kategori (Publik & Terakselerasi Redis)
- `GET /api/v1/categories` — Daftar kategori aktif.
- `GET /api/v1/categories/:slug` — Detail kategori berdasarkan slug.
- `GET /api/v1/products` — Filter katalog (`page`, `limit`, `search`, `category_id`, `sort`, `min_price`, `max_price`).
- `GET /api/v1/products/:slug` — Detail produk lengkap (menggunakan Redis cache-aside).

### Keranjang Belanja (Tersinkronisasi DB + Redis)
- `GET /api/v1/cart` — Ambil ringkasan keranjang (total item, total berat, total harga).
- `POST /api/v1/cart/items` — Tambah produk ke keranjang (dengan validasi stok).
- `PUT /api/v1/cart/items/:product_id` — Ubah kuantitas produk.
- `DELETE /api/v1/cart/items/:product_id` — Hapus item dari keranjang.
- `DELETE /api/v1/cart` — Kosongkan seluruh keranjang.

### Checkout & Transaksi (100% Delivery Model)
- `POST /api/v1/orders/estimate` — Hitung estimasi ongkos kirim dan total tagihan.
- `POST /api/v1/orders` — Eksekusi checkout transaksi (kunci stok, snapshot harga, buat order).
- `GET /api/v1/orders` — Riwayat pesanan pelanggan.
- `GET /api/v1/orders/:code` — Rincian invoice pesanan dan status pengiriman.

### Pembayaran Midtrans Snap & Webhook Callback
- `POST /api/v1/payments/:order_code/snap` — Buat Snap Token & URL Midtrans.
- `POST /api/v1/payments/midtrans-callback` — Endpoint Webhook Midtrans (Verifikasi SHA-512, update status pesanan, pencatatan `payment_logs` JSONB, restock otomatis jika gagal/batal).

### Manajemen Super Admin (`RequireRole("Super Admin")`)
- `POST /api/v1/admin/categories` — Tambah kategori baru.
- `PUT /api/v1/admin/categories/:id` — Ubah kategori.
- `DELETE /api/v1/admin/categories/:id` — Hapus kategori.
- `POST /api/v1/admin/products` — Tambah produk baru (stok, harga, varian, berat).
- `PUT /api/v1/admin/products/:id` — Ubah produk.
- `DELETE /api/v1/admin/products/:id` — Hapus produk.
- `GET /api/v1/admin/orders` — Daftar seluruh pesanan dari semua pelanggan.
- `PUT /api/v1/admin/orders/:id/status` — Ubah status pesanan (`PROCESSING`, `ON_DELIVERY` + input nomor resi, `COMPLETED`, `CANCELLED` + otomatis restock).

---

## ⚡ 6. PENGUJIAN BEBAN (k6 LOAD TESTING - 1.000 VUs)

Script pengujian terletak pada [`k6/load_test.js`](file:///c:/xampp/htdocs/eCommerce/k6/load_test.js) yang mensimulasikan trafik 1.000 Virtual Users serentak.

### Menjalankan Load Test via Docker:
```bash
docker run --rm -i --network="host" grafana/k6 run - < k6/load_test.js
```

Atau jika k6 terpasang secara lokal:
```bash
k6 run k6/load_test.js
```

### Kriteria Kelulusan Performa:
- **`http_req_duration p(95) < 200ms`** (Didukung query caching di Redis).
- **`http_req_failed < 1%`** (Tingkat error di bawah 1% pada puncak 1.000 VU).
- **`product_list_duration p(95) < 100ms`** (Akselerasi cache-aside katalog produk).

---

## 📂 7. STRUKTUR DIREKTORI PROYEK

```text
eCommerce/
├── backend/
│   ├── cmd/api/main.go                 # Entrypoint Go API & Graceful Shutdown
│   ├── config/config.go                # Loader konfigurasi Environment
│   ├── internal/
│   │   ├── domain/                     # Entities, Repository Interfaces & Errors
│   │   ├── usecase/                    # Business Logic Clean Architecture
│   │   ├── repository/
│   │   │   ├── postgres/               # Implementasi PostgreSQL & ACID Transactions
│   │   │   └── redis/                  # Implementasi Redis Catalog & Cart Cache
│   │   ├── delivery/http/
│   │   │   ├── handler/                # Fiber HTTP Request Handlers
│   │   │   ├── middleware/             # JWT & RBAC "Super Admin" Middlewares
│   │   │   └── router.go               # Pengelompokan Rute API v1
│   │   └── pkg/
│   │       ├── hash/                   # Bcrypt Hasher
│   │       ├── jwt/                    # JWT Token Generator & Validator
│   │       ├── midtrans/               # Midtrans Snap Client & SHA-512 Verifier
│   │       └── response/               # Standardized JSON Response Formatter
│   ├── migrations/                     # 11 Tabel DDL & Seeder SQL Files
│   ├── Dockerfile                      # Multi-Stage Go Build
│   └── go.mod
├── frontend/
│   ├── assets/css/main.css             # Tailwind CSS & Design System
│   ├── components/                     # Navbar, CartDrawer, ProductCard, AdminSidebar, Footer
│   ├── composables/useApi.ts           # Unified Fetching Client with JWT
│   ├── layouts/                        # Default & Admin Layouts
│   ├── pages/                          # Customer & Admin Dashboard Pages
│   ├── plugins/midtrans.client.ts      # Client Plugin Snap.js
│   ├── stores/                         # Pinia Auth & Cart Stores
│   ├── types/index.ts                  # TypeScript Strict Contracts
│   ├── Dockerfile                      # Multi-Stage Nuxt Nitro Build
│   ├── nuxt.config.ts
│   └── package.json
├── k6/
│   └── load_test.js                    # k6 Benchmark 1.000 VUs
├── docker-compose.yml                  # Postgres, Redis, Go Backend, Nuxt Frontend
├── Makefile                            # Orchestration Shortcuts
└── README.md
```

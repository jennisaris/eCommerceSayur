# Refactoring Progress

Tanggal audit dan refactoring: 2026-11-09

## Fase 1 — Audit

- Audit source code, database, arsitektur, performa, dan keamanan dilakukan sebelum perubahan.
- Baseline backend: `go test ./...` dan `go vet ./...` lulus.
- Baseline frontend: production build lulus; pemeriksaan TypeScript menemukan konfigurasi `vite.server.host` yang tidak valid.
- Backup database sebelum migration:
  `C:\Users\ARIEFK~1\AppData\Local\Temp\ecommerce-audit-backups\20261108\ecommerce_db_before_refactor.dump`
- Backup tambahan sebelum auth session migration:
  `C:\Users\ARIEFK~1\AppData\Local\Temp\ecommerce-audit-backups\20261108\ecommerce_db_before_auth_sessions.dump`

## Fase 2 — Refactoring

### Critical

- Menghapus kredensial demo dari UI, README, dan script pengujian.
- Seeder development tidak lagi otomatis dipasang ke container production; `RUN_SEEDER` production dikunci `false`.
- Expiration JWT production diturunkan menjadi default 2 jam.
- Migration database tidak lagi dieksekusi sebagai init script container dan sekarang gagal-fast saat error.

### High

- Otorisasi role admin memuat role terbaru dari database, sehingga pencabutan role tidak bergantung pada claim JWT lama.
- Menambahkan idempotensi webhook pembayaran berbasis `event_key` SHA-256 dan unique partial index.
- Memvalidasi nominal, mata uang, dan transaction ID webhook sebelum memperbarui pembayaran.
- Restock order cancellation dilindungi lock dan tidak mengulang restock jika order sudah cancelled.
- Menambahkan index foreign key, index query utama, dan unique index gateway payment.

### Medium

- Payload webhook tidak lagi ditulis penuh ke log; hanya ukuran payload dan error internal yang dicatat.
- Response error webhook dibuat generik agar detail internal tidak dikirim ke client.
- Pagination product dan order dibatasi maksimum 100 item per halaman dengan default yang aman.
- Akses produk inactive dipindahkan ke route admin terproteksi; katalog publik selalu membatasi status `ACTIVE`.
- Menghapus properti Vite yang tidak valid agar TypeScript frontend dapat diperiksa.

### Sisa High/Medium yang kemudian ditutup

- Token akses dipindahkan dari `localStorage` ke HttpOnly/Secure cookie dengan access-token berumur pendek.
- Refresh token dirotasi dan direvoke melalui tabel `auth_sessions`; logout mencabut refresh session.
- CSRF double-submit token diterapkan untuk seluruh request mutating berbasis cookie.
- Password registration sekarang wajib 12-72 karakter dengan huruf besar, huruf kecil, angka, dan simbol.
- Production mewajibkan verifikasi email melalui SMTP; token verifikasi disimpan sebagai hash dan memiliki masa berlaku 24 jam.
- Error 500 pada handler dipetakan ke pesan generik dan detail hanya dicatat di server.
- Production config menolak DB superuser, SSL database disabled, JWT default, dan SMTP verification yang belum dikonfigurasi.
- Docker runtime backend dan frontend dijalankan sebagai non-root; Nginx menambahkan CSP dan security headers.

## Fase 3 — Testing & Validasi

- `go test ./...` — lulus.
- `go vet ./...` — lulus.
- `npx tsc --noEmit` — lulus.
- `npm.cmd run build` — lulus.
- Migration `000003_security_indexes` dan `000004_payment_webhook_idempotency` dijalankan transactional dan diuji ulang secara idempotent.
- Migration `000005_auth_sessions` dijalankan transactional dan diuji ulang secara idempotent.
- Row count lokal setelah migration tetap konsisten: users 4, products 10, orders 6, payments 6, payment_logs 1.
- Smoke test lokal cookie auth lulus: CSRF bootstrap, register, `/auth/me`, refresh rotation, logout, dan penolakan session setelah logout.
- Tidak ada deployment production.

## Prasyarat operasional deployment

- Sediakan sertifikat TLS/domain pada reverse proxy sebelum membuka port HTTPS production.
- Sediakan `DB_USER` least-privilege, `DB_SSLMODE=require` atau mode verifikasi sertifikat, dan seluruh secret production.
- Sediakan `SMTP_HOST`, `SMTP_USERNAME`, `SMTP_PASSWORD`, dan `SMTP_FROM`; aplikasi production akan menolak start jika kosong.
- Jalankan integration/E2E audit dengan credential melalui environment variable setelah service production-like tersedia.

Prasyarat tersebut tidak dapat diisi otomatis tanpa akses ke domain, sertifikat, credential database, dan credential SMTP milik deployment. Tidak ada deployment production yang dilakukan.

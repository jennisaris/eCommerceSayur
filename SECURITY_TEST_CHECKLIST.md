# Security Verification & Test Checklist (16 Poin)

Dokumen ini berisi daftar uji keamanan lengkap untuk seluruh **16 poin keamanan** BookVerse eCommerce beserta cakupan tes otomatis dan langkah uji manual (termasuk pengujian menggunakan 2 akun berbeda).

---

## Cara Menjalankan Seluruh Tes Otomatis

```powershell
# 1. Backend: Unit Test + Security E2E Integration Test + Static Analysis
cd c:\xampp\htdocs\eCommerce\backend
go test -v ./...
go vet ./...

# 2. Backend: Vulnerability Scan (Go modules & stdlib)
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# 3. Frontend: TypeScript Check & Dependency Audit
cd c:\xampp\htdocs\eCommerce\frontend
npx.cmd tsc --noEmit
npm.cmd audit
```

---

## Daftar Tes 16 Poin Keamanan

| # | Poin Keamanan | Implementasi / Bukti | Tes Otomatis | Langkah Uji Manual |
|---|---|---|---|---|
| **1** | **Kunci API & Secret** | Secret (`MIDTRANS_SERVER_KEY`, `JWT_SECRET`, DB/Redis/SMTP) hanya di server; `.dockerignore` ditambahkan di `backend/` & `frontend/`; production fail-fast menolak placeholder secret; fallback dummy dihapus di `nuxt.config.ts`. | `backend/internal/config/config_test.go` (`TestLoadConfigRejectsPlaceholderMidtransKeysInProduction`, `TestLoadConfigRejectsDefaultDBAndEmptyRedisPasswordInProduction`) | Inspect bundle browser (DevTools → Sources/Network): pastikan tidak ada `Mid-server-*`, `JWT_SECRET`, atau password DB/SMTP. Pastikan prefix Midtrans di `.env` lokal sesuai Sandbox (`SB-Mid-*`) sebelum deploy. |
| **2** | **Login, Session, Logout, Reset Password, Kedaluwarsa** | Cookie `HttpOnly`, `SameSite=Lax`, `Secure` (wajib `true` di production); access token JWT 15 menit + refresh token 30 hari dengan rotasi di tabel `auth_sessions`; logout mencabut sesi; fitur `/forgot-password` & `/reset-password` mencabut seluruh sesi lama user. | `backend/internal/usecase/auth_password_test.go` (`TestPasswordResetFlowAndSessionRevocation`), `backend/internal/delivery/http/security_integration_test.go` (`TestSecurityE2E_LoginSessionLogoutAndPasswordReset`) | Daftar akun → login → cek cookie di DevTools (`HttpOnly=✓`, `SameSite=Lax`) → klik **Lupa Password?** → reset password → pastikan password lama ditolak dan sesi browser lama gugur. |
| **3** | **Akses Data (IDOR) & RBAC Admin** | Kepemilikan `cart`, `orders`, dan `payments` dicek di server per `user_id`; `RequireRole` & `GetOrderByCode` memeriksa role `Super Admin` langsung ke PostgreSQL (bukan claim JWT basi); katalog publik terkunci aktif; route admin terpisah di `/api/v1/admin/*`. | `backend/internal/usecase/access_control_test.go`, `backend/internal/delivery/http/middleware/role_middleware_test.go`, `security_integration_test.go` (`TestSecurityE2E_TwoAccountDataIsolationAndAdminRBAC`) | **Tes 2 Akun Berbeda:** Login Akun A (browser normal) → buat pesanan `INV-...`. Login Akun B (incognito) → coba buka `/orders/INV-...` atau panggil `GET /api/v1/orders/INV-...` & `POST /api/v1/payments/INV-.../snap` → harus **`403 Forbidden`**. Coba akses `/api/v1/admin/orders` dari Akun B → harus **`403 Forbidden`**. |
| **4** | **Password Hashing** | Seluruh password di-hash dengan `golang.org/x/crypto/bcrypt` cost `12`; validasi batas `72` UTF-8 byte; field `User.Password` ber-tag `json:"-"`. | `backend/internal/pkg/hash/hash_test.go`, `backend/internal/usecase/auth_password_test.go` (`TestHasStrongPassword`) | Cek tabel `users` di PostgreSQL: kolom `password` harus selalu diawali `$2a$10$` / `$2a$12$` (tidak pernah plain text). |
| **5** | **Validasi Input di Server** | Validator terpusat berbasis `validator/v10` (`safe_url`, `slug`, `phone`, `latitude`, `longitude`, batas angka & panjang field, whitelist query `sort`/`status`, validasi ID `> 0`, stok kumulatif cart, Fiber `BodyLimit` 1 MB). | `backend/internal/pkg/validate/validate_test.go`, `backend/internal/delivery/http/handler/validation_test.go` | Kirim `GET /api/v1/products?min_price=-10` → `400`; kirim `POST /api/v1/admin/products` dengan `image: "javascript:alert(1)"` → `400`; kirim `POST /api/v1/cart/items` dengan `product_id: -1` atau `quantity: 5000` → `400`. |
| **6** | **Escape Output (XSS)** | Nol pemakaian `v-html`/`innerHTML`/`dangerouslySetInnerHTML` di kode aplikasi; seluruh teks memakai auto-escape Vue `{{ ... }}`; URL redirect Midtrans & preview gambar admin dibatasi melalui `frontend/utils/safeUrl.ts`. | `backend/internal/pkg/validate/validate_test.go` (`TestSafeURLValidation`), `npx tsc --noEmit` | Input nama/alamat/catatan pesanan berisi `<script>alert(1)</script>` → pastikan tampil sebagai teks biasa (tidak ada popup alert). |
| **7** | **Query Database (SQL Injection)** | 100% query memakai GORM ORM + parameterized placeholders (`?`); wildcard `ILIKE` di-escape (`ESCAPE '\'`). | `backend/internal/repository/postgres/like_test.go` | Kirim `GET /api/v1/products?search=' OR '1'='1` dan `?search=%25` → tidak menyebabkan error SQL dan tidak mengembalikan seluruh baris. |
| **8** | **Rate Limit (`429`)** | Fiber limiter aktif di login/register/refresh/verify-email (`10/1m`), reset password (`5/15m`), order/estimate (`30/1m`), Midtrans Snap (`15/1m` per user/IP); `EnableTrustedProxyCheck` aktif; edge `limit_req_zone` di `nginx.conf`. | `backend/internal/delivery/http/rate_limit_test.go` (`TestIPRateLimiterReturns429WhenLimitExceeded`, `TestUserOrIPRateLimiterScopesPerUser`) | Kirim >10 percobaan login gagal berturut-turut dalam 1 menit → respons ke-11 harus **`429 Too Many Requests`**. |
| **9** | **Security Headers & HTTPS** | Fiber `helmet` + Nitro `server/middleware/security-headers.ts` + `nginx.conf` memasang `CSP`, `HSTS` (production), `X-Frame-Options`, `X-Content-Type-Options: nosniff`, `Referrer-Policy`; production config menolak `CLIENT_BASE_URL` non-`https://`. | `backend/internal/delivery/http/security_headers_test.go`, `backend/internal/config/config_test.go` (`TestLoadConfigRejectsNonHTTPSClientURLInProduction`) | Cek header response (`curl -I http://localhost:8080/health` & `http://127.0.0.1:3000`); saat deploy production pasang sertifikat TLS & aktifkan redirect `301 https://` di `nginx.conf`. |
| **10** | **CSRF & Larangan Mutasi via GET** | Double-submit `HttpOnly` + `SameSite=Lax` cookie + header `X-CSRF-Token`, verifikasi `Origin`/`Referer` + blokir `Sec-Fetch-Site: cross-site`; `VerifyEmail` dipindah ke `POST`. | `backend/internal/delivery/http/middleware/csrf_middleware_test.go`, `security_integration_test.go` (`TestSecurityE2E_PaymentWebhookAndCSRFProtection`) | Kirim `POST /api/v1/orders` tanpa `X-CSRF-Token` atau dengan header `Origin: https://evil.example` → **`403 Forbidden`**. Kirim `GET /api/v1/auth/verify-email?token=abc` → **`405 Method Not Allowed`**. |
| **11** | **Upload File** | Tidak ada fitur upload file aktif (`N/A`); batas ukuran body diperketat (`1 MB` Fiber, `2 MB` Nginx). | — | Verifikasi tidak ada endpoint multipart upload yang terbuka. |
| **12** | **Harga, Diskon, Ongkir & Webhook** | Harga, diskon, berat, ongkir, dan total dihitung atomik di dalam transaksi DB `FOR UPDATE`; verifikasi SHA-512 Midtrans constant-time (`subtle.ConstantTimeCompare`), idempotensi `event_key`, validasi nominal/mata uang, dan state transition guard. | `backend/internal/pkg/midtrans/midtrans_test.go`, `backend/internal/usecase/order_shipping_test.go`, `backend/internal/usecase/payment_webhook_test.go`, `security_integration_test.go` | Kirim webhook Midtrans dengan `signature_key` palsu → `401`; ubah `gross_amount` → `400`; pastikan `expire` terlambat tidak membatalkan order `PAID`. |
| **13** | **Library & Dependency Audit** | 7 modul Go rentan diperbarui (`0` kerentanan modul pihak ketiga di `govulncheck`); `lucide-vue-next` dihapus; `npm audit fix` menutup seluruh kerentanan High & Low (`0 critical`, `0 high`). | `govulncheck ./...`, `npm audit`, `go test ./...`, `npx tsc --noEmit` | Update Go host ke `1.27.4+` untuk menutup patch stdlib Go; evaluasi migrasi Nuxt 4 di jadwal terpisah bila ingin menutup 8 moderate `nuxt 3.x`. |
| **14** | **Riwayat Git & `.gitignore`** | `.env` masuk `.gitignore` dan tidak pernah ter-commit; nol kunci API live pihak ketiga di seluruh riwayat Git. | Scan `git rev-list --all` | Opsional: jalankan `gitleaks detect --source .` secara berkala. Jangan pernah gunakan secret default/akun demo seeder di production. |
| **15** | **Logging, Monitoring & Alerting** | Structured JSON audit log (`log/slog`) di `backend/internal/pkg/seclog` (email di-mask, tanpa password/token) + deteksi & alert aktivitas mencurigakan (≥5 login gagal / 5 menit) + dukungan `SECURITY_ALERT_WEBHOOK_URL`. | `backend/internal/pkg/seclog/seclog_test.go` | Lakukan 5x login gagal berturut-turut → pastikan muncul log JSON `category:"security_alert"` di console server tanpa password/token. |
| **16** | **Tes Keamanan Otomatis & Manual** | Full-router security E2E integration tests untuk Login, Akses Data (2 akun berbeda + Admin), Pembayaran/Webhook, dan CSRF/Origin. | `backend/internal/delivery/http/security_integration_test.go` | Jalankan `go test -v ./internal/delivery/http -run TestSecurityE2E` setelah setiap perubahan fitur/keamanan. |

---

## Prasyarat Manual Sebelum Go-Live Production
1. Pasang sertifikat TLS domain & aktifkan blok port `443` + redirect `301 https://` di `docker/nginx/nginx.conf` dan `docker-compose.prod.yml`.
2. Buat role PostgreSQL non-superuser least-privilege (`DB_USER` bukan `postgres`) dan pastikan `DB_SSLMODE=require` / `verify-full`.
3. Isi semua secret unik production (`JWT_SECRET` ≥ 32 karakter acak, `DB_PASSWORD`, `REDIS_PASSWORD`, `MIDTRANS_SERVER_KEY`/`CLIENT_KEY` production, SMTP, dan opsional `SECURITY_ALERT_WEBHOOK_URL`).
4. Perbarui instalasi Go host ke `1.27.4+` agar patch stdlib Go terbaru aktif.

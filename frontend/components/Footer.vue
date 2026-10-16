<template>
  <footer class="bg-gray-100 text-gray-700 border-t border-gray-200 pt-12 pb-8 mt-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 pb-8 border-b border-gray-200">
        <!-- Brand Info -->
        <div class="space-y-4 md:col-span-1">
          <NuxtLink to="/" class="inline-block">
            <img
              src="/assets/images/logo/freshcart-logo.svg"
              alt="Sayur FreshCart"
              class="h-8 w-auto"
            />
          </NuxtLink>
          <p class="text-xs text-gray-500 leading-relaxed">
            Pusat belanja sayuran segar, buah organik, dan bahan makanan 24 Jam Non-Stop dengan 100% Home Delivery Model langsung ke alamat Anda.
          </p>
          <div class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-green-50 border border-green-200 text-green-700 text-xs font-bold">
            <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
            Siap Kirim 24/7 (Non-Stop)
          </div>
        </div>

        <!-- Categories Links -->
        <div>
          <h4 class="text-sm font-bold text-gray-900 mb-3">Kategori Produk</h4>
          <ul class="space-y-2 text-xs text-gray-500">
            <li><NuxtLink to="/products" class="hover:text-green-600 transition-colors">Sayuran & Bumbu Dapur</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-green-600 transition-colors">Buah-buahan Segar</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-green-600 transition-colors">Susu, Roti & Telur</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-green-600 transition-colors">Daging, Ayam & Ikan</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-green-600 transition-colors">Beras, Minyak & Bumbu</NuxtLink></li>
          </ul>
        </div>

        <!-- Delivery & Service Policy -->
        <div>
          <h4 class="text-sm font-bold text-gray-900 mb-3">Layanan Pengiriman</h4>
          <ul class="space-y-2 text-xs text-gray-500">
            <li class="flex items-center gap-1.5"><i class="ti ti-check text-green-600"></i> 100% Model Home Delivery</li>
            <li class="flex items-center gap-1.5"><i class="ti ti-check text-green-600"></i> Bebas Antre Tanpa Opsi Pickup</li>
            <li class="flex items-center gap-1.5"><i class="ti ti-check text-green-600"></i> Kurir Terlatih & Berpendingin</li>
            <li class="flex items-center gap-1.5"><i class="ti ti-check text-green-600"></i> Garansi Sayur Segar 100%</li>
          </ul>
        </div>

        <!-- Payment Partners & Credentials -->
        <div>
          <h4 class="text-sm font-bold text-gray-900 mb-3">Metode Pembayaran</h4>
          <p class="text-xs text-gray-500 mb-3">
            Didukung oleh Midtrans Payment Gateway resmi: QRIS, BCA, Mandiri, BNI, BRI Virtual Account, dan Kartu Kredit.
          </p>
          <div class="flex items-center flex-wrap gap-2 mb-4">
            <img src="/assets/images/payment/visa.svg" alt="Visa" class="h-6" />
            <img src="/assets/images/payment/mastercard.svg" alt="MasterCard" class="h-6" />
            <img src="/assets/images/payment/paypal.svg" alt="Paypal" class="h-6" />
            <img src="/assets/images/payment/american-express.svg" alt="Amex" class="h-6" />
          </div>
          <div class="p-2.5 rounded-lg bg-gray-200/60 text-[11px] text-gray-600 font-mono">
            Demo Super Admin:<br />
            admin@ecommerce.com / Admin123!
          </div>
        </div>
      </div>

      <div class="pt-6 flex flex-col sm:flex-row items-center justify-between text-xs text-gray-500 gap-4">
        <p>&copy; 2026 Sayur FreshCart E-Commerce. All rights reserved.</p>
        <div class="flex items-center gap-4">
          <span class="hover:text-green-600 cursor-pointer">Privacy Policy</span>
          <span>&bull;</span>
          <span class="hover:text-green-600 cursor-pointer">Terms of Service</span>
          <span>&bull;</span>
          <span class="text-gray-400">Powered by Golang Clean Architecture & Nuxt 3</span>
        </div>
      </div>
    </div>
  </footer>
</template>

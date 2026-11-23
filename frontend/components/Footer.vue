<template>
  <footer class="bg-slate-900 text-slate-300 border-t border-slate-800 pt-14 pb-8 mt-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 pb-10 border-b border-slate-800">
        <!-- Brand Info -->
        <div class="space-y-4 md:col-span-1">
          <NuxtLink to="/" class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-xl bg-indigo-600 flex items-center justify-center text-white text-base">
              📖
            </div>
            <span class="font-extrabold text-lg text-white tracking-tight">
              Book<span class="text-indigo-400">Verse</span>
            </span>
          </NuxtLink>
          <p class="text-xs text-slate-400 leading-relaxed">
            Toko buku online modern terlengkap dan terpercaya. Menghadirkan ribuan koleksi novel fiksi, buku self-improvement, bisnis, akademik, dan manga original bergaransi langsung dari penerbit resmi.
          </p>
          <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-indigo-950/80 border border-indigo-700/50 text-indigo-300 text-xs font-semibold">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            100% Buku Segel Original
          </div>
        </div>

        <!-- Categories Links -->
        <div>
          <h4 class="text-sm font-bold text-white mb-3 tracking-wide">Kategori Pilihan</h4>
          <ul class="space-y-2 text-xs text-slate-400">
            <li><NuxtLink to="/products" class="hover:text-indigo-400 transition-colors">Fiksi & Sastra Populer</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-indigo-400 transition-colors">Pengembangan Diri & Karir</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-indigo-400 transition-colors">Bisnis, Investasi & Keuangan</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-indigo-400 transition-colors">Teknologi, AI & Programming</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:text-indigo-400 transition-colors">Komik & Manga Berlisensi</NuxtLink></li>
          </ul>
        </div>

        <!-- Delivery & Service Policy -->
        <div>
          <h4 class="text-sm font-bold text-white mb-3 tracking-wide">Jaminan Layanan</h4>
          <ul class="space-y-2.5 text-xs text-slate-400">
            <li class="flex items-center gap-2"><i class="ti ti-shield-check text-emerald-400 text-sm"></i> Garansi 100% Original & Bersegel</li>
            <li class="flex items-center gap-2"><i class="ti ti-package text-emerald-400 text-sm"></i> Packing Aman Bubble Wrap + Kardus</li>
            <li class="flex items-center gap-2"><i class="ti ti-truck-delivery text-emerald-400 text-sm"></i> Pengiriman Ekspres Seluruh Indonesia</li>
            <li class="flex items-center gap-2"><i class="ti ti-receipt text-emerald-400 text-sm"></i> Invoice Pembelian Resmi Terbit Otomatis</li>
          </ul>
        </div>

        <!-- Payment Partners & Credentials -->
        <div>
          <h4 class="text-sm font-bold text-white mb-3 tracking-wide">Metode Pembayaran</h4>
          <p class="text-xs text-slate-400 mb-3">
            Didukung oleh Midtrans Payment Gateway resmi: QRIS, GoPay, OVO, ShopeePay, Virtual Account Bank (BCA, Mandiri, BNI, BRI), dan Kartu Kredit.
          </p>
          <div class="flex items-center flex-wrap gap-2 mb-4 bg-slate-800/60 p-2.5 rounded-xl border border-slate-700/50">
            <img src="/assets/images/payment/visa.svg" alt="Visa" class="h-5" />
            <img src="/assets/images/payment/mastercard.svg" alt="MasterCard" class="h-5" />
            <img src="/assets/images/payment/paypal.svg" alt="Paypal" class="h-5" />
            <img src="/assets/images/payment/american-express.svg" alt="Amex" class="h-5" />
          </div>
        </div>
      </div>

      <div class="pt-6 flex flex-col sm:flex-row items-center justify-between text-xs text-slate-500 gap-4">
        <p>&copy; 2026 BookVerse Indonesia. Hak Cipta Dilindungi Undang-Undang.</p>
        <div class="flex items-center gap-4">
          <span class="hover:text-indigo-400 cursor-pointer">Kebijakan Privasi</span>
          <span>&bull;</span>
          <span class="hover:text-indigo-400 cursor-pointer">Syarat & Ketentuan</span>
          <span>&bull;</span>
          <span class="text-slate-400 font-medium">Powered by Golang Clean Architecture & Nuxt 3</span>
        </div>
      </div>
    </div>
  </footer>
</template>

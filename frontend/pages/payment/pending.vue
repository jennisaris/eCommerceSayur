<template>
  <div class="max-w-2xl mx-auto px-4 py-16 sm:py-24 text-center">
    <div class="w-20 h-20 rounded-full bg-amber-100 text-amber-600 flex items-center justify-center mx-auto mb-6 shadow-lg shadow-amber-500/20">
      <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </div>

    <span class="inline-block px-3 py-1 text-xs font-bold bg-amber-50 text-amber-700 border border-amber-200 rounded-full mb-3">
      Menunggu Pembayaran
    </span>

    <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 mb-2">
      Instruksi Pembayaran Dikirim
    </h1>
    <p class="text-xs sm:text-sm text-slate-600 max-w-md mx-auto mb-6 leading-relaxed">
      Selesaikan pembayaran pesanan <span class="font-bold text-slate-900">{{ orderCode }}</span> sesuai nomor Virtual Account atau kode QRIS sebelum batas waktu habis agar stok tetap aman.
    </p>

    <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
      <NuxtLink
        :to="`/orders/${orderCode}`"
        class="w-full sm:w-auto px-6 py-3 bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs rounded-xl shadow-md transition-all"
      >
        Lihat Detail & Bayar
      </NuxtLink>
      <NuxtLink
        to="/"
        class="w-full sm:w-auto px-6 py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs rounded-xl transition-all"
      >
        Kembali ke Beranda
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const orderCode = computed(() => (route.query.order_code as string) || '')
</script>

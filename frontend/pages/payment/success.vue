<template>
  <div class="max-w-2xl mx-auto px-4 py-16 sm:py-24 text-center">
    <div class="w-20 h-20 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center mx-auto mb-6 shadow-lg shadow-emerald-500/20">
      <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
      </svg>
    </div>

    <span class="inline-block px-3 py-1 text-xs font-bold bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-full mb-3">
      Pembayaran Sukses Terverifikasi
    </span>

    <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 mb-2">
      Terima Kasih Atas Pesanan Anda!
    </h1>
    <p class="text-xs sm:text-sm text-slate-600 max-w-md mx-auto mb-6 leading-relaxed">
      Pesanan dengan kode <span class="font-bold text-slate-900">{{ orderCode }}</span> telah berhasil dibayar. Tim kami segera menyiapkan paket Anda untuk pengantaran.
    </p>

    <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
      <NuxtLink
        :to="`/orders/${orderCode}`"
        class="w-full sm:w-auto px-6 py-3 bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs rounded-xl shadow-md transition-all"
      >
        Lacak Status Pesanan
      </NuxtLink>
      <NuxtLink
        to="/products"
        class="w-full sm:w-auto px-6 py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs rounded-xl transition-all"
      >
        Lanjut Belanja
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const orderCode = computed(() => (route.query.order_code as string) || '')
</script>

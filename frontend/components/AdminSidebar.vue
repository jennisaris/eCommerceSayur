<template>
  <aside class="w-64 bg-slate-900 text-white min-h-screen flex flex-col justify-between p-4 flex-shrink-0">
    <div class="space-y-6">
      <!-- Admin Header with FreshCart Logo -->
      <div class="flex flex-col gap-2 px-2 py-3 border-b border-slate-800">
        <NuxtLink to="/" class="flex items-center gap-2">
          <img
            src="/assets/images/logo/freshcart-logo.svg"
            alt="FreshCart Logo"
            class="h-7 w-auto bg-white p-1 rounded-md"
          />
        </NuxtLink>
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold text-gray-300">Admin Console</span>
          <span class="text-[10px] text-green-400 font-bold px-1.5 py-0.5 rounded bg-green-950 border border-green-800 uppercase tracking-wider">
            Super Admin
          </span>
        </div>
      </div>

      <!-- Navigation Menu -->
      <nav class="space-y-1 text-xs">
        <NuxtLink
          to="/admin"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg font-medium text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
          active-class="bg-green-600 text-white hover:bg-green-600 font-bold"
        >
          <i class="ti ti-dashboard text-base"></i>
          <span>Dashboard Overview</span>
        </NuxtLink>

        <NuxtLink
          to="/admin/categories"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg font-medium text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
          active-class="bg-green-600 text-white hover:bg-green-600 font-bold"
        >
          <i class="ti ti-category text-base"></i>
          <span>Categories</span>
        </NuxtLink>

        <NuxtLink
          to="/admin/products"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg font-medium text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
          active-class="bg-green-600 text-white hover:bg-green-600 font-bold"
        >
          <i class="ti ti-package text-base"></i>
          <span>Products</span>
        </NuxtLink>

        <NuxtLink
          to="/admin/orders"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg font-medium text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
          active-class="bg-green-600 text-white hover:bg-green-600 font-bold"
        >
          <i class="ti ti-shopping-cart text-base"></i>
          <span>Orders & Delivery</span>
        </NuxtLink>
      </nav>
    </div>

    <!-- Bottom Actions -->
    <div class="space-y-2 pt-4 border-t border-slate-800 text-xs">
      <NuxtLink
        to="/"
        class="flex items-center gap-2.5 px-3 py-2 rounded-lg font-medium text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
      >
        <i class="ti ti-arrow-left text-sm"></i>
        <span>Kembali ke Toko</span>
      </NuxtLink>

      <button
        @click="authStore.logout"
        class="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg font-medium text-red-400 hover:text-red-300 hover:bg-red-500/10 transition-colors"
      >
        <i class="ti ti-logout text-sm"></i>
        <span>Keluar (Logout)</span>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
</script>

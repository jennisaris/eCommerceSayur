<template>
  <header class="sticky top-0 z-40 w-full bg-white border-b border-slate-200 transition-all shadow-xs">
    <!-- Top Announcement Bar -->
    <div class="bg-indigo-950 text-indigo-100 text-[11px] py-1.5 px-4 text-center font-medium flex items-center justify-center gap-2">
      <span>📚 <strong>Pesta Literasi Nasional:</strong> Diskon hingga 35% untuk Buku Best Seller & Bebas Ongkir se-Indonesia!</span>
      <NuxtLink to="/products" class="underline font-bold text-amber-300 hover:text-amber-200 ml-1">Jelajahi Sekarang &rarr;</NuxtLink>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-20 gap-4">
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center gap-3 group flex-shrink-0">
          <div class="w-10 h-10 rounded-2xl bg-indigo-600 flex items-center justify-center text-white shadow-md shadow-indigo-600/30 text-xl font-bold group-hover:scale-105 transition-transform">
            📖
          </div>
          <div class="flex flex-col">
            <span class="font-extrabold text-xl tracking-tight text-slate-900 group-hover:text-indigo-600 transition-colors">
              Book<span class="text-indigo-600">Verse</span>
            </span>
            <span class="text-[10px] text-slate-500 font-semibold tracking-wider uppercase -mt-0.5">
              Toko Buku Original
            </span>
          </div>
          <span class="hidden lg:inline-block px-2.5 py-0.5 text-[10px] font-extrabold bg-indigo-50 text-indigo-700 rounded-full border border-indigo-200 uppercase tracking-wider ml-1">
            100% Original
          </span>
        </NuxtLink>

        <!-- Search Bar with Tabler Icon -->
        <div class="hidden md:flex flex-1 max-w-xl mx-4">
          <form @submit.prevent="handleSearch" class="relative w-full">
            <input
              v-model="searchQuery"
              type="search"
              placeholder="Cari judul buku, nama penulis (e.g. James Clear), penerbit, atau ISBN..."
              class="w-full pl-10 pr-4 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-600 transition-all placeholder:text-slate-400"
            />
            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M10 10m-7 0a7 7 0 1 0 14 0a7 7 0 1 0 -14 0" />
                <path d="M21 21l-6 -6" />
              </svg>
            </div>
          </form>
        </div>

        <!-- Navigation Links & Actions -->
        <div class="flex items-center gap-4 sm:gap-6">
          <NuxtLink
            to="/products"
            class="hidden sm:inline-flex text-xs font-bold text-slate-700 hover:text-indigo-600 transition-colors"
          >
            Katalog Buku
          </NuxtLink>

          <!-- Super Admin Button if logged in as Admin -->
          <NuxtLink
            v-if="authStore.isSuperAdmin"
            to="/admin"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-bold bg-amber-500 text-white rounded-lg hover:bg-amber-600 shadow-xs transition-all"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Dashboard Admin
          </NuxtLink>

          <!-- Cart Bag Trigger Button (with Tabler Shopping Bag Icon) -->
          <button
            @click="cartStore.toggleDrawer"
            type="button"
            class="relative p-2 rounded-full text-slate-700 hover:text-indigo-600 hover:bg-slate-100 transition-all focus:outline-none"
            aria-label="Keranjang Buku"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
              <path d="M6.331 8h11.339a2 2 0 0 1 1.977 2.304l-1.255 8.152a3 3 0 0 1 -2.966 2.544h-6.852a3 3 0 0 1 -2.965 -2.544l-1.255 -8.152a2 2 0 0 1 1.977 -2.304z" />
              <path d="M9 11v-5a3 3 0 0 1 6 0v5" />
            </svg>
            <!-- Badge -->
            <span
              v-if="cartStore.totalItems > 0"
              class="absolute -top-1 -right-1 bg-indigo-600 text-white text-[11px] font-extrabold h-5 min-w-[20px] px-1 rounded-full flex items-center justify-center ring-2 ring-white"
            >
              {{ cartStore.totalItems }}
            </span>
          </button>

          <!-- User Menu / Auth Buttons -->
          <div v-if="authStore.isAuthenticated" class="relative group">
            <button class="flex items-center gap-2 p-1 rounded-full hover:bg-slate-100 transition-colors">
              <img
                :src="authStore.user?.photo || '/assets/images/avatar/avatar-1.jpg'"
                alt="Avatar"
                class="w-8 h-8 rounded-full object-cover ring-2 ring-indigo-600/30"
              />
              <span class="hidden md:inline-block text-xs font-semibold text-slate-700 max-w-[120px] truncate">
                {{ authStore.user?.name }}
              </span>
            </button>
            <!-- Dropdown -->
            <div class="absolute right-0 mt-1 w-52 bg-white border border-slate-200 rounded-2xl shadow-xl py-1.5 hidden group-hover:block transition-all z-50">
              <div class="px-4 py-2 border-b border-slate-100">
                <p class="text-xs font-bold text-slate-900 truncate">{{ authStore.user?.name }}</p>
                <p class="text-[11px] text-slate-500 truncate">{{ authStore.user?.email }}</p>
              </div>
              <NuxtLink to="/orders" class="block px-4 py-2 text-xs text-slate-700 hover:bg-slate-50 font-medium">
                Riwayat Pesanan Buku
              </NuxtLink>
              <button
                @click="authStore.logout"
                class="w-full text-left px-4 py-2 text-xs text-rose-600 hover:bg-rose-50 font-semibold"
              >
                Keluar (Logout)
              </button>
            </div>
          </div>

          <div v-else class="flex items-center gap-2">
            <NuxtLink
              to="/login"
              class="px-3 py-1.5 text-xs font-bold text-slate-700 hover:text-indigo-600 transition-colors"
            >
              Sign In
            </NuxtLink>
            <NuxtLink
              to="/register"
              class="px-3.5 py-1.5 text-xs font-bold bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 shadow-xs transition-all"
            >
              Sign Up
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const cartStore = useCartStore()
const router = useRouter()
const searchQuery = ref('')

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      path: '/products',
      query: { search: searchQuery.value.trim() },
    })
  }
}

onMounted(() => {
  if (authStore.isAuthenticated) {
    cartStore.fetchCart()
  }
})
</script>

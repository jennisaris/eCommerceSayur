<template>
  <div class="space-y-12 pb-16">
    <!-- Hero Slider Section (Sayur FreshCart Hero) -->
    <section class="mt-4">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
          class="relative rounded-2xl overflow-hidden bg-cover bg-center shadow-sm"
          style="background-image: url('/assets/images/slider/slide-1.jpg'); min-height: 440px;"
        >
          <!-- Subtle Dark/Green Tint for Contrast -->
          <div class="absolute inset-0 bg-gradient-to-r from-black/50 via-black/20 to-transparent"></div>

          <div class="relative z-10 py-16 px-8 sm:px-12 lg:py-24 max-w-xl flex flex-col justify-center">
            <span class="inline-flex items-center gap-1.5 px-3 py-1 text-xs font-bold rounded-lg bg-yellow-400 text-gray-900 w-fit mb-4 shadow-sm">
              <i class="ti ti-flame text-sm"></i> Opening Sale Discount 50%
            </span>

            <h1 class="text-3xl sm:text-5xl font-extrabold text-white leading-tight mb-4 drop-shadow-sm">
              SuperMarket For <br />
              <span class="text-green-400">Fresh Grocery</span>
            </h1>

            <p class="text-white/90 text-sm sm:text-base leading-relaxed mb-6 font-medium drop-shadow-sm">
              Model pengiriman 100% Home Delivery buka 24/7 non-stop. Belanja sayur segar, buah organik, dan bahan dapur kapan saja langsung diantar ke depan pintu Anda.
            </p>

            <div class="flex flex-wrap gap-3 items-center">
              <NuxtLink
                to="/products"
                class="btn inline-flex items-center gap-2 bg-green-600 text-white font-bold text-sm px-6 py-3 rounded-xl hover:bg-green-700 shadow-lg shadow-green-600/30 transition-all"
              >
                <span>Shop Now</span>
                <i class="ti ti-arrow-right text-base"></i>
              </NuxtLink>

              <NuxtLink
                to="/products?sort=price_asc"
                class="btn inline-flex items-center gap-2 bg-white/90 hover:bg-white text-gray-800 font-bold text-sm px-6 py-3 rounded-xl shadow transition-all backdrop-blur-sm"
              >
                <span>Promo Spesial</span>
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Featured Categories Section -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h2 class="text-xl sm:text-2xl font-bold text-gray-900">Featured Categories</h2>
          <p class="text-xs text-gray-500 mt-0.5">Pilih kategori sayur dan bahan makanan berkualitas untuk keluarga Anda</p>
        </div>
        <NuxtLink to="/products" class="text-xs font-bold text-green-600 hover:text-green-700 flex items-center gap-1">
          <span>Lihat Semua</span>
          <i class="ti ti-chevron-right text-sm"></i>
        </NuxtLink>
      </div>

      <!-- Categories Grid -->
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
        <CategoryCard
          v-for="cat in displayCategories"
          :key="cat.id"
          :category="cat"
        />
      </div>
    </section>

    <!-- Double Promo Banners -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Banner 1: Fruits & Vegetables -->
        <div
          class="rounded-2xl p-8 bg-cover bg-center shadow-xs relative overflow-hidden flex flex-col justify-between min-h-[220px]"
          style="background-image: url('/assets/images/banner/grocery-banner.png')"
        >
          <div class="relative z-10 max-w-xs space-y-2">
            <h3 class="text-2xl font-extrabold text-gray-900 leading-snug">
              Fruits & Vegetables
            </h3>
            <p class="text-sm text-gray-600">
              Get Upto <span class="font-extrabold text-gray-900">30%</span> Off
            </p>
          </div>
          <div class="relative z-10 mt-6">
            <NuxtLink
              to="/products"
              class="btn inline-flex items-center gap-1.5 bg-gray-900 hover:bg-gray-800 text-white font-bold text-xs px-4 py-2.5 rounded-lg transition-colors"
            >
              <span>Shop Now</span>
              <i class="ti ti-arrow-right text-xs"></i>
            </NuxtLink>
          </div>
        </div>

        <!-- Banner 2: Freshly Baked Buns -->
        <div
          class="rounded-2xl p-8 bg-cover bg-center shadow-xs relative overflow-hidden flex flex-col justify-between min-h-[220px]"
          style="background-image: url('/assets/images/banner/grocery-banner-2.jpg')"
        >
          <div class="relative z-10 max-w-xs space-y-2">
            <h3 class="text-2xl font-extrabold text-gray-900 leading-snug">
              Freshly Baked Buns
            </h3>
            <p class="text-sm text-gray-600">
              Get Upto <span class="font-extrabold text-gray-900">25%</span> Off
            </p>
          </div>
          <div class="relative z-10 mt-6">
            <NuxtLink
              to="/products"
              class="btn inline-flex items-center gap-1.5 bg-gray-900 hover:bg-gray-800 text-white font-bold text-xs px-4 py-2.5 rounded-lg transition-colors"
            >
              <span>Shop Now</span>
              <i class="ti ti-arrow-right text-xs"></i>
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <!-- Popular Products Grid (Sayur Katalog) -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h2 class="text-xl sm:text-2xl font-bold text-gray-900">Popular Products</h2>
          <p class="text-xs text-gray-500 mt-0.5">Produk sayur dan kebutuhan pokok paling banyak dicari</p>
        </div>

        <!-- Category Filter Tabs -->
        <div class="flex items-center gap-2 overflow-x-auto pb-1 text-xs">
          <button
            @click="selectedCategory = 0"
            :class="selectedCategory === 0 ? 'bg-green-600 text-white font-bold' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
            class="px-3 py-1.5 rounded-lg transition-colors whitespace-nowrap"
          >
            Semua
          </button>
          <button
            v-for="cat in categories.slice(0, 4)"
            :key="cat.id"
            @click="selectedCategory = cat.id"
            :class="selectedCategory === cat.id ? 'bg-green-600 text-white font-bold' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
            class="px-3 py-1.5 rounded-lg transition-colors whitespace-nowrap"
          >
            {{ cat.name }}
          </button>
        </div>
      </div>

      <!-- Loading skeleton -->
      <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
        <div v-for="n in 10" :key="n" class="h-80 rounded-lg bg-gray-200 animate-pulse"></div>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredProducts.length === 0" class="text-center py-16 bg-white rounded-2xl border border-gray-200 p-8">
        <i class="ti ti-basket-cancel text-4xl text-gray-400 mb-2 inline-block"></i>
        <h3 class="text-base font-bold text-gray-800">Belum ada produk untuk kategori ini</h3>
        <p class="text-xs text-gray-500 mt-1">Silakan pilih kategori lain atau lihat seluruh katalog.</p>
      </div>

      <!-- Products Grid -->
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
        <ProductCard
          v-for="prod in filteredProducts"
          :key="prod.id"
          :product="prod"
        />
      </div>
    </section>

    <!-- Daily Best Sells Highlight Banner Section -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="relative rounded-2xl overflow-hidden border border-gray-200 bg-white shadow-xs p-6 sm:p-8 flex flex-col lg:flex-row items-center justify-between gap-8">
        <div class="space-y-4 max-w-lg">
          <span class="inline-block px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider rounded bg-green-100 text-green-700">
            Paling Direkomendasikan
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-gray-900 leading-tight">
            100% Organik & Segar Langsung Dari Petani Lokal
          </h2>
          <p class="text-xs sm:text-sm text-gray-600 leading-relaxed">
            Semua produk sayuran dan bahan dapur dipilih melalui seleksi kualitas ketat setiap hari. Nikmati kesegaran maksimal dengan jaminan kirim cepat sampai ke dapur Anda.
          </p>
          <div class="flex items-center gap-6 pt-2">
            <div>
              <div class="text-xl sm:text-2xl font-black text-green-600">24/7</div>
              <div class="text-[11px] text-gray-500 font-medium">Layanan Buka</div>
            </div>
            <div class="h-8 w-px bg-gray-200"></div>
            <div>
              <div class="text-xl sm:text-2xl font-black text-green-600">100%</div>
              <div class="text-[11px] text-gray-500 font-medium">Home Delivery</div>
            </div>
            <div class="h-8 w-px bg-gray-200"></div>
            <div>
              <div class="text-xl sm:text-2xl font-black text-green-600">Instant</div>
              <div class="text-[11px] text-gray-500 font-medium">Midtrans Snap</div>
            </div>
          </div>
        </div>

        <div class="w-full lg:w-auto flex justify-center">
          <img
            src="/assets/images/banner/banner-deal.jpg"
            alt="Fresh Grocery Deals"
            class="rounded-xl max-h-72 object-cover shadow-sm"
          />
        </div>
      </div>
    </section>

    <!-- 4 Key Value Pillars (Sayur FreshCart Signature) -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 pt-4">
        <!-- Feature 1 -->
        <div class="flex items-start gap-4 p-5 rounded-xl bg-white border border-gray-200 shadow-xs">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-lg bg-green-50">
            <img src="/assets/images/icons/clock.svg" alt="10 Minute Delivery" class="w-7 h-7" />
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-gray-900">Pengiriman 24 Jam Non-Stop</h3>
            <p class="text-xs text-gray-500 leading-relaxed">
              100% Model Home Delivery. Pesanan langsung diantar ke alamat Anda kapan saja.
            </p>
          </div>
        </div>

        <!-- Feature 2 -->
        <div class="flex items-start gap-4 p-5 rounded-xl bg-white border border-gray-200 shadow-xs">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-lg bg-yellow-50">
            <img src="/assets/images/icons/gift.svg" alt="Best Prices & Offers" class="w-7 h-7" />
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-gray-900">Harga Terbaik & Promo</h3>
            <p class="text-xs text-gray-500 leading-relaxed">
              Harga terjangkau langsung dari pemasok dengan promo diskon reguler dan hemat.
            </p>
          </div>
        </div>

        <!-- Feature 3 -->
        <div class="flex items-start gap-4 p-5 rounded-xl bg-white border border-gray-200 shadow-xs">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-lg bg-blue-50">
            <img src="/assets/images/icons/package.svg" alt="Wide Assortment" class="w-7 h-7" />
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-gray-900">Pilihan Produk Lengkap</h3>
            <p class="text-xs text-gray-500 leading-relaxed">
              Sayuran, buah, daging, ikan, bumbu dapur, hingga produk harian tersedia lengkap.
            </p>
          </div>
        </div>

        <!-- Feature 4 -->
        <div class="flex items-start gap-4 p-5 rounded-xl bg-white border border-gray-200 shadow-xs">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-lg bg-purple-50">
            <img src="/assets/images/icons/refresh-cw.svg" alt="Easy Returns" class="w-7 h-7" />
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-gray-900">Midtrans Checkout Aman</h3>
            <p class="text-xs text-gray-500 leading-relaxed">
              Verifikasi signature otomatis dengan dukungan QRIS, Virtual Account, & Kartu Kredit.
            </p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { Category, Product } from '~/types'

const api = useApi()
const categories = ref<Category[]>([])
const products = ref<Product[]>([])
const isLoading = ref(true)
const selectedCategory = ref<number>(0)

// Default sample categories if DB is not seeded yet
const fallbackCategories: Category[] = [
  { id: 1, name: 'Dairy, Bread & Eggs', slug: 'dairy-bread-eggs', icon: '/assets/images/category/category-dairy-bread-eggs.jpg', is_active: true },
  { id: 2, name: 'Snack & Munchies', slug: 'snack-munchies', icon: '/assets/images/category/category-snack-munchies.jpg', is_active: true },
  { id: 3, name: 'Bakery & Biscuits', slug: 'bakery-biscuits', icon: '/assets/images/category/category-bakery-biscuits.jpg', is_active: true },
  { id: 4, name: 'Instant Food', slug: 'instant-food', icon: '/assets/images/category/category-instant-food.jpg', is_active: true },
  { id: 5, name: 'Tea, Coffee & Drinks', slug: 'tea-coffee-drinks', icon: '/assets/images/category/category-tea-coffee-drinks.jpg', is_active: true },
  { id: 6, name: 'Atta, Rice & Dal', slug: 'atta-rice-dal', icon: '/assets/images/category/category-atta-rice-dal.jpg', is_active: true },
]

const displayCategories = computed(() => {
  return categories.value.length > 0 ? categories.value.slice(0, 6) : fallbackCategories
})

const filteredProducts = computed(() => {
  if (selectedCategory.value === 0) return products.value
  return products.value.filter(p => p.category_id === selectedCategory.value)
})

onMounted(async () => {
  try {
    const [catRes, prodRes] = await Promise.all([
      api.get<Category[]>('/categories', { active: 'true' }),
      api.get<Product[]>('/products', { limit: 10, status: 'ACTIVE' }),
    ])
    if (catRes.data && catRes.data.length > 0) categories.value = catRes.data
    if (prodRes.data && prodRes.data.length > 0) products.value = prodRes.data
  } catch (err) {
    console.error('Failed to load homepage data:', err)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="space-y-12 pb-16">
    <!-- Hero Slider Section (BookVerse Luxury Library Hero) -->
    <section class="mt-4">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
          class="relative rounded-3xl overflow-hidden bg-cover bg-center shadow-xl border border-slate-800"
          style="background-image: url('/images/slider/bookstore-hero.jpg'); min-height: 480px;"
        >
          <!-- Elegant Gradient Overlay -->
          <div class="absolute inset-0 bg-gradient-to-r from-slate-950/90 via-slate-950/60 to-transparent"></div>

          <div class="relative z-10 py-16 px-8 sm:px-14 lg:py-24 max-w-2xl flex flex-col justify-center">
            <span class="inline-flex items-center gap-2 px-3.5 py-1 text-xs font-bold rounded-full bg-amber-400 text-slate-950 w-fit mb-4 shadow-sm tracking-wide">
              <span>📚</span> Pesta Literasi & Book Fair 2026
            </span>

            <h1 class="text-3xl sm:text-5xl font-extrabold text-white leading-tight mb-4 drop-shadow-md">
              Jelajahi Dunia Lewat <br />
              <span class="text-amber-400 font-serif italic">Buku Original Terbaik</span>
            </h1>

            <p class="text-slate-200 text-sm sm:text-base leading-relaxed mb-8 font-normal drop-shadow-sm">
              Temukan ribuan karya novel fiksi best seller, buku pengembangan diri, bisnis, teknologi, dan komik berlisensi resmi. 100% buku baru bersegel dengan pengemasan aman bubble wrap tebal dan pengiriman cepat ke seluruh Indonesia.
            </p>

            <div class="flex flex-wrap gap-4 items-center">
              <NuxtLink
                to="/products"
                class="btn inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm px-7 py-3.5 rounded-xl shadow-lg shadow-indigo-600/30 transition-all transform hover:-translate-y-0.5"
              >
                <span>Beli Buku Sekarang</span>
                <i class="ti ti-arrow-right text-base"></i>
              </NuxtLink>

              <NuxtLink
                to="/products?category=fiksi-sastra"
                class="btn inline-flex items-center gap-2 bg-white/10 hover:bg-white/20 text-white border border-white/20 font-bold text-sm px-6 py-3.5 rounded-xl transition-all backdrop-blur-md"
              >
                <span>Novel Best Seller</span>
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
          <h2 class="text-xl sm:text-2xl font-extrabold text-slate-900 tracking-tight">Kategori Buku Pilihan</h2>
          <p class="text-xs text-slate-500 mt-0.5">Pilih genre bacaan favorit Anda dari fiksi, motivasi, bisnis, hingga sains modern</p>
        </div>
        <NuxtLink to="/products" class="text-xs font-bold text-indigo-600 hover:text-indigo-800 flex items-center gap-1">
          <span>Lihat Semua Kategori</span>
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
        <!-- Banner 1: Novel & Sastra -->
        <div class="rounded-2xl p-8 bg-gradient-to-br from-indigo-900 via-indigo-950 to-slate-950 text-white shadow-md relative overflow-hidden flex flex-col justify-between min-h-[220px] border border-indigo-800/50">
          <div class="relative z-10 max-w-sm space-y-2">
            <span class="text-[11px] font-bold text-amber-300 uppercase tracking-wider">Koleksi Fiksi & Sastra</span>
            <h3 class="text-2xl font-extrabold leading-snug">
              Novel Terlaris & Edisi Kolektor
            </h3>
            <p class="text-xs text-indigo-200">
              Dapatkan diskon spesial hingga <span class="font-bold text-white text-sm">30%</span> untuk novel pemenang penghargaan.
            </p>
          </div>
          <div class="relative z-10 mt-6">
            <NuxtLink
              to="/products"
              class="inline-flex items-center gap-1.5 bg-amber-400 hover:bg-amber-500 text-slate-950 font-bold text-xs px-5 py-2.5 rounded-xl transition-all shadow-sm"
            >
              <span>Jelajahi Novel</span>
              <i class="ti ti-arrow-right text-xs"></i>
            </NuxtLink>
          </div>
        </div>

        <!-- Banner 2: Self Improvement & Bisnis -->
        <div class="rounded-2xl p-8 bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-950 text-white shadow-md relative overflow-hidden flex flex-col justify-between min-h-[220px] border border-slate-700/50">
          <div class="relative z-10 max-w-sm space-y-2">
            <span class="text-[11px] font-bold text-emerald-400 uppercase tracking-wider">Pengembangan Diri & Finansial</span>
            <h3 class="text-2xl font-extrabold leading-snug">
              Investasi Terbaik untuk Pikiran Anda
            </h3>
            <p class="text-xs text-slate-300">
              Hemat hingga <span class="font-bold text-white text-sm">25%</span> untuk buku psikologi, habit, dan strategi bisnis.
            </p>
          </div>
          <div class="relative z-10 mt-6">
            <NuxtLink
              to="/products"
              class="inline-flex items-center gap-1.5 bg-white hover:bg-slate-100 text-slate-950 font-bold text-xs px-5 py-2.5 rounded-xl transition-all shadow-sm"
            >
              <span>Lihat Buku Bisnis</span>
              <i class="ti ti-arrow-right text-xs"></i>
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <!-- Popular Products Grid (Book Catalog) -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h2 class="text-xl sm:text-2xl font-extrabold text-slate-900 tracking-tight">Buku Paling Populer & Best Seller</h2>
          <p class="text-xs text-slate-500 mt-0.5">Buku-buku pilihan favorit pembaca dengan rating tertinggi dan ulasan memuaskan</p>
        </div>

        <!-- Category Filter Tabs -->
        <div class="flex items-center gap-2 overflow-x-auto pb-1 text-xs">
          <button
            @click="selectedCategory = 0"
            :class="selectedCategory === 0 ? 'bg-indigo-600 text-white font-bold shadow-xs' : 'bg-slate-100 text-slate-700 hover:bg-slate-200'"
            class="px-3.5 py-1.5 rounded-lg transition-colors whitespace-nowrap"
          >
            Semua Buku
          </button>
          <button
            v-for="cat in categories.slice(0, 4)"
            :key="cat.id"
            @click="selectedCategory = cat.id"
            :class="selectedCategory === cat.id ? 'bg-indigo-600 text-white font-bold shadow-xs' : 'bg-slate-100 text-slate-700 hover:bg-slate-200'"
            class="px-3.5 py-1.5 rounded-lg transition-colors whitespace-nowrap"
          >
            {{ cat.name }}
          </button>
        </div>
      </div>

      <!-- Loading skeleton -->
      <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
        <div v-for="n in 10" :key="n" class="h-80 rounded-2xl bg-slate-200 animate-pulse"></div>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredProducts.length === 0" class="text-center py-16 bg-white rounded-2xl border border-slate-200 p-8 shadow-xs">
        <i class="ti ti-book-off text-4xl text-slate-400 mb-2 inline-block"></i>
        <h3 class="text-base font-bold text-slate-800">Belum ada buku untuk kategori ini</h3>
        <p class="text-xs text-slate-500 mt-1">Silakan pilih kategori lain atau lihat seluruh katalog buku kami.</p>
      </div>

      <!-- Products Grid -->
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5">
        <ProductCard
          v-for="prod in filteredProducts"
          :key="prod.id"
          :product="prod"
        />
      </div>
    </section>

    <!-- Daily Best Sells Highlight Banner Section -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="relative rounded-3xl overflow-hidden border border-slate-200 bg-white shadow-xs p-6 sm:p-10 flex flex-col lg:flex-row items-center justify-between gap-8">
        <div class="space-y-4 max-w-lg">
          <span class="inline-block px-3 py-1 text-[11px] font-bold uppercase tracking-wider rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200/60">
            Jaminan Kualitas BookVerse
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-slate-900 leading-tight">
            100% Buku Original Langsung dari Penerbit Terpercaya
          </h2>
          <p class="text-xs sm:text-sm text-slate-600 leading-relaxed">
            Setiap buku dikirim dalam kondisi baru bersegel dari penerbit terkemuka (Gramedia, Penerbit Buku Kompas, Mizan, Bentang Pustaka, Elex Media, dll.). Kami menjamin kepuasan membaca Anda dengan jaminan ganti baru jika buku cacat produksi.
          </p>
          <div class="flex items-center gap-6 pt-3">
            <div>
              <div class="text-xl sm:text-2xl font-black text-indigo-600">10.000+</div>
              <div class="text-[11px] text-slate-500 font-medium">Judul Buku</div>
            </div>
            <div class="h-8 w-px bg-slate-200"></div>
            <div>
              <div class="text-xl sm:text-2xl font-black text-indigo-600">100%</div>
              <div class="text-[11px] text-slate-500 font-medium">Original & Segel</div>
            </div>
            <div class="h-8 w-px bg-slate-200"></div>
            <div>
              <div class="text-xl sm:text-2xl font-black text-indigo-600">24/7</div>
              <div class="text-[11px] text-slate-500 font-medium">Layanan Kirim</div>
            </div>
          </div>
        </div>

        <div class="w-full lg:w-auto flex justify-center">
          <div class="relative p-3 bg-gradient-to-tr from-amber-100 to-indigo-100 rounded-2xl shadow-inner">
            <img
              src="/images/products/book-atomic-habits.jpg"
              alt="Atomic Habits Bestseller"
              class="rounded-xl max-h-72 object-cover shadow-lg border border-white"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- 4 Key Value Pillars (BookVerse Signature) -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 pt-4">
        <!-- Feature 1 -->
        <div class="flex items-start gap-4 p-5 rounded-2xl bg-white border border-slate-200/80 shadow-xs hover:shadow-md transition-shadow">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-xl bg-indigo-50 text-indigo-600 text-xl font-bold">
            🛡️
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-slate-900">100% Buku Original</h3>
            <p class="text-xs text-slate-500 leading-relaxed">
              Bebas buku bajakan. Semua koleksi didapatkan resmi dari penerbit terdaftar.
            </p>
          </div>
        </div>

        <!-- Feature 2 -->
        <div class="flex items-start gap-4 p-5 rounded-2xl bg-white border border-slate-200/80 shadow-xs hover:shadow-md transition-shadow">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-xl bg-amber-50 text-amber-600 text-xl font-bold">
            📦
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-slate-900">Packing Anti-Tekuk</h3>
            <p class="text-xs text-slate-500 leading-relaxed">
              Setiap buku dibungkus bubble wrap berlapis dan kardus tebal agar tiba mulus.
            </p>
          </div>
        </div>

        <!-- Feature 3 -->
        <div class="flex items-start gap-4 p-5 rounded-2xl bg-white border border-slate-200/80 shadow-xs hover:shadow-md transition-shadow">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-xl bg-emerald-50 text-emerald-600 text-xl font-bold">
            🚚
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-slate-900">Kirim Cepat 24/7</h3>
            <p class="text-xs text-slate-500 leading-relaxed">
              Pesanan diproses setiap hari dengan pengiriman langsung ke alamat rumah Anda.
            </p>
          </div>
        </div>

        <!-- Feature 4 -->
        <div class="flex items-start gap-4 p-5 rounded-2xl bg-white border border-slate-200/80 shadow-xs hover:shadow-md transition-shadow">
          <div class="w-12 h-12 flex-shrink-0 flex items-center justify-center rounded-xl bg-purple-50 text-purple-600 text-xl font-bold">
            💳
          </div>
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-slate-900">Pembayaran Mudah</h3>
            <p class="text-xs text-slate-500 leading-relaxed">
              Bayar instan via QRIS, Virtual Account seluruh bank, e-Wallet, atau Kartu Kredit.
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
  { id: 1, name: 'Fiksi & Sastra', slug: 'fiksi-sastra', icon: '/assets/images/category/category-fruits-vegetables.jpg', is_active: true },
  { id: 2, name: 'Pengembangan Diri', slug: 'self-improvement', icon: '/assets/images/category/category-dairy-bread-eggs.jpg', is_active: true },
  { id: 3, name: 'Bisnis & Finansial', slug: 'bisnis-investasi', icon: '/assets/images/category/category-snack-munchies.jpg', is_active: true },
  { id: 4, name: 'Teknologi & Sains', slug: 'teknologi-sains', icon: '/assets/images/category/category-bakery-biscuits.jpg', is_active: true },
  { id: 5, name: 'Komik & Manga', slug: 'komik-manga', icon: '/assets/images/category/category-chicken-meat-fish.jpg', is_active: true },
  { id: 6, name: 'Anak & Edukasi', slug: 'anak-edukasi', icon: '/assets/images/category/category-tea-coffee-drinks.jpg', is_active: true },
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

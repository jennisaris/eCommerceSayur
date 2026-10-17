<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-10">
    <!-- Breadcrumbs -->
    <nav class="flex items-center gap-2 text-xs text-gray-500 mb-6">
      <NuxtLink to="/" class="hover:text-indigo-600 flex items-center gap-1">
        <i class="ti ti-home text-sm"></i>
        <span>Home</span>
      </NuxtLink>
      <span>/</span>
      <span class="text-gray-800 font-semibold">Koleksi Buku</span>
    </nav>

    <!-- Header Section -->
    <div class="mb-8">
      <h1 class="text-2xl sm:text-3xl font-extrabold text-gray-900">Katalog Pustaka & Buku Pilihan</h1>
      <p class="text-xs sm:text-sm text-gray-500 mt-1">
        Eksplorasi ribuan judul buku fiksi, non-fiksi, sains, teknologi, dan manga orisinal dengan jaminan kualitas penerbit resmi.
      </p>
    </div>

    <!-- Layout: Left Sidebar Filter + Right Product Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8 items-start">
      <!-- Left Sidebar Filters -->
      <aside class="space-y-6 bg-white p-5 rounded-xl border border-gray-200">
        <!-- Search Input -->
        <div>
          <h3 class="text-sm font-bold text-gray-900 mb-3 flex items-center gap-2">
            <i class="ti ti-search text-indigo-600"></i> Cari Judul / Penulis
          </h3>
          <div class="relative">
            <input
              v-model="filters.search"
              @keyup.enter="fetchProducts"
              type="text"
              placeholder="Ketik judul buku, penulis, ISBN..."
              class="w-full pl-9 pr-3 py-2 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600"
            />
            <i class="ti ti-search text-gray-400 absolute left-3 top-2.5 text-sm"></i>
          </div>
        </div>

        <!-- Categories Filter List -->
        <div>
          <h3 class="text-sm font-bold text-gray-900 mb-3 flex items-center gap-2">
            <i class="ti ti-books text-indigo-600"></i> Genre & Kategori
          </h3>
          <ul class="space-y-1.5 text-xs text-gray-600">
            <li>
              <button
                @click="selectCategory(0)"
                :class="filters.category_id === 0 ? 'text-indigo-600 font-bold bg-indigo-50' : 'hover:text-indigo-600 hover:bg-gray-50'"
                class="w-full text-left px-3 py-2 rounded-lg flex items-center justify-between transition-colors"
              >
                <span>Semua Genre</span>
                <i class="ti ti-chevron-right text-xs"></i>
              </button>
            </li>
            <li v-for="cat in categories" :key="cat.id">
              <button
                @click="selectCategory(cat.id)"
                :class="filters.category_id === cat.id ? 'text-indigo-600 font-bold bg-indigo-50' : 'hover:text-indigo-600 hover:bg-gray-50'"
                class="w-full text-left px-3 py-2 rounded-lg flex items-center justify-between transition-colors"
              >
                <span>{{ cat.name }}</span>
                <i class="ti ti-chevron-right text-xs"></i>
              </button>
            </li>
          </ul>
        </div>

        <!-- Price / Ordering Sort -->
        <div>
          <h3 class="text-sm font-bold text-gray-900 mb-3 flex items-center gap-2">
            <i class="ti ti-arrows-sort text-indigo-600"></i> Urutkan Buku
          </h3>
          <select
            v-model="filters.sort"
            @change="fetchProducts"
            class="w-full px-3 py-2 text-xs bg-gray-50 border border-gray-300 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-indigo-400"
          >
            <option value="latest">Terbitan Terbaru</option>
            <option value="price_asc">Harga Terendah</option>
            <option value="price_desc">Harga Tertinggi</option>
          </select>
        </div>

        <button
          @click="resetFilters"
          type="button"
          class="w-full py-2 px-3 text-xs font-semibold text-gray-700 hover:text-red-600 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors flex items-center justify-center gap-1.5"
        >
          <i class="ti ti-refresh text-xs"></i>
          <span>Reset Filter</span>
        </button>

        <!-- Book Protection Banner in Sidebar -->
        <div class="p-4 rounded-xl bg-indigo-50 border border-indigo-100 text-center space-y-2">
          <i class="ti ti-shield-check text-2xl text-indigo-600"></i>
          <h4 class="text-xs font-bold text-indigo-950">100% Original & Aman</h4>
          <p class="text-[11px] text-indigo-700 leading-relaxed">Semua buku dibungkus ekstra double bubble wrap & kardus sudut anti-penyok.</p>
        </div>
      </aside>

      <!-- Right Main: Products List -->
      <main class="lg:col-span-3 space-y-6">
        <!-- Top Toolbar -->
        <div class="bg-white p-4 rounded-xl border border-gray-200 flex flex-wrap items-center justify-between gap-4">
          <div class="text-xs text-gray-600">
            Menampilkan <span class="font-bold text-gray-900">{{ products.length }}</span> buku
            <span v-if="filters.search" class="text-indigo-600 font-semibold"> untuk "{{ filters.search }}"</span>
          </div>

          <div class="flex items-center gap-3">
            <button
              @click="fetchProducts"
              class="btn px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-semibold shadow-xs flex items-center gap-1"
            >
              <i class="ti ti-filter text-xs"></i>
              <span>Terapkan</span>
            </button>
          </div>
        </div>

        <!-- Products Grid -->
        <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <div v-for="n in 8" :key="n" class="h-80 rounded-lg bg-gray-200 animate-pulse"></div>
        </div>

        <div v-else-if="products.length === 0" class="text-center py-20 bg-white rounded-2xl border border-gray-200 p-8">
          <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-400">
            <i class="ti ti-book-off text-3xl"></i>
          </div>
          <h3 class="text-base font-bold text-gray-800">Tidak ada buku ditemukan</h3>
          <p class="text-xs text-gray-500 mt-1">Coba gunakan judul buku lain atau ubah filter kategori genre.</p>
          <button
            @click="resetFilters"
            class="mt-4 px-4 py-2 bg-indigo-600 text-white text-xs font-bold rounded-lg hover:bg-indigo-700 transition-all shadow-xs"
          >
            Reset Semua Filter
          </button>
        </div>

        <div v-else class="space-y-8">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-4 gap-4">
            <ProductCard
              v-for="product in products"
              :key="product.id"
              :product="product"
            />
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 pt-4">
            <button
              @click="changePage(currentPage - 1)"
              :disabled="currentPage <= 1"
              class="px-3 py-1.5 text-xs font-semibold bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:pointer-events-none transition-colors"
            >
              Sebelumnya
            </button>
            <span class="text-xs font-bold text-gray-700 px-3">
              Halaman {{ currentPage }} dari {{ totalPages }}
            </span>
            <button
              @click="changePage(currentPage + 1)"
              :disabled="currentPage >= totalPages"
              class="px-3 py-1.5 text-xs font-semibold bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:pointer-events-none transition-colors"
            >
              Selanjutnya
            </button>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Category, Product, ApiResponse } from '~/types'

const route = useRoute()
const router = useRouter()
const api = useApi()

const categories = ref<Category[]>([])
const products = ref<Product[]>([])
const isLoading = ref(true)

const currentPage = ref(1)
const totalPages = ref(1)
const limit = 12

const filters = reactive({
  search: (route.query.search as string) || '',
  category_id: Number(route.query.category_id) || 0,
  sort: (route.query.sort as string) || 'latest',
})

const selectCategory = (catId: number) => {
  filters.category_id = catId
  currentPage.value = 1
  fetchProducts()
}

const fetchCategories = async () => {
  try {
    const res = await api.get<Category[]>('/categories', { active: 'true' })
    if (res.data) categories.value = res.data
  } catch (err) {
    console.error('Failed to load categories:', err)
  }
}

const fetchProducts = async () => {
  isLoading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      limit,
      status: 'ACTIVE',
    }

    if (filters.search) params.search = filters.search
    if (filters.category_id > 0) params.category_id = filters.category_id
    if (filters.sort) params.sort = filters.sort

    // Sync query string
    router.replace({
      query: {
        ...(filters.search ? { search: filters.search } : {}),
        ...(filters.category_id ? { category_id: filters.category_id } : {}),
        ...(filters.sort !== 'latest' ? { sort: filters.sort } : {}),
        ...(currentPage.value > 1 ? { page: currentPage.value } : {}),
      },
    })

    const res = await api.get<Product[]>('/products', params)
    if (res.data) {
      products.value = res.data
      if (res.meta) {
        totalPages.value = res.meta.total_pages || 1
      }
    }
  } catch (err) {
    console.error('Failed to load products:', err)
  } finally {
    isLoading.value = false
  }
}

const resetFilters = () => {
  filters.search = ''
  filters.category_id = 0
  filters.sort = 'latest'
  currentPage.value = 1
  fetchProducts()
}

const changePage = (page: number) => {
  currentPage.value = page
  fetchProducts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  await Promise.all([fetchCategories(), fetchProducts()])
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-10">
    <!-- Breadcrumb -->
    <nav class="flex items-center gap-2 text-xs text-gray-500 mb-6">
      <NuxtLink to="/" class="hover:text-green-600 flex items-center gap-1">
        <i class="ti ti-home text-sm"></i>
        <span>Home</span>
      </NuxtLink>
      <span>/</span>
      <NuxtLink to="/products" class="hover:text-green-600">Shop</NuxtLink>
      <span>/</span>
      <span class="text-gray-800 font-semibold truncate max-w-xs">{{ product?.name || 'Detail Produk' }}</span>
    </nav>

    <div v-if="isLoading" class="grid grid-cols-1 md:grid-cols-2 gap-10 animate-pulse">
      <div class="aspect-square bg-gray-200 rounded-xl"></div>
      <div class="space-y-4">
        <div class="h-6 w-32 bg-gray-200 rounded"></div>
        <div class="h-10 w-full bg-gray-200 rounded"></div>
        <div class="h-8 w-48 bg-gray-200 rounded"></div>
        <div class="h-32 w-full bg-gray-200 rounded"></div>
      </div>
    </div>

    <div v-else-if="!product" class="text-center py-20 bg-white rounded-2xl border border-gray-200 p-8">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-100 flex items-center justify-center text-gray-400">
        <i class="ti ti-package-off text-3xl"></i>
      </div>
      <h2 class="text-xl font-bold text-gray-800">Produk tidak ditemukan</h2>
      <p class="text-xs text-gray-500 mt-2">Produk yang Anda cari mungkin telah dinonaktifkan atau dihapus.</p>
      <NuxtLink to="/products" class="mt-4 inline-block px-5 py-2.5 bg-green-600 text-white rounded-lg text-xs font-bold hover:bg-green-700 shadow-xs transition-all">
        Kembali ke Katalog
      </NuxtLink>
    </div>

    <!-- FreshCart Shop Single View -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-10 items-start">
      <!-- Product Image Display -->
      <div class="sticky top-24 bg-white p-6 rounded-2xl border border-gray-200 shadow-xs flex items-center justify-center">
        <div class="relative w-full aspect-square flex items-center justify-center overflow-hidden">
          <img
            :src="product.image || '/assets/images/products/product-single-img-1.jpg'"
            :alt="product.name"
            class="max-w-full max-h-[420px] object-contain hover:scale-105 transition-transform duration-300"
          />
          <div v-if="hasDiscount" class="absolute top-2 left-2 z-10">
            <span class="inline-block px-2.5 py-1 text-xs font-bold rounded bg-red-600 text-white">
              -{{ discountPercentage }}% Off
            </span>
          </div>
        </div>
      </div>

      <!-- Product Purchase Options & Details -->
      <div class="space-y-6">
        <div>
          <span class="text-xs font-bold text-green-600 uppercase tracking-wider mb-2 inline-block">
            {{ product.category?.name || 'Sayur & Bahan Makanan' }}
          </span>
          <h1 class="text-2xl sm:text-3xl font-extrabold text-gray-900 leading-snug">
            {{ product.name }}
          </h1>

          <!-- Ratings from FreshCart template -->
          <div class="flex items-center gap-2 mt-2 text-xs">
            <div class="flex items-center text-yellow-500">
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-half-filled"></i>
            </div>
            <span class="text-green-600 font-bold">(32 reviews)</span>
            <span class="text-gray-300">&bull;</span>
            <span :class="product.stock > 0 ? 'text-green-600 font-bold' : 'text-red-500 font-bold'">
              {{ product.stock > 0 ? 'In Stock' : 'Out of Stock' }}
            </span>
          </div>
        </div>

        <!-- Pricing Section -->
        <div class="p-5 rounded-xl bg-gray-50 border border-gray-200 flex items-center justify-between">
          <div>
            <div class="text-2xl sm:text-3xl font-black text-gray-900">
              {{ formatIDR(effectivePrice) }}
            </div>
            <div v-if="hasDiscount" class="flex items-center gap-2 mt-1">
              <span class="text-xs text-gray-400 line-through">{{ formatIDR(product.regular_price) }}</span>
              <span class="text-[11px] font-bold px-2 py-0.5 rounded bg-red-100 text-red-700">
                Hemat {{ discountPercentage }}%
              </span>
            </div>
          </div>

          <div class="text-right">
            <span
              :class="product.stock > 0 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
              class="inline-block px-3 py-1 text-xs font-bold rounded-full"
            >
              {{ product.stock > 0 ? `Tersedia: ${product.stock} ${product.unit}` : 'Stok Habis' }}
            </span>
            <div class="text-[11px] text-gray-500 mt-1">Berat: {{ product.weight }} gram</div>
          </div>
        </div>

        <!-- Description -->
        <div>
          <h3 class="text-xs font-bold uppercase tracking-wider text-gray-500 mb-2">Deskripsi Produk</h3>
          <p class="text-xs sm:text-sm text-gray-600 leading-relaxed whitespace-pre-line">
            {{ product.description || 'Produk sayuran segar pilihan terbaik, dipetik langsung dari kebun mitra terpercaya dengan standar kebersihan tinggi.' }}
          </p>
        </div>

        <!-- Quantity & Add to Cart Controls -->
        <div class="space-y-4 pt-2 border-t border-gray-200">
          <div class="flex items-center gap-4">
            <span class="text-xs font-bold text-gray-700">Jumlah Beli:</span>
            <div class="flex items-center border border-gray-300 rounded-lg overflow-hidden bg-white">
              <button
                @click="qty = Math.max(1, qty - 1)"
                :disabled="qty <= 1"
                class="px-3 py-1.5 text-xs text-gray-600 hover:bg-gray-100 disabled:opacity-40 transition-colors border-r border-gray-200"
              >
                -
              </button>
              <input
                v-model.number="qty"
                type="number"
                min="1"
                :max="product.stock"
                class="w-12 text-center text-xs font-bold text-gray-900 focus:outline-none"
              />
              <button
                @click="qty = Math.min(product.stock, qty + 1)"
                :disabled="qty >= product.stock"
                class="px-3 py-1.5 text-xs text-gray-600 hover:bg-gray-100 disabled:opacity-40 transition-colors border-l border-gray-200"
              >
                +
              </button>
            </div>
            <span class="text-xs text-gray-500">Maks. {{ product.stock }} unit</span>
          </div>

          <div class="flex flex-col sm:flex-row gap-3 pt-2">
            <button
              @click="handleAddToCart"
              :disabled="product.stock <= 0 || isAdding"
              type="button"
              class="btn flex-1 py-3 px-6 bg-green-600 hover:bg-green-700 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-green-600/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50"
            >
              <i v-if="isAdding" class="ti ti-loader-2 animate-spin text-base"></i>
              <i v-else class="ti ti-plus text-base"></i>
              <span>{{ isAdding ? 'Menambahkan...' : 'Tambah ke Keranjang' }}</span>
            </button>

            <button
              @click="handleBuyNow"
              :disabled="product.stock <= 0"
              type="button"
              class="btn py-3 px-6 bg-gray-900 hover:bg-gray-800 text-white font-bold text-xs sm:text-sm rounded-xl transition-all disabled:opacity-50 flex items-center justify-center gap-1.5"
            >
              <span>Beli Sekarang</span>
              <i class="ti ti-arrow-right text-sm"></i>
            </button>
          </div>
        </div>

        <!-- FreshCart Product Specification Table -->
        <div class="pt-4 border-t border-gray-200">
          <table class="w-full text-xs text-left">
            <tbody class="divide-y divide-gray-100">
              <tr>
                <td class="py-2.5 text-gray-500 w-1/3">Kode SKU</td>
                <td class="py-2.5 font-mono font-semibold text-gray-800">{{ product.sku || 'SAYUR-PRD-' + product.id }}</td>
              </tr>
              <tr>
                <td class="py-2.5 text-gray-500">Kategori</td>
                <td class="py-2.5 font-semibold text-gray-800">{{ product.category?.name || 'Sayur' }}</td>
              </tr>
              <tr>
                <td class="py-2.5 text-gray-500">Model Pengiriman</td>
                <td class="py-2.5 font-semibold text-green-700">100% Home Delivery (24 Jam Non-Stop)</td>
              </tr>
              <tr>
                <td class="py-2.5 text-gray-500">Jaminan Kualitas</td>
                <td class="py-2.5 text-gray-800">100% Sayur Segar & Higienis</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Product } from '~/types'

const route = useRoute()
const router = useRouter()
const api = useApi()
const cartStore = useCartStore()

const product = ref<Product | null>(null)
const isLoading = ref(true)
const qty = ref(1)
const isAdding = ref(false)

const effectivePrice = computed(() => {
  if (!product.value) return 0
  if (product.value.sale_price > 0 && product.value.sale_price < product.value.regular_price) {
    return product.value.sale_price
  }
  return product.value.regular_price
})

const hasDiscount = computed(() => {
  if (!product.value) return false
  return product.value.sale_price > 0 && product.value.sale_price < product.value.regular_price
})

const discountPercentage = computed(() => {
  if (!hasDiscount.value || !product.value) return 0
  const diff = product.value.regular_price - product.value.sale_price
  return Math.round((diff / product.value.regular_price) * 100)
})

const handleAddToCart = async () => {
  if (!product.value) return
  isAdding.value = true
  try {
    await cartStore.addToCart(product.value.id, qty.value)
    cartStore.openDrawer()
  } finally {
    isAdding.value = false
  }
}

const handleBuyNow = async () => {
  if (!product.value) return
  await cartStore.addToCart(product.value.id, qty.value)
  router.push('/checkout')
}

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
}

onMounted(async () => {
  const slug = route.params.slug as string
  try {
    const res = await api.get<Product>(`/products/slug/${slug}`)
    if (res.data) {
      product.value = res.data
    }
  } catch (err) {
    console.error('Failed to load product details:', err)
  } finally {
    isLoading.value = false
  }
})
</script>

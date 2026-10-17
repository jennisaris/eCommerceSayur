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
      <!-- Book Cover Image Display -->
      <div class="sticky top-24 bg-white p-6 rounded-3xl border border-slate-200 shadow-sm flex items-center justify-center">
        <div class="relative w-full aspect-[3/4] max-w-sm flex items-center justify-center overflow-hidden rounded-2xl bg-slate-50 p-2 shadow-inner">
          <img
            :src="product.image || '/images/products/book-atomic-habits.jpg'"
            :alt="product.name"
            class="max-w-full max-h-[460px] object-cover rounded-xl shadow-md hover:scale-105 transition-transform duration-300"
          />
          <div v-if="hasDiscount" class="absolute top-4 left-4 z-10">
            <span class="inline-block px-3 py-1 text-xs font-extrabold rounded-lg bg-rose-600 text-white shadow-sm">
              Hemat {{ discountPercentage }}%
            </span>
          </div>
        </div>
      </div>

      <!-- Book Purchase Options & Details -->
      <div class="space-y-6">
        <div>
          <span class="text-xs font-bold text-indigo-600 uppercase tracking-wider mb-2 inline-block">
            {{ product.category?.name || 'Katalog Buku' }}
          </span>
          <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 leading-snug">
            {{ product.name }}
          </h1>

          <div class="flex items-center gap-2 mt-2 text-xs">
            <div class="flex items-center text-amber-500">
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-filled"></i>
              <i class="ti ti-star-half-filled"></i>
            </div>
            <span class="text-indigo-600 font-bold">(48 ulasan pembaca)</span>
            <span class="text-slate-300">&bull;</span>
            <span :class="product.stock > 0 ? 'text-emerald-600 font-bold' : 'text-rose-500 font-bold'">
              {{ product.stock > 0 ? 'Tersedia (Ready Stock)' : 'Stok Habis' }}
            </span>
          </div>
        </div>

        <!-- Pricing Section -->
        <div class="p-5 rounded-2xl bg-slate-50 border border-slate-200 flex items-center justify-between">
          <div>
            <div class="text-2xl sm:text-3xl font-black text-slate-900">
              {{ formatIDR(effectivePrice) }}
            </div>
            <div v-if="hasDiscount" class="flex items-center gap-2 mt-1">
              <span class="text-xs text-slate-400 line-through">{{ formatIDR(product.regular_price) }}</span>
              <span class="text-[11px] font-bold px-2 py-0.5 rounded-md bg-rose-100 text-rose-700">
                Hemat {{ discountPercentage }}%
              </span>
            </div>
          </div>

          <div class="text-right">
            <span
              :class="product.stock > 0 ? 'bg-emerald-100 text-emerald-800' : 'bg-rose-100 text-rose-800'"
              class="inline-block px-3 py-1 text-xs font-bold rounded-full"
            >
              {{ product.stock > 0 ? `Tersedia: ${product.stock} ${product.unit}` : 'Stok Habis' }}
            </span>
            <div class="text-[11px] text-slate-500 mt-1 font-medium">Bobot Buku: {{ product.weight }} gram</div>
          </div>
        </div>

        <!-- Book Synopsis & Description -->
        <div>
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-500 mb-2">Sinopsis & Keterangan Buku</h3>
          <p class="text-xs sm:text-sm text-slate-600 leading-relaxed whitespace-pre-line bg-white p-4 rounded-xl border border-slate-200/80">
            {{ product.description || 'Buku pilihan original berkualitas tinggi langsung dari penerbit resmi dengan jaminan isi lengkap dan segel pabrik.' }}
          </p>
        </div>

        <!-- Quantity & Add to Cart Controls -->
        <div class="space-y-4 pt-2 border-t border-slate-200">
          <div class="flex items-center gap-4">
            <span class="text-xs font-bold text-slate-700">Jumlah Eksemplar:</span>
            <div class="flex items-center border border-slate-300 rounded-xl overflow-hidden bg-white shadow-xs">
              <button
                @click="qty = Math.max(1, qty - 1)"
                :disabled="qty <= 1"
                class="px-3.5 py-1.5 text-xs text-slate-600 hover:bg-slate-100 disabled:opacity-40 transition-colors border-r border-slate-200"
              >
                -
              </button>
              <input
                v-model.number="qty"
                type="number"
                min="1"
                :max="product.stock"
                class="w-12 text-center text-xs font-bold text-slate-900 focus:outline-none"
              />
              <button
                @click="qty = Math.min(product.stock, qty + 1)"
                :disabled="qty >= product.stock"
                class="px-3.5 py-1.5 text-xs text-slate-600 hover:bg-slate-100 disabled:opacity-40 transition-colors border-l border-slate-200"
              >
                +
              </button>
            </div>
            <span class="text-xs text-slate-500">Maks. {{ product.stock }} eksemplar</span>
          </div>

          <div class="flex flex-col sm:flex-row gap-3 pt-2">
            <button
              @click="handleAddToCart"
              :disabled="product.stock <= 0 || isAdding"
              type="button"
              class="btn flex-1 py-3 px-6 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-indigo-600/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50"
            >
              <i v-if="isAdding" class="ti ti-loader-2 animate-spin text-base"></i>
              <i v-else class="ti ti-shopping-cart-plus text-base"></i>
              <span>{{ isAdding ? 'Menambahkan...' : 'Tambah ke Keranjang' }}</span>
            </button>

            <button
              @click="handleBuyNow"
              :disabled="product.stock <= 0"
              type="button"
              class="btn py-3 px-6 bg-slate-900 hover:bg-slate-800 text-white font-bold text-xs sm:text-sm rounded-xl transition-all disabled:opacity-50 flex items-center justify-center gap-1.5"
            >
              <span>Beli Sekarang</span>
              <i class="ti ti-arrow-right text-sm"></i>
            </button>
          </div>
        </div>

        <!-- Bookverse Specification Table -->
        <div class="pt-4 border-t border-slate-200">
          <table class="w-full text-xs text-left">
            <tbody class="divide-y divide-slate-100">
              <tr>
                <td class="py-2.5 text-slate-500 w-1/3">Format Edisi</td>
                <td class="py-2.5 font-semibold text-slate-800">{{ product.variant_name || 'Buku Cetak Fisik' }}</td>
              </tr>
              <tr>
                <td class="py-2.5 text-slate-500">Kategori Genre</td>
                <td class="py-2.5 font-semibold text-slate-800">{{ product.category?.name || 'Buku Pilihan' }}</td>
              </tr>
              <tr>
                <td class="py-2.5 text-slate-500">Standar Pengemasan</td>
                <td class="py-2.5 font-semibold text-indigo-700">Bubble Wrap Tebal Berlapis + Kardus Aman Anti-Tekuk</td>
              </tr>
              <tr>
                <td class="py-2.5 text-slate-500">Keaslian Buku</td>
                <td class="py-2.5 text-emerald-700 font-bold">100% Original Langsung dari Penerbit Resmi</td>
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

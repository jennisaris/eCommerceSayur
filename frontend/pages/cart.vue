<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-10">
    <!-- Breadcrumbs -->
    <nav class="flex items-center gap-2 text-xs text-gray-500 mb-6">
      <NuxtLink to="/" class="hover:text-green-600 flex items-center gap-1">
        <i class="ti ti-home text-sm"></i>
        <span>Home</span>
      </NuxtLink>
      <span>/</span>
      <NuxtLink to="/products" class="hover:text-green-600">Shop</NuxtLink>
      <span>/</span>
      <span class="text-gray-800 font-semibold">Shop Cart</span>
    </nav>

    <div class="mb-8">
      <h1 class="text-2xl sm:text-3xl font-extrabold text-gray-900">Keranjang Belanja</h1>
      <p class="text-xs sm:text-sm text-gray-500 mt-1">Periksa daftar pesanan sayuran dan bahan dapur Anda sebelum checkout.</p>
    </div>

    <!-- Empty state -->
    <div v-if="cartStore.items.length === 0" class="text-center py-20 bg-white rounded-2xl border border-gray-200 p-8 max-w-xl mx-auto shadow-xs">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-100 flex items-center justify-center text-gray-400">
        <i class="ti ti-shopping-cart-off text-3xl"></i>
      </div>
      <h3 class="text-base font-bold text-gray-800">Keranjang Belanja Anda Kosong</h3>
      <p class="text-xs text-gray-500 mt-1">Anda belum menambahkan produk apa pun ke keranjang belanja.</p>
      <NuxtLink
        to="/products"
        class="mt-6 inline-flex items-center gap-2 px-6 py-3 bg-green-600 text-white rounded-xl text-xs font-bold hover:bg-green-700 transition-all shadow-sm"
      >
        <span>Mulai Belanja Sekarang</span>
        <i class="ti ti-arrow-right text-xs"></i>
      </NuxtLink>
    </div>

    <!-- Cart with Items -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
      <!-- Items Table/List -->
      <div class="lg:col-span-2 space-y-4">
        <!-- Delivery info alert -->
        <div class="p-3.5 bg-green-50 border border-green-200 rounded-xl text-xs text-green-900 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <i class="ti ti-truck text-base text-green-600"></i>
            <span>Layanan <strong>100% Home Delivery</strong> siap antar 24/7 non-stop</span>
          </div>
          <span class="font-bold text-green-700">24 Jam</span>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-xs">
          <div class="p-4 border-b border-gray-100 flex items-center justify-between bg-gray-50/50">
            <span class="text-xs font-bold text-gray-800">Daftar Produk ({{ cartStore.totalItems }} Item)</span>
            <button
              @click="confirmClearCart"
              type="button"
              class="text-xs text-red-600 hover:text-red-800 font-semibold flex items-center gap-1 transition-colors"
            >
              <i class="ti ti-trash text-xs"></i>
              <span>Kosongkan Keranjang</span>
            </button>
          </div>

          <div class="divide-y divide-gray-100 p-4 space-y-4">
            <div
              v-for="item in cartStore.items"
              :key="item.id"
              class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 pt-4 first:pt-0"
            >
              <div class="flex items-center gap-4">
                <img
                  :src="item.product?.image || '/assets/images/products/product-img-1.jpg'"
                  :alt="item.product?.name"
                  class="w-16 h-16 sm:w-20 sm:h-20 rounded-lg object-contain bg-gray-50 border border-gray-200 flex-shrink-0"
                />
                <div>
                  <h3 class="text-xs sm:text-sm font-bold text-gray-900 line-clamp-1">
                    {{ item.product?.name }}
                  </h3>
                  <div class="text-xs text-gray-500 mt-1">
                    {{ formatIDR(item.product?.sale_price && item.product.sale_price > 0 ? item.product.sale_price : item.product?.regular_price || 0) }} / unit
                  </div>
                  <div class="text-[11px] text-gray-400 mt-0.5">Berat: {{ item.product?.weight || 0 }} gram</div>
                </div>
              </div>

              <!-- Controls -->
              <div class="flex items-center justify-between w-full sm:w-auto gap-6">
                <div class="flex items-center border border-gray-300 rounded-lg overflow-hidden bg-white">
                  <button
                    @click="cartStore.updateQuantity(item.product_id, item.quantity - 1)"
                    type="button"
                    class="px-3 py-1 text-xs text-gray-600 hover:bg-gray-100 border-r border-gray-200 transition-colors"
                  >
                    -
                  </button>
                  <span class="px-3.5 py-1 text-xs font-bold text-gray-900">{{ item.quantity }}</span>
                  <button
                    @click="cartStore.updateQuantity(item.product_id, item.quantity + 1)"
                    type="button"
                    class="px-3 py-1 text-xs text-gray-600 hover:bg-gray-100 border-l border-gray-200 transition-colors"
                  >
                    +
                  </button>
                </div>

                <div class="text-right">
                  <div class="text-xs sm:text-sm font-bold text-gray-900">
                    {{ formatIDR((item.product?.sale_price && item.product.sale_price > 0 ? item.product.sale_price : item.product?.regular_price || 0) * item.quantity) }}
                  </div>
                  <button
                    @click="cartStore.removeItem(item.product_id)"
                    class="text-[11px] text-red-500 hover:text-red-700 font-medium inline-flex items-center gap-0.5 mt-1"
                  >
                    <i class="ti ti-trash text-xs"></i>
                    <span>Hapus</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-between items-center pt-2">
          <NuxtLink
            to="/products"
            class="btn px-4 py-2 bg-white text-gray-800 border border-gray-300 rounded-lg text-xs font-semibold hover:bg-gray-50 flex items-center gap-1.5 transition-colors"
          >
            <i class="ti ti-arrow-left text-xs"></i>
            <span>Lanjut Belanja</span>
          </NuxtLink>
        </div>
      </div>

      <!-- Summary Sidebar -->
      <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-xs space-y-5 sticky top-24">
        <h3 class="text-base font-bold text-gray-900 border-b border-gray-100 pb-3">Ringkasan Belanja</h3>

        <div class="space-y-3 text-xs">
          <div class="flex justify-between text-gray-600">
            <span>Subtotal Produk</span>
            <span class="font-bold text-gray-900">{{ formatIDR(cartStore.totalPrice) }}</span>
          </div>

          <div class="flex justify-between text-gray-600">
            <span>Total Estimasi Berat</span>
            <span class="font-bold text-gray-900">{{ (cartStore.totalWeight / 1000).toFixed(2) }} kg</span>
          </div>

          <div class="flex justify-between text-gray-600">
            <span>Pengiriman</span>
            <span class="text-green-600 font-semibold">Dihitung saat checkout</span>
          </div>

          <div class="pt-3 border-t border-gray-200 flex justify-between items-center text-sm">
            <span class="font-bold text-gray-900">Total Sementara</span>
            <span class="text-base font-black text-green-600">{{ formatIDR(cartStore.totalPrice) }}</span>
          </div>
        </div>

        <NuxtLink
          to="/checkout"
          class="btn w-full py-3.5 px-4 bg-green-600 hover:bg-green-700 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-green-600/25 flex items-center justify-center gap-2 transition-all"
        >
          <span>Lanjut ke Checkout</span>
          <i class="ti ti-arrow-right text-sm"></i>
        </NuxtLink>

        <div class="p-3 bg-gray-50 rounded-lg border border-gray-200 text-center">
          <div class="text-[11px] text-gray-500 flex items-center justify-center gap-1">
            <i class="ti ti-shield-check text-green-600 text-sm"></i>
            <span>Transaksi Dijamin Aman & Idempoten</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const cartStore = useCartStore()

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
}

const confirmClearCart = async () => {
  if (confirm('Apakah Anda yakin ingin mengosongkan semua produk dari keranjang belanja?')) {
    await cartStore.clearCart()
  }
}

onMounted(() => {
  cartStore.fetchCart()
})
</script>

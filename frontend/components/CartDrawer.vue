<template>
  <div v-if="cartStore.isDrawerOpen" class="fixed inset-0 z-50 overflow-hidden">
    <!-- Backdrop -->
    <div
      @click="cartStore.closeDrawer"
      class="absolute inset-0 bg-gray-900/40 backdrop-blur-xs transition-opacity"
    ></div>

    <div class="fixed inset-y-0 right-0 max-w-full flex pl-10">
      <div class="w-screen max-w-md bg-white shadow-2xl flex flex-col">
        <!-- Drawer Header -->
        <div class="p-5 border-b border-gray-200 flex items-center justify-between bg-white">
          <div class="flex items-center gap-2">
            <i class="ti ti-shopping-bag text-xl text-indigo-600"></i>
            <h2 class="text-base font-bold text-gray-900">Keranjang Buku</h2>
            <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200">
              {{ cartStore.totalItems }}
            </span>
          </div>
          <button
            @click="cartStore.closeDrawer"
            class="p-2 text-gray-500 hover:text-gray-800 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <i class="ti ti-x text-lg"></i>
          </button>
        </div>

        <!-- Packaging Protection Alert Banner -->
        <div class="bg-indigo-50 text-indigo-900 text-xs px-5 py-2.5 border-b border-indigo-100 flex items-center justify-between">
          <span class="flex items-center gap-1.5 font-medium">
            <i class="ti ti-shield-check text-indigo-600 text-sm"></i>
            Kemasan aman dengan bubble wrap & kardus tebal!
          </span>
          <NuxtLink to="/checkout" @click="cartStore.closeDrawer" class="text-indigo-700 font-bold hover:underline">
            Checkout
          </NuxtLink>
        </div>

        <!-- Cart Items List -->
        <div class="flex-1 overflow-y-auto p-5 space-y-4">
          <div v-if="cartStore.items.length === 0" class="text-center py-16">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-400">
              <i class="ti ti-book-off text-3xl"></i>
            </div>
            <p class="text-sm font-semibold text-gray-800">Keranjang Buku Anda masih kosong</p>
            <p class="text-xs text-gray-500 mt-1">Eksplorasi ribuan judul buku fiksi, non-fiksi, dan teknologi di katalog.</p>
            <NuxtLink
              to="/products"
              @click="cartStore.closeDrawer"
              class="mt-4 inline-block px-4 py-2 text-xs font-bold bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-all shadow-xs"
            >
              Cari Buku Sekarang
            </NuxtLink>
          </div>

          <div
            v-for="item in cartStore.items"
            :key="item.id"
            class="flex gap-3.5 p-3 rounded-lg border border-gray-200 hover:border-indigo-600 bg-white transition-all"
          >
            <!-- Product Image (Book Cover Aspect Ratio) -->
            <img
              :src="item.product?.image || '/assets/images/products/product-img-1.jpg'"
              alt="Book Cover"
              class="w-14 h-20 rounded-md object-cover bg-gray-50 flex-shrink-0 shadow-xs border border-gray-100"
            />

            <!-- Details -->
            <div class="flex-1 flex flex-col justify-between min-w-0">
              <div>
                <h4 class="text-xs sm:text-sm font-semibold text-gray-900 truncate">
                  {{ item.product?.name }}
                </h4>
                <div class="text-[11px] text-gray-500 mt-0.5 flex items-center gap-2">
                  <span class="font-bold text-gray-900">
                    {{ formatIDR(item.product?.sale_price && item.product.sale_price > 0 ? item.product.sale_price : item.product?.regular_price || 0) }}
                  </span>
                  <span class="text-gray-300">&bull;</span>
                  <span>{{ item.product?.weight || 0 }}g</span>
                </div>
              </div>

              <!-- Quantity Controls & Delete -->
              <div class="flex items-center justify-between mt-2 pt-2 border-t border-gray-100">
                <div class="flex items-center border border-gray-300 rounded-lg overflow-hidden bg-white">
                  <button
                    @click="cartStore.updateQuantity(item.product_id, item.quantity - 1)"
                    type="button"
                    class="px-2.5 py-0.5 text-xs text-gray-600 hover:bg-gray-100 transition-colors border-r border-gray-200"
                  >
                    -
                  </button>
                  <span class="px-2.5 py-0.5 text-xs font-bold text-gray-800">
                    {{ item.quantity }}
                  </span>
                  <button
                    @click="cartStore.updateQuantity(item.product_id, item.quantity + 1)"
                    type="button"
                    class="px-2.5 py-0.5 text-xs text-gray-600 hover:bg-gray-100 transition-colors border-l border-gray-200"
                  >
                    +
                  </button>
                </div>

                <button
                  @click="cartStore.removeItem(item.product_id)"
                  type="button"
                  class="text-xs font-medium text-gray-500 hover:text-red-600 transition-colors flex items-center gap-1"
                >
                  <i class="ti ti-trash text-sm"></i>
                  <span>Hapus</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Drawer Footer -->
        <div v-if="cartStore.items.length > 0" class="p-5 border-t border-gray-200 bg-gray-50 space-y-3">
          <div class="space-y-1.5 text-xs">
            <div class="flex justify-between text-gray-500">
              <span>Total Estimasi Berat</span>
              <span class="font-semibold text-gray-700">{{ (cartStore.totalWeight / 1000).toFixed(2) }} kg</span>
            </div>
            <div class="flex justify-between text-gray-500">
              <span>Subtotal Belanja</span>
              <span class="text-sm font-bold text-gray-900">{{ formatIDR(cartStore.totalPrice) }}</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2 pt-1">
            <NuxtLink
              to="/cart"
              @click="cartStore.closeDrawer"
              class="btn py-2.5 px-3 bg-white hover:bg-gray-100 text-gray-800 border border-gray-300 text-xs font-bold rounded-lg text-center transition-all"
            >
              Lihat Keranjang
            </NuxtLink>
            <NuxtLink
              to="/checkout"
              @click="cartStore.closeDrawer"
              class="btn py-2.5 px-3 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-lg shadow-sm text-center flex items-center justify-center gap-1 transition-all"
            >
              <span>Checkout</span>
              <i class="ti ti-arrow-right text-xs"></i>
            </NuxtLink>
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
</script>

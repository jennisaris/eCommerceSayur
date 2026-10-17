<template>
  <div class="relative rounded-2xl break-words border bg-white border-slate-200 transition-all duration-300 hover:border-indigo-500 hover:shadow-xl hover:-translate-y-1 card-product flex flex-col justify-between h-full group overflow-hidden">
    <div class="p-4 flex-auto flex flex-col justify-between">
      <!-- Book Cover Image & Badges -->
      <div class="text-center relative flex justify-center mb-3">
        <!-- Discount / Status Badge -->
        <div v-if="hasDiscount" class="absolute top-1 left-1 z-10">
          <span class="inline-block px-2 py-0.5 text-center font-extrabold text-[10px] rounded-md bg-rose-600 text-white shadow-xs">
            -{{ discountPercentage }}%
          </span>
        </div>
        <div v-else-if="product.stock <= 0" class="absolute top-1 left-1 z-10">
          <span class="inline-block px-2 py-0.5 text-center font-bold text-[10px] rounded-md bg-slate-600 text-white shadow-xs">
            Stok Habis
          </span>
        </div>
        <div v-else class="absolute top-1 left-1 z-10">
          <span class="inline-block px-2 py-0.5 text-center font-bold text-[10px] rounded-md bg-indigo-600 text-white shadow-xs">
            Original
          </span>
        </div>

        <NuxtLink :to="`/products/${product.slug}`" class="block w-full overflow-hidden aspect-[3/4] flex items-center justify-center bg-slate-50 rounded-xl p-1 shadow-inner">
          <img
            :src="product.image || '/images/products/book-atomic-habits.jpg'"
            :alt="product.name"
            class="w-full h-full object-cover rounded-lg group-hover:scale-105 transition-transform duration-300 shadow-sm"
            loading="lazy"
          />
        </NuxtLink>

        <!-- Hover Quick Actions -->
        <div class="absolute w-full bottom-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center gap-2">
          <NuxtLink
            :to="`/products/${product.slug}`"
            class="h-8 w-8 bg-white/95 shadow-md inline-flex items-center justify-center rounded-xl text-slate-700 hover:bg-indigo-600 hover:text-white transition-all backdrop-blur-xs"
            title="Sinopsis & Detail"
          >
            <i class="ti ti-eye text-sm"></i>
          </NuxtLink>
        </div>
      </div>

      <!-- Book Meta & Info -->
      <div class="flex flex-col gap-1.5 mt-auto">
        <div class="flex items-center justify-between text-[11px] text-slate-500">
          <span class="font-medium truncate max-w-[120px]">{{ product.category?.name || 'Buku Pilihan' }}</span>
          <span v-if="product.variant_name" class="px-1.5 py-0.2 bg-slate-100 text-slate-600 rounded text-[10px] font-medium truncate max-w-[90px]">
            {{ product.variant_name }}
          </span>
        </div>

        <h3 class="text-xs sm:text-sm font-bold text-slate-800 line-clamp-2 hover:text-indigo-600 transition-colors leading-snug">
          <NuxtLink :to="`/products/${product.slug}`">
            {{ product.name }}
          </NuxtLink>
        </h3>

        <!-- Rating Stars -->
        <div class="flex items-center gap-1 text-amber-500 text-xs mt-0.5">
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-half-filled"></i>
          <span class="text-slate-500 text-[10px] ml-1 font-semibold">4.9 (48 ulasan)</span>
        </div>

        <!-- Price & Add Button -->
        <div class="flex justify-between items-center mt-2 pt-2.5 border-t border-slate-100">
          <div>
            <span class="text-slate-900 font-extrabold text-xs sm:text-sm block">
              {{ formatIDR(effectivePrice) }}
            </span>
            <span v-if="hasDiscount" class="line-through text-slate-400 text-[11px]">
              {{ formatIDR(product.regular_price) }}
            </span>
          </div>

          <div>
            <button
              @click="handleAddToCart"
              :disabled="product.stock <= 0 || isAdding"
              type="button"
              class="btn inline-flex items-center gap-1 bg-indigo-600 text-white disabled:opacity-50 disabled:pointer-events-none hover:bg-indigo-700 active:bg-indigo-800 text-xs py-1.5 px-3 rounded-xl font-bold transition-all shadow-xs"
            >
              <i v-if="isAdding" class="ti ti-loader-2 animate-spin text-xs"></i>
              <i v-else class="ti ti-shopping-cart-plus text-xs"></i>
              <span>{{ isAdding ? '...' : '+ Beli' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Product } from '~/types'

const props = defineProps<{
  product: Product
}>()

const cartStore = useCartStore()
const isAdding = ref(false)

const effectivePrice = computed(() => {
  if (props.product.sale_price > 0 && props.product.sale_price < props.product.regular_price) {
    return props.product.sale_price
  }
  return props.product.regular_price
})

const hasDiscount = computed(() => {
  return props.product.sale_price > 0 && props.product.sale_price < props.product.regular_price
})

const discountPercentage = computed(() => {
  if (!hasDiscount.value) return 0
  const diff = props.product.regular_price - props.product.sale_price
  return Math.round((diff / props.product.regular_price) * 100)
})

const handleAddToCart = async () => {
  isAdding.value = true
  try {
    await cartStore.addToCart(props.product.id, 1)
  } finally {
    setTimeout(() => {
      isAdding.value = false
    }, 400)
  }
}

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
}
</script>

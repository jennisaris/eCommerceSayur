<template>
  <div class="relative rounded-lg break-words border bg-white border-gray-300 transition duration-75 hover:transition hover:duration-500 ease-in-out hover:border-green-600 hover:shadow-md card-product flex flex-col justify-between h-full group">
    <div class="p-4 flex-auto flex flex-col justify-between">
      <!-- Image & Badges -->
      <div class="text-center relative flex justify-center mb-3">
        <!-- Discount / Status Badge -->
        <div v-if="hasDiscount" class="absolute top-0 left-0 z-10">
          <span class="inline-block px-1.5 py-0.5 text-center font-bold text-xs align-baseline leading-none rounded bg-red-600 text-white">
            -{{ discountPercentage }}%
          </span>
        </div>
        <div v-else-if="product.stock <= 0" class="absolute top-0 left-0 z-10">
          <span class="inline-block px-1.5 py-0.5 text-center font-bold text-xs align-baseline leading-none rounded bg-gray-600 text-white">
            Habis
          </span>
        </div>
        <div v-else class="absolute top-0 left-0 z-10">
          <span class="inline-block px-1.5 py-0.5 text-center font-bold text-xs align-baseline leading-none rounded bg-green-600 text-white">
            Segar
          </span>
        </div>

        <NuxtLink :to="`/products/${product.slug}`" class="block w-full overflow-hidden aspect-square flex items-center justify-center">
          <img
            :src="product.image || '/assets/images/products/product-img-1.jpg'"
            :alt="product.name"
            class="w-full h-auto max-h-44 object-contain group-hover:scale-105 transition-transform duration-300"
            loading="lazy"
          />
        </NuxtLink>

        <!-- Hover Quick Actions -->
        <div class="absolute w-full bottom-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center gap-2">
          <NuxtLink
            :to="`/products/${product.slug}`"
            class="h-8 w-8 bg-white shadow-md inline-flex items-center justify-center rounded-lg text-gray-700 hover:bg-green-600 hover:text-white transition-colors"
            title="Lihat Detail"
          >
            <i class="ti ti-eye text-base"></i>
          </NuxtLink>
          <button
            type="button"
            class="h-8 w-8 bg-white shadow-md inline-flex items-center justify-center rounded-lg text-gray-700 hover:bg-green-600 hover:text-white transition-colors"
            title="Suka"
          >
            <i class="ti ti-heart text-base"></i>
          </button>
        </div>
      </div>

      <!-- Product Meta & Info -->
      <div class="flex flex-col gap-2 mt-auto">
        <div class="flex items-center justify-between text-xs text-gray-500">
          <span>{{ product.category?.name || 'Sayuran & Segar' }}</span>
          <span v-if="product.weight" class="text-gray-400">{{ product.weight }}g</span>
        </div>

        <h3 class="text-sm font-semibold text-gray-800 line-clamp-1 hover:text-green-600 transition-colors">
          <NuxtLink :to="`/products/${product.slug}`">
            {{ product.name }}
          </NuxtLink>
        </h3>

        <!-- Rating Stars -->
        <div class="flex items-center gap-1 text-yellow-500 text-xs">
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-filled"></i>
          <i class="ti ti-star-half-filled"></i>
          <span class="text-gray-500 text-[11px] ml-1 font-medium">4.8 (24)</span>
        </div>

        <!-- Price & Add Button -->
        <div class="flex justify-between items-center mt-2 pt-2 border-t border-gray-100">
          <div>
            <span class="text-gray-900 font-bold text-sm block">
              {{ formatIDR(effectivePrice) }}
            </span>
            <span v-if="hasDiscount" class="line-through text-gray-400 text-xs">
              {{ formatIDR(product.regular_price) }}
            </span>
          </div>

          <div>
            <button
              @click="handleAddToCart"
              :disabled="product.stock <= 0 || isAdding"
              type="button"
              class="btn inline-flex items-center gap-x-1 bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:bg-green-700 active:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-400 text-xs py-1 px-2.5 rounded-lg font-semibold transition-all shadow-xs"
            >
              <i v-if="isAdding" class="ti ti-loader-2 animate-spin text-xs"></i>
              <i v-else class="ti ti-plus text-xs"></i>
              <span>{{ isAdding ? '...' : 'Add' }}</span>
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

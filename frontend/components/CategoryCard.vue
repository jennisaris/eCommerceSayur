<template>
  <NuxtLink :to="`/products?category_id=${category.id}`" class="block">
    <div class="relative rounded-xl break-words border bg-white border-gray-200 transition duration-150 hover:duration-300 ease-in-out hover:border-indigo-600 hover:shadow-md p-3 text-center group h-full flex flex-col justify-center items-center">
      <div class="py-3 text-center flex flex-col items-center">
        <div class="w-16 h-16 rounded-xl bg-indigo-50 flex items-center justify-center mb-3 group-hover:scale-110 transition-transform duration-300">
          <i class="ti ti-books text-2xl text-indigo-600"></i>
        </div>
        <div class="text-xs sm:text-sm font-semibold text-gray-800 group-hover:text-indigo-600 transition-colors line-clamp-1">
          {{ category.name }}
        </div>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import type { Category } from '~/types'

defineProps<{
  category: Category
}>()
</script>

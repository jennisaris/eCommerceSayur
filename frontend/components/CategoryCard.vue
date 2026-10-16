<template>
  <NuxtLink :to="`/products?category_id=${category.id}`" class="block">
    <div class="relative rounded-lg break-words border bg-white border-gray-300 transition duration-75 hover:transition hover:duration-500 ease-in-out hover:border-green-600 hover:shadow-md p-3 text-center group h-full flex flex-col justify-center items-center">
      <div class="py-3 text-center flex flex-col items-center">
        <img
          :src="category.icon || '/assets/images/category/category-dairy-bread-eggs.jpg'"
          :alt="category.name"
          class="mb-3 mx-auto w-20 h-20 object-contain group-hover:scale-105 transition-transform duration-300"
        />
        <div class="text-xs sm:text-sm font-semibold text-gray-800 group-hover:text-green-600 transition-colors line-clamp-1">
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

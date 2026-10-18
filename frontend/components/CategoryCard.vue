<template>
  <NuxtLink :to="`/products?category_id=${category.id}`" class="block group h-full">
    <div class="relative rounded-2xl break-words border bg-white border-slate-200/90 transition duration-200 ease-in-out hover:border-indigo-600 hover:shadow-lg hover:-translate-y-1 p-4 text-center h-full flex flex-col justify-center items-center shadow-xs">
      <div class="text-center flex flex-col items-center w-full">
        <div class="w-16 h-16 sm:w-20 sm:h-20 rounded-2xl p-1 flex items-center justify-center mb-3 group-hover:scale-105 transition-transform duration-300">
          <img
            :src="categoryImage"
            :alt="category.name"
            class="w-full h-full object-contain filter drop-shadow-xs"
          />
        </div>
        <div class="text-xs sm:text-sm font-bold text-slate-800 group-hover:text-indigo-600 transition-colors line-clamp-1 w-full">
          {{ category.name }}
        </div>
        <span class="text-[11px] text-slate-400 mt-1 group-hover:text-indigo-600 font-semibold transition-colors flex items-center gap-0.5">
          <span>Koleksi Buku</span>
          <i class="ti ti-chevron-right text-xs"></i>
        </span>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import type { Category } from '~/types'

const props = defineProps<{
  category: Category
}>()

// Smart mapping to our bespoke high-res category illustration logos
const categoryImage = computed(() => {
  if (props.category.icon && props.category.icon.startsWith('/images/category/')) {
    return props.category.icon
  }
  const slug = (props.category.slug || '').toLowerCase()
  if (slug.includes('fiksi') || slug.includes('sastra')) return '/images/category/cat-fiksi.svg'
  if (slug.includes('self') || slug.includes('motivasi') || slug.includes('pengembangan')) return '/images/category/cat-self-improvement.svg'
  if (slug.includes('bisnis') || slug.includes('finansial') || slug.includes('investasi')) return '/images/category/cat-bisnis.svg'
  if (slug.includes('tekno') || slug.includes('coding') || slug.includes('sains')) return '/images/category/cat-teknologi.svg'
  if (slug.includes('komik') || slug.includes('manga') || slug.includes('grafis')) return '/images/category/cat-manga.svg'
  if (slug.includes('anak') || slug.includes('edukasi') || slug.includes('pelajar')) return '/images/category/cat-anak.svg'
  return '/images/category/cat-fiksi.svg'
})
</script>

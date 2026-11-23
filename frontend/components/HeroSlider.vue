<template>
  <section class="mt-4">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div
        class="relative min-h-[480px] overflow-hidden rounded-3xl border border-slate-800 bg-slate-950 shadow-xl"
        role="region"
        aria-roledescription="carousel"
        aria-label="Promo BookVerse"
        @mouseenter="pauseSlider"
        @mouseleave="startSlider"
      >
        <div
          v-for="(slide, index) in slides"
          :key="slide.image"
          class="absolute inset-0 transition-all duration-700 ease-out"
          :class="index === currentSlide ? 'opacity-100 translate-x-0' : 'pointer-events-none opacity-0 translate-x-4'"
          :aria-hidden="index !== currentSlide"
        >
          <img :src="slide.image" :alt="slide.imageAlt" class="absolute inset-0 h-full w-full object-cover" />
          <div class="absolute inset-0 bg-gradient-to-r from-slate-950/95 via-slate-950/65 to-slate-950/10"></div>

          <div class="relative z-10 flex min-h-[480px] max-w-2xl flex-col justify-center px-8 py-16 sm:px-14 lg:py-24">
            <span class="mb-4 inline-flex w-fit items-center gap-2 rounded-full bg-amber-400 px-3.5 py-1 text-xs font-bold tracking-wide text-slate-950 shadow-sm">
              <i class="ti ti-sparkles text-xs"></i>
              {{ slide.eyebrow }}
            </span>

            <h1 class="mb-4 text-3xl font-extrabold leading-tight text-white drop-shadow-md sm:text-5xl">
              {{ slide.title }} <br />
              <span class="font-serif italic text-amber-400">{{ slide.highlight }}</span>
            </h1>

            <p class="mb-8 max-w-xl text-sm font-normal leading-relaxed text-slate-200 drop-shadow-sm sm:text-base">
              {{ slide.description }}
            </p>

            <div class="flex flex-wrap items-center gap-4">
              <NuxtLink
                :to="slide.primaryLink"
                class="btn inline-flex items-center gap-2 rounded-xl bg-indigo-600 px-7 py-3.5 text-sm font-bold text-white shadow-lg shadow-indigo-600/30 transition-all hover:-translate-y-0.5 hover:bg-indigo-700"
              >
                <span>{{ slide.primaryLabel }}</span>
                <i class="ti ti-arrow-right text-base"></i>
              </NuxtLink>

              <NuxtLink
                :to="slide.secondaryLink"
                class="btn inline-flex items-center gap-2 rounded-xl border border-white/20 bg-white/10 px-6 py-3.5 text-sm font-bold text-white backdrop-blur-md transition-all hover:bg-white/20"
              >
                {{ slide.secondaryLabel }}
              </NuxtLink>
            </div>
          </div>
        </div>

        <button
          type="button"
          aria-label="Slide sebelumnya"
          class="absolute left-4 top-1/2 z-20 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-white/20 bg-slate-950/35 text-white backdrop-blur-md transition hover:bg-white/20"
          @click="previousSlide"
        >
          <i class="ti ti-chevron-left text-xl"></i>
        </button>
        <button
          type="button"
          aria-label="Slide berikutnya"
          class="absolute right-4 top-1/2 z-20 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-white/20 bg-slate-950/35 text-white backdrop-blur-md transition hover:bg-white/20"
          @click="nextSlide"
        >
          <i class="ti ti-chevron-right text-xl"></i>
        </button>

        <div class="absolute bottom-6 left-1/2 z-20 flex -translate-x-1/2 items-center gap-2" role="tablist" aria-label="Pilih slide">
          <button
            v-for="(slide, index) in slides"
            :key="`indicator-${slide.image}`"
            type="button"
            :aria-label="`Buka slide ${index + 1}`"
            :aria-selected="index === currentSlide"
            role="tab"
            class="h-2 rounded-full transition-all"
            :class="index === currentSlide ? 'w-8 bg-amber-400' : 'w-2 bg-white/50 hover:bg-white/80'"
            @click="goToSlide(index)"
          ></button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

type HeroSlide = {
  image: string
  imageAlt: string
  eyebrow: string
  title: string
  highlight: string
  description: string
  primaryLabel: string
  primaryLink: string
  secondaryLabel: string
  secondaryLink: string
}

const slides: HeroSlide[] = [
  {
    image: '/images/slider/bookstore-hero-modern.png',
    imageAlt: 'Tumpukan buku premium di sudut baca modern dengan pencahayaan hangat',
    eyebrow: 'Pesta Literasi & Book Fair 2026',
    title: 'Jelajahi Dunia Lewat',
    highlight: 'Buku Original Terbaik',
    description: 'Temukan ribuan karya novel, pengembangan diri, bisnis, teknologi, dan komik berlisensi resmi dengan pengemasan aman serta pengiriman cepat ke seluruh Indonesia.',
    primaryLabel: 'Beli Buku Sekarang',
    primaryLink: '/products',
    secondaryLabel: 'Novel Best Seller',
    secondaryLink: '/products?category=fiksi-sastra',
  },
  {
    image: '/images/slider/bookstore-hero-collection.png',
    imageAlt: 'Koleksi buku pilihan BookVerse',
    eyebrow: 'Koleksi Pilihan Pembaca',
    title: 'Temukan Bacaan yang',
    highlight: 'Menginspirasi Anda',
    description: 'Dari cerita yang menghibur sampai wawasan yang mengubah cara pandang, pilih buku terbaik untuk menemani setiap perjalanan Anda.',
    primaryLabel: 'Lihat Koleksi',
    primaryLink: '/products',
    secondaryLabel: 'Kategori Buku',
    secondaryLink: '/products',
  },
  {
    image: '/images/slider/bookstore-hero-reading.png',
    imageAlt: 'Buku dan meja baca dengan pencahayaan elegan',
    eyebrow: 'Promo Pengembangan Diri',
    title: 'Investasi Terbaik untuk',
    highlight: 'Pikiran Anda',
    description: 'Hemat untuk buku psikologi, habit, analisis investasi, dan strategi bisnis modern yang membantu Anda terus bertumbuh.',
    primaryLabel: 'Jelajahi Promo',
    primaryLink: '/products',
    secondaryLabel: 'Buku Bisnis',
    secondaryLink: '/products?category=bisnis-investasi',
  },
]

const currentSlide = ref(0)
let sliderTimer: ReturnType<typeof setInterval> | undefined

function nextSlide() {
  currentSlide.value = (currentSlide.value + 1) % slides.length
}

function previousSlide() {
  currentSlide.value = (currentSlide.value - 1 + slides.length) % slides.length
}

function goToSlide(index: number) {
  currentSlide.value = index
  startSlider()
}

function pauseSlider() {
  if (sliderTimer) clearInterval(sliderTimer)
  sliderTimer = undefined
}

function startSlider() {
  pauseSlider()
  sliderTimer = setInterval(nextSlide, 6000)
}

onMounted(startSlider)
onUnmounted(pauseSlider)
</script>

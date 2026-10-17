import dns from 'node:dns'

// Force IPv4 first to prevent IPv6 ::1 loopback mismatch (426 Upgrade Required)
dns.setDefaultResultOrder('ipv4first')

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-10-01',
  devtools: { enabled: false },
  devServer: {
    host: '127.0.0.1',
    port: 3000,
  },
  vite: {
    server: {
      host: '127.0.0.1',
      strictPort: true,
      hmr: {
        protocol: 'ws',
        host: '127.0.0.1',
        port: 3000,
        clientPort: 3000,
      },
    },
  },
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
  ],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
      midtransClientKey: process.env.NUXT_PUBLIC_MIDTRANS_CLIENT_KEY || 'SB-Mid-client-TEST-KEY-SAMPLE',
      midtransSnapUrl: process.env.NUXT_PUBLIC_MIDTRANS_SNAP_URL || 'https://app.sandbox.midtrans.com/snap/snap.js',
    }
  },
  app: {
    head: {
      title: 'BookVerse - Toko Buku Online Terlengkap & Pengiriman Cepat 24/7',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Platform toko buku online modern dengan koleksi buku fiksi, non-fiksi, bisnis, akademik, dan komik original terlengkap. Pengiriman cepat dan bergaransi resmi se-Indonesia.' }
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@300;400;500;600;700;800&family=Inter:wght@400;500;600;700&display=swap' },
        { rel: 'stylesheet', href: 'https://cdn.jsdelivr.net/npm/@tabler/icons-webfont@2.46.0/tabler-icons.min.css' },
        { rel: 'stylesheet', href: '/assets/css/theme.min.css' }
      ]
    }
  }
})

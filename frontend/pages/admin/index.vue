<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-black text-slate-900 tracking-tight">Ringkasan Operasional Toko</h1>
      <p class="text-xs text-slate-500 mt-1">Monitoring status transaksi, perputaran produk, dan logistik pengantaran 24/7.</p>
    </div>

    <!-- Stat Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="p-5 rounded-2xl bg-white border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="p-3 rounded-xl bg-brand-50 text-brand-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
          </svg>
        </div>
        <div>
          <span class="text-xs font-semibold text-slate-400">Total Pesanan</span>
          <div class="text-xl font-black text-slate-900">{{ stats.totalOrders }}</div>
        </div>
      </div>

      <div class="p-5 rounded-2xl bg-white border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="p-3 rounded-xl bg-emerald-50 text-emerald-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div>
          <span class="text-xs font-semibold text-slate-400">Total Pendapatan</span>
          <div class="text-xl font-black text-slate-900">{{ formatIDR(stats.totalRevenue) }}</div>
        </div>
      </div>

      <div class="p-5 rounded-2xl bg-white border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="p-3 rounded-xl bg-purple-50 text-purple-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1M5 17a2 2 0 104 0m-4 0a2 2 0 114 0m6 0a2 2 0 104 0m-4 0a2 2 0 114 0" />
          </svg>
        </div>
        <div>
          <span class="text-xs font-semibold text-slate-400">Pengiriman Aktif</span>
          <div class="text-xl font-black text-slate-900">{{ stats.activeDeliveries }}</div>
        </div>
      </div>

      <div class="p-5 rounded-2xl bg-white border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="p-3 rounded-xl bg-amber-50 text-amber-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
        <div>
          <span class="text-xs font-semibold text-slate-400">Total Katalog Produk</span>
          <div class="text-xl font-black text-slate-900">{{ stats.totalProducts }}</div>
        </div>
      </div>
    </div>

    <!-- Quick Shortcuts & Recent Orders -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200 p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="text-sm font-bold text-slate-900">Pesanan Masuk Terbaru</h3>
          <NuxtLink to="/admin/orders" class="text-xs font-semibold text-brand-600 hover:underline">
            Lihat Semua Pesanan
          </NuxtLink>
        </div>

        <div v-if="recentOrders.length === 0" class="text-center py-10 text-xs text-slate-400">
          Belum ada pesanan terbaru
        </div>

        <div v-else class="divide-y divide-slate-100">
          <div
            v-for="order in recentOrders"
            :key="order.id"
            class="py-3 flex items-center justify-between text-xs"
          >
            <div>
              <span class="font-bold font-mono text-slate-900">{{ order.order_code }}</span>
              <p class="text-slate-500 mt-0.5">{{ order.recipient_name }} - {{ order.shipping_address }}</p>
            </div>
            <div class="text-right">
              <span class="font-black text-slate-900">{{ formatIDR(order.total_amount) }}</span>
              <div class="text-[10px] font-bold mt-0.5 text-brand-600 uppercase">{{ order.status }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Action Navigation -->
      <div class="bg-white rounded-2xl border border-slate-200 p-6 space-y-4 h-fit">
        <h3 class="text-sm font-bold text-slate-900 border-b border-slate-100 pb-3">Aksi Cepat Admin</h3>
        <div class="space-y-2">
          <NuxtLink
            to="/admin/products"
            class="w-full flex items-center justify-between p-3 rounded-xl bg-slate-50 hover:bg-brand-50 hover:text-brand-700 text-xs font-bold text-slate-700 transition-colors"
          >
            <span>+ Tambah / Kelola Produk</span>
            <span>&rarr;</span>
          </NuxtLink>
          <NuxtLink
            to="/admin/categories"
            class="w-full flex items-center justify-between p-3 rounded-xl bg-slate-50 hover:bg-brand-50 hover:text-brand-700 text-xs font-bold text-slate-700 transition-colors"
          >
            <span>+ Tambah / Kelola Kategori</span>
            <span>&rarr;</span>
          </NuxtLink>
          <NuxtLink
            to="/admin/orders"
            class="w-full flex items-center justify-between p-3 rounded-xl bg-slate-50 hover:bg-brand-50 hover:text-brand-700 text-xs font-bold text-slate-700 transition-colors"
          >
            <span>Proses Pengiriman Pesanan</span>
            <span>&rarr;</span>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Order, Product } from '~/types'

definePageMeta({
  layout: 'admin',
})

const api = useApi()

const stats = reactive({
  totalOrders: 0,
  totalRevenue: 0,
  activeDeliveries: 0,
  totalProducts: 0,
})

const recentOrders = ref<Order[]>([])

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
}

onMounted(async () => {
  try {
    const [ordersRes, prodsRes] = await Promise.all([
      api.get<Order[]>('/admin/orders', { limit: 10 }),
      api.get<Product[]>('/admin/products', { limit: 1, status: 'all' }),
    ])

    if (ordersRes.data) {
      recentOrders.value = ordersRes.data
      stats.totalOrders = ordersRes.meta?.total_items || ordersRes.data.length

      stats.totalRevenue = ordersRes.data.reduce((acc, curr) => {
        if (curr.status === 'PAID' || curr.status === 'COMPLETED' || curr.status === 'PROCESSING' || curr.status === 'ON_DELIVERY') {
          return acc + curr.total_amount
        }
        return acc
      }, 0)

      stats.activeDeliveries = ordersRes.data.filter(
        (o) => o.status === 'PROCESSING' || o.status === 'ON_DELIVERY'
      ).length
    }

    if (prodsRes.meta) {
      stats.totalProducts = prodsRes.meta.total_items
    }
  } catch (err) {
    console.error('Failed to load admin stats:', err)
  }
})
</script>

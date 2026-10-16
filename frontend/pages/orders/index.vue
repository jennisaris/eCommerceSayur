<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
    <div class="mb-8">
      <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900">Riwayat Pesanan</h1>
      <p class="text-xs sm:text-sm text-slate-500 mt-1">Daftar transaksi belanja dan status pengantaran 100% Delivery pesanan Anda.</p>
    </div>

    <!-- Status Tabs -->
    <div class="flex items-center gap-2 overflow-x-auto pb-4 mb-6 border-b border-slate-200">
      <button
        v-for="tab in statusTabs"
        :key="tab.value"
        @click="filterStatus(tab.value)"
        :class="activeStatus === tab.value ? 'bg-brand-600 text-white shadow-xs' : 'bg-white text-slate-600 hover:bg-slate-100 border border-slate-200'"
        class="px-4 py-2 rounded-xl text-xs font-bold transition-all flex-shrink-0"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Orders List -->
    <div v-if="isLoading" class="space-y-4">
      <div v-for="n in 3" :key="n" class="h-40 rounded-2xl bg-slate-200 animate-pulse"></div>
    </div>

    <div v-else-if="orders.length === 0" class="text-center py-20 bg-white rounded-3xl border border-slate-200 p-8">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-100 flex items-center justify-center text-slate-400">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
      </div>
      <h3 class="text-base font-bold text-slate-800">Belum ada pesanan</h3>
      <p class="text-xs text-slate-500 mt-1">Pesanan Anda akan ditampilkan di sini setelah melakukan checkout.</p>
      <NuxtLink to="/products" class="mt-4 inline-block px-5 py-2.5 bg-brand-600 text-white rounded-xl text-xs font-bold">
        Mulai Belanja
      </NuxtLink>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="order in orders"
        :key="order.id"
        class="bg-white rounded-2xl border border-slate-200 p-5 hover:border-slate-300 transition-all shadow-xs flex flex-col md:flex-row items-start md:items-center justify-between gap-4"
      >
        <div class="space-y-2">
          <div class="flex items-center gap-3">
            <span class="text-xs font-black text-brand-600 font-mono">{{ order.order_code }}</span>
            <span :class="getStatusBadgeClass(order.status)" class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold border">
              {{ formatStatus(order.status) }}
            </span>
          </div>

          <div class="text-xs text-slate-500">
            <span>Dipesan pada: {{ formatDate(order.created_at) }}</span>
            <span class="mx-2 text-slate-300">•</span>
            <span>Penerima: <strong class="text-slate-700">{{ order.recipient_name }}</strong></span>
          </div>

          <p class="text-xs text-slate-600 line-clamp-1">
            Alamat: {{ order.shipping_address }}
          </p>
        </div>

        <div class="flex items-center justify-between w-full md:w-auto gap-6 border-t md:border-t-0 pt-3 md:pt-0 border-slate-100">
          <div class="text-left md:text-right">
            <div class="text-[11px] text-slate-400">Total Tagihan</div>
            <div class="text-base font-black text-slate-900">{{ formatIDR(order.total_amount) }}</div>
          </div>

          <NuxtLink
            :to="`/orders/${order.order_code}`"
            class="px-4 py-2 bg-slate-100 hover:bg-brand-50 hover:text-brand-700 border border-slate-200 rounded-xl text-xs font-bold transition-all"
          >
            Lihat Invoice
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Order } from '~/types'

const api = useApi()
const authStore = useAuthStore()
const router = useRouter()

const orders = ref<Order[]>([])
const isLoading = ref(true)
const activeStatus = ref('')

const statusTabs = [
  { label: 'Semua Pesanan', value: '' },
  { label: 'Menunggu Bayar', value: 'PENDING_PAYMENT' },
  { label: 'Dibayar', value: 'PAID' },
  { label: 'Diproses', value: 'PROCESSING' },
  { label: 'Sedang Dikirim', value: 'ON_DELIVERY' },
  { label: 'Selesai', value: 'COMPLETED' },
  { label: 'Dibatalkan', value: 'CANCELLED' },
]

const fetchOrders = async () => {
  isLoading.value = true
  try {
    const params: Record<string, any> = { limit: 50 }
    if (activeStatus.value) params.status = activeStatus.value

    const res = await api.get<Order[]>('/orders', params)
    if (res.data) orders.value = res.data
  } catch (err) {
    console.error('Failed to load orders:', err)
  } finally {
    isLoading.value = false
  }
}

const filterStatus = (val: string) => {
  activeStatus.value = val
  fetchOrders()
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'PAID':
      return 'bg-emerald-50 text-emerald-700 border-emerald-200'
    case 'PROCESSING':
      return 'bg-blue-50 text-blue-700 border-blue-200'
    case 'ON_DELIVERY':
      return 'bg-purple-50 text-purple-700 border-purple-200'
    case 'COMPLETED':
      return 'bg-teal-50 text-teal-700 border-teal-200'
    case 'CANCELLED':
      return 'bg-rose-50 text-rose-700 border-rose-200'
    default:
      return 'bg-amber-50 text-amber-700 border-amber-200'
  }
}

const formatStatus = (status: string) => {
  switch (status) {
    case 'PENDING_PAYMENT':
      return 'Menunggu Pembayaran'
    case 'PAID':
      return 'Pembayaran Diterima'
    case 'PROCESSING':
      return 'Sedang Dikemas'
    case 'ON_DELIVERY':
      return 'Dalam Pengantaran'
    case 'COMPLETED':
      return 'Pesanan Selesai'
    case 'CANCELLED':
      return 'Dibatalkan'
    default:
      return status
  }
}

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
}

const formatDate = (val: string) => {
  if (!val) return ''
  return new Date(val).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  fetchOrders()
})
</script>

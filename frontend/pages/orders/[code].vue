<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
    <!-- Breadcrumb -->
    <div class="flex items-center gap-2 text-xs text-slate-500 mb-6">
      <NuxtLink to="/orders" class="hover:text-brand-600">Riwayat Pesanan</NuxtLink>
      <span>/</span>
      <span class="text-slate-800 font-semibold">{{ order?.order_code || 'Detail Invoice' }}</span>
    </div>

    <div v-if="isLoading" class="h-96 rounded-3xl bg-slate-200 animate-pulse"></div>

    <div v-else-if="!order" class="text-center py-20 bg-white rounded-3xl border border-slate-200 p-8">
      <h3 class="text-base font-bold text-slate-800">Pesanan Tidak Ditemukan</h3>
      <NuxtLink to="/orders" class="mt-4 inline-block px-5 py-2.5 bg-brand-600 text-white rounded-xl text-xs font-bold">
        Kembali ke Daftar Pesanan
      </NuxtLink>
    </div>

    <div v-else class="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
      <!-- Invoice Header -->
      <div class="p-6 sm:p-8 border-b border-slate-100 bg-slate-50/50 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">Invoice Resmi</span>
          <h1 class="text-xl sm:text-2xl font-black text-slate-900 mt-1 font-mono">
            {{ order.order_code }}
          </h1>
          <p class="text-xs text-slate-500 mt-1">Dibuat pada: {{ formatDate(order.created_at) }}</p>
        </div>

        <div class="text-left sm:text-right">
          <span :class="getStatusBadgeClass(order.status)" class="inline-block px-3 py-1 rounded-full text-xs font-extrabold border">
            {{ formatStatus(order.status) }}
          </span>
          <div v-if="order.tracking_number" class="mt-2 text-xs text-purple-700 bg-purple-50 px-2.5 py-1 rounded-lg border border-purple-200">
            Resi: <strong class="font-mono">{{ order.tracking_number }}</strong>
          </div>
        </div>
      </div>

      <div class="p-6 sm:p-8 space-y-8">
        <!-- Delivery Details (100% Delivery Model) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 p-5 rounded-2xl bg-slate-50 border border-slate-100">
          <div>
            <h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">Penerima Paket (100% Delivery)</h3>
            <p class="text-sm font-bold text-slate-900">{{ order.recipient_name }}</p>
            <p class="text-xs text-slate-600 mt-0.5">{{ order.recipient_phone }}</p>
            <p class="text-xs text-slate-600 mt-2 leading-relaxed">
              <strong>Alamat:</strong> {{ order.shipping_address }}
            </p>
          </div>

          <div>
            <h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">Informasi Pembayaran</h3>
            <p class="text-xs text-slate-600">
              Metode: <strong class="text-slate-800">{{ order.payment?.payment_method || 'Midtrans Snap' }}</strong>
            </p>
            <p class="text-xs text-slate-600 mt-1">
              Status Midtrans: <strong class="text-slate-800 uppercase">{{ order.payment?.payment_status || 'PENDING' }}</strong>
            </p>
            <p v-if="order.payment?.payment_gateway_id" class="text-xs text-slate-600 mt-1">
              Gateway ID: <code class="text-[11px] bg-slate-200 px-1 rounded">{{ order.payment.payment_gateway_id }}</code>
            </p>
            <p v-if="order.remarks" class="text-xs text-slate-500 mt-2 italic">
              "Catatan: {{ order.remarks }}"
            </p>
          </div>
        </div>

        <!-- Snapshotted Items Table (CRITICAL RULE: Never recalculated from products) -->
        <div>
          <h3 class="text-sm font-bold text-slate-900 mb-4">Rincian Item (Snapshot Harga Saat Checkout)</h3>
          <div class="border border-slate-200 rounded-2xl overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-slate-50 text-slate-500 font-bold border-b border-slate-200">
                <tr>
                  <th class="p-3.5">Nama Produk</th>
                  <th class="p-3.5 text-center">Jumlah</th>
                  <th class="p-3.5 text-right">Harga Satuan (Snapshot)</th>
                  <th class="p-3.5 text-right">Subtotal</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="item in order.items" :key="item.id">
                  <td class="p-3.5 font-semibold text-slate-900">{{ item.product_name }}</td>
                  <td class="p-3.5 text-center font-bold text-slate-700">{{ item.quantity }}x</td>
                  <td class="p-3.5 text-right text-slate-600">{{ formatIDR(item.price) }}</td>
                  <td class="p-3.5 text-right font-bold text-slate-900">{{ formatIDR(item.subtotal) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Cost Breakdown -->
        <div class="border-t border-slate-100 pt-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div v-if="order.status === 'PENDING_PAYMENT'" class="w-full sm:w-auto">
            <button
              @click="handlePayNow"
              :disabled="isPaying"
              class="w-full sm:w-auto px-6 py-3.5 bg-brand-600 hover:bg-brand-700 text-white font-extrabold text-xs sm:text-sm rounded-xl shadow-lg shadow-brand-500/25 flex items-center justify-center gap-2 transition-all"
            >
              <svg v-if="isPaying" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
              </svg>
              <span>Bayar Sekarang (Midtrans Snap)</span>
            </button>
          </div>

          <div class="w-full sm:w-72 space-y-2 text-xs ml-auto">
            <div class="flex justify-between text-slate-600">
              <span>Subtotal Item</span>
              <span class="font-bold text-slate-800">{{ formatIDR(calculateItemsSubtotal) }}</span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span>Ongkos Kirim (100% Delivery)</span>
              <span class="font-bold text-brand-600">{{ formatIDR(order.shipping_fee) }}</span>
            </div>
            <div class="flex justify-between text-sm pt-2 border-t border-slate-200">
              <span class="font-black text-slate-900">Total Pembayaran</span>
              <span class="font-black text-brand-600 text-base">{{ formatIDR(order.total_amount) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Order } from '~/types'

const route = useRoute()
const router = useRouter()
const api = useApi()

const order = ref<Order | null>(null)
const isLoading = ref(true)
const isPaying = ref(false)

const calculateItemsSubtotal = computed(() => {
  if (!order.value?.items) return 0
  return order.value.items.reduce((sum, it) => sum + it.subtotal, 0)
})

const fetchOrder = async () => {
  const code = route.params.code as string
  isLoading.value = true
  try {
    const res = await api.get<Order>(`/orders/${code}`)
    if (res.data) order.value = res.data
  } catch (err) {
    console.error('Failed to load order invoice:', err)
  } finally {
    isLoading.value = false
  }
}

const handlePayNow = async () => {
  if (!order.value) return
  isPaying.value = true
  try {
    const snapRes = await api.post<{ snap_token: string; redirect_url: string }>(
      `/payments/${order.value.order_code}/snap`
    )
    const snapToken = snapRes.data?.snap_token

    if (typeof window !== 'undefined' && window.snap && snapToken) {
      window.snap.pay(snapToken, {
        onSuccess: () => {
          router.push(`/payment/success?order_code=${order.value?.order_code}`)
        },
        onPending: () => {
          router.push(`/payment/pending?order_code=${order.value?.order_code}`)
        },
        onError: () => {
          router.push(`/payment/failed?order_code=${order.value?.order_code}`)
        },
        onClose: () => {
          fetchOrder()
        },
      })
    } else if (snapRes.data?.redirect_url) {
      window.location.href = snapRes.data.redirect_url
    }
  } catch (err: any) {
    alert(err.message || 'Gagal memicu pembayaran Midtrans')
  } finally {
    isPaying.value = false
  }
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
  fetchOrder()
})
</script>

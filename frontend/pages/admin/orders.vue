<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">Pesanan Masuk & Pengiriman</h1>
        <p class="text-xs text-slate-500 mt-1">Kelola status proses pesanan, input nomor resi kurir, dan pengantaran 100% Delivery.</p>
      </div>

      <!-- Status Filter -->
      <select
        v-model="statusFilter"
        @change="fetchOrders"
        class="px-3.5 py-2 text-xs bg-white border border-slate-200 rounded-xl font-semibold text-slate-700 focus:outline-none focus:border-brand-500"
      >
        <option value="">Semua Status Transaksi</option>
        <option value="PENDING_PAYMENT">Menunggu Pembayaran</option>
        <option value="PAID">Sudah Dibayar (Siap Proses)</option>
        <option value="PROCESSING">Sedang Dikemas (Packing)</option>
        <option value="ON_DELIVERY">Dalam Pengiriman Kurir</option>
        <option value="COMPLETED">Selesai Diterima</option>
        <option value="CANCELLED">Dibatalkan</option>
      </select>
    </div>

    <!-- Orders Table -->
    <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-xs">
      <div v-if="isLoading" class="p-8 text-center text-xs text-slate-400">
        Memuat data pesanan...
      </div>

      <div v-else-if="orders.length === 0" class="p-12 text-center text-xs text-slate-500">
        Tidak ada pesanan untuk filter status ini.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-slate-50 text-slate-500 font-bold border-b border-slate-200 uppercase tracking-wider">
            <tr>
              <th class="p-4">Kode Invoice</th>
              <th class="p-4">Penerima & Alamat (100% Delivery)</th>
              <th class="p-4">Total Tagihan</th>
              <th class="p-4">Resi Kurir</th>
              <th class="p-4">Status Pesanan</th>
              <th class="p-4 text-right">Ubah Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="order in orders" :key="order.id" class="hover:bg-slate-50/50">
              <td class="p-4">
                <span class="font-bold font-mono text-slate-900 block">{{ order.order_code }}</span>
                <span class="text-[10px] text-slate-400">{{ formatDate(order.created_at) }}</span>
              </td>
              <td class="p-4 max-w-xs">
                <div class="font-bold text-slate-900">{{ order.recipient_name }} ({{ order.recipient_phone }})</div>
                <div class="text-slate-500 text-[11px] truncate mt-0.5">{{ order.shipping_address }}</div>
              </td>
              <td class="p-4">
                <div class="font-black text-slate-900">{{ formatIDR(order.total_amount) }}</div>
                <div class="text-[10px] text-slate-400">Ongkir: {{ formatIDR(order.shipping_fee) }}</div>
              </td>
              <td class="p-4">
                <span v-if="order.tracking_number" class="font-mono text-xs font-bold text-purple-700 bg-purple-50 px-2 py-0.5 rounded border border-purple-200">
                  {{ order.tracking_number }}
                </span>
                <span v-else class="text-slate-400 italic text-[11px]">Belum dikirim</span>
              </td>
              <td class="p-4">
                <span :class="getStatusBadgeClass(order.status)" class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold border">
                  {{ formatStatus(order.status) }}
                </span>
              </td>
              <td class="p-4 text-right">
                <button
                  @click="openStatusModal(order)"
                  class="px-3 py-1.5 bg-slate-100 hover:bg-brand-50 hover:text-brand-700 text-slate-700 rounded-xl text-xs font-bold border border-slate-200 transition-colors"
                >
                  Kelola Status
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Update Status Modal -->
    <div v-if="isModalOpen && selectedOrder" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="w-full max-w-md bg-white rounded-3xl p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="text-sm font-bold text-slate-900">
            Perbarui Status Pesanan
          </h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600 text-sm">&times;</button>
        </div>

        <div class="text-xs text-slate-600 bg-slate-50 p-3 rounded-xl border border-slate-100 space-y-1">
          <div>Kode: <strong class="font-mono text-slate-900">{{ selectedOrder.order_code }}</strong></div>
          <div>Penerima: <strong class="text-slate-900">{{ selectedOrder.recipient_name }}</strong></div>
          <div>Status Saat Ini: <strong class="text-brand-600">{{ selectedOrder.status }}</strong></div>
        </div>

        <form @submit.prevent="submitStatusUpdate" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Status Baru</label>
            <select
              v-model="newStatus"
              required
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
            >
              <option value="PROCESSING">PROCESSING (Sedang Dikemas)</option>
              <option value="ON_DELIVERY">ON_DELIVERY (Dalam Pengantaran Kurir)</option>
              <option value="COMPLETED">COMPLETED (Selesai)</option>
              <option value="CANCELLED">CANCELLED (Batalkan & Restock Stok)</option>
            </select>
          </div>

          <div v-if="newStatus === 'ON_DELIVERY'">
            <label class="block text-xs font-bold text-slate-700 mb-1">
              Nomor Resi Pengiriman <span class="text-rose-500">*</span>
            </label>
            <input
              v-model="trackingNumber"
              type="text"
              required
              placeholder="Contoh: JNE-9876543210 atau KURIR-01"
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500 font-mono"
            />
          </div>

          <div v-if="newStatus === 'CANCELLED'" class="p-3 bg-rose-50 border border-rose-200 rounded-xl text-xs text-rose-700 leading-relaxed">
            ⚠️ Membatalkan pesanan ini akan secara otomatis mengembalikan (restock) seluruh stok item ke database produk.
          </div>

          <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100">
            <button
              @click="isModalOpen = false"
              type="button"
              class="px-4 py-2 text-xs font-semibold text-slate-600 hover:bg-slate-100 rounded-xl"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-5 py-2 text-xs font-bold bg-brand-600 hover:bg-brand-700 text-white rounded-xl shadow-md"
            >
              Simpan Perubahan
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Order } from '~/types'

definePageMeta({
  layout: 'admin',
})

const api = useApi()
const orders = ref<Order[]>([])
const isLoading = ref(true)
const statusFilter = ref('')

const isModalOpen = ref(false)
const selectedOrder = ref<Order | null>(null)
const newStatus = ref('PROCESSING')
const trackingNumber = ref('')

const fetchOrders = async () => {
  isLoading.value = true
  try {
    const params: Record<string, any> = { limit: 100 }
    if (statusFilter.value) params.status = statusFilter.value

    const res = await api.get<Order[]>('/admin/orders', params)
    if (res.data) orders.value = res.data
  } catch (err) {
    console.error('Failed to load admin orders:', err)
  } finally {
    isLoading.value = false
  }
}

const openStatusModal = (order: Order) => {
  selectedOrder.value = order
  newStatus.value = order.status === 'PAID' ? 'PROCESSING' : order.status === 'PROCESSING' ? 'ON_DELIVERY' : 'COMPLETED'
  trackingNumber.value = order.tracking_number || ''
  isModalOpen.value = true
}

const submitStatusUpdate = async () => {
  if (!selectedOrder.value) return
  try {
    await api.put(`/admin/orders/${selectedOrder.value.id}/status`, {
      status: newStatus.value,
      tracking_number: trackingNumber.value,
    })
    isModalOpen.value = false
    await fetchOrders()
  } catch (err: any) {
    alert(err.message || 'Gagal memperbarui status pesanan')
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
      return 'Menunggu Bayar'
    case 'PAID':
      return 'Sudah Dibayar'
    case 'PROCESSING':
      return 'Sedang Dikemas'
    case 'ON_DELIVERY':
      return 'Pengantaran'
    case 'COMPLETED':
      return 'Selesai'
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
  fetchOrders()
})
</script>

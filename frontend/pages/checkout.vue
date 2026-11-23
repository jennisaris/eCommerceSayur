<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-10">
    <!-- Breadcrumbs -->
    <nav class="flex items-center gap-2 text-xs text-gray-500 mb-6">
      <NuxtLink to="/" class="hover:text-indigo-600 flex items-center gap-1">
        <i class="ti ti-home text-sm"></i>
        <span>Home</span>
      </NuxtLink>
      <span>/</span>
      <NuxtLink to="/cart" class="hover:text-indigo-600">Keranjang</NuxtLink>
      <span>/</span>
      <span class="text-gray-800 font-semibold">Checkout Buku</span>
    </nav>

    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl sm:text-3xl font-extrabold text-gray-900">Checkout & Pengiriman Buku</h1>
      <p class="text-xs sm:text-sm text-gray-500 mt-1">
        Lengkapi detail alamat pengiriman buku Anda dengan proteksi kemasan kardus sudut & bubble wrap tebal, bayar praktis via Midtrans Snap.
      </p>
    </div>

    <!-- Empty cart state -->
    <div v-if="cartStore.items.length === 0" class="text-center py-20 bg-white rounded-2xl border border-gray-200 p-8">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-400">
        <i class="ti ti-book-off text-3xl"></i>
      </div>
      <h3 class="text-base font-bold text-gray-800">Keranjang belanja Anda kosong</h3>
      <p class="text-xs text-gray-500 mt-1">Silakan pilih buku terlebih dahulu sebelum melakukan checkout.</p>
      <NuxtLink to="/products" class="mt-4 inline-block px-5 py-2.5 bg-indigo-600 text-white rounded-lg text-xs font-bold hover:bg-indigo-700 shadow-xs transition-all">
        Kembali ke Katalog Buku
      </NuxtLink>
    </div>

    <!-- Checkout Form Grid -->
    <form v-else @submit.prevent="handleCheckout" class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
      <!-- Left Column: Delivery Form -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Packaging & Delivery Protection Card -->
        <div class="p-4 rounded-xl bg-indigo-50 border border-indigo-200 flex items-start gap-3.5">
          <div class="w-10 h-10 rounded-lg bg-indigo-600 text-white flex items-center justify-center flex-shrink-0 shadow-xs">
            <i class="ti ti-shield-check text-xl"></i>
          </div>
          <div class="text-xs">
            <div class="font-bold text-indigo-950 text-sm">Standar Pengiriman Buku Anti-Tekuk</div>
            <p class="text-indigo-800 mt-0.5 leading-relaxed">
              Setiap buku dikemas rapat dengan lapisan double bubble wrap serta pelindung karton tebal di setiap sudut, memastikan buku sampai dalam kondisi 100% mulus dan sempurna.
            </p>
          </div>
        </div>

        <!-- Form Card: Delivery Details -->
        <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-xs space-y-5">
          <h2 class="text-base font-bold text-gray-900 border-b border-gray-100 pb-3 flex items-center gap-2">
            <i class="ti ti-map-pin text-indigo-600 text-lg"></i>
            <span>Data Penerima & Alamat Pengiriman</span>
          </h2>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-gray-700 mb-1.5">
                Nama Lengkap Penerima <span class="text-red-500">*</span>
              </label>
              <input
                v-model="form.recipient_name"
                type="text"
                required
                placeholder="Contoh: Budi Santoso"
                class="w-full px-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-700 mb-1.5">
                Nomor Telepon / WhatsApp <span class="text-red-500">*</span>
              </label>
              <input
                v-model="form.recipient_phone"
                type="tel"
                required
                placeholder="Contoh: 081234567890"
                class="w-full px-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1.5">
              Alamat Lengkap Pengantaran <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="form.shipping_address"
              required
              rows="3"
              placeholder="Sebutkan nama jalan, nomor rumah, RT/RW, kelurahan, kecamatan, kota, kode pos, dan patokan tujuan..."
              class="w-full px-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
            ></textarea>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1.5">
              Catatan untuk Kurir (Opsional)
            </label>
            <input
              v-model="form.remarks"
              type="text"
              placeholder="Contoh: Titip di pos sekuriti jika gerbang terkunci"
              class="w-full px-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
            />
          </div>
        </div>

        <!-- Payment Gateway Selection Info -->
        <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-xs space-y-4">
          <h2 class="text-base font-bold text-gray-900 border-b border-gray-100 pb-3 flex items-center gap-2">
            <i class="ti ti-credit-card text-indigo-600 text-lg"></i>
            <span>Metode Pembayaran (Midtrans Snap)</span>
          </h2>
          <div class="p-4 rounded-lg bg-gray-50 border border-gray-200 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <i class="ti ti-shield-lock text-2xl text-indigo-600"></i>
              <div class="text-xs">
                <div class="font-bold text-gray-900">Midtrans Pop-up Snap Gateway</div>
                <div class="text-gray-500">QRIS, BCA, Mandiri, BNI, BRI Virtual Account, dan Kartu Kredit</div>
              </div>
            </div>
            <div class="flex items-center gap-1.5">
              <img src="/assets/images/payment/visa.svg" alt="Visa" class="h-5" />
              <img src="/assets/images/payment/mastercard.svg" alt="MC" class="h-5" />
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Order Details Summary -->
      <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-xs space-y-6 sticky top-24">
        <h3 class="text-base font-bold text-gray-900 border-b border-gray-100 pb-3 flex items-center justify-between">
          <span>Detail Pesanan</span>
          <span class="text-xs font-normal text-gray-500">{{ cartStore.totalItems }} Buku</span>
        </h3>

        <!-- Purchased items list -->
        <div class="max-h-60 overflow-y-auto divide-y divide-gray-100 space-y-3 pr-1">
          <div
            v-for="item in cartStore.items"
            :key="item.id"
            class="flex items-center justify-between text-xs pt-3 first:pt-0"
          >
            <div class="flex items-center gap-3 min-w-0">
              <img
                :src="item.product?.image || '/images/products/book-atomic-habits.jpg'"
                alt="Product"
                class="w-10 h-14 object-cover rounded-md bg-gray-50 border border-gray-200 flex-shrink-0"
              />
              <div class="min-w-0">
                <div class="font-semibold text-gray-900 truncate max-w-[140px]">
                  {{ item.product?.name }}
                </div>
                <div class="text-[11px] text-gray-400">
                  {{ item.quantity }} x {{ formatIDR(item.product?.sale_price && item.product.sale_price > 0 ? item.product.sale_price : item.product?.regular_price || 0) }}
                </div>
              </div>
            </div>
            <span class="font-bold text-gray-900 flex-shrink-0">
              {{ formatIDR((item.product?.sale_price && item.product.sale_price > 0 ? item.product.sale_price : item.product?.regular_price || 0) * item.quantity) }}
            </span>
          </div>
        </div>

        <!-- Cost Breakdown -->
        <div class="space-y-2.5 text-xs pt-4 border-t border-gray-100">
          <div class="flex justify-between text-gray-600">
            <span>Subtotal Buku</span>
            <span class="font-semibold text-gray-900">{{ formatIDR(cartStore.totalPrice) }}</span>
          </div>
          <div class="flex justify-between text-gray-600">
            <span>Total Berat Pengiriman</span>
            <span class="font-semibold text-gray-900">{{ (cartStore.totalWeight / 1000).toFixed(2) }} kg</span>
          </div>
          <div class="flex justify-between text-gray-600">
            <span>Ongkos Kirim Kurir</span>
            <span class="font-semibold text-indigo-600">{{ formatIDR(shippingFee) }}</span>
          </div>
          <div class="flex justify-between pt-3 border-t border-gray-200 text-sm">
            <span class="font-extrabold text-gray-900">Total Tagihan</span>
            <span class="text-base font-extrabold text-indigo-600">{{ formatIDR(totalBill) }}</span>
          </div>
        </div>

        <!-- Submit Button -->
        <button
          :disabled="isSubmitting || cartStore.items.length === 0"
          type="submit"
          class="w-full py-3.5 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm rounded-xl shadow-md shadow-indigo-600/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50"
        >
          <i v-if="isSubmitting" class="ti ti-loader-2 animate-spin text-base"></i>
          <span v-else class="flex items-center gap-1.5">
            <span>Bayar Sekarang (Midtrans Snap)</span>
            <i class="ti ti-arrow-right text-base"></i>
          </span>
        </button>

        <!-- Security guarantee footer -->
        <div class="flex items-center justify-center gap-1.5 text-[11px] text-gray-400">
          <i class="ti ti-shield-check text-indigo-600 text-sm"></i>
          <span>Enkripsi SSL 256-Bit & Idempoten Midtrans</span>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import type { Order, CheckoutEstimateResponse } from '~/types'

const authStore = useAuthStore()
const cartStore = useCartStore()
const api = useApi()
const router = useRouter()

const form = reactive({
  recipient_name: authStore.user?.name || '',
  recipient_phone: authStore.user?.phone || '',
  shipping_address: authStore.user?.address || '',
  remarks: '',
})

const isSubmitting = ref(false)
const shippingFee = ref(15000)

const totalBill = computed(() => {
  return cartStore.totalPrice + shippingFee.value
})

const fetchEstimate = async () => {
  try {
    const res = await api.post<CheckoutEstimateResponse>('/orders/estimate', {
      shipping_address: form.shipping_address,
    })
    if (res.data) {
      shippingFee.value = res.data.shipping_fee
    }
  } catch (err) {
    console.error('Failed to estimate shipping fee:', err)
  }
}

const handleCheckout = async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }

  isSubmitting.value = true
  try {
    // 1. Create order in backend (locks stock, price snapshotting, PENDING_PAYMENT)
    const orderRes = await api.post<Order>('/orders', form)
    const createdOrder = orderRes.data

    if (!createdOrder || !createdOrder.order_code) {
      throw new Error('Gagal memproses pesanan')
    }

    // Order is created and DB cart is cleared; sync cart immediately so UI never stays stale
    await cartStore.fetchCart()

    // 2. Request Midtrans Snap Token
    let snapRes: { data?: { snap_token: string; redirect_url: string } }
    try {
      snapRes = await api.post<{ snap_token: string; redirect_url: string }>(
        '/payments/snap',
        { order_code: createdOrder.order_code }
      )
    } catch (snapErr: any) {
      alert(snapErr.message || 'Pesanan sudah dibuat, tetapi gagal memuat Midtrans. Silakan lanjutkan pembayaran dari halaman pesanan.')
      router.push(`/orders/${createdOrder.order_code}`)
      return
    }

    const snapToken = snapRes.data?.snap_token

    if (!snapToken) {
      router.push(`/orders/${createdOrder.order_code}`)
      return
    }

    // 3. Trigger Midtrans Snap.js Modal Popup
    if (typeof window !== 'undefined' && window.snap) {
      window.snap.pay(snapToken, {
        onSuccess: (result: any) => {
          console.log('Payment success:', result)
          router.push(`/payment/success?order_code=${createdOrder.order_code}`)
        },
        onPending: (result: any) => {
          console.log('Payment pending:', result)
          router.push(`/payment/pending?order_code=${createdOrder.order_code}`)
        },
        onError: (result: any) => {
          console.log('Payment error:', result)
          router.push(`/payment/failed?order_code=${createdOrder.order_code}`)
        },
        onClose: () => {
          router.push(`/orders/${createdOrder.order_code}`)
        },
      })
    } else {
      // Fallback: Redirect to Midtrans Redirect URL if Snap JS not loaded
      if (snapRes.data?.redirect_url && isAllowedMidtransUrl(snapRes.data.redirect_url)) {
        window.location.href = snapRes.data.redirect_url
      } else {
        router.push(`/orders/${createdOrder.order_code}`)
      }
    }
  } catch (err: any) {
    alert(err.message || 'Terjadi kesalahan saat memproses pesanan')
  } finally {
    isSubmitting.value = false
  }
}

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
}

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  await cartStore.fetchCart()
  await fetchEstimate()
})
</script>

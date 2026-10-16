<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">Manajemen Katalog Produk</h1>
        <p class="text-xs text-slate-500 mt-1">Atur ketersediaan stok, harga promo, varian, dan berat pengiriman produk.</p>
      </div>

      <button
        @click="openCreateModal"
        class="px-4 py-2.5 bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs rounded-xl shadow-md transition-all flex items-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span>Tambah Produk Baru</span>
      </button>
    </div>

    <!-- Table Card -->
    <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-xs">
      <div v-if="isLoading" class="p-8 text-center text-xs text-slate-400">
        Memuat data produk...
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-slate-50 text-slate-500 font-bold border-b border-slate-200 uppercase tracking-wider">
            <tr>
              <th class="p-4">Produk</th>
              <th class="p-4">Kategori</th>
              <th class="p-4">Harga & Promo</th>
              <th class="p-4">Stok & Satuan</th>
              <th class="p-4">Berat</th>
              <th class="p-4">Status</th>
              <th class="p-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="prod in products" :key="prod.id" class="hover:bg-slate-50/50">
              <td class="p-4">
                <div class="flex items-center gap-3">
                  <img
                    :src="prod.image || 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=80&h=80&fit=crop'"
                    alt="Product"
                    class="w-10 h-10 rounded-lg object-cover bg-slate-100 border border-slate-200 flex-shrink-0"
                  />
                  <div>
                    <span class="font-bold text-slate-900 line-clamp-1 max-w-xs">{{ prod.name }}</span>
                    <span v-if="prod.variant_name" class="text-[10px] text-slate-400">{{ prod.variant_name }}</span>
                  </div>
                </div>
              </td>
              <td class="p-4 font-semibold text-slate-700">{{ prod.category?.name || '-' }}</td>
              <td class="p-4">
                <div class="font-black text-slate-900">
                  {{ formatIDR(prod.sale_price > 0 ? prod.sale_price : prod.regular_price) }}
                </div>
                <div v-if="prod.sale_price > 0" class="text-[10px] text-slate-400 line-through">
                  {{ formatIDR(prod.regular_price) }}
                </div>
              </td>
              <td class="p-4 font-bold" :class="prod.stock <= 5 ? 'text-rose-600' : 'text-slate-800'">
                {{ prod.stock }} {{ prod.unit }}
              </td>
              <td class="p-4 text-slate-600">{{ prod.weight }}g</td>
              <td class="p-4">
                <span
                  :class="prod.status === 'ACTIVE' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-slate-100 text-slate-500 border-slate-200'"
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold border"
                >
                  {{ prod.status }}
                </span>
              </td>
              <td class="p-4 text-right space-x-2">
                <button
                  @click="openEditModal(prod)"
                  class="px-2.5 py-1 text-xs font-semibold text-brand-600 hover:bg-brand-50 rounded-lg transition-colors"
                >
                  Edit
                </button>
                <button
                  @click="deleteProduct(prod.id)"
                  class="px-2.5 py-1 text-xs font-semibold text-rose-600 hover:bg-rose-50 rounded-lg transition-colors"
                >
                  Hapus
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form (Create / Edit) -->
    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm overflow-y-auto">
      <div class="w-full max-w-2xl bg-white rounded-3xl p-6 shadow-2xl space-y-4 my-8">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="text-sm font-bold text-slate-900">
            {{ isEditing ? 'Edit Data Produk' : 'Tambah Produk Baru' }}
          </h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600 text-sm">&times;</button>
        </div>

        <form @submit.prevent="submitForm" class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Produk</label>
              <input
                v-model="form.name"
                @input="generateSlug"
                type="text"
                required
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Kategori</label>
              <select
                v-model.number="form.category_id"
                required
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              >
                <option value="" disabled>Pilih Kategori</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Slug URL</label>
              <input
                v-model="form.slug"
                type="text"
                required
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl font-mono focus:outline-none focus:border-brand-500"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Varian (Opsional)</label>
              <input
                v-model="form.variant_name"
                type="text"
                placeholder="Misal: Warna Hitam, Size 42"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Harga Normal (IDR)</label>
              <input
                v-model.number="form.regular_price"
                type="number"
                required
                min="1000"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Harga Diskon (0 jika normal)</label>
              <input
                v-model.number="form.sale_price"
                type="number"
                min="0"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Jumlah Stok</label>
              <input
                v-model.number="form.stock"
                type="number"
                required
                min="0"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Berat (Gram - Ongkir)</label>
              <input
                v-model.number="form.weight"
                type="number"
                required
                min="1"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Satuan Barang</label>
              <input
                v-model="form.unit"
                type="text"
                placeholder="pcs, unit, kg, pasang"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Status Penjualan</label>
              <select
                v-model="form.status"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              >
                <option value="ACTIVE">ACTIVE</option>
                <option value="INACTIVE">INACTIVE</option>
                <option value="OUT_OF_STOCK">OUT_OF_STOCK</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">URL Foto Produk</label>
            <div class="flex items-center gap-3">
              <input
                v-model="form.image"
                type="text"
                placeholder="https://... atau /assets/images/products/..."
                class="flex-1 px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
              <img
                v-if="form.image"
                :src="form.image"
                alt="Preview"
                class="w-9 h-9 rounded-lg object-cover bg-slate-100 border border-slate-200 flex-shrink-0"
                @error="(e: any) => e.target.style.display = 'none'"
              />
            </div>
            <p class="text-[10px] text-slate-400 mt-1">Dapat menggunakan URL web (https://) atau file aset lokal (/assets/images/products/...)</p>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Deskripsi Lengkap</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
            ></textarea>
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
              Simpan Produk
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Category, Product } from '~/types'

definePageMeta({
  layout: 'admin',
})

const api = useApi()
const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const isLoading = ref(true)

const isModalOpen = ref(false)
const isEditing = ref(false)
const editId = ref<number | null>(null)

const form = reactive({
  category_id: 1,
  name: '',
  slug: '',
  image: '',
  description: '',
  regular_price: 100000,
  sale_price: 0,
  unit: 'pcs',
  weight: 500,
  stock: 20,
  variant_name: '',
  status: 'ACTIVE',
})

const fetchCategories = async () => {
  try {
    const res = await api.get<Category[]>('/categories')
    if (res.data) categories.value = res.data
  } catch (err) {
    console.error('Failed to load categories:', err)
  }
}

const fetchProducts = async () => {
  isLoading.value = true
  try {
    const res = await api.get<Product[]>('/products', { limit: 100, status: 'all' })
    if (res.data) products.value = res.data
  } catch (err) {
    console.error('Failed to load products:', err)
  } finally {
    isLoading.value = false
  }
}

const generateSlug = () => {
  if (!isEditing.value) {
    form.slug = form.name
      .toLowerCase()
      .trim()
      .replace(/[^\w\s-]/g, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '')
  }
}

const openCreateModal = () => {
  isEditing.value = false
  editId.value = null
  form.category_id = categories.value[0]?.id || 1
  form.name = ''
  form.slug = ''
  form.image = 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=800&fit=crop'
  form.description = ''
  form.regular_price = 150000
  form.sale_price = 0
  form.unit = 'pcs'
  form.weight = 500
  form.stock = 50
  form.variant_name = ''
  form.status = 'ACTIVE'
  isModalOpen.value = true
}

const openEditModal = (prod: Product) => {
  isEditing.value = true
  editId.value = prod.id
  form.category_id = prod.category_id
  form.name = prod.name
  form.slug = prod.slug
  form.image = prod.image || ''
  form.description = prod.description || ''
  form.regular_price = prod.regular_price
  form.sale_price = prod.sale_price
  form.unit = prod.unit
  form.weight = prod.weight
  form.stock = prod.stock
  form.variant_name = prod.variant_name || ''
  form.status = prod.status
  isModalOpen.value = true
}

const submitForm = async () => {
  try {
    if (isEditing.value && editId.value) {
      await api.put(`/admin/products/${editId.value}`, form)
    } else {
      await api.post('/admin/products', form)
    }
    isModalOpen.value = false
    await fetchProducts()
  } catch (err: any) {
    alert(err.message || 'Gagal menyimpan data produk')
  }
}

const deleteProduct = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus produk ini?')) return
  try {
    await api.delete(`/admin/products/${id}`)
    await fetchProducts()
  } catch (err: any) {
    alert(err.message || 'Gagal menghapus produk')
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
  await fetchCategories()
  await fetchProducts()
})
</script>

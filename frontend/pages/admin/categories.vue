<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">Manajemen Kategori</h1>
        <p class="text-xs text-slate-500 mt-1">Kelola struktur kategori produk bertingkat untuk katalog toko.</p>
      </div>

      <button
        @click="openCreateModal"
        class="px-4 py-2.5 bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs rounded-xl shadow-md transition-all flex items-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span>Tambah Kategori Baru</span>
      </button>
    </div>

    <!-- Table Card -->
    <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-xs">
      <div v-if="isLoading" class="p-8 text-center text-xs text-slate-400">
        Memuat data kategori...
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-slate-50 text-slate-500 font-bold border-b border-slate-200 uppercase tracking-wider">
            <tr>
              <th class="p-4">Icon & Kategori</th>
              <th class="p-4">Slug Rute</th>
              <th class="p-4">Deskripsi</th>
              <th class="p-4">Status</th>
              <th class="p-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="cat in categories" :key="cat.id" class="hover:bg-slate-50/50">
              <td class="p-4">
                <div class="flex items-center gap-3">
                  <img
                    :src="cat.icon || 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=80&h=80&fit=crop'"
                    alt="Icon"
                    class="w-9 h-9 rounded-lg object-cover bg-slate-100 border border-slate-200"
                  />
                  <div>
                    <span class="font-bold text-slate-900">{{ cat.name }}</span>
                    <span v-if="cat.children && cat.children.length > 0" class="block text-[10px] text-brand-600">
                      {{ cat.children.length }} sub-kategori
                    </span>
                  </div>
                </div>
              </td>
              <td class="p-4 font-mono text-slate-500">{{ cat.slug }}</td>
              <td class="p-4 text-slate-600 max-w-xs truncate">{{ cat.description || '-' }}</td>
              <td class="p-4">
                <span
                  :class="cat.status ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-slate-100 text-slate-500 border-slate-200'"
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold border"
                >
                  {{ cat.status ? 'Aktif' : 'Non-Aktif' }}
                </span>
              </td>
              <td class="p-4 text-right space-x-2">
                <button
                  @click="openEditModal(cat)"
                  class="px-2.5 py-1 text-xs font-semibold text-brand-600 hover:bg-brand-50 rounded-lg transition-colors"
                >
                  Edit
                </button>
                <button
                  @click="deleteCategory(cat.id)"
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
    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="w-full max-w-lg bg-white rounded-3xl p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="text-sm font-bold text-slate-900">
            {{ isEditing ? 'Edit Kategori' : 'Tambah Kategori Baru' }}
          </h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600 text-sm">&times;</button>
        </div>

        <form @submit.prevent="submitForm" class="space-y-3.5">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Nama Kategori</label>
            <input
              v-model="form.name"
              @input="generateSlug"
              type="text"
              required
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Slug Rute URL</label>
            <input
              v-model="form.slug"
              type="text"
              required
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl font-mono focus:outline-none focus:border-brand-500"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">URL Gambar Icon</label>
            <div class="flex items-center gap-3">
              <input
                v-model="form.icon"
                type="text"
                placeholder="https://... atau /assets/images/category/..."
                class="flex-1 px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
              />
              <img
                v-if="toSafeAssetUrl(form.icon)"
                :src="toSafeAssetUrl(form.icon)"
                alt="Preview"
                class="w-8 h-8 rounded-lg object-cover bg-slate-100 border border-slate-200 flex-shrink-0"
                @error="(e: any) => e.target.style.display = 'none'"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Deskripsi</label>
            <textarea
              v-model="form.description"
              rows="2"
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:border-brand-500"
            ></textarea>
          </div>

          <div class="flex items-center gap-2 pt-1">
            <input
              v-model="form.status"
              type="checkbox"
              id="catStatus"
              class="rounded border-slate-300 text-brand-600"
            />
            <label for="catStatus" class="text-xs font-semibold text-slate-700">Kategori Aktif & Tampil di Toko</label>
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
              Simpan Kategori
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Category } from '~/types'

definePageMeta({
  layout: 'admin',
})

const api = useApi()
const categories = ref<Category[]>([])
const isLoading = ref(true)

const isModalOpen = ref(false)
const isEditing = ref(false)
const editId = ref<number | null>(null)

const form = reactive({
  name: '',
  slug: '',
  icon: '',
  description: '',
  status: true,
})

const fetchCategories = async () => {
  isLoading.value = true
  try {
    const res = await api.get<Category[]>('/admin/categories')
    if (res.data) categories.value = res.data
  } catch (err) {
    console.error('Failed to load categories:', err)
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
  form.name = ''
  form.slug = ''
  form.icon = 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=100&h=100&fit=crop'
  form.description = ''
  form.status = true
  isModalOpen.value = true
}

const openEditModal = (cat: Category) => {
  isEditing.value = true
  editId.value = cat.id
  form.name = cat.name
  form.slug = cat.slug
  form.icon = cat.icon || ''
  form.description = cat.description || ''
  form.status = cat.status
  isModalOpen.value = true
}

const submitForm = async () => {
  try {
    if (isEditing.value && editId.value) {
      await api.put(`/admin/categories/${editId.value}`, form)
    } else {
      await api.post('/admin/categories', form)
    }
    isModalOpen.value = false
    await fetchCategories()
  } catch (err: any) {
    alert(err.message || 'Gagal menyimpan kategori')
  }
}

const deleteCategory = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus kategori ini?')) return
  try {
    await api.delete(`/admin/categories/${id}`)
    await fetchCategories()
  } catch (err: any) {
    alert(err.message || 'Gagal menghapus kategori')
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
    <div class="mb-6 flex items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-extrabold text-slate-900">Profil Saya</h1>
        <p class="text-xs text-slate-500 mt-1">
          Perbarui nama, alamat email, nomor HP, dan alamat pengiriman rumah Anda.
        </p>
      </div>
      <NuxtLink
        to="/orders"
        class="px-4 py-2 text-xs font-bold text-indigo-700 bg-indigo-50 hover:bg-indigo-100 rounded-xl border border-indigo-200 transition-colors"
      >
        Riwayat Pesanan
      </NuxtLink>
    </div>

    <div v-if="isLoadingProfile" class="bg-white rounded-2xl border border-slate-200 p-8 text-center text-xs text-slate-500">
      Memuat profil Anda...
    </div>

    <div v-else class="bg-white rounded-2xl border border-slate-200 shadow-xs p-6 sm:p-8 space-y-6">
      <div v-if="errorMessage" class="p-3 bg-rose-50 text-rose-700 border border-rose-200 text-xs rounded-xl font-medium">
        {{ errorMessage }}
      </div>

      <div v-if="successMessage" class="p-3 bg-emerald-50 text-emerald-700 border border-emerald-200 text-xs rounded-xl font-medium">
        {{ successMessage }}
      </div>

      <form @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">Nama Lengkap</label>
          <input
            v-model="form.name"
            type="text"
            required
            minlength="2"
            maxlength="100"
            placeholder="Nama lengkap Anda"
            class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">Alamat Email</label>
            <input
              v-model="form.email"
              type="email"
              required
              maxlength="100"
              placeholder="nama@email.com"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">Nomor HP / WhatsApp</label>
            <input
              v-model="form.phone"
              type="tel"
              maxlength="25"
              placeholder="081234567890"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
            />
          </div>
        </div>

        <div v-if="isEmailChanged" class="p-4 bg-amber-50 border border-amber-200 rounded-xl space-y-2">
          <label class="block text-xs font-bold text-amber-900">
            Konfirmasi Kata Sandi Saat Ini (Wajib untuk mengubah email)
          </label>
          <input
            v-model="form.current_password"
            type="password"
            required
            maxlength="72"
            placeholder="Masukkan kata sandi akun Anda saat ini"
            class="w-full px-3.5 py-2.5 text-xs bg-white border border-amber-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-amber-400 focus:border-amber-600 transition-all"
          />
          <p class="text-[11px] text-amber-800">
            Demi keamanan akun, perubahan alamat email memerlukan verifikasi kata sandi saat ini.
          </p>
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">Alamat Rumah / Pengiriman Utama</label>
          <textarea
            v-model="form.address"
            rows="3"
            maxlength="500"
            placeholder="Jl. Nama Jalan No. XX, Kelurahan, Kecamatan, Kota, Kode Pos"
            class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
          ></textarea>
          <p class="text-[11px] text-slate-400 mt-1">
            Alamat ini otomatis diisi saat Anda melakukan checkout pesanan buku.
          </p>
        </div>

        <div class="pt-3 flex flex-wrap items-center justify-between gap-3 border-t border-slate-100">
          <NuxtLink to="/forgot-password" class="text-xs font-semibold text-indigo-600 hover:underline">
            Ingin ubah kata sandi? Reset Password &rarr;
          </NuxtLink>

          <div class="flex items-center gap-2.5">
            <NuxtLink
              to="/"
              class="px-5 py-2.5 bg-white hover:bg-slate-50 text-slate-700 font-bold text-xs rounded-xl border border-slate-300 transition-colors"
            >
              Beranda
            </NuxtLink>
            <button
              :disabled="isSaving"
              type="submit"
              class="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-xl shadow-md shadow-indigo-600/20 flex items-center gap-2 transition-all disabled:opacity-50"
            >
              <span v-if="isSaving">Menyimpan...</span>
              <span v-else>Simpan Perubahan</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()

const isLoadingProfile = ref(true)
const isSaving = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const originalEmail = ref('')

const form = reactive({
  name: '',
  email: '',
  phone: '',
  address: '',
  current_password: '',
})

const isEmailChanged = computed(() => {
  return form.email.trim().toLowerCase() !== originalEmail.value.trim().toLowerCase()
})

const syncFormFromUser = () => {
  if (!authStore.user) return
  form.name = authStore.user.name || ''
  form.email = authStore.user.email || ''
  form.phone = authStore.user.phone || ''
  form.address = authStore.user.address || ''
  form.current_password = ''
  originalEmail.value = authStore.user.email || ''
}

const handleSave = async () => {
  errorMessage.value = ''
  successMessage.value = ''
  isSaving.value = true
  try {
    const payload: {
      name: string
      email: string
      phone?: string
      address?: string
      current_password?: string
    } = {
      name: form.name.trim(),
      email: form.email.trim(),
      phone: form.phone.trim(),
      address: form.address.trim(),
    }
    if (isEmailChanged.value) {
      payload.current_password = form.current_password
    }

    const res = await authStore.updateProfile(payload)
    syncFormFromUser()
    successMessage.value = res.message || 'Profil berhasil diperbarui.'
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal menyimpan perubahan profil.'
  } finally {
    isSaving.value = false
  }
}

onMounted(async () => {
  await authStore.fetchProfile()
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  syncFormFromUser()
  isLoadingProfile.value = false
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 sm:py-16">
    <div class="flex justify-center">
      <div class="w-full max-w-md bg-white p-8 rounded-2xl border border-gray-200 shadow-sm space-y-6">
        <div class="space-y-1">
          <h1 class="text-2xl font-extrabold text-gray-900">Atur Ulang Kata Sandi</h1>
          <p class="text-xs text-gray-500">
            Buat kata sandi baru minimal 12 karakter dengan kombinasi huruf besar, huruf kecil, angka, dan simbol.
          </p>
        </div>

        <div v-if="errorMessage" class="p-3 bg-red-50 text-red-700 border border-red-200 text-xs rounded-xl font-medium flex items-center gap-2">
          <i class="ti ti-alert-triangle text-base"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <div v-if="successMessage" class="p-3 bg-emerald-50 text-emerald-700 border border-emerald-200 text-xs rounded-xl font-medium space-y-2">
          <div class="flex items-center gap-2">
            <i class="ti ti-check text-base"></i>
            <span>{{ successMessage }}</span>
          </div>
          <div>
            <NuxtLink to="/login" class="font-bold text-indigo-600 hover:underline">
              Lanjut ke Halaman Login
            </NuxtLink>
          </div>
        </div>

        <form v-if="!successMessage" @submit.prevent="handleReset" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1.5">Kata Sandi Baru</label>
            <div class="relative">
              <input
                v-model="newPassword"
                :type="showPassword ? 'text' : 'password'"
                required
                minlength="12"
                maxlength="72"
                placeholder="Minimal 12 karakter kuat"
                class="w-full pl-9 pr-9 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
              />
              <i class="ti ti-lock text-gray-400 absolute left-3 top-3 text-sm"></i>
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-2.5 text-gray-400 hover:text-gray-600"
              >
                <i :class="showPassword ? 'ti ti-eye-off' : 'ti ti-eye'" class="text-sm"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1.5">Konfirmasi Kata Sandi Baru</label>
            <div class="relative">
              <input
                v-model="confirmPassword"
                :type="showPassword ? 'text' : 'password'"
                required
                minlength="12"
                maxlength="72"
                placeholder="Ulangi kata sandi baru"
                class="w-full pl-9 pr-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
              />
              <i class="ti ti-lock-check text-gray-400 absolute left-3 top-3 text-sm"></i>
            </div>
          </div>

          <button
            :disabled="isLoading || !token"
            type="submit"
            class="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-indigo-600/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50"
          >
            <i v-if="isLoading" class="ti ti-loader-2 animate-spin text-base"></i>
            <span v-else>Simpan Kata Sandi Baru</span>
          </button>
        </form>

        <div class="text-center text-xs text-gray-500 pt-3 border-t border-gray-100">
          <NuxtLink to="/login" class="font-bold text-indigo-600 hover:underline">
            Kembali ke Sign In
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()

const token = computed(() => String(route.query.token || '').trim())
const newPassword = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const isLoading = ref(false)
const errorMessage = ref(token.value ? '' : 'Tautan reset password tidak valid atau token tidak ditemukan.')
const successMessage = ref('')

const handleReset = async () => {
  errorMessage.value = ''
  if (!token.value) {
    errorMessage.value = 'Token reset password tidak ditemukan.'
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = 'Konfirmasi kata sandi tidak cocok.'
    return
  }

  isLoading.value = true
  try {
    const res = await api.post('/auth/reset-password', {
      token: token.value,
      new_password: newPassword.value,
    })
    successMessage.value = res.message || 'Password berhasil diperbarui. Silakan login kembali.'
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal mengatur ulang kata sandi.'
  } finally {
    isLoading.value = false
  }
}
</script>

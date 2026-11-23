<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 sm:py-16">
    <div class="flex justify-center">
      <div class="w-full max-w-md bg-white p-8 rounded-2xl border border-gray-200 shadow-sm space-y-6">
        <div class="space-y-1">
          <h1 class="text-2xl font-extrabold text-gray-900">Lupa Kata Sandi</h1>
          <p class="text-xs text-gray-500">
            Masukkan alamat email terdaftar Anda. Kami akan mengirimkan tautan untuk mengatur ulang kata sandi.
          </p>
        </div>

        <div v-if="errorMessage" class="p-3 bg-red-50 text-red-700 border border-red-200 text-xs rounded-xl font-medium flex items-center gap-2">
          <i class="ti ti-alert-triangle text-base"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <div v-if="successMessage" class="p-3 bg-emerald-50 text-emerald-700 border border-emerald-200 text-xs rounded-xl font-medium flex items-center gap-2">
          <i class="ti ti-check text-base"></i>
          <span>{{ successMessage }}</span>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1.5">Alamat Email</label>
            <div class="relative">
              <input
                v-model="email"
                type="email"
                required
                placeholder="nama@email.com"
                class="w-full pl-9 pr-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-indigo-600 transition-all"
              />
              <i class="ti ti-mail text-gray-400 absolute left-3 top-3 text-sm"></i>
            </div>
          </div>

          <button
            :disabled="isLoading"
            type="submit"
            class="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-indigo-600/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50"
          >
            <i v-if="isLoading" class="ti ti-loader-2 animate-spin text-base"></i>
            <span v-else>Kirim Tautan Reset</span>
          </button>
        </form>

        <div class="text-center text-xs text-gray-500 pt-3 border-t border-gray-100">
          Ingat kata sandi Anda?
          <NuxtLink to="/login" class="font-bold text-indigo-600 hover:underline ml-1">
            Kembali ke Sign In
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const api = useApi()

const email = ref('')
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleSubmit = async () => {
  isLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const res = await api.post('/auth/forgot-password', { email: email.value })
    successMessage.value = res.message || 'Jika email terdaftar, tautan reset password telah dikirim.'
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal memproses permintaan reset password'
  } finally {
    isLoading.value = false
  }
}
</script>

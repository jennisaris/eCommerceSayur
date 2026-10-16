<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 sm:py-16">
    <div class="flex flex-wrap justify-center items-center gap-8 lg:gap-16">
      <!-- Illustration (FreshCart signup-g.svg) -->
      <div class="w-full md:w-5/12 lg:w-1/2 flex justify-center order-2 lg:order-1">
        <img
          src="/assets/images/svg-graphics/signup-g.svg"
          alt="Sign up to FreshCart"
          class="max-w-sm sm:max-w-md w-full h-auto drop-shadow-sm"
        />
      </div>

      <!-- Sign Up Form -->
      <div class="w-full md:w-7/12 lg:w-5/12 max-w-md order-1 lg:order-2 bg-white p-8 rounded-2xl border border-gray-200 shadow-sm space-y-5">
        <div class="space-y-1">
          <h1 class="text-2xl font-extrabold text-gray-900">Get Started with FreshCart</h1>
          <p class="text-xs text-gray-500">
            Daftar akun sekarang dan nikmati belanja kebutuhan sayur segar 24/7 diantar langsung ke rumah.
          </p>
        </div>

        <!-- Error Alert -->
        <div v-if="errorMessage" class="p-3 bg-red-50 text-red-700 border border-red-200 text-xs rounded-xl font-medium flex items-center gap-2">
          <i class="ti ti-alert-triangle text-base"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <form @submit.prevent="handleRegister" class="space-y-3.5">
          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1">Nama Lengkap</label>
            <div class="relative">
              <input
                v-model="form.name"
                type="text"
                required
                placeholder="Nama Lengkap Anda"
                class="w-full pl-9 pr-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-400 focus:border-green-600 transition-all"
              />
              <i class="ti ti-user text-gray-400 absolute left-3 top-3 text-sm"></i>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1">Alamat Email</label>
            <div class="relative">
              <input
                v-model="form.email"
                type="email"
                required
                placeholder="nama@email.com"
                class="w-full pl-9 pr-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-400 focus:border-green-600 transition-all"
              />
              <i class="ti ti-mail text-gray-400 absolute left-3 top-3 text-sm"></i>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1">Kata Sandi</label>
            <div class="relative">
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                minlength="6"
                placeholder="Minimal 6 karakter"
                class="w-full pl-9 pr-9 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-400 focus:border-green-600 transition-all"
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
            <label class="block text-xs font-bold text-gray-700 mb-1">Nomor Telepon / WhatsApp</label>
            <div class="relative">
              <input
                v-model="form.phone"
                type="tel"
                placeholder="081234567890"
                class="w-full pl-9 pr-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-400 focus:border-green-600 transition-all"
              />
              <i class="ti ti-phone text-gray-400 absolute left-3 top-3 text-sm"></i>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1">Alamat Lengkap Pengiriman</label>
            <textarea
              v-model="form.address"
              rows="2"
              placeholder="Alamat rumah atau kantor untuk 100% Home Delivery..."
              class="w-full px-3 py-2 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-400 focus:border-green-600 transition-all"
            ></textarea>
          </div>

          <button
            :disabled="isLoading"
            type="submit"
            class="w-full py-3 px-4 bg-green-600 hover:bg-green-700 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-green-600/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50 mt-2"
          >
            <i v-if="isLoading" class="ti ti-loader-2 animate-spin text-base"></i>
            <span v-else>Daftar Sekarang</span>
          </button>
        </form>

        <div class="text-center text-xs text-gray-500 pt-3 border-t border-gray-100">
          Sudah memiliki akun?
          <NuxtLink to="/login" class="font-bold text-green-600 hover:underline ml-1">
            Sign In di Sini
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()

const form = reactive({
  name: '',
  email: '',
  password: '',
  phone: '',
  address: '',
})

const showPassword = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')

const handleRegister = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    await authStore.register(form)
    router.push('/')
  } catch (err: any) {
    errorMessage.value = err.message || 'Pendaftaran gagal. Pastikan email belum terdaftar.'
  } finally {
    isLoading.value = false
  }
}
</script>

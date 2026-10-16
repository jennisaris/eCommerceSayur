<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 sm:py-16">
    <div class="flex flex-wrap justify-center items-center gap-8 lg:gap-16">
      <!-- Illustration (FreshCart signin-g.svg) -->
      <div class="w-full md:w-5/12 lg:w-1/2 flex justify-center order-2 lg:order-1">
        <img
          src="/assets/images/svg-graphics/signin-g.svg"
          alt="Sign in to FreshCart"
          class="max-w-sm sm:max-w-md w-full h-auto drop-shadow-sm"
        />
      </div>

      <!-- Sign In Form -->
      <div class="w-full md:w-7/12 lg:w-5/12 max-w-md order-1 lg:order-2 bg-white p-8 rounded-2xl border border-gray-200 shadow-sm space-y-6">
        <div class="space-y-1">
          <h1 class="text-2xl font-extrabold text-gray-900">Sign In to FreshCart</h1>
          <p class="text-xs text-gray-500">
            Selamat datang kembali! Masuk untuk memesan sayuran segar dan cek status pesanan.
          </p>
        </div>

        <!-- Quick Demo Login Buttons -->
        <div class="p-3 bg-gray-50 rounded-xl border border-gray-200 space-y-2">
          <span class="text-[10px] font-bold text-gray-500 uppercase tracking-wider block text-center">Akun Uji Coba Cepat (Demo Fill)</span>
          <div class="grid grid-cols-2 gap-2">
            <button
              @click="fillAdmin"
              type="button"
              class="px-2.5 py-1.5 text-xs font-bold bg-amber-50 text-amber-700 border border-amber-200 hover:bg-amber-100 rounded-lg transition-colors flex items-center justify-center gap-1"
            >
              <i class="ti ti-shield-lock text-xs"></i>
              <span>Super Admin</span>
            </button>
            <button
              @click="fillCustomer"
              type="button"
              class="px-2.5 py-1.5 text-xs font-bold bg-green-50 text-green-700 border border-green-200 hover:bg-green-100 rounded-lg transition-colors flex items-center justify-center gap-1"
            >
              <i class="ti ti-user text-xs"></i>
              <span>Customer</span>
            </button>
          </div>
        </div>

        <!-- Error Alert -->
        <div v-if="errorMessage" class="p-3 bg-red-50 text-red-700 border border-red-200 text-xs rounded-xl font-medium flex items-center gap-2">
          <i class="ti ti-alert-triangle text-base"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-gray-700 mb-1.5">Alamat Email</label>
            <div class="relative">
              <input
                v-model="email"
                type="email"
                required
                placeholder="nama@email.com"
                class="w-full pl-9 pr-3.5 py-2.5 text-xs bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-400 focus:border-green-600 transition-all"
              />
              <i class="ti ti-mail text-gray-400 absolute left-3 top-3 text-sm"></i>
            </div>
          </div>

          <div>
            <div class="flex items-center justify-between mb-1.5">
              <label class="block text-xs font-bold text-gray-700">Kata Sandi</label>
              <span class="text-[11px] text-green-600 hover:underline cursor-pointer">Lupa Password?</span>
            </div>
            <div class="relative">
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                required
                placeholder="••••••••"
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

          <button
            :disabled="isLoading"
            type="submit"
            class="w-full py-3 px-4 bg-green-600 hover:bg-green-700 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-green-600/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50"
          >
            <i v-if="isLoading" class="ti ti-loader-2 animate-spin text-base"></i>
            <span v-else>Sign In</span>
          </button>
        </form>

        <div class="text-center text-xs text-gray-500 pt-3 border-t border-gray-100">
          Belum memiliki akun?
          <NuxtLink to="/register" class="font-bold text-green-600 hover:underline ml-1">
            Sign Up Sekarang
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')

const fillAdmin = () => {
  email.value = 'admin@ecommerce.com'
  password.value = 'Admin123!'
  errorMessage.value = ''
}

const fillCustomer = () => {
  email.value = 'budi@customer.com'
  password.value = 'Admin123!'
  errorMessage.value = ''
}

const handleLogin = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    await authStore.login(email.value, password.value)
    if (authStore.isSuperAdmin) {
      router.push('/admin')
    } else {
      router.push('/')
    }
  } catch (err: any) {
    errorMessage.value = err.message || 'Email atau kata sandi tidak valid'
  } finally {
    isLoading.value = false
  }
}
</script>

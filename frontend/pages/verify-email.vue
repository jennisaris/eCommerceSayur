<template>
  <main class="min-h-screen flex items-center justify-center bg-slate-50 px-4">
    <section class="w-full max-w-md rounded-2xl bg-white p-8 text-center shadow-sm">
      <h1 class="text-xl font-bold text-slate-900">Verifikasi Email</h1>
      <p class="mt-3 text-sm text-slate-600">{{ message }}</p>
      <NuxtLink to="/login" class="mt-6 inline-block rounded-lg bg-indigo-600 px-5 py-3 text-sm font-semibold text-white">
        Ke halaman login
      </NuxtLink>
    </section>
  </main>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()
const message = ref('Memverifikasi email...')

onMounted(async () => {
  const token = String(route.query.token || '')
  if (!token) {
    message.value = 'Token verifikasi tidak ditemukan.'
    return
  }

  try {
    await api.post('/auth/verify-email', { token })
    message.value = 'Email berhasil diverifikasi. Silakan login.'
  } catch {
    message.value = 'Token verifikasi tidak valid atau sudah kedaluwarsa.'
  }
})
</script>

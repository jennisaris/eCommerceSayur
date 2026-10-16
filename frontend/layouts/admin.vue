<template>
  <div class="min-h-screen flex bg-slate-100">
    <!-- Desktop Sidebar -->
    <AdminSidebar class="hidden md:flex" />

    <!-- Main Admin Content Area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-y-auto">
      <header class="h-16 bg-white border-b border-slate-200 px-6 flex items-center justify-between">
        <h1 class="text-sm font-bold text-slate-800">Panel Manajemen E-Commerce</h1>
        <div class="flex items-center gap-3">
          <span class="text-xs font-semibold text-slate-500">Super Admin Mode</span>
          <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 ring-4 ring-emerald-100"></div>
        </div>
      </header>

      <main class="p-6 md:p-8">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  if (!authStore.isSuperAdmin) {
    alert('Akses khusus Super Admin!')
    router.push('/')
  }
})
</script>

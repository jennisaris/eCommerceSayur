import { defineStore } from 'pinia'
import type { CartSummary } from '~/types'

export const useCartStore = defineStore('cart', () => {
  const summary = ref<CartSummary>({
    cart: { id: 0, user_id: 0, items: [] },
    total_items: 0,
    total_weight: 0,
    total_price: 0,
  })
  const isDrawerOpen = ref(false)
  const isLoading = ref(false)

  const items = computed(() => summary.value.cart.items || [])
  const totalItems = computed(() => summary.value.total_items)
  const totalPrice = computed(() => summary.value.total_price)
  const totalWeight = computed(() => summary.value.total_weight)

  const openDrawer = () => (isDrawerOpen.value = true)
  const closeDrawer = () => (isDrawerOpen.value = false)
  const toggleDrawer = () => (isDrawerOpen.value = !isDrawerOpen.value)

  const fetchCart = async () => {
    const authStore = useAuthStore()
    if (!authStore.isAuthenticated) {
      summary.value = {
        cart: { id: 0, user_id: 0, items: [] },
        total_items: 0,
        total_weight: 0,
        total_price: 0,
      }
      return
    }

    isLoading.value = true
    try {
      const api = useApi()
      const res = await api.get<CartSummary>('/cart')
      if (res.data) {
        summary.value = res.data
      }
    } catch (err) {
      console.error('Failed to fetch cart:', err)
    } finally {
      isLoading.value = false
    }
  }

  const addToCart = async (productId: number, quantity: number = 1) => {
    const authStore = useAuthStore()
    if (!authStore.isAuthenticated) {
      const router = useRouter()
      router.push('/login')
      return false
    }

    isLoading.value = true
    try {
      const api = useApi()
      const res = await api.post<CartSummary>('/cart/items', {
        product_id: productId,
        quantity,
      })
      if (res.data) {
        summary.value = res.data
        openDrawer()
      }
      return true
    } catch (err: any) {
      alert(err.message || 'Gagal menambahkan produk ke keranjang')
      return false
    } finally {
      isLoading.value = false
    }
  }

  const updateQuantity = async (productId: number, quantity: number) => {
    if (quantity <= 0) {
      return removeItem(productId)
    }

    isLoading.value = true
    try {
      const api = useApi()
      const res = await api.put<CartSummary>(`/cart/items/${productId}`, {
        quantity,
      })
      if (res.data) {
        summary.value = res.data
      }
    } catch (err: any) {
      alert(err.message || 'Gagal mengubah jumlah barang')
    } finally {
      isLoading.value = false
    }
  }

  const removeItem = async (productId: number) => {
    isLoading.value = true
    try {
      const api = useApi()
      const res = await api.delete<CartSummary>(`/cart/items/${productId}`)
      if (res.data) {
        summary.value = res.data
      }
    } catch (err: any) {
      alert(err.message || 'Gagal menghapus barang dari keranjang')
    } finally {
      isLoading.value = false
    }
  }

  const clearCart = async () => {
    isLoading.value = true
    try {
      const api = useApi()
      await api.delete('/cart')
      summary.value = {
        cart: { id: 0, user_id: 0, items: [] },
        total_items: 0,
        total_weight: 0,
        total_price: 0,
      }
    } catch (err: any) {
      console.error('Failed to clear cart:', err)
    } finally {
      isLoading.value = false
    }
  }

  return {
    summary,
    items,
    totalItems,
    totalPrice,
    totalWeight,
    isDrawerOpen,
    isLoading,
    openDrawer,
    closeDrawer,
    toggleDrawer,
    fetchCart,
    addToCart,
    updateQuantity,
    removeItem,
    clearCart,
  }
})

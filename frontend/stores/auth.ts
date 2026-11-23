import { defineStore } from 'pinia'
import type { User } from '~/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const initialized = ref(false)
  const isAuthenticated = computed(() => !!user.value)
  const isSuperAdmin = computed(() => {
    return user.value?.roles?.some((r) => r.name === 'Super Admin') ?? false
  })

  const setUser = (newUser: User) => {
    user.value = newUser
  }

  const login = async (email: string, password: string) => {
    const api = useApi()
    const res = await api.post<{ user: User }>('/auth/login', {
      email,
      password,
    })
    setUser(res.data.user)
    initialized.value = true
    return res.data
  }

  const register = async (payload: {
    name: string
    email: string
    password: string
    phone?: string
    address?: string
  }) => {
    const api = useApi()
    const res = await api.post<{ user: User; verification_required?: boolean }>('/auth/register', payload)
    if (!res.data.verification_required) {
      setUser(res.data.user)
      initialized.value = true
    }
    return res.data
  }

  const fetchProfile = async () => {
    if (initialized.value) return user.value
    initialized.value = true
    const api = useApi()
    try {
      const res = await api.get<User>('/auth/me')
      setUser(res.data)
      return res.data
    } catch {
      user.value = null
      return null
    }
  }

  const updateProfile = async (payload: {
    name: string
    email: string
    phone?: string
    address?: string
    current_password?: string
  }) => {
    const api = useApi()
    const res = await api.put<User>('/auth/me', payload)
    setUser(res.data)
    return res
  }

  const clearSession = () => {
    user.value = null
    initialized.value = true
  }

  const logout = async (redirect = true) => {
    const api = useApi()
    try {
      await api.post('/auth/logout')
    } catch {
      // Clear client state even if the server session is already unavailable.
    }
    clearSession()
    if (!redirect) return
    const router = useRouter()
    router.push('/login')
  }

  return {
    user,
    isAuthenticated,
    isSuperAdmin,
    login,
    register,
    fetchProfile,
    updateProfile,
    clearSession,
    logout,
  }
})

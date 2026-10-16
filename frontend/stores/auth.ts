import { defineStore } from 'pinia'
import type { User } from '~/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => !!token.value)
  const isSuperAdmin = computed(() => {
    return user.value?.roles?.some((r) => r.name === 'Super Admin') ?? false
  })

  // Initialize from localStorage in client
  if (process.client) {
    const savedToken = localStorage.getItem('auth_token')
    if (savedToken) {
      token.value = savedToken
      const savedUser = localStorage.getItem('auth_user')
      if (savedUser) {
        try {
          user.value = JSON.parse(savedUser)
        } catch {}
      }
    }
  }

  const setAuth = (newToken: string, newUser: User) => {
    token.value = newToken
    user.value = newUser
    if (process.client) {
      localStorage.setItem('auth_token', newToken)
      localStorage.setItem('auth_user', JSON.stringify(newUser))
    }
  }

  const login = async (email: string, password: string) => {
    const api = useApi()
    const res = await api.post<{ token: string; user: User }>('/auth/login', {
      email,
      password,
    })
    setAuth(res.data.token, res.data.user)
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
    const res = await api.post<{ token: string; user: User }>('/auth/register', payload)
    setAuth(res.data.token, res.data.user)
    return res.data
  }

  const fetchProfile = async () => {
    if (!token.value) return null
    const api = useApi()
    try {
      const res = await api.get<User>('/auth/me')
      user.value = res.data
      if (process.client) {
        localStorage.setItem('auth_user', JSON.stringify(res.data))
      }
      return res.data
    } catch {
      logout()
      return null
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    if (process.client) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
    }
    const router = useRouter()
    router.push('/login')
  }

  return {
    token,
    user,
    isAuthenticated,
    isSuperAdmin,
    login,
    register,
    fetchProfile,
    logout,
  }
})

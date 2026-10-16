import type { ApiResponse } from '~/types'

export const useApi = () => {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const request = async <T>(
    endpoint: string,
    options: {
      method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
      body?: any
      params?: Record<string, any>
      headers?: Record<string, string>
    } = {}
  ): Promise<ApiResponse<T>> => {
    const url = `${config.public.apiBase}${endpoint}`

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      ...(options.headers || {}),
    }

    if (authStore.token) {
      headers['Authorization'] = `Bearer ${authStore.token}`
    }

    try {
      const response = await $fetch<ApiResponse<T>>(url, {
        method: options.method || 'GET',
        body: options.body,
        params: options.params,
        headers,
      })
      return response
    } catch (err: any) {
      const errorMsg =
        err?.data?.message || err?.message || 'Terjadi kesalahan pada server'
      throw new Error(errorMsg)
    }
  }

  return {
    get: <T>(endpoint: string, params?: Record<string, any>) =>
      request<T>(endpoint, { method: 'GET', params }),
    post: <T>(endpoint: string, body?: any) =>
      request<T>(endpoint, { method: 'POST', body }),
    put: <T>(endpoint: string, body?: any) =>
      request<T>(endpoint, { method: 'PUT', body }),
    delete: <T>(endpoint: string) =>
      request<T>(endpoint, { method: 'DELETE' }),
  }
}

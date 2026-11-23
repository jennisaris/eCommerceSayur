import type { ApiResponse } from '~/types'

const resolveApiBase = (configuredBase: string): string => {
  if (process.client && typeof window !== 'undefined') {
    const pageHost = window.location.hostname
    if (pageHost === '127.0.0.1' && configuredBase.startsWith('http://localhost:')) {
      return configuredBase.replace('http://localhost:', 'http://127.0.0.1:')
    }
    if (pageHost === 'localhost' && configuredBase.startsWith('http://127.0.0.1:')) {
      return configuredBase.replace('http://127.0.0.1:', 'http://localhost:')
    }
  }
  return configuredBase
}

export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = resolveApiBase(String(config.public.apiBase || ''))
  const authStore = useAuthStore()
  let csrfToken: string | null = null

  const unsafeMethods = new Set(['POST', 'PUT', 'PATCH', 'DELETE'])

  const ensureCsrfToken = async () => {
    if (csrfToken) return csrfToken
    const response = await $fetch<ApiResponse<{ token: string }>>(
      `${apiBase}/auth/csrf`,
      { method: 'GET', credentials: 'include' },
    )
    csrfToken = response.data.token
    return csrfToken
  }

  const request = async <T>(
    endpoint: string,
    options: {
      method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
      body?: any
      params?: Record<string, any>
      headers?: Record<string, string>
    } = {},
    retryAfterRefresh = false,
  ): Promise<ApiResponse<T>> => {
    const url = `${apiBase}${endpoint}`

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      ...(options.headers || {}),
    }

    const method = options.method || 'GET'
    if (unsafeMethods.has(method)) {
      headers['X-CSRF-Token'] = await ensureCsrfToken()
    }

    try {
      const response = await $fetch<ApiResponse<T>>(url, {
        method,
        body: options.body,
        params: options.params,
        headers,
        credentials: 'include',
      })
      return response
    } catch (err: any) {
      const statusCode = err?.response?.status || err?.statusCode
      const canRefresh = !retryAfterRefresh &&
        statusCode === 401 &&
        !endpoint.endsWith('/auth/login') &&
        !endpoint.endsWith('/auth/register') &&
        !endpoint.endsWith('/auth/refresh') &&
        !endpoint.endsWith('/auth/logout')

      if (canRefresh) {
        try {
          await request('/auth/refresh', { method: 'POST' }, true)
          return request<T>(endpoint, options, true)
        } catch {
          authStore.clearSession()
        }
      }

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

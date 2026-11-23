import { isAllowedMidtransUrl } from '~/utils/safeUrl'

declare global {
  interface Window {
    snap?: {
      pay: (
        token: string,
        options: {
          onSuccess?: (result: any) => void
          onPending?: (result: any) => void
          onError?: (result: any) => void
          onClose?: () => void
        }
      ) => void
    }
  }
}

export default defineNuxtPlugin(() => {
  if (process.client) {
    const config = useRuntimeConfig()
    const clientKey = String(config.public.midtransClientKey || '').trim()
    const snapUrl = String(config.public.midtransSnapUrl || '').trim()
    if (!clientKey || !isAllowedMidtransUrl(snapUrl)) {
      return
    }
    const script = document.createElement('script')
    script.src = snapUrl
    script.setAttribute('data-client-key', clientKey)
    script.async = true
    document.head.appendChild(script)
  }
})

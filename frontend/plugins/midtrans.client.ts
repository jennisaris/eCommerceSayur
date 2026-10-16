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
    const script = document.createElement('script')
    script.src = config.public.midtransSnapUrl as string
    script.setAttribute('data-client-key', config.public.midtransClientKey as string)
    script.async = true
    document.head.appendChild(script)
  }
})

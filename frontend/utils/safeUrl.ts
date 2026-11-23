const ALLOWED_MIDTRANS_HOSTS = new Set([
  'app.sandbox.midtrans.com',
  'app.midtrans.com',
])

export function isAllowedMidtransUrl(rawUrl: string | null | undefined): boolean {
  const trimmed = String(rawUrl || '').trim()
  if (!trimmed) return false
  try {
    const parsed = new URL(trimmed)
    return parsed.protocol === 'https:' && ALLOWED_MIDTRANS_HOSTS.has(parsed.hostname.toLowerCase())
  } catch {
    return false
  }
}

export function toSafeAssetUrl(rawUrl: string | null | undefined): string {
  const trimmed = String(rawUrl || '').trim()
  if (!trimmed) return ''
  if (/[\s\u0000-\u001F\u007F]/.test(trimmed)) return ''
  if (trimmed.startsWith('/')) {
    return trimmed.startsWith('//') ? '' : trimmed
  }
  try {
    const parsed = new URL(trimmed)
    if (parsed.protocol === 'https:' || parsed.protocol === 'http:') {
      return parsed.toString()
    }
  } catch {
    return ''
  }
  return ''
}

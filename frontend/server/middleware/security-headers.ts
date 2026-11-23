import { defineEventHandler, setResponseHeaders } from 'h3'

const PRODUCTION_CSP = [
  "default-src 'self'",
  "base-uri 'self'",
  "frame-ancestors 'self'",
  "object-src 'none'",
  "script-src 'self' 'unsafe-inline' https://app.sandbox.midtrans.com https://app.midtrans.com",
  "connect-src 'self' https://app.sandbox.midtrans.com https://app.midtrans.com https://api.sandbox.midtrans.com https://api.midtrans.com",
  "img-src 'self' data: blob: https:",
  "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdn.jsdelivr.net",
  "font-src 'self' data: https://fonts.gstatic.com",
  "frame-src 'self' https://app.sandbox.midtrans.com https://app.midtrans.com",
].join('; ')

export default defineEventHandler((event) => {
  const headers: Record<string, string> = {
    'X-Frame-Options': 'SAMEORIGIN',
    'X-Content-Type-Options': 'nosniff',
    'Referrer-Policy': 'strict-origin-when-cross-origin',
    'Permissions-Policy': 'camera=(), geolocation=(), microphone=()',
  }

  if (process.env.NODE_ENV === 'production') {
    headers['Strict-Transport-Security'] = 'max-age=31536000; includeSubDomains'
    headers['Content-Security-Policy'] = PRODUCTION_CSP
  }

  setResponseHeaders(event, headers)
})

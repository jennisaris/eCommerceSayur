export interface Role {
  id: number
  name: string
}

export interface User {
  id: number
  name: string
  email: string
  phone?: string
  address?: string
  photo?: string
  lat?: string
  lng?: string
  is_verified: boolean
  roles?: Role[]
  created_at: string
}

export interface Category {
  id: number
  parent_id?: number | null
  name: string
  slug: string
  icon?: string
  status: boolean
  description?: string
  children?: Category[]
}

export interface Product {
  id: number
  category_id: number
  category?: Category
  parent_id?: number | null
  name: string
  slug: string
  image: string
  description: string
  regular_price: number
  sale_price: number
  unit: string
  weight: number
  stock: number
  variant_name?: string
  status: 'ACTIVE' | 'INACTIVE' | 'OUT_OF_STOCK'
  created_at: string
}

export interface CartItem {
  id: number
  cart_id: number
  product_id: number
  product?: Product
  quantity: number
}

export interface Cart {
  id: number
  user_id: number
  items: CartItem[]
}

export interface CartSummary {
  cart: Cart
  total_items: number
  total_weight: number
  total_price: number
}

export type OrderStatus =
  | 'PENDING_PAYMENT'
  | 'PAID'
  | 'PROCESSING'
  | 'ON_DELIVERY'
  | 'COMPLETED'
  | 'CANCELLED'

export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  product?: Product
  product_name: string
  quantity: number
  price: number
  subtotal: number
}

export interface Payment {
  id: number
  order_id: number
  payment_method?: string
  payment_status: string
  payment_gateway_id?: string
  gross_amount: number
  payment_url?: string
  snap_token?: string
}

export interface Order {
  id: number
  order_code: string
  buyer_id: number
  buyer?: User
  status: OrderStatus
  recipient_name: string
  recipient_phone: string
  shipping_address: string
  shipping_lat?: string
  shipping_lng?: string
  shipping_fee: number
  tracking_number?: string
  total_amount: number
  remarks?: string
  items?: OrderItem[]
  payment?: Payment
  created_at: string
}

export interface ApiResponse<T = any> {
  success: boolean
  message: string
  data: T
  meta?: {
    current_page: number
    per_page: number
    total_items: number
    total_pages: number
  }
}

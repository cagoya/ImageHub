// User related types
export interface User {
  id: number
  username: string
  email: string
  parse_exif: boolean
  is_admin: boolean
  is_suspended: boolean
  suspended_until?: string | null
  created_time: string
  updated_time: string
}

export interface UserListResponse {
  users: User[]
  total: number
}

// Image related types
export interface Image {
  image_id: number
  image_name: string
  image_path: string
  ai_tags: string[]
  manual_tags: string[]
}

export interface ImageItem {
  image_id: number
  image_name: string
  image_path: string
  ai_tags: string[]
  manual_tags: string[]
}

export interface ImageDetail {
  id: number
  name: string
  path: string
  captured_time?: string
  latitude?: number
  longitude?: number
  formatted_address?: string
  province?: string
  city?: string
  district?: string
  width_pixel?: number
  height_pixel?: number
  user_id: number
  created_time: string
  updated_time: string
  tags?: string[]
}

// Pagination related types
export interface PaginationState {
  currentPage: number
  pageSize: number
  total: number
}

// API response types
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// Form types
export interface LoginForm {
  name: string
  password: string
}

export interface RegisterForm {
  name: string
  email: string
  password: string
}

export interface ForgetPasswordForm {
  username: string
  email: string
  newPassword: string
}

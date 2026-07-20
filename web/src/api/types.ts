export interface Document {
  id: string
  original_name: string
  file_size: number
  content_type: string
  status: 'pending' | 'ocr-processing' | 'done' | 'failed' | string
  ocr_result?: string
  thumbnail_key?: string
  created_at: string
}

export interface LoginResponse {
  message: string
  token: string
}

export interface RegisterResponse {
  message: string
  user_id: string
}

export interface ApiError {
  error?: string
  response?: string
}


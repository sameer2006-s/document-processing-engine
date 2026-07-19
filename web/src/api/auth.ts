import { apiFetch } from './client'
import type { LoginResponse, RegisterResponse } from './types'

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  email: string
  password: string
  firstName: string
  lastName: string
}

export function login(input: LoginInput) {
  return apiFetch<LoginResponse>('/login', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function register(input: RegisterInput) {
  return apiFetch<RegisterResponse>('/register', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

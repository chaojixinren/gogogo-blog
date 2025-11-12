import api from './api'
import type { User } from '@/types'

export interface AuthResponse {
  token: string
  user: User
}

export const login = async (payload: {
  username: string
  password: string
}): Promise<AuthResponse> => {
  const { data } = await api.post<AuthResponse>('/auth/login', payload)
  return data
}

export const register = async (payload: {
  username: string
  password: string
  email?: string
  displayName?: string
}): Promise<AuthResponse> => {
  const { data } = await api.post<AuthResponse>('/auth/register', payload)
  return data
}

export const fetchProfile = async (): Promise<User> => {
  const { data } = await api.get<{ user: User }>('/me')
  return data.user
}

import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { User } from '@/types'
import * as authService from '@/services/auth'

const TOKEN_STORAGE_KEY = 'auth_token'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isInitialized = ref(false)
  const isProcessing = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => Boolean(user.value && token.value))

  const setSession = (authToken: string | null, profile: User | null) => {
    token.value = authToken
    user.value = profile
    if (authToken) {
      localStorage.setItem(TOKEN_STORAGE_KEY, authToken)
    } else {
      localStorage.removeItem(TOKEN_STORAGE_KEY)
    }
  }

  const initialize = async () => {
    if (isInitialized.value) return

    const storedToken = localStorage.getItem(TOKEN_STORAGE_KEY)
    if (!storedToken) {
      isInitialized.value = true
      return
    }

    token.value = storedToken
    try {
      const profile = await authService.fetchProfile()
      user.value = profile
    } catch (err) {
      console.error('Failed to restore session', err)
      setSession(null, null)
    } finally {
      isInitialized.value = true
    }
  }

  const login = async (payload: {
    username: string
    password: string
  }) => {
    isProcessing.value = true
    error.value = null
    try {
      const response = await authService.login(payload)
      setSession(response.token, response.user)
      return response.user
    } catch (err) {
      error.value = 'Login failed'
      throw err
    } finally {
      isProcessing.value = false
    }
  }

  const register = async (payload: {
    username: string
    password: string
    email?: string
    displayName?: string
  }) => {
    isProcessing.value = true
    error.value = null
    try {
      const response = await authService.register(payload)
      setSession(response.token, response.user)
      return response.user
    } catch (err) {
      error.value = 'Register failed'
      throw err
    } finally {
      isProcessing.value = false
    }
  }

  const logout = async () => {
    setSession(null, null)
  }

  const refreshProfile = async () => {
    if (!token.value) return
    try {
      const profile = await authService.fetchProfile()
      user.value = profile
    } catch (err) {
      console.error('Failed to refresh profile', err)
      setSession(null, null)
    }
  }

  return {
    user,
    token,
    isProcessing,
    isInitialized,
    error,
    isAuthenticated,
    initialize,
    login,
    register,
    logout,
    refreshProfile,
    setSession,
  }
})

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || null)
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthed = computed(() => !!token.value && !!user.value)
  
  // Ролевые геттеры
  const role = computed(() => user.value?.role || '')
  const isManager = computed(() => user.value?.role === 'manager')
  const isEngineer = computed(() => user.value?.role === 'engineer')
  const isLead = computed(() => user.value?.role === 'lead')
  const isViewer = computed(() => user.value?.role === 'viewer')

  const setAuth = (newToken, newUser) => {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const clearAuth = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  const register = async (credentials) => {
    try {
      const response = await authApi.register(credentials)
      setAuth(response.token, response.user)
      return response
    } catch (error) {
      throw error
    }
  }

  const login = async (credentials) => {
    try {
      const response = await authApi.login(credentials)
      setAuth(response.token, response.user)
      return response
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    clearAuth()
  }

  const checkAuth = async () => {
    if (!token.value) return false
    
    try {
      const response = await authApi.me()
      user.value = response
      return true
    } catch (error) {
      clearAuth()
      return false
    }
  }

  return {
    token,
    user,
    isAuthed,
    role,
    isManager,
    isEngineer,
    isLead,
    isViewer,
    register,
    login,
    logout,
    checkAuth
  }
})

import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import router from '../router'

const baseURL = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

const http = axios.create({
  baseURL,
  timeout: 10000
})

// Request interceptor
http.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
http.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default http

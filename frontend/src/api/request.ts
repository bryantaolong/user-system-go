import axios from 'axios'
import { Message } from '@arco-design/web-vue'
import router from '@/router'

const instance = axios.create({
  timeout: 10000,
})

// Request interceptor
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor
instance.interceptors.response.use(
  (response) => {
    // Pass blob through for file downloads
    if (response.config.responseType === 'blob') {
      return response
    }
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      Message.error('认证信息失效，请重新登录')
      router.push('/login')
    }
    return Promise.reject(error.response?.data || error)
  }
)

export default instance

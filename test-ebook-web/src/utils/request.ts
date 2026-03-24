import axios from 'axios'
import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    // 如果 code 不是 200，则报错
    if (res.code !== 200) {
      ElMessage.error(res.message || 'Error')
      
      // 401: 未登录或 Token 过期
      if (res.code === 401) {
        localStorage.removeItem('token')
        sessionStorage.removeItem('token')
        window.location.href = '/login'
      }
      return Promise.reject(new Error(res.message || 'Error'))
    }
    return res.data
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      ElMessage.error('登录状态已失效，请重新登录')
      localStorage.removeItem('token')
      sessionStorage.removeItem('token')
      window.location.href = '/login'
    } else {
      ElMessage.error(error.message || 'Network Error')
    }
    return Promise.reject(error)
  }
)

export default service

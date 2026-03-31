import axios from 'axios'
import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import Cookies from 'js-cookie'

declare module 'axios' {
  export interface AxiosRequestConfig {
    silent?: boolean
  }
}

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = Cookies.get('token')
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
      if (!response.config.silent) {
        ElMessage.error(res.msg || 'Error')
      }
      
      // 401: 未登录或 Token 过期
      if (res.code === 401) {
        Cookies.remove('token')
        window.location.href = '/login'
      }
      return Promise.reject(new Error(res.msg || 'Error'))
    }
    return res.data
  },
  (error) => {
    const isSilent = error.config?.silent
    if (error.response && error.response.status === 401) {
      if (!isSilent) {
        ElMessage.error('登录状态已失效，请重新登录')
      }
      Cookies.remove('token')
      window.location.href = '/login'
    } else {
      if (!isSilent) {
        // 优先显示后端返回的业务错误信息
        const backendMsg = error.response?.data?.msg || error.response?.data?.message
        ElMessage.error(backendMsg || error.message || 'Network Error')
      }
    }
    return Promise.reject(error)
  }
)

export default service

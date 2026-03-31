import { defineStore } from 'pinia'
import { login, getInfo } from '@/api/auth'
import { ref } from 'vue'
import Cookies from 'js-cookie'
export interface User {
  id: number
  username: string
  role: string
  theme: string
  is_active: boolean
  permissions: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(Cookies.get('token') || '')
  const user = ref<User | null>(null)
  const isLoggedIn = ref(!!token.value)

  async function handleLogin(form: any, remember: boolean) {
    const res: any = await login(form)
    token.value = res.token
    user.value = res.user
    isLoggedIn.value = true
    
    if (remember) {
      Cookies.set('token', res.token, { expires: 0.5 })
    } else {
      Cookies.set('token', res.token)
    }
  }

  async function fetchUser() {
    if (!token.value) return
    const res: any = await getInfo()
    user.value = res
  }

  function logout() {
    token.value = ''
    user.value = null
    isLoggedIn.value = false
    Cookies.remove('token')
  }

  return {
    token,
    user,
    isLoggedIn,
    handleLogin,
    fetchUser,
    logout,
  }
})

import { defineStore } from 'pinia'
import { login, getInfo } from '@/api/auth'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || sessionStorage.getItem('token') || '')
  const user = ref<any>(null)
  const isLoggedIn = ref(!!token.value)

  async function handleLogin(form: any, remember: boolean) {
    const res: any = await login(form)
    token.value = res.token
    user.value = res.user
    isLoggedIn.value = true
    
    const storage = remember ? localStorage : sessionStorage
    storage.setItem('token', res.token)
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
    localStorage.removeItem('token')
    sessionStorage.removeItem('token')
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

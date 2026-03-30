import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCategories } from '@/api/category'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<any[]>([])
  const loading = ref(false)
  const isLoaded = ref(false)

  const fetchCategories = async (force = false) => {
    if (isLoaded.value && !force) return
    loading.value = true
    try {
      const res: any = await getCategories()
      categories.value = res
      isLoaded.value = true
    } catch (error) {
      console.error('Failed to fetch categories:', error)
    } finally {
      loading.value = false
    }
  }

  return {
    categories,
    loading,
    isLoaded,
    fetchCategories
  }
})

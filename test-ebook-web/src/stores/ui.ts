import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useResponsive } from '@/composables/useResponsive'

export const useUiStore = defineStore('ui', () => {
  const { isMobile, isTablet } = useResponsive()
  
  const isSidebarCollapsed = ref(false)
  const isDrawerVisible = ref(false)
  
  // 自动处理侧边栏折叠：平板端默认折叠
  watch(isTablet, (val) => {
    if (val) isSidebarCollapsed.value = true
    else if (!isMobile.value) isSidebarCollapsed.value = false
  }, { immediate: true })
  
  // 切换折叠状态
  const toggleSidebar = () => {
    isSidebarCollapsed.value = !isSidebarCollapsed.value
  }
  
  // 切换抽屉显示
  const toggleDrawer = () => {
    isDrawerVisible.value = !isDrawerVisible.value
  }
  
  const closeDrawer = () => {
    isDrawerVisible.value = false
  }
  
  return {
    isSidebarCollapsed,
    isDrawerVisible,
    toggleSidebar,
    toggleDrawer,
    closeDrawer
  }
})

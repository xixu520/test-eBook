import { useBreakpoints } from '@vueuse/core'
import { computed } from 'vue'

const breakpointsElementPlus = {
  xs: 768,
  sm: 768,
  md: 992,
  lg: 1200,
  xl: 1920,
}

export function useResponsive() {
  const breakpoints = useBreakpoints(breakpointsElementPlus)
  
  const isMobile = breakpoints.smaller('sm') // < 768px
  const isTablet = breakpoints.between('sm', 'md') // 768px - 992px
  const isDesktop = breakpoints.greaterOrEqual('md') // >= 992px
  
  const deviceType = computed(() => {
    if (isMobile.value) return 'mobile'
    if (isTablet.value) return 'tablet'
    return 'desktop'
  })
  
  return {
    isMobile,
    isTablet,
    isDesktop,
    deviceType
  }
}

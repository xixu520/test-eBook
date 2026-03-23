<template>
  <div class="main-layout">
    <TopNavBar />
    <div class="main-container">
      <!-- 桌面端侧边栏 -->
      <SideMenu v-if="!isMobile" />
      
      <!-- 移动端抽屉侧边栏 -->
      <el-drawer
        v-if="isMobile"
        v-model="ui.isDrawerVisible"
        direction="ltr"
        size="280px"
        :with-header="false"
        class="mobile-drawer"
      >
        <SideMenu />
      </el-drawer>
      
      <div class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import TopNavBar from './TopNavBar.vue'
import SideMenu from './SideMenu.vue'
import { useResponsive } from '@/composables/useResponsive'
import { useUiStore } from '@/stores/ui'

const { isMobile } = useResponsive()
const ui = useUiStore()
</script>

<style scoped lang="scss">
.main-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  
  .main-container {
    flex: 1;
    display: flex;
    overflow: hidden;
    position: relative;
    
    .main-content {
      flex: 1;
      padding: 0;
      overflow-y: auto;
      background-color: #f5f7fa;
      transition: all 0.3s;
    }
  }
}

:deep(.mobile-drawer) {
  .el-drawer__body {
    padding: 0;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

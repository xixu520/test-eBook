<template>
  <div class="main-layout">
    <TopNavBar />
    <div class="main-container">
      <SideMenu />
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
    
    .main-content {
      flex: 1;
      padding: 0;
      overflow-y: auto;
      background-color: #f5f7fa;
    }
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

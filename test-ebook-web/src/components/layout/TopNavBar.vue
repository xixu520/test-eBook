<template>
  <div class="top-nav-bar">
    <!-- 公告横幅 -->
    <el-alert
      v-if="announcement && showAnnouncement"
      :title="announcement.content"
      type="info"
      center
      closable
      @close="handleCloseAnnouncement"
      class="announcement-bar"
    />
    
    <div class="header-content">
      <div class="left-section">
        <img src="@/assets/logo.png" alt="Logo" class="logo" />
        <span class="site-name">{{ siteName }}</span>
      </div>
      
      <div class="center-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索标准号或名称 (Ctrl+K)"
          class="search-input"
          :prefix-icon="Search"
          @keyup.enter="handleSearch"
          ref="searchInputRef"
        />
      </div>
      
      <div class="right-section">
        <el-switch
          v-model="isDark"
          class="theme-switch"
          inline-prompt
          :active-icon="Moon"
          :inactive-icon="Sunny"
          @change="toggleTheme"
        />
        
        <el-dropdown trigger="click">
          <div class="user-info">
            <el-avatar :size="32" icon="UserFilled" />
            <span class="username">{{ auth.user?.username }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>个人中心</el-dropdown-item>
              <el-dropdown-item v-if="['admin', 'editor'].includes(auth.user?.role || '')" @click="goToAdmin">管理后台</el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Search, Sunny, Moon, UserFilled } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { getActiveAnnouncement } from '@/api/document'

const auth = useAuthStore()
const router = useRouter()

const siteName = ref('建筑标准文件管理系统')
const searchKeyword = ref('')
const showAnnouncement = ref(!sessionStorage.getItem('announcement-closed'))
const announcement = ref<any>(null)
const isDark = ref(false)
const searchInputRef = ref()

onMounted(async () => {
  try {
    const res: any = await getActiveAnnouncement()
    announcement.value = res
  } catch (error) {
    console.error(error)
  }
  
  window.addEventListener('keydown', handleGlobalKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeyDown)
})

const handleGlobalKeyDown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchInputRef.value?.focus()
  }
}

const handleCloseAnnouncement = () => {
  sessionStorage.setItem('announcement-closed', 'true')
}

const handleSearch = () => {
  router.push({ path: '/', query: { keyword: searchKeyword.value } })
}

const toggleTheme = (val: any) => {
  document.documentElement.setAttribute('data-theme', val ? 'dark' : 'light')
  // 这里可以调用后端 API 保存主题偏好
}

const goToAdmin = () => {
  router.push('/admin')
}

const handleLogout = () => {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped lang="scss">
.top-nav-bar {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  position: sticky;
  top: 0;
  z-index: 1000;
  
  .announcement-bar {
    border-radius: 0;
  }
  
  .header-content {
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    
    .left-section {
      display: flex;
      align-items: center;
      gap: 10px;
      
      .logo {
        height: 32px;
      }
      
      .site-name {
        font-size: 18px;
        font-weight: bold;
        color: #303133;
      }
    }
    
    .center-section {
      flex: 1;
      max-width: 500px;
      margin: 0 20px;
      
      .search-input {
        :deep(.el-input__wrapper) {
          border-radius: 20px;
        }
      }
    }
    
    .right-section {
      display: flex;
      align-items: center;
      gap: 20px;
      
      .user-info {
        display: flex;
        align-items: center;
        gap: 8px;
        cursor: pointer;
        padding: 4px 8px;
        border-radius: 4px;
        transition: background-color 0.2s;
        
        &:hover {
          background-color: #f5f7fa;
        }
        
        .username {
          font-size: 14px;
          color: #606266;
        }
      }
    }
  }
}
</style>

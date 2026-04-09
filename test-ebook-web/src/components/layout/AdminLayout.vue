<template>
  <div class="admin-layout">
    <div class="admin-sidebar" :class="{ collapsed: isCollapsed }">
      <div class="sidebar-header">
        <img src="@/assets/logo.png" alt="Logo" class="logo" v-if="!isCollapsed" />
        <span class="title" v-if="!isCollapsed">管理后台</span>
        <el-icon class="toggle" @click="isCollapsed = !isCollapsed">
          <Fold v-if="!isCollapsed" /><Expand v-else />
        </el-icon>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        class="admin-menu"
        :collapse="isCollapsed"
        router
      >
        <el-menu-item index="/admin">
          <el-icon><DataBoard /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>
        
        <el-menu-item index="/admin/categories" v-if="hasAccess(['admin', 'editor'])">
          <el-icon><Folder /></el-icon>
          <template #title>分类管理</template>
        </el-menu-item>

        <el-menu-item index="/admin/field-config" v-if="hasAccess(['admin', 'editor'])">
          <el-icon><Setting /></el-icon>
          <template #title>属性管理</template>
        </el-menu-item>
        
        <el-menu-item index="/admin/ocr" v-if="hasAccess(['admin', 'editor'])">
          <el-icon><Cpu /></el-icon>
          <template #title>OCR 任务</template>
        </el-menu-item>
        
        <el-menu-item index="/admin/documents">
          <el-icon><Document /></el-icon>
          <template #title>文档管理</template>
        </el-menu-item>
        
        <el-menu-item index="/admin/recycle" v-if="hasAccess(['admin'])">
          <el-icon><Delete /></el-icon>
          <template #title>回收站</template>
        </el-menu-item>
        
        <el-menu-item index="/admin/audit" v-if="hasAccess(['admin'])">
          <el-icon><Memo /></el-icon>
          <template #title>审计日志</template>
        </el-menu-item>
        
        <el-menu-item index="/admin/users" v-if="hasAccess(['admin'])">
          <el-icon><User /></el-icon>
          <template #title>用户管理</template>
        </el-menu-item>

        <el-menu-item index="/admin/settings" v-if="hasAccess(['admin'])">
          <el-icon><Setting /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>
        
        <div class="menu-footer">
          <el-menu-item @click="router.push('/')">
            <el-icon><Back /></el-icon>
            <template #title>返回主页</template>
          </el-menu-item>
        </div>
      </el-menu>
    </div>
    
    <div class="admin-main">
      <div class="admin-header">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>管理后台</el-breadcrumb-item>
          <el-breadcrumb-item>{{ currentRouteTitle }}</el-breadcrumb-item>
        </el-breadcrumb>
        
        <div class="user-actions">
          <el-dropdown trigger="click">
            <div class="avatar-wrap">
              <el-avatar :size="32"><UserFilled /></el-avatar>
              <span class="name">{{ auth.user?.username }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
      <div class="admin-content">
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
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  DataBoard, Folder, Document, User, Back, Fold, Expand, UserFilled, Cpu, Setting, Delete, Memo 
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const isCollapsed = ref(false)
const activeMenu = computed(() => route.path)
const currentRouteTitle = computed(() => route.meta.title || '')

const hasAccess = (roles: string[]) => {
  return auth.user && roles.includes(auth.user.role)
}

const handleLogout = () => {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped lang="scss">
.admin-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  
  .admin-sidebar {
    width: 220px;
    background-color: #001529;
    color: #fff;
    transition: width 0.3s;
    display: flex;
    flex-direction: column;
    
    &.collapsed {
      width: 64px;
    }
    
    .sidebar-header {
      height: 60px;
      display: flex;
      align-items: center;
      padding: 0 16px;
      gap: 12px;
      overflow: hidden;
      white-space: nowrap;
      
      .logo {
        height: 32px;
      }
      
      .title {
        font-weight: bold;
        font-size: 16px;
      }
      
      .toggle {
        margin-left: auto;
        cursor: pointer;
        font-size: 18px;
        &:hover { color: #409eff; }
      }
    }
    
    .el-menu {
      border-right: none;
      background-color: transparent;
      flex: 1;
      
      :deep(.el-menu-item) {
        color: rgba(255, 255, 255, 0.65);
        &:hover {
          color: #fff;
          background-color: #1890ff;
        }
        &.is-active {
          background-color: #1890ff;
          color: #fff;
        }
      }
    }
    
    .menu-footer {
      border-top: 1px solid rgba(255, 255, 255, 0.1);
      margin-top: auto;
    }
  }
  
  .admin-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    background-color: #f0f2f5;
    overflow: hidden;
    
    .admin-header {
      height: 60px;
      background-color: #fff;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 24px;
      box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
      z-index: 10;
      
      .user-actions {
        .avatar-wrap {
          display: flex;
          align-items: center;
          gap: 8px;
          cursor: pointer;
          .name { font-size: 14px; color: #666; }
        }
      }
    }
    
    .admin-content {
      flex: 1;
      padding: 24px;
      overflow-y: auto;
    }
  }
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>

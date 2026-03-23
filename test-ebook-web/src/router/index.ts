import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import MainLayout from '@/components/layout/MainLayout.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/LoginPage.vue'),
      meta: { title: '登录' }
    },
    {
      path: '/',
      component: MainLayout,
      children: [
        {
          path: '',
          name: 'Home',
          component: () => import('@/views/home/HomePage.vue'),
          meta: { requiresAuth: true, title: '首页' }
        }
      ]
    },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true, roles: ['admin', 'editor'] },
      children: [
        {
          path: '',
          name: 'AdminDashboard',
          component: () => import('@/views/admin/dashboard/DashboardPage.vue'),
          meta: { title: '仪表盘' }
        },
        {
          path: 'categories',
          name: 'AdminCategories',
          component: () => import('@/views/admin/category/CategoryPage.vue'),
          meta: { title: '分类管理' }
        },
        {
          path: 'ocr',
          name: 'AdminOcr',
          component: () => import('@/views/admin/ocr/OcrTaskPage.vue'),
          meta: { title: 'OCR 任务' }
        },
        {
          path: 'users',
          name: 'AdminUsers',
          component: () => import('@/views/admin/user/UserListPage.vue'),
          meta: { title: '用户管理', requiresAdmin: true }
        },
        {
          path: 'documents',
          name: 'AdminDocuments',
          component: () => import('@/views/admin/document/AdminDocumentsPage.vue'),
          meta: { title: '文档管理' }
        },
        {
          path: 'recycle',
          name: 'AdminRecycle',
          component: () => import('@/views/admin/recycle/RecycleBinPage.vue'),
          meta: { title: '回收站' }
        },
        {
          path: 'audit',
          name: 'AdminAudit',
          component: () => import('@/views/admin/audit/AuditLogPage.vue'),
          meta: { title: '审计日志' }
        },
        {
          path: 'settings',
          name: 'AdminSettings',
          component: () => import('@/views/admin/settings/SettingsPage.vue'),
          meta: { title: '系统设置', requiresAdmin: true }
        }
      ]
    }
  ]
})

router.beforeEach(async (to, _from, next) => {
  const auth = useAuthStore()
  
  // 处理需要登录的路由
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!auth.token) {
      return next('/login')
    }
    
    // 如果有 token 但无用户信息，尝试拉取
    if (!auth.user) {
      try {
        await auth.fetchUser()
      } catch (error) {
        auth.logout()
        return next('/login')
      }
    }

    // 处理角色访问控制
    const roles = to.matched.find(record => record.meta.roles)?.meta.roles as string[] | undefined
    if (roles && !roles.includes(auth.user?.role || '')) {
      return next('/')
    }

    // 处理管理员权限 (针对明确标记为 requiresAdmin 的路由)
    if (to.matched.some(record => record.meta.requiresAdmin) && auth.user?.role !== 'admin') {
      return next('/')
    }
  }
  
  next()
})

export default router

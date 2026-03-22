<template>
  <div class="login-container">
    <el-card class="login-card" shadow="always">
      <div class="login-header">
        <img src="@/assets/logo.png" alt="Logo" v-if="hasLogo" />
        <h2 v-else>电子书管理系统</h2>
      </div>
      
      <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top">
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="loginForm.username" 
            placeholder="请输入用户名" 
            :prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="请输入密码" 
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>
        
        <div class="login-options">
          <el-checkbox v-model="rememberMe">记住密码</el-checkbox>
          <el-link type="primary" :underline="false">忘记密码？</el-link>
        </div>
        
        <el-button 
          type="primary" 
          class="login-button" 
          :loading="loading" 
          @click="handleLogin"
        >
          登录
        </el-button>
        
        <div class="register-link">
          还没有账号？<el-link type="primary">立即注册</el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const auth = useAuthStore()
const router = useRouter()

const hasLogo = ref(true)
const loading = ref(false)
const rememberMe = ref(false)
const loginFormRef = ref()

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, message: '密码不能少于6位', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await auth.handleLogin(loginForm, rememberMe.value)
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        console.error(error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  
  .login-card {
    width: 420px;
    border-radius: 12px;
    
    .login-header {
      text-align: center;
      margin-bottom: 30px;
      
      h2 {
        color: #303133;
        margin: 0;
      }
      
      img {
        width: 120px;
      }
    }
    
    .login-options {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
    }
    
    .login-button {
      width: 100%;
      height: 40px;
      font-size: 16px;
    }
    
    .register-link {
      text-align: center;
      margin-top: 20px;
      font-size: 14px;
      color: #606266;
    }
  }
}
</style>

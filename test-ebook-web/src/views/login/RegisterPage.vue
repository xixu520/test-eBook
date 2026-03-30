<template>
  <div class="register-container">
    <el-card class="register-card" shadow="always">
      <div class="register-header">
        <img src="@/assets/logo.png" alt="Logo" v-if="hasLogo" />
        <h2 v-else>用户注册</h2>
      </div>
      
      <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" label-position="top" @submit.prevent="handleRegister">
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="registerForm.username" 
            placeholder="请输入用户名" 
            :prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="registerForm.password" 
            type="password" 
            placeholder="请输入密码" 
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="registerForm.confirmPassword" 
            type="password" 
            placeholder="请再次输入密码" 
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>
        
        <el-button 
          type="primary" 
          class="register-button" 
          :loading="loading" 
          native-type="submit"
        >
          立即注册
        </el-button>
        
        <div class="login-link">
          已有账号？<el-link type="primary" @click="router.push('/login')">立即登录</el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { User, Lock } from '@element-plus/icons-vue'
import { register } from '@/api/auth'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const hasLogo = ref(true)
const loading = ref(false)
const registerFormRef = ref()

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const validatePass2 = (_rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validatePass2, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await register({
          username: registerForm.username,
          password: registerForm.password
        })
        ElMessage.success('注册成功，请登录')
        router.push('/login')
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
.register-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  
  .register-card {
    width: 100%;
    max-width: 420px;
    margin: 0 20px;
    border-radius: 12px;
    
    .register-header {
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
    
    .register-button {
      width: 100%;
      height: 40px;
      font-size: 16px;
      margin-top: 10px;
    }
    
    .login-link {
      text-align: center;
      margin-top: 20px;
      font-size: 14px;
      color: #606266;
    }
  }
}
</style>

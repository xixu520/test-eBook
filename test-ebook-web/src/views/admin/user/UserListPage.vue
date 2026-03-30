<template>
  <div class="user-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" :icon="Plus" @click="handleAdd">新增用户</el-button>
        </div>
      </template>

      <el-table :data="userList" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column label="角色">
          <template #default="{ row }">
            <el-tag :type="roleTagType(row.role)">{{ roleName(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_active"
              @change="(val: boolean) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间">
          <template #default="{ row }">
             {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="handleResetPwd(row)">重置密码</el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="row.role !== 'admin'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 用户编辑/新增弹窗 -->
    <el-dialog
      v-model="formVisible"
      :title="isEdit ? '编辑用户' : '新增用户'"
      width="450px"
      destroy-on-close
    >
      <el-form :model="userForm" :rules="userRules" ref="userFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" :disabled="isEdit" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="userForm.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="标注员" value="editor" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="submitUserForm" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 密码重置弹窗 -->
    <el-dialog
      v-model="pwdVisible"
      title="重置密码"
      width="400px"
      destroy-on-close
    >
      <el-form :model="pwdForm" :rules="pwdRules" ref="pwdFormRef" label-width="100px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="pwdForm.password" type="password" show-password placeholder="请输入新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pwdVisible = false">取消</el-button>
        <el-button type="primary" @click="submitResetPwd" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { getUsers, updateUserStatus, deleteUser, createUser, updateUser, resetPassword } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const userList = ref([])

// 表单控制
const formVisible = ref(false)
const isEdit = ref(false)
const userFormRef = ref()
const userForm = reactive({
  id: 0,
  username: '',
  password: '',
  role: 'user'
})

const userRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

// 密码重置控制
const pwdVisible = ref(false)
const pwdFormRef = ref()
const pwdForm = reactive({
  id: 0,
  password: ''
})

const pwdRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ]
}

const roleTagType = (role: string) => {
  switch (role) {
    case 'admin': return 'danger'
    case 'editor': return 'warning'
    default: return 'info'
  }
}

const roleName = (role: string) => {
  switch (role) {
    case 'admin': return '管理员'
    case 'editor': return '标注员'
    default: return '普通用户'
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getUsers({ page: 1, page_size: 100 })
    userList.value = res.list || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  userForm.username = ''
  userForm.password = ''
  userForm.role = 'user'
  formVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  userForm.id = row.id
  userForm.username = row.username
  userForm.role = row.role
  formVisible.value = true
}

const handleResetPwd = (row: any) => {
  pwdForm.id = row.id
  pwdForm.password = ''
  pwdVisible.value = true
}

const submitUserForm = async () => {
  if (!userFormRef.value) return
  await userFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          await updateUser(userForm.id, { role: userForm.role })
          ElMessage.success('用户权限已更新')
        } else {
          await createUser(userForm)
          ElMessage.success('用户创建成功')
        }
        formVisible.value = false
        loadData()
      } catch (e) {
      } finally {
        submitting.value = false
      }
    }
  })
}

const submitResetPwd = async () => {
  if (!pwdFormRef.value) return
  await pwdFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitting.value = true
      try {
        await resetPassword(pwdForm.id, { password: pwdForm.password })
        ElMessage.success('密码重置成功')
        pwdVisible.value = false
      } catch (e) {
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleStatusChange = async (row: any, val: boolean) => {
  try {
    await updateUserStatus(row.id, val)
    ElMessage.success(`用户 ${row.username} 已${val ? '启用' : '禁用'}`)
  } catch (error) {
    row.is_active = !val // 失败回滚
    console.error(error)
  }
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除用户 ${row.username} 吗？`, '警告', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteUser(row.id)
      ElMessage.success('用户已删除')
      loadData()
    } catch (error) {
      console.error(error)
    }
  })
}

const formatDate = (date: any) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(loadData)
</script>

<style scoped lang="scss">
.user-list-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
  }
}
</style>


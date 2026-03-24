<template>
  <div class="user-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" :icon="Plus">新增用户</el-button>
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
        <el-table-column prop="last_login" label="最近登录" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEditRole(row)">修改角色</el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="row.role !== 'admin'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 修改角色对话框 -->
    <el-dialog v-model="roleDialogVisible" title="修改用户角色" width="400px">
      <el-form :model="roleForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="roleForm.username" disabled />
        </el-form-item>
        <el-form-item label="核心角色">
          <el-select v-model="roleForm.role" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="标注员 (编辑)" value="editor" />
            <el-option label="普通用户 (只读)" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRoleChange" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { getUsers, updateUserStatus, deleteUser, updateUserRole } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const userList = ref([])
const roleDialogVisible = ref(false)
const submitting = ref(false)

const roleForm = reactive({
  id: 0,
  username: '',
  role: ''
})

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

const handleStatusChange = async (row: any, val: boolean) => {
  try {
    await updateUserStatus(row.id, val)
    ElMessage.success(`用户 ${row.username} 已${val ? '启用' : '禁用'}`)
  } catch (error) {
    row.is_active = !val // 失败回滚
    console.error(error)
  }
}

const handleEditRole = (row: any) => {
  roleForm.id = row.id
  roleForm.username = row.username
  roleForm.role = row.role
  roleDialogVisible.value = true
}

const submitRoleChange = async () => {
  submitting.value = true
  try {
    await updateUserRole(roleForm.id, roleForm.role)
    ElMessage.success('角色更新成功')
    roleDialogVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
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

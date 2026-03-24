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
            <el-button link type="danger" @click="handleDelete(row)" v-if="row.role !== 'admin'">删除</el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="row.role !== 'admin'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { getUsers, updateUserStatus, deleteUser } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const userList = ref([])

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

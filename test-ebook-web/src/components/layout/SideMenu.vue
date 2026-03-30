<template>
  <div class="side-menu" :class="{ collapsed: ui.isSidebarCollapsed, 'is-mobile': isMobile }">
    <div v-if="!isMobile" class="collapse-trigger" @click="ui.toggleSidebar">
      <el-icon><Fold v-if="!ui.isSidebarCollapsed" /><Expand v-else /></el-icon>
    </div>
    
    <el-menu
      :default-active="activeCategory"
      class="el-menu-vertical"
      :collapse="!isMobile && ui.isSidebarCollapsed"
      @select="handleSelect"
    >
      <el-menu-item index="all">
        <el-icon><Files /></el-icon>
        <template #title>全部文件</template>
      </el-menu-item>
      
      <el-sub-menu v-for="cat in categoryTree" :key="cat.ID" :index="cat.ID.toString()">
        <template #title>
          <el-icon><Folder /></el-icon>
          <span>{{ cat.name }}</span>
        </template>
        
        <el-menu-item 
          v-for="sub in cat.children" 
          :key="sub.ID" 
          :index="sub.ID.toString()"
        >
          {{ sub.name }}
          <span class="count">{{ sub.doc_count }}</span>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { Folder, Files, Fold, Expand } from '@element-plus/icons-vue'
import { useCategoryStore } from '@/stores/category'
import { useRouter, useRoute } from 'vue-router'
import { useUiStore } from '@/stores/ui'
import { useResponsive } from '@/composables/useResponsive'

const router = useRouter()
const route = useRoute()
const ui = useUiStore()
const { isMobile } = useResponsive()

const categoryStore = useCategoryStore()
const activeCategory = computed(() => (route.query.category_id as string) || 'all')

onMounted(() => {
  categoryStore.fetchCategories()
})

const categoryTree = computed(() => {
  return categoryStore.categories
})

const handleSelect = (index: string) => {
  if (index === 'all') {
    router.push({ path: '/', query: {} })
  } else {
    router.push({ path: '/', query: { category_id: index } })
  }
  
  if (isMobile.value) {
    ui.closeDrawer()
  }
}

// 导出供模板使用
defineExpose({
  categoryTree,
  activeCategory,
  handleSelect
})
</script>

<style scoped lang="scss">
.side-menu {
  width: 240px;
  height: calc(100vh - 56px);
  border-right: 1px solid #dcdfe6;
  background-color: #fff;
  transition: width 0.3s;
  display: flex;
  flex-direction: column;
  
  &.collapsed {
    width: 64px;
  }
  
  &.is-mobile {
    width: 100%;
    height: 100%;
  }
  
  .collapse-trigger {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: #909399;
    border-bottom: 1px solid #f2f6fc;
    
    &:hover {
      background-color: #f5f7fa;
      color: #303133;
    }
  }
  
  .el-menu-vertical {
    border-right: none;
    flex: 1;
    overflow-y: auto;
  }
  
  .count {
    font-size: 12px;
    color: #909399;
    margin-left: auto;
    background-color: #f0f2f5;
    padding: 2px 6px;
    border-radius: 10px;
  }
}
</style>

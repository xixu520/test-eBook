<template>
  <el-config-provider>
    <router-view />
  </el-config-provider>
</template>

<script setup lang="ts">
import { ElConfigProvider } from 'element-plus'
</script>

<style>
body {
  margin: 0;
  padding: 0;
}
#app {
  height: 100vh;
}
</style>

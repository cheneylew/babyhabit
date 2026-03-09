<template>
  <div id="app">
    <router-view v-if="$route.meta.requiresAuth && !isAuthenticated">
      <el-result
        icon="error"
        title="请先登录"
        sub-title="该页面需要登录才能访问"
      >
        <template #extra>
          <el-button type="primary" @click="$router.push('/login')">去登录</el-button>
        </template>
      </el-result>
    </router-view>
    <router-view v-else />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useUserStore } from './store/user'

const userStore = useUserStore()
const isAuthenticated = computed(() => userStore.isAuthenticated)
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Arial, sans-serif;
  background-color: #f5f7fa;
}

#app {
  min-height: 100vh;
}
</style>
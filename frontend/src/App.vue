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
  max-width: 100%;
}

html, body {
  font-family: Arial, sans-serif;
  background-color: #f5f7fa;
  overflow-x: hidden;
  width: 100%;
  position: relative;
}

#app {
  min-height: 100vh;
  width: 100%;
  overflow-x: hidden;
}

/* 防止图片和其他元素超出容器 */
img, video, iframe {
  max-width: 100%;
  height: auto;
}

/* 防止长文本导致横向滚动 */
p, h1, h2, h3, h4, h5, h6 {
  word-wrap: break-word;
  overflow-wrap: break-word;
}
</style>
<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">
          <h2>BabyHabit 登录</h2>
        </div>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item>
          <div class="form-actions">
            <el-button type="primary" @click="login" :loading="isLoading">登录</el-button>
            <el-button @click="$router.push('/register')">注册</el-button>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useUserStore } from '../store/user'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const router = useRouter()
const loginFormRef = ref(null)
const isLoading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const login = async () => {
  if (!loginFormRef.value.validate()) return
  
  isLoading.value = true
  try {
    const result = await userStore.login(loginForm.username, loginForm.password)
    ElMessage.success('登录成功')
    // 根据用户类型跳转页面
    if (result.user.user_type === 1) {
      router.push('/admin')
    } else {
      router.push('/home')
    }
  } catch (error) {
    ElMessage.error(userStore.error || '登录失败')
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  border-radius: 10px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-card {
    width: 90%;
    max-width: 400px;
  }
  
  .el-form {
    padding: 0 15px;
  }
  
  .el-button {
    width: 90px;
  }
}

.login-header {
  text-align: center;
  color: #333;
}

.login-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.el-button {
  width: 100px;
  margin-right: 10px;
  display: inline-block;
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}
</style>
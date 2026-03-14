<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <div class="register-header">
          <h2>BabyHabit 注册</h2>
        </div>
      </template>
      <el-form :model="registerForm" :rules="rules" ref="registerFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="registerForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="registerForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="registerForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="用户类型" prop="user_type">
          <el-radio-group v-model="registerForm.user_type">
            <el-radio label="1">父母</el-radio>
            <el-radio label="2">小孩</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="registerForm.user_type === 2" label="父母ID" prop="parent_id">
          <el-input v-model="registerForm.parent_id" placeholder="请输入父母账号ID" />
        </el-form-item>
        <el-form-item>
          <div class="form-actions">
            <el-button type="primary" @click="register" :loading="isLoading">注册</el-button>
            <el-button @click="$router.push('/login')">登录</el-button>
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
const registerFormRef = ref(null)
const isLoading = ref(false)

const registerForm = reactive({
  username: '',
  password: '',
  name: '',
  phone: '',
  email: '',
  user_type: 1,
  parent_id: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, message: '密码长度至少6位', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  user_type: [{ required: true, message: '请选择用户类型', trigger: 'change' }],
  parent_id: [{ required: registerForm.user_type === '2', message: '请输入父母账号ID', trigger: 'blur' }]
}

const register = async () => {
  if (!registerFormRef.value.validate()) return
  
  isLoading.value = true
  try {
    await userStore.register({
      ...registerForm,
      user_type: parseInt(registerForm.user_type),
      parent_id: registerForm.user_type === '2' ? parseInt(registerForm.parent_id) : 0
    })
    ElMessage.success('注册成功')
    router.push('/home')
  } catch (error) {
    ElMessage.error(userStore.error || '注册失败')
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.register-card {
  width: 500px;
  border-radius: 10px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-card {
    width: 90%;
    max-width: 500px;
  }
  
  .el-form {
    padding: 0 15px;
  }
  
  .el-button {
    width: 90px;
  }
}

.register-header {
  text-align: center;
  color: #333;
}

.register-header h2 {
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
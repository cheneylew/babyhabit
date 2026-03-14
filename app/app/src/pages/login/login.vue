<template>
  <view class="login-container">
    <view class="login-card">
      <view class="login-header">
        <image src="../../static/logo.png" class="logo" mode="aspectFit"></image>
        <h2>BabyHabit</h2>
        <p>让好习惯成为孩子的终身财富</p>
      </view>
      <form @submit.prevent="login">
        <view class="form-item">
          <view class="input-wrapper">
            <text class="input-icon">👤</text>
            <input 
              type="text" 
              v-model="loginForm.username" 
              placeholder="请输入用户名" 
              class="input"
              :class="{ 'input-focused': isUsernameFocused }"
              @focus="isUsernameFocused = true"
              @blur="isUsernameFocused = false"
            />
          </view>
        </view>
        <view class="form-item">
          <view class="input-wrapper">
            <text class="input-icon">🔒</text>
            <input 
              type="password" 
              v-model="loginForm.password" 
              placeholder="请输入密码" 
              class="input"
              :class="{ 'input-focused': isPasswordFocused }"
              @focus="isPasswordFocused = true"
              @blur="isPasswordFocused = false"
            />
            <text 
              class="password-toggle" 
              @click="togglePasswordVisibility"
            >
              {{ showPassword ? '👁️' : '👁️‍🗨️' }}
            </text>
          </view>
        </view>
        <view class="form-actions">
          <button 
            type="primary" 
            @click="login" 
            :loading="isLoading"
            class="login-button"
          >
            登录
          </button>
          <button 
            type="default" 
            @click="navigateToRegister"
            class="register-button"
          >
            注册
          </button>
        </view>
      </form>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useUserStore } from '../../store/user'

const userStore = useUserStore()
const isLoading = ref(false)
const isUsernameFocused = ref(false)
const isPasswordFocused = ref(false)
const showPassword = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const login = async () => {
  if (!loginForm.username || !loginForm.password) {
    uni.showToast({ title: '请输入用户名和密码', icon: 'none' })
    return
  }
  
  isLoading.value = true
  try {
    const result = await userStore.login(loginForm.username, loginForm.password)
    uni.showToast({ title: '登录成功', icon: 'success' })
    // 根据用户类型跳转页面
    if (result.user.user_type === 1) {
      uni.redirectTo({ url: '/pages/admin/admin' })
    } else {
      uni.redirectTo({ url: '/pages/home/home' })
    }
  } catch (error) {
    uni.showToast({ title: userStore.error || '登录失败', icon: 'none' })
  } finally {
    isLoading.value = false
  }
}

const navigateToRegister = () => {
  uni.navigateTo({ url: '/pages/register/register' })
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background-color: #fff;
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  padding: 40px 25px;
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  border-radius: 20px;
  background-color: #f5f5f5;
  padding: 10px;
}

.login-header h2 {
  margin: 0 0 10px 0;
  font-size: 28px;
  font-weight: 700;
  color: #333;
}

.login-header p {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.form-item {
  margin-bottom: 20px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  border-radius: 12px;
  border: 2px solid #e0e0e0;
  transition: all 0.3s ease;
  background-color: #f9f9f9;
}

.input-wrapper:focus-within {
  border-color: #409eff;
  background-color: #fff;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

.input-icon {
  position: absolute;
  left: 15px;
  font-size: 18px;
  color: #999;
  z-index: 1;
}

.input {
  flex: 1;
  padding: 16px 15px 16px 45px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  background: transparent;
  color: #333;
}

.input::placeholder {
  color: #999;
}

.password-toggle {
  position: absolute;
  right: 15px;
  font-size: 18px;
  cursor: pointer;
  z-index: 1;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-top: 40px;
}

button {
  padding: 16px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.login-button {
  background: linear-gradient(135deg, #409eff 0%, #667eea 100%);
  color: #fff;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
}

.register-button {
  background-color: #fff;
  color: #409eff;
  border: 2px solid #409eff;
}

.register-button:hover {
  background-color: #f0f9ff;
  transform: translateY(-2px);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .login-card {
    padding: 30px 20px;
  }
  
  .login-header h2 {
    font-size: 24px;
  }
  
  .input {
    padding: 14px 15px 14px 45px;
    font-size: 15px;
  }
  
  button {
    padding: 14px;
    font-size: 15px;
  }
}
</style>
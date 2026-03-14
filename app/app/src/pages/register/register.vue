<template>
  <view class="register-container">
    <view class="register-card">
      <view class="register-header">
        <image src="../../static/logo.png" class="logo" mode="aspectFit"></image>
        <h2>注册账号</h2>
        <p>加入 BabyHabit 培养孩子好习惯</p>
      </view>
      <form @submit.prevent="register">
        <view class="form-item">
          <view class="input-wrapper">
            <text class="input-icon">👤</text>
            <input 
              type="text" 
              v-model="registerForm.username" 
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
              v-model="registerForm.password" 
              placeholder="请输入密码" 
              class="input"
              :class="{ 'input-focused': isPasswordFocused }"
              @focus="isPasswordFocused = true"
              @blur="isPasswordFocused = false"
            />
            <text 
              class="password-toggle" 
              @click="togglePasswordVisibility('password')"
            >
              {{ showPassword ? '👁️' : '👁️‍🗨️' }}
            </text>
          </view>
        </view>
        <view class="form-item">
          <view class="input-wrapper">
            <text class="input-icon">🔒</text>
            <input 
              type="password" 
              v-model="registerForm.confirmPassword" 
              placeholder="请确认密码" 
              class="input"
              :class="{ 'input-focused': isConfirmPasswordFocused }"
              @focus="isConfirmPasswordFocused = true"
              @blur="isConfirmPasswordFocused = false"
            />
            <text 
              class="password-toggle" 
              @click="togglePasswordVisibility('confirmPassword')"
            >
              {{ showConfirmPassword ? '👁️' : '👁️‍🗨️' }}
            </text>
          </view>
        </view>
        <view class="form-item">
          <view class="input-wrapper">
            <text class="input-icon">📧</text>
            <input 
              type="email" 
              v-model="registerForm.email" 
              placeholder="请输入邮箱" 
              class="input"
              :class="{ 'input-focused': isEmailFocused }"
              @focus="isEmailFocused = true"
              @blur="isEmailFocused = false"
            />
          </view>
        </view>
        <view class="form-actions">
          <button 
            type="primary" 
            @click="register" 
            :loading="isLoading"
            class="register-button"
          >
            注册
          </button>
          <button 
            type="default" 
            @click="navigateToLogin"
            class="login-button"
          >
            登录
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
const isConfirmPasswordFocused = ref(false)
const isEmailFocused = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: ''
})

const register = async () => {
  if (!registerForm.username || !registerForm.password || !registerForm.confirmPassword || !registerForm.email) {
    uni.showToast({ title: '请填写所有字段', icon: 'none' })
    return
  }
  
  if (registerForm.password !== registerForm.confirmPassword) {
    uni.showToast({ title: '两次输入的密码不一致', icon: 'none' })
    return
  }
  
  isLoading.value = true
  try {
    const result = await userStore.register({
      username: registerForm.username,
      password: registerForm.password,
      email: registerForm.email
    })
    uni.showToast({ title: '注册成功', icon: 'success' })
    // 注册成功后跳转到首页
    uni.redirectTo({ url: '/pages/home/home' })
  } catch (error) {
    uni.showToast({ title: userStore.error || '注册失败', icon: 'none' })
  } finally {
    isLoading.value = false
  }
}

const navigateToLogin = () => {
  uni.navigateTo({ url: '/pages/login/login' })
}

const togglePasswordVisibility = (field) => {
  if (field === 'password') {
    showPassword.value = !showPassword.value
  } else if (field === 'confirmPassword') {
    showConfirmPassword.value = !showConfirmPassword.value
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 400px;
  background-color: #fff;
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  padding: 30px 25px;
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

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo {
  width: 60px;
  height: 60px;
  margin: 0 auto 15px;
  border-radius: 15px;
  background-color: #f5f5f5;
  padding: 8px;
}

.register-header h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 700;
  color: #333;
}

.register-header p {
  margin: 0;
  font-size: 13px;
  color: #666;
}

.form-item {
  margin-bottom: 15px;
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
  padding: 14px 15px 14px 45px;
  border: none;
  border-radius: 12px;
  font-size: 15px;
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
  gap: 12px;
  margin-top: 30px;
}

button {
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.register-button {
  background: linear-gradient(135deg, #409eff 0%, #667eea 100%);
  color: #fff;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}

.register-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
}

.login-button {
  background-color: #fff;
  color: #409eff;
  border: 2px solid #409eff;
}

.login-button:hover {
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
  .register-card {
    padding: 25px 20px;
  }
  
  .register-header h2 {
    font-size: 22px;
  }
  
  .input {
    padding: 12px 15px 12px 45px;
    font-size: 14px;
  }
  
  button {
    padding: 12px;
    font-size: 14px;
  }
}
</style>
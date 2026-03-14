<template>
  <view class="admin-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <h3>管理员控制台</h3>
      </view>
      <view class="header-right">
        <text class="user-name">{{ userStore.user?.name }}</text>
        <button type="button" style="float: right;" plain @click="logout" class="logout-button">
          退出登录
        </button>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="admin-content">
      <view class="admin-cards">
        <view class="card" @click="navigateToChildren">
          <view class="card-icon">👶</view>
          <view class="card-title">小孩管理</view>
          <view class="card-description">添加、编辑和管理小孩账号</view>
        </view>
        <view class="card" @click="navigateToHabits">
          <view class="card-icon">📋</view>
          <view class="card-title">习惯管理</view>
          <view class="card-description">创建和管理习惯任务</view>
        </view>
        <view class="card" @click="navigateToRewards">
          <view class="card-icon">🎁</view>
          <view class="card-title">奖励管理</view>
          <view class="card-description">设置和管理奖励物品</view>
        </view>
        <view class="card" @click="navigateToExchanges">
          <view class="card-icon">📦</view>
          <view class="card-title">兑换管理</view>
          <view class="card-description">处理和管理兑换请求</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { useUserStore } from '../../store/user'

const userStore = useUserStore()

// 导航到小孩管理页面
const navigateToChildren = () => {
  uni.navigateTo({ url: '/pages/admin/children/children' })
}

// 导航到习惯管理页面
const navigateToHabits = () => {
  uni.navigateTo({ url: '/pages/admin/habits/habits' })
}

// 导航到奖励管理页面
const navigateToRewards = () => {
  uni.navigateTo({ url: '/pages/admin/rewards/rewards' })
}

// 导航到兑换管理页面
const navigateToExchanges = () => {
  uni.navigateTo({ url: '/pages/admin/exchanges/exchanges' })
}

// 登出
const logout = () => {
  userStore.logout()
  uni.redirectTo({ url: '/pages/login/login' })
}
</script>

<style scoped>
.admin-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.header-left h3 {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-name {
  font-size: 14px;
  color: #606266;
}

.logout-button {
  font-size: 16px;
  padding: 6px 10px;
}

.admin-content {
  padding-top: 20px;
}

.admin-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.card {
  padding: 30px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
}

.card-icon {
  font-size: 48px;
  margin-bottom: 16px;
  display: block;
}

.card-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.card-description {
  font-size: 14px;
  color: #666;
  line-height: 1.4;
}

button {
  border: none;
  border-radius: 4px;
  font-size: 14px;
  padding: 8px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

button:hover {
  transform: translateY(-1px);
}

button[plain] {
  background-color: transparent;
  border: 1px solid;
}

button[plain][type="primary"] {
  color: #409eff;
  border-color: #409eff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-container {
    padding: 15px;
  }
  
  .admin-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .header-left h3 {
    font-size: 20px;
  }
  
  .header-right {
    width: 100%;
    justify-content: space-between;
  }
  
  .admin-cards {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .card {
    padding: 25px 15px;
  }
  
  .card-icon {
    font-size: 40px;
    margin-bottom: 12px;
  }
  
  .card-title {
    font-size: 18px;
  }
  
  .card-description {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .card {
    padding: 20px 12px;
  }
  
  .card-icon {
    font-size: 36px;
  }
  
  .card-title {
    font-size: 16px;
  }
  
  .card-description {
    font-size: 12px;
  }
}
</style>
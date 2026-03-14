<template>
  <view class="rewards-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>奖励管理</h3>
      </view>
      <view class="header-right">
        <text class="user-name">{{ userStore.user?.name }}</text>
        <button type="primary" plain @click="logout" class="logout-button">
          🚪
        </button>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="rewards-content">
      <button type="primary" @click="navigateToAddReward" class="add-button">
        <text class="add-icon">+</text>
        添加奖励
      </button>
      <view v-if="rewards.length > 0" class="rewards-list">
        <view v-for="reward in rewards" :key="reward.id" class="reward-item">
          <view class="reward-info">
            <view class="reward-header">
              <view class="reward-name">{{ reward.name }}</view>
              <view class="reward-status" :class="reward.status === 1 ? 'status-active' : 'status-inactive'">
                {{ reward.status === 1 ? '上架' : '下架' }}
              </view>
            </view>
            <view class="reward-details">
              <text class="detail-item">积分: {{ reward.points_required }}</text>
              <text class="detail-item">库存: {{ reward.stock }}</text>
              <text class="detail-item">描述: {{ reward.description || '无' }}</text>
              <text class="detail-item">创建时间: {{ reward.created_at }}</text>
            </view>
          </view>
          <view class="reward-actions">
            <button type="primary" size="small" @click="navigateToEditReward(reward)" class="action-button primary">编辑</button>
            <button type="danger" size="small" @click="deleteReward(reward.id)" class="action-button danger">删除</button>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">🎁</text>
        <text>暂无奖励</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../store/user'
import api from '../../../api'

const userStore = useUserStore()

// 奖励相关
const rewards = ref([])

// 导航返回
const navigateBack = () => {
  uni.navigateBack()
}

// 导航到添加奖励页面
const navigateToAddReward = () => {
  uni.navigateTo({ url: '/pages/admin/rewards/edit/edit' })
}

// 导航到编辑奖励页面
const navigateToEditReward = (reward) => {
  uni.navigateTo({ url: `/pages/admin/rewards/edit/edit?rewardId=${reward.id}` })
}

// 加载奖励列表
const loadRewards = async () => {
  try {
    const response = await api.get('/api/admin/rewards')
    rewards.value = response.data.items || []
  } catch (error) {
    console.error('Failed to load rewards:', error)
  }
}

// 删除奖励
const deleteReward = async (id) => {
  try {
    await uni.showModal({
      title: '删除奖励',
      content: '确定要删除这个奖励吗？',
      success: async (res) => {
        if (res.confirm) {
          await api.delete(`/api/admin/rewards/${id}`)
          await loadRewards()
          uni.showToast({ title: '删除成功', icon: 'success' })
        }
      }
    })
  } catch (error) {
    console.error('Failed to delete reward:', error)
    uni.showToast({ title: '删除失败', icon: 'none' })
  }
}

// 登出
const logout = () => {
  userStore.logout()
  uni.redirectTo({ url: '/pages/login/login' })
}

onMounted(async () => {
  await loadRewards()
})
</script>

<style scoped>
.rewards-container {
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

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.back-btn {
  font-size: 16px;
  color: #409eff;
  background-color: transparent;
  border: none;
}

.header-left h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
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

.rewards-content {
  padding-top: 20px;
}

.add-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #409eff 0%, #667eea 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transition: all 0.3s ease;
}

.add-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.4);
}

.add-icon {
  font-size: 20px;
  font-weight: bold;
}

.rewards-list {
  margin-top: 20px;
}

.reward-item {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 12px;
  margin-bottom: 15px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.reward-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.reward-info {
  flex: 1;
  min-width: 0;
}

.reward-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.reward-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.reward-details {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.detail-item {
  font-size: 14px;
  color: #666;
}

.reward-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
}

.action-button {
  font-size: 12px;
  padding: 6px 12px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-button:hover {
  transform: translateY(-1px);
}

.action-button.primary {
  background-color: #409eff;
  color: #fff;
}

.action-button.danger {
  background-color: #f56c6c;
  color: #fff;
}

.action-button.info {
  background-color: #909399;
  color: #fff;
}

.reward-status, .exchange-status {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-active {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-inactive {
  background-color: #fef0f0;
  color: #f56c6c;
}

.status-pending {
  background-color: #ecf5ff;
  color: #409eff;
}

.status-canceled {
  background-color: #f5f7fa;
  color: #909399;
}

.empty-state {
  padding: 60px 0;
  text-align: center;
  color: #999;
}

.empty-icon {
  display: block;
  font-size: 48px;
  margin-bottom: 16px;
}

/* 对话框样式 */
.dialog-content {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  width: 80%;
  max-width: 500px;
}

.dialog-header {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
}

.dialog-body {
  margin-bottom: 20px;
}

.form-item {
  margin-bottom: 15px;
}

.form-item label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.form-item input, .form-item textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
}

.form-item textarea {
  min-height: 80px;
  resize: vertical;
}

.radio-group {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

/* 兑换记录列表样式 */
.exchanges-list {
  max-height: 400px;
  overflow-y: auto;
}

.exchange-item {
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.exchange-child {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.exchange-date, .exchange-points {
  font-size: 14px;
  color: #666;
}

.exchange-points {
  color: #f56c6c;
  font-weight: 500;
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

button[type="primary"] {
  background-color: #409eff;
  color: #fff;
}

button[type="default"] {
  background-color: #f0f0f0;
  color: #333;
}

button[type="danger"] {
  background-color: #f56c6c;
  color: #fff;
}

button[size="small"] {
  font-size: 12px;
  padding: 6px 12px;
  align-self: flex-start;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .rewards-container {
    padding: 15px;
  }
  
  .admin-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .header-left {
    width: 100%;
  }
  
  .header-right {
    width: 100%;
    justify-content: space-between;
  }
  
  .reward-item {
    flex-direction: column;
    align-items: stretch;
  }
  
  .reward-actions {
    flex-direction: row;
    flex-wrap: wrap;
    margin-top: 15px;
  }
  
  .action-button {
    flex: 1;
    min-width: 80px;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .add-button {
    padding: 10px 16px;
    font-size: 14px;
  }
  
  .reward-item {
    padding: 15px;
  }
  
  .reward-name {
    font-size: 16px;
  }
  
  .detail-item {
    font-size: 13px;
  }
  
  .action-button {
    font-size: 11px;
    padding: 5px 10px;
  }
}
</style>
<template>
  <view class="reward-redemptions-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>奖励兑换记录</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="reward-redemptions-content">
      <view v-if="rewardRedemptions.length > 0" class="redemptions-list">
        <view v-for="redemption in rewardRedemptions" :key="redemption.id" class="redemption-item">
          <view class="redemption-child">{{ redemption.child.name }}</view>
          <view class="redemption-reward">{{ redemption.reward.name }}</view>
          <view class="redemption-points">{{ redemption.reward.points_required }} 积分</view>
          <view class="redemption-time">{{ redemption.redemption_time }}</view>
          <view class="redemption-status" :class="redemption.status === 1 ? 'status-active' : 'status-inactive'">
            {{ redemption.status === 1 ? '已兑换' : '已取消' }}
          </view>
          <button 
            type="danger" 
            size="small" 
            @click="cancelRedemption(redemption)"
            :disabled="redemption.status !== 1"
            class="cancel-button"
          >
            取消
          </button>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">🎯</text>
        <text>暂无兑换记录</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'

const userStore = useUserStore()

// 兑换记录相关
const rewardRedemptions = ref([])

// 导航返回
const navigateBack = () => {
  uni.navigateBack()
}

// 加载奖励兑换记录
const loadRewardRedemptions = async () => {
  try {
    const response = await api.get('/admin/reward-redemptions')
    rewardRedemptions.value = response.data.redemptions || []
  } catch (error) {
    console.error('Failed to load reward redemptions:', error)
    uni.showToast({ title: '加载兑换记录失败', icon: 'none' })
  }
}

// 取消奖励兑换
const cancelRedemption = async (redemption) => {
  try {
    await uni.showModal({
      title: '取消兑换',
      content: `确定要取消 ${redemption.child.name} 对 "${redemption.reward.name}" 的兑换吗？`,
      success: async (res) => {
        if (res.confirm) {
          await api.delete(`/admin/reward-redemptions/${redemption.id}`)
          uni.showToast({ title: '取消成功', icon: 'success' })
          // 重新加载列表
          await loadRewardRedemptions()
        }
      }
    })
  } catch (error) {
    console.error('Failed to cancel redemption:', error)
    uni.showToast({ title: '取消失败', icon: 'none' })
  }
}

onMounted(async () => {
  await loadRewardRedemptions()
})
</script>

<style scoped>
.reward-redemptions-container {
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

.reward-redemptions-content {
  padding-top: 20px;
}

.redemptions-list {
  margin-top: 20px;
}

.redemption-item {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 12px;
  margin-bottom: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.redemption-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.redemption-child {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.redemption-reward {
  font-size: 16px;
  font-weight: 500;
  color: #409eff;
  margin-bottom: 8px;
}

.redemption-points {
  font-size: 14px;
  color: #67c23a;
  font-weight: 500;
  margin-bottom: 8px;
}

.redemption-time {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.redemption-status {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 12px;
}

.status-active {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-inactive {
  background-color: #fef0f0;
  color: #f56c6c;
}

.cancel-button {
  background-color: #f56c6c;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.4);
}

.cancel-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .reward-redemptions-container {
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
  
  .redemption-item {
    padding: 15px;
  }
  
  .redemption-child {
    font-size: 16px;
  }
}
</style>
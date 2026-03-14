<template>
  <view class="reward-assignments-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>奖励分配记录</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="reward-assignments-content">
      <view v-if="rewardAssignments.length > 0" class="assignments-list">
        <view v-for="assignment in rewardAssignments" :key="assignment.id" class="assignment-item">
          <view class="assignment-child">{{ assignment.child.name }}</view>
          <view class="assignment-reward">{{ assignment.reward.name }}</view>
          <view class="assignment-points">{{ assignment.reward.points_required }} 积分</view>
          <view class="assignment-time">{{ assignment.assign_time }}</view>
          <view class="assignment-status" :class="assignment.status === 1 ? 'status-active' : 'status-inactive'">
            {{ assignment.status === 1 ? '已分配' : '已撤销' }}
          </view>
          <button 
            type="danger" 
            size="small" 
            @click="revokeAssignment(assignment)"
            :disabled="assignment.status !== 1"
            class="revoke-button"
          >
            撤销
          </button>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">🎁</text>
        <text>暂无分配记录</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'

const userStore = useUserStore()

// 分配记录相关
const rewardAssignments = ref([])

// 导航返回
const navigateBack = () => {
  uni.navigateBack()
}

// 加载奖励分配记录
const loadRewardAssignments = async () => {
  try {
    const response = await api.get('/admin/reward-assignments')
    rewardAssignments.value = response.data.assignments || []
  } catch (error) {
    console.error('Failed to load reward assignments:', error)
    uni.showToast({ title: '加载分配记录失败', icon: 'none' })
  }
}

// 撤销奖励分配
const revokeAssignment = async (assignment) => {
  try {
    await uni.showModal({
      title: '撤销奖励',
      content: `确定要撤销分配给 ${assignment.child.name} 的奖励 "${assignment.reward.name}" 吗？`,
      success: async (res) => {
        if (res.confirm) {
          await api.delete(`/admin/reward-assignments/${assignment.id}`)
          uni.showToast({ title: '撤销成功', icon: 'success' })
          // 重新加载列表
          await loadRewardAssignments()
        }
      }
    })
  } catch (error) {
    console.error('Failed to revoke assignment:', error)
    uni.showToast({ title: '撤销失败', icon: 'none' })
  }
}

onMounted(async () => {
  await loadRewardAssignments()
})
</script>

<style scoped>
.reward-assignments-container {
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

.reward-assignments-content {
  padding-top: 20px;
}

.assignments-list {
  margin-top: 20px;
}

.assignment-item {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 12px;
  margin-bottom: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.assignment-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.assignment-child {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.assignment-reward {
  font-size: 16px;
  font-weight: 500;
  color: #409eff;
  margin-bottom: 8px;
}

.assignment-points {
  font-size: 14px;
  color: #67c23a;
  font-weight: 500;
  margin-bottom: 8px;
}

.assignment-time {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.assignment-status {
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

.revoke-button {
  background-color: #f56c6c;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.revoke-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.4);
}

.revoke-button:disabled {
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
  .reward-assignments-container {
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
  
  .assignment-item {
    padding: 15px;
  }
  
  .assignment-child {
    font-size: 16px;
  }
}
</style>
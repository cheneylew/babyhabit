<template>
  <view class="assigned-habits-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>已分配习惯</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="assigned-habits-content">
      <view class="child-info">
        <text class="child-label">小孩:</text>
        <text class="child-name">{{ currentChild?.name }}</text>
      </view>
      
      <view v-if="assignedHabits.length > 0" class="assigned-habits-list">
        <view v-for="assignment in assignedHabits" :key="assignment.id" class="assigned-habit-item">
          <view class="assigned-habit-name">{{ assignment.habit.name }}</view>
          <view class="assigned-habit-points">{{ assignment.habit.reward_points }} 积分</view>
          <view class="assigned-habit-time">{{ assignment.assign_time }}</view>
          <button type="danger" size="small" @click="deleteAssignedHabit(assignment)" class="delete-button">删除</button>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">📋</text>
        <text>暂无已分配习惯</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'

const userStore = useUserStore()

// 已分配习惯相关
const assignedHabits = ref([])
const currentChild = ref(null)

// 获取路由参数
const getRouteParams = () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  return currentPage.options
}

// 导航返回
const navigateBack = () => {
  uni.navigateBack()
}

// 加载已分配习惯列表
const loadAssignedHabits = async (childId) => {
  try {
    console.log('Loading assigned habits for child:', childId)
    const response = await api.get('/api/admin/assigned-habits', {
      params: { child_id: childId }
    })
    console.log('Response data:', response.data)
    assignedHabits.value = response.data.assignments || []
    console.log('Assigned habits:', assignedHabits.value)
  } catch (error) {
    console.error('Failed to load assigned habits:', error)
  }
}

// 删除已分配习惯
const deleteAssignedHabit = async (assignment) => {
  try {
    await uni.showModal({
      title: '删除习惯',
      content: `确定要删除分配给 ${currentChild.value.name} 的习惯 "${assignment.habit.name}" 吗？`,
      success: async (res) => {
        if (res.confirm) {
          await api.delete(`/api/admin/habit-assignments/${assignment.id}`)
          uni.showToast({ title: '删除成功', icon: 'success' })
          // 重新加载列表
          await loadAssignedHabits(currentChild.value.id)
        }
      }
    })
  } catch (error) {
    console.error('Failed to delete assigned habit:', error)
    uni.showToast({ title: '删除失败', icon: 'none' })
  }
}

onMounted(async () => {
  const params = getRouteParams()
  console.log('Route params:', params)
  if (params.childId) {
    // 加载小孩信息
    try {
      console.log('Loading child info for ID:', params.childId)
      const response = await api.get(`/api/admin/children/${params.childId}`)
      console.log('Child info response:', response.data)
      currentChild.value = response.data.child
      // 加载已分配习惯
      await loadAssignedHabits(params.childId)
    } catch (error) {
      console.error('Failed to load child info:', error)
      uni.showToast({ title: '加载小孩信息失败', icon: 'none' })
    }
  } else {
    console.error('No childId in route params')
  }
})
</script>

<style scoped>
.assigned-habits-container {
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

.assigned-habits-content {
  padding-top: 20px;
}

.child-info {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f0f9eb;
  border-radius: 8px;
}

.child-label {
  font-size: 14px;
  font-weight: 500;
  color: #666;
}

.child-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.assigned-habits-list {
  margin-top: 20px;
}

.assigned-habit-item {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 12px;
  margin-bottom: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.assigned-habit-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.assigned-habit-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.assigned-habit-points {
  font-size: 14px;
  color: #67c23a;
  font-weight: 500;
  margin-bottom: 8px;
}

.assigned-habit-time {
  font-size: 14px;
  color: #666;
  margin-bottom: 12px;
}

.delete-button {
  background-color: #f56c6c;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.delete-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.4);
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
  .assigned-habits-container {
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
  
  .assigned-habit-item {
    padding: 15px;
  }
  
  .assigned-habit-name {
    font-size: 16px;
  }
}
</style>
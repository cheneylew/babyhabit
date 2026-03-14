<template>
  <view class="habit-assignments-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>习惯分配记录</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="habit-assignments-content">
      <view class="habit-info">
        <text class="habit-label">习惯:</text>
        <text class="habit-name">{{ currentHabit?.name }}</text>
      </view>
      
      <view v-if="habitAssignments.length > 0" class="assignments-list">
        <view v-for="assignment in habitAssignments" :key="assignment.id" class="assignment-item">
          <view class="assignment-child">{{ assignment.child.name }}</view>
          <view class="assignment-time">{{ assignment.assign_time }}</view>
          <button type="danger" size="small" @click="deleteAssignment(assignment)" class="delete-button">删除</button>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">👥</text>
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
const habitAssignments = ref([])
const currentHabit = ref(null)

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

// 加载分配记录
const loadHabitAssignments = async (habitId) => {
  try {
    const response = await api.get('/api/admin/assigned-habits', {
      params: { habit_id: habitId }
    })
    habitAssignments.value = response.data.assignments || []
  } catch (error) {
    console.error('Failed to load assignments:', error)
    uni.showToast({ title: '加载分配记录失败', icon: 'none' })
  }
}

// 删除习惯分配
const deleteAssignment = async (assignment) => {
  try {
    await uni.showModal({
      title: '删除分配',
      content: `确定要删除分配给 ${assignment.child.name} 的习惯吗？`,
      success: async (res) => {
        if (res.confirm) {
          await api.delete(`/api/admin/habit-assignments/${assignment.id}`)
          uni.showToast({ title: '删除成功', icon: 'success' })
          // 重新加载列表
          await loadHabitAssignments(currentHabit.value.id)
        }
      }
    })
  } catch (error) {
    console.error('Failed to delete assignment:', error)
    uni.showToast({ title: '删除失败', icon: 'none' })
  }
}

onMounted(async () => {
  const params = getRouteParams()
  if (params.habitId) {
    // 加载习惯信息
    try {
      const response = await api.get(`/api/admin/habits/${params.habitId}`)
      currentHabit.value = response.data.habit
      // 加载分配记录
      await loadHabitAssignments(params.habitId)
    } catch (error) {
      console.error('Failed to load habit info:', error)
      uni.showToast({ title: '加载习惯信息失败', icon: 'none' })
    }
  }
})
</script>

<style scoped>
.habit-assignments-container {
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

.habit-assignments-content {
  padding-top: 20px;
}

.habit-info {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f0f9eb;
  border-radius: 8px;
}

.habit-label {
  font-size: 14px;
  font-weight: 500;
  color: #666;
}

.habit-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
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

.assignment-time {
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
  .habit-assignments-container {
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
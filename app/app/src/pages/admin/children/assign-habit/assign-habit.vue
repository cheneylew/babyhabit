<template>
  <view class="assign-habit-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>分配习惯</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="assign-habit-content">
      <view class="form-container">
        <view class="child-info">
          <text class="child-label">小孩:</text>
          <text class="child-name">{{ currentChild?.name }}</text>
        </view>
        
        <view class="form-item">
          <label>选择习惯</label>
          <picker @change="handleHabitChange" :value="habitIndex" :range="habits" range-key="name">
            <view class="picker">
              {{ selectedHabit?.name || '请选择习惯' }}
            </view>
          </picker>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="action-buttons">
        <button type="default" @click="navigateBack" class="cancel-button">取消</button>
        <button type="primary" @click="saveAssignHabit" class="save-button" :disabled="!selectedHabit">保存</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'

const userStore = useUserStore()

// 习惯相关
const habits = ref([])
const habitIndex = ref(0)
const currentChild = ref(null)

// 计算属性：当前选中的习惯
const selectedHabit = computed(() => {
  return habits.value[habitIndex.value] || null
})

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

// 加载习惯列表
const loadHabits = async () => {
  try {
    const response = await api.get('/api/admin/habits')
    habits.value = response.data.habits || []
  } catch (error) {
    console.error('Failed to load habits:', error)
    uni.showToast({ title: '加载习惯失败', icon: 'none' })
  }
}

// 处理习惯选择变化
const handleHabitChange = (e) => {
  habitIndex.value = e.detail.value
}

// 保存分配习惯
const saveAssignHabit = async () => {
  try {
    await api.post('/api/admin/habits/assign', {
      habit_id: selectedHabit.value.id,
      child_id: currentChild.value.id
    })
    uni.showToast({ title: '习惯分配成功', icon: 'success' })
    uni.navigateBack()
  } catch (error) {
    console.error('Failed to assign habit:', error)
    uni.showToast({ title: '习惯分配失败', icon: 'none' })
  }
}

onMounted(async () => {
  const params = getRouteParams()
  if (params.childId) {
    // 加载小孩信息
    try {
      const response = await api.get(`/api/admin/children/${params.childId}`)
      currentChild.value = response.data.child
    } catch (error) {
      console.error('Failed to load child info:', error)
      uni.showToast({ title: '加载小孩信息失败', icon: 'none' })
    }
  }
  await loadHabits()
})
</script>

<style scoped>
.assign-habit-container {
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

.assign-habit-content {
  padding-top: 20px;
}

.form-container {
  background-color: #f9f9f9;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
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

.form-item {
  margin-bottom: 20px;
}

.form-item label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.picker {
  padding: 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  font-size: 14px;
  color: #606266;
  background-color: #fff;
}

.action-buttons {
  display: flex;
  gap: 15px;
  justify-content: flex-end;
}

.cancel-button,
.save-button {
  padding: 12px 24px;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.cancel-button {
  background-color: #f0f0f0;
  color: #333;
  border: none;
}

.save-button {
  background: linear-gradient(135deg, #409eff 0%, #667eea 100%);
  color: #fff;
  border: none;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.save-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.cancel-button:hover,
.save-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .assign-habit-container {
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
  
  .form-container {
    padding: 15px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .cancel-button,
  .save-button {
    width: 100%;
    text-align: center;
  }
}
</style>
<template>
  <view class="habit-records-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>习惯打卡记录</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="habit-records-content">
      <view class="habit-info">
        <text class="habit-label">习惯:</text>
        <text class="habit-name">{{ currentHabit?.name }}</text>
      </view>
      
      <view class="date-range">
        <text class="date-label">时间范围:</text>
        <text class="date-value">{{ startDate }} 至 {{ endDate }}</text>
      </view>
      
      <view v-if="habitCheckinRecords.length > 0" class="records-list">
        <view v-for="record in habitCheckinRecords" :key="record.id" class="record-item">
          <view class="record-child">{{ record.child.name }}</view>
          <view class="record-date">{{ record.checkin_date }} {{ record.checkin_time }}</view>
          <view class="record-type">{{ record.checkin_type === 1 ? '正常' : '补卡' }}</view>
          <view class="record-points">+{{ record.points_rewarded }} 积分</view>
          <view class="record-status" :class="record.status === 1 ? 'status-active' : 'status-inactive'">
            {{ record.status === 1 ? '正常' : '已回退' }}
          </view>
          <button 
            type="danger" 
            size="small" 
            @click="rollbackCheckin(record)"
            :disabled="record.status !== 1"
            class="rollback-button"
          >
            回退
          </button>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">📅</text>
        <text>暂无打卡记录</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'
import dayjs from 'dayjs'

const userStore = useUserStore()

// 打卡记录相关
const habitCheckinRecords = ref([])
const currentHabit = ref(null)
const startDate = ref('')
const endDate = ref('')

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

// 加载习惯打卡记录
const loadHabitCheckinRecords = async (habitId) => {
  try {
    const start = dayjs().subtract(30, 'day').format('YYYY-MM-DD')
    const end = dayjs().format('YYYY-MM-DD')
    startDate.value = start
    endDate.value = end
    
    const response = await api.get('/api/admin/child/checkin-records', {
      params: {
        habit_id: habitId,
        start_date: start,
        end_date: end
      }
    })
    habitCheckinRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load checkin records:', error)
    uni.showToast({ title: '加载打卡记录失败', icon: 'none' })
  }
}

// 回退打卡
const rollbackCheckin = async (record) => {
  try {
    await uni.showModal({
      title: '回退打卡',
      content: `确定要回退 ${record.child.name} 的这次打卡吗？回退后将扣除 ${record.points_rewarded} 积分。`,
      success: async (res) => {
        if (res.confirm) {
          await api.post('/api/admin/checkin/rollback', {
            checkin_id: record.id,
            reason: '管理员手动回退'
          })
          uni.showToast({ title: '回退成功', icon: 'success' })
          // 重新加载列表
          await loadHabitCheckinRecords(currentHabit.value.id)
        }
      }
    })
  } catch (error) {
    console.error('Failed to rollback checkin:', error)
    uni.showToast({ title: '回退失败', icon: 'none' })
  }
}

onMounted(async () => {
  const params = getRouteParams()
  if (params.habitId) {
    // 加载习惯信息
    try {
      const response = await api.get(`/api/admin/habits/${params.habitId}`)
      currentHabit.value = response.data.habit
      // 加载打卡记录
      await loadHabitCheckinRecords(params.habitId)
    } catch (error) {
      console.error('Failed to load habit info:', error)
      uni.showToast({ title: '加载习惯信息失败', icon: 'none' })
    }
  }
})
</script>

<style scoped>
.habit-records-container {
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

.habit-records-content {
  padding-top: 20px;
}

.habit-info {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
  padding: 15px;
  background-color: #f0f9eb;
  border-radius: 8px;
}

.date-range {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #ecf5ff;
  border-radius: 8px;
}

.habit-label,
.date-label {
  font-size: 14px;
  font-weight: 500;
  color: #666;
}

.habit-name,
.date-value {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.records-list {
  margin-top: 20px;
}

.record-item {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 12px;
  margin-bottom: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.record-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.record-child {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.record-date,
.record-type {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.record-points {
  font-size: 14px;
  color: #67c23a;
  font-weight: 500;
  margin-bottom: 8px;
}

.record-status {
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

.rollback-button {
  background-color: #f56c6c;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.rollback-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.4);
}

.rollback-button:disabled {
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
  .habit-records-container {
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
  
  .record-item {
    padding: 15px;
  }
  
  .record-child {
    font-size: 16px;
  }
}
</style>
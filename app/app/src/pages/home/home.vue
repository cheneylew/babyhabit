<template>
  <view class="home-container">
    <!-- 用户信息卡片 -->
    <view class="user-info-card">
      <view class="user-info">
        <view class="user-avatar">
          <image :src="user?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" mode="aspectFill"></image>
          <text v-if="!user?.avatar" class="avatar-text">{{ user?.name?.charAt(0) || 'U' }}</text>
        </view>
        <view class="user-details">
          <h3>{{ user?.name }}</h3>
          <p>{{ user?.user_type === 1 ? '父母账号' : '小孩账号' }}</p>
        </view>
        <button type="danger" plain @click="logout" class="logout-button">
          🚪
        </button>
      </view>
      <view class="points-section">
        <button type="primary" size="small" class="points-badge">
          💰 我的积分 {{ user?.points_balance || 0 }}
        </button>
        <button 
          v-if="user?.user_type === 2" 
          type="warning" 
          size="small" 
          @click="navigateToChildRewards"
          class="rewards-button"
        >
          🎁 积分商城
        </button>
      </view>
    </view>

    <!-- 快捷操作 -->
    <view class="quick-actions">
      <button type="primary" @click="navigateToCheckin" class="action-button checkin-button">
        <text class="action-icon">📅</text>
        <text class="action-text">去打卡</text>
      </button>
      <button 
        v-if="user?.user_type === 1" 
        type="success" 
        @click="navigateToAdmin" 
        class="action-button admin-button"
      >
        <text class="action-icon">⚙️</text>
        <text class="action-text">管理后台</text>
      </button>
    </view>

    <!-- 今日待打卡 -->
    <view class="habits-card">
      <view class="card-header">
        <h3>今日待打卡</h3>
      </view>
      <view v-if="habits.length > 0" class="habits-list">
        <view 
          v-for="habit in habits" 
          :key="habit.id" 
          class="habit-item" 
          :class="{ 'checked-in': isCheckedIn(habit.id) }"
          @click="navigateToCheckin"
        >
          <view class="habit-icon" :class="{ 'checked-in': isCheckedIn(habit.id) }">
            <image v-if="habit.icon" :src="habit.icon" mode="aspectFill"></image>
            <view v-else class="default-icon" :class="{ 'checked-in': isCheckedIn(habit.id) }">{{ habit.name.charAt(0) }}</view>
          </view>
          <view class="habit-details">
            <h4 :class="{ 'checked-in': isCheckedIn(habit.id) }">{{ habit.name }}</h4>
            <p>{{ habit.description }}</p>
            <view class="habit-meta">
              <span class="habit-time">{{ habit.checkin_time_start.substring(0, 5) }} - {{ habit.checkin_time_end.substring(0, 5) }}</span>
              <span class="reward-points">{{ habit.reward_points }} 积分</span>
            </view>
          </view>
          <view class="habit-status">
            <view v-if="isCheckedIn(habit.id)" class="status-badge checked">已打卡</view>
            <view v-else class="status-badge pending">待打卡</view>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text>🎉 暂无待打卡习惯</text>
      </view>
    </view>

    <!-- 统计信息 -->
    <view class="stats-card">
      <view class="card-header">
        <h3>打卡统计</h3>
      </view>
      <view class="stats-grid">
        <view class="stat-item" @click="navigateToCheckin">
          <view class="stat-icon">📊</view>
          <view class="stat-value">{{ todayCheckins }}</view>
          <view class="stat-label">今日打卡</view>
        </view>
        <view class="stat-item" @click="navigateToCheckin">
          <view class="stat-icon">📈</view>
          <view class="stat-value">{{ totalCheckins }}</view>
          <view class="stat-label">累计打卡</view>
        </view>
        <view class="stat-item" @click="navigateToCheckin">
          <view class="stat-icon">🔥</view>
          <view class="stat-value">{{ currentStreak }}</view>
          <view class="stat-label">连续打卡</view>
        </view>
        <view class="stat-item" @click="user?.user_type === 2 ? navigateToChildRewards : navigateToAdmin">
          <view class="stat-icon">💰</view>
          <view class="stat-value">{{ totalPoints }}</view>
          <view class="stat-label">总积分</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../../store/user'
import api from '../../api'
import dayjs from 'dayjs'

const userStore = useUserStore()
const user = computed(() => userStore.user)

// 导航方法
const navigateToCheckin = () => {
  uni.navigateTo({ url: '/pages/checkin/checkin' })
}

const navigateToChildRewards = () => {
  uni.navigateTo({ url: '/pages/child-rewards/child-rewards' })
}

const navigateToAdmin = () => {
  uni.navigateTo({ url: '/pages/admin/admin' })
}

// 登出
const logout = () => {
  userStore.logout()
  uni.redirectTo({ url: '/pages/login/login' })
}

const habits = ref([])
const checkinRecords = ref([])
const todayCheckins = ref(0)
const totalCheckins = ref(0)
const currentStreak = ref(0)
const totalPoints = ref(0)

const loadHabits = async () => {
  try {
    const response = await api.get('/api/habits')
    habits.value = response.data.habits
  } catch (error) {
    console.error('Failed to load habits:', error)
  }
}

const loadCheckinRecords = async () => {
  try {
    const startDate = dayjs().subtract(30, 'day').format('YYYY-MM-DD')
    const endDate = dayjs().format('YYYY-MM-DD')
    const response = await api.get('/api/checkin/records', {
      params: { start_date: startDate, end_date: endDate }
    })
    checkinRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load checkin records:', error)
  }
}

const isCheckedIn = (habitId) => {
  const today = dayjs().format('YYYY-MM-DD')
  return checkinRecords.value.some(record => {
    const recordDate = dayjs(record.checkin_date).format('YYYY-MM-DD')
    return record.habit_id === habitId && recordDate === today
  })
}

const loadStats = async () => {
  try {
    const today = dayjs().format('YYYY-MM-DD')
    // 今日打卡
    todayCheckins.value = checkinRecords.value.filter(record => {
      const recordDate = dayjs(record.checkin_date).format('YYYY-MM-DD')
      return recordDate === today
    }).length
    
    // 累计打卡
    totalCheckins.value = checkinRecords.value.length
    
    // 连续打卡（取所有习惯中最大的连续打卡天数）
    const streakPromises = habits.value.map(habit => 
      api.get(`/api/checkin/streak/${habit.id}`).then(res => res.data.streak?.current_streak || 0)
    )
    const streaks = await Promise.all(streakPromises)
    currentStreak.value = Math.max(...streaks, 0)
    
    // 总积分
    totalPoints.value = user.value?.points_balance || 0
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

onMounted(async () => {
  if (!user.value) {
    await userStore.getUserInfo()
  }
  await loadHabits()
  await loadCheckinRecords()
  await loadStats()
})
</script>

<style scoped>
.home-container {
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e9f2 100%);
  min-height: 100vh;
}

/* 用户信息卡片 */
.user-info-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  padding: 25px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  margin-bottom: 20px;
  color: white;
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
}

.user-avatar {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  overflow: hidden;
  background-color: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3px solid rgba(255, 255, 255, 0.3);
  flex-shrink: 0;
  position: relative;
}

.user-avatar image {
  width: 100%;
  height: 100%;
}

.avatar-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 28px;
  font-weight: bold;
  color: #fff;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-details h3 {
  margin: 0 0 5px 0;
  font-size: 20px;
  font-weight: 700;
  color: white;
}

.user-details p {
  margin: 0 0 12px 0;
  color: rgba(255, 255, 255, 0.8);
}

.logout-button {
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 10px 15px;
  font-size: 18px;
  color: white;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.logout-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.05);
}

.points-section {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.points-badge, .rewards-button {
  border-radius: 20px;
  font-size: 14px;
  padding: 8px 16px;
  font-weight: 600;
  transition: all 0.3s ease;
  border: none;
}

.points-badge {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.rewards-button {
  background: rgba(255, 255, 255, 0.9);
  color: #764ba2;
}

.rewards-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 快捷操作 */
.quick-actions {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.action-button {
  flex: 1;
  border-radius: 15px;
  padding: 18px;
  font-size: 15px;
  font-weight: 600;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  border: none;
  color: white;
}

.action-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.action-icon {
  font-size: 24px;
}

.action-text {
  font-size: 14px;
}

.checkin-button {
  background: linear-gradient(135deg, #409eff 0%, #667eea 100%);
}

.admin-button {
  background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);
}

/* 习惯卡片区域 */
.habits-card, .stats-card {
  background-color: #fff;
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  transition: all 0.3s ease;
  animation: fadeIn 0.5s ease-out 0.2s both;
}

.habits-card:hover, .stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.12);
}

.card-header {
  margin-bottom: 15px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-header h3::before {
  content: '📋';
  font-size: 18px;
}

.stats-card .card-header h3::before {
  content: '📊';
}

.habits-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.habit-item {
  display: flex;
  align-items: center;
  padding: 18px;
  background-color: #f9f9f9;
  border-radius: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
  border: 2px solid transparent;
}

.habit-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  border-color: #e6f7ff;
}

.habit-item.checked-in {
  background-color: #f0f9ff;
  border-color: #91d5ff;
  opacity: 0.8;
}

.habit-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  overflow: hidden;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.habit-icon.checked-in {
  background-color: #e6f7ff;
  box-shadow: 0 0 0 4px rgba(24, 144, 255, 0.1);
  opacity: 0.8;
}

.habit-icon image {
  width: 100%;
  height: 100%;
}

.default-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
  flex-shrink: 0;
}

.default-icon.checked-in {
  background: linear-gradient(135deg, #91d5ff 0%, #409eff 100%);
}

.habit-details {
  flex: 1;
  min-width: 0;
}

.habit-details h4 {
  margin: 0 0 6px 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  transition: all 0.3s ease;
}

.habit-details h4.checked-in {
  color: #1890ff;
  text-decoration: line-through;
}

.habit-details p {
  margin: 0 0 8px 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
}

.habit-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.habit-time {
  color: #999;
  font-size: 13px;
  background: #f0f0f0;
  padding: 4px 10px;
  border-radius: 12px;
}

.reward-points {
  background-color: #f0f9eb;
  color: #67c23a;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
}

.habit-status {
  margin-left: 15px;
}

.status-badge {
  padding: 6px 12px;
  border-radius: 15px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.pending {
  background: #fff7e6;
  color: #fa8c16;
}

.status-badge.checked {
  background: #f6ffed;
  color: #52c41a;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
  color: #999;
  font-size: 16px;
  background: #f9f9f9;
  border-radius: 12px;
  margin-top: 10px;
}

/* 统计信息区域 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.stat-item {
  background: linear-gradient(135deg, #f9f9f9 0%, #f0f0f0 100%);
  border-radius: 15px;
  padding: 20px;
  text-align: center;
  transition: all 0.3s ease;
  cursor: pointer;
  border: 2px solid transparent;
}

.stat-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
  border-color: #e6f7ff;
}

.stat-icon {
  font-size: 28px;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #667eea;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 13px;
  color: #666;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 15px;
  }
  
  .user-details {
    text-align: center;
  }
  
  .points-section {
    justify-content: center;
    width: 100%;
  }
  
  .quick-actions {
    flex-direction: column;
  }
  
  .action-button {
    flex-direction: row;
    justify-content: center;
    gap: 10px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .home-container {
    padding: 15px;
  }
  
  .user-info-card {
    padding: 20px;
  }
  
  .user-avatar {
    width: 60px;
    height: 60px;
  }
  
  .user-details h3 {
    font-size: 18px;
  }
  
  .habit-item {
    padding: 15px;
  }
  
  .habit-icon {
    width: 48px;
    height: 48px;
  }
  
  .default-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }
  
  .stats-grid {
    gap: 12px;
  }
  
  .stat-item {
    padding: 15px;
  }
  
  .stat-value {
    font-size: 24px;
  }
}
</style>
<template>
  <div class="home-container">
    <el-card class="user-info-card">
      <div class="user-info">
        <el-avatar :size="80" :src="user?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'">
          {{ user?.name?.charAt(0) || 'U' }}
        </el-avatar>
        <div class="user-details">
          <h3>{{ user?.name }}</h3>
          <p>{{ user?.user_type === 1 ? '父母账号' : '小孩账号' }}</p>
          <div class="points-section">
            <el-badge :value="user?.points_balance || 0" type="success" class="points-badge">
              <el-button type="primary" size="small">我的积分</el-button>
            </el-badge>
            <el-button 
              v-if="user?.user_type === 2" 
              type="warning" 
              size="small" 
              @click="$router.push('/child-rewards')"
              style="margin-left: 10px;"
            >
              积分商城
            </el-button>
            <el-button 
              v-if="user?.user_type === 2" 
              type="info" 
              size="small" 
              @click="$router.push('/vocabulary')"
              style="margin-left: 10px;"
            >
              艾宾浩斯单词记忆
            </el-button>
          </div>
        </div>
        <div class="user-actions">
          <el-button v-if="user?.user_type === 1" type="success" plain @click="$router.push('/admin')">管理后台</el-button>
          <el-button type="danger" plain @click="logout">登出</el-button>
        </div>
      </div>
    </el-card>

    <el-card class="habits-card">
      <template #header>
        <div class="card-header">
          <h3>今日待打卡</h3>
          <el-button type="primary" size="small" @click="$router.push('/checkin')">去打卡</el-button>
        </div>
      </template>
      <div v-if="todayHabits.length > 0" class="habits-list">
        <el-card v-for="habit in todayHabits" :key="habit.id" class="habit-item" :class="{ 'checked-in': isCheckedIn(habit.id) }">
          <div class="habit-info">
            <div class="habit-icon" :class="{ 'checked-in': isCheckedIn(habit.id) }">
              <img v-if="habit.icon" :src="habit.icon" alt="habit icon" />
              <div v-else class="default-icon" :class="{ 'checked-in': isCheckedIn(habit.id) }">{{ habit.name.charAt(0) }}</div>
            </div>
            <div class="habit-details">
              <h4 :class="{ 'checked-in': isCheckedIn(habit.id) }">{{ habit.name }}</h4>
              <p>{{ habit.description }}</p>
              <div class="habit-time">
                <span>{{ habit.checkin_time_start.substring(0, 5) }} - {{ habit.checkin_time_end.substring(0, 5) }}</span>
                <el-tag size="small" type="success">{{ habit.reward_points }} 积分</el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </div>
      <div v-else class="empty-state">
        <el-empty description="暂无待打卡习惯" />
      </div>
    </el-card>

    <el-card class="stats-card">
      <template #header>
        <h3>统计信息</h3>
      </template>
      <div class="stats-grid">
        <el-statistic :value="todayCheckins" title="今日打卡" />
        <el-statistic :value="totalCheckins" title="累计打卡" />
        <el-statistic :value="currentStreak" title="连续打卡" />
        <el-statistic :value="totalPoints" title="总积分" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import api from '../api'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()
const user = computed(() => userStore.user)

// 登出
const logout = () => {
  userStore.logout()
  router.push('/login')
}

const habits = ref([])
const checkinRecords = ref([])
const todayCheckins = ref(0)
const totalCheckins = ref(0)
const currentStreak = ref(0)
const totalPoints = ref(0)

// 过滤出今天可打卡的习惯并排序（和打卡页保持一致）
const todayHabits = computed(() => {
  const today = dayjs()
  
  // 过滤出今天可打卡的习惯
  const filteredHabits = habits.value.filter(habit => {
    // 周期性习惯：检查是否在允许的星期内
    if (habit.schedule_type === 2) {
      if (!habit.schedule_detail || habit.schedule_detail.length === 0) {
        return false
      }
      
      // schedule_detail 可能是数组也可能是逗号分隔的字符串
      let allowedDays = []
      if (Array.isArray(habit.schedule_detail)) {
        allowedDays = habit.schedule_detail
      } else if (typeof habit.schedule_detail === 'string') {
        allowedDays = habit.schedule_detail.split(',').map(d => parseInt(d))
      }
      
      const todayWeekday = today.day() // 0-6, 0 表示周日
      
      if (!allowedDays.includes(todayWeekday)) {
        return false
      }
    }
    
    return true
  })
  
  // 排序：未打卡的排在前面，已打卡的排在后面
  // 未打卡的按积分倒序排列
  return filteredHabits.sort((a, b) => {
    const aChecked = isCheckedIn(a.id)
    const bChecked = isCheckedIn(b.id)
    
    // 如果一个已打卡一个未打卡，未打卡的排在前面
    if (aChecked !== bChecked) {
      return aChecked ? 1 : -1
    }
    
    // 如果都未打卡，按积分倒序排列
    if (!aChecked && !bChecked) {
      return b.reward_points - a.reward_points
    }
    
    // 如果都已打卡，保持原有顺序
    return 0
  })
})

const loadHabits = async () => {
  try {
    const response = await api.get('/habits')
    habits.value = response.data.habits
  } catch (error) {
    console.error('Failed to load habits:', error)
  }
}

const loadCheckinRecords = async () => {
  try {
    const startDate = dayjs().subtract(30, 'day').format('YYYY-MM-DD')
    const endDate = dayjs().format('YYYY-MM-DD')
    const response = await api.get('/checkin/records', {
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
      api.get(`/checkin/streak/${habit.id}`).then(res => res.data.streak?.current_streak || 0)
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
  max-width: 1200px;
  margin: 0 auto;
}

.user-info-card {
  margin-bottom: 15px;
  padding: 12px;
}

/* 用户信息区域 */
.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.user-avatar {
  flex-shrink: 0;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-details h3 {
  margin: 0 0 5px 0;
  font-size: 20px;
  font-weight: 600;
}

.user-details p {
  margin: 0 0 12px 0;
  color: #666;
}

/* 积分区域 */
.points-section {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.points-badge {
  margin: 0;
}

.points-badge .el-button {
  font-size: 14px;
  padding: 6px 16px;
}

.points-section .el-button {
  font-size: 14px;
  padding: 6px 16px;
}

.user-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
}

.user-actions .el-button {
  font-size: 14px;
  padding: 6px 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .user-avatar {
    align-self: center;
  }
  
  .user-details {
    width: 100%;
    text-align: center;
  }
  
  .points-section {
    justify-content: center;
    width: 100%;
    gap: 8px;
  }
  
  .user-actions {
    width: 100%;
    flex-direction: row;
    justify-content: center;
    gap: 8px;
  }
  
  .user-actions .el-button {
    flex: 1;
    max-width: 120px;
  }
}

@media (max-width: 480px) {
  .user-details h3 {
    font-size: 18px;
  }
  
  .points-section {
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }
  
  .points-section .el-button {
    width: 100%;
    max-width: 200px;
  }
  
  .user-actions {
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }
  
  .user-actions .el-button {
    width: 100%;
    max-width: 200px;
  }
}

/* 习惯卡片区域 */
.habits-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
}

.habits-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.habit-item {
  transition: all 0.3s ease;
}

.habit-item.checked-in {
  opacity: 0.6;
  filter: grayscale(80%);
}

.habit-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.habit-info {
  display: flex;
  align-items: flex-start;
  gap: 15px;
}

.habit-icon {
  flex-shrink: 0;
}

.habit-icon.checked-in {
  opacity: 0.5;
}

.habit-icon img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
}

.default-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
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
  background: linear-gradient(135deg, #999 0%, #666 100%);
}

.habit-details {
  flex: 1;
  min-width: 0;
}

.habit-details h4 {
  margin: 0 0 5px 0;
  font-size: 16px;
}

.habit-details h4.checked-in {
  color: #999;
}

.habit-details p {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
}

.habit-time {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.habit-time span {
  color: #999;
  font-size: 14px;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}

/* 统计信息区域 */
.stats-card {
  margin-bottom: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.el-statistic {
  text-align: center;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
}

.el-statistic__value {
  font-size: 24px;
  font-weight: bold;
  color: #667eea;
}

.el-statistic__label {
  font-size: 14px;
  color: #666;
  margin-top: 10px;
}
</style>
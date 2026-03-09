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
      <div v-if="habits.length > 0" class="habits-list">
        <el-card v-for="habit in habits" :key="habit.id" class="habit-item" :class="{ 'checked-in': isCheckedIn(habit.id) }">
          <div class="habit-info">
            <div class="habit-icon" :class="{ 'checked-in': isCheckedIn(habit.id) }">
              <img v-if="habit.icon" :src="habit.icon" alt="habit icon" />
              <div v-else class="default-icon" :class="{ 'checked-in': isCheckedIn(habit.id) }">{{ habit.name.charAt(0) }}</div>
            </div>
            <div class="habit-details">
              <h4 :class="{ 'checked-in': isCheckedIn(habit.id) }">{{ habit.name }}</h4>
              <p>{{ habit.description }}</p>
              <div class="habit-time">
                <span>{{ habit.checkin_time_start }} - {{ habit.checkin_time_end }}</span>
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
  margin-bottom: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  padding: 20px;
}

.user-details {
  margin-left: 20px;
  flex: 1;
}

.user-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-left: 20px;
}

.user-details h3 {
  margin: 0 0 10px 0;
  font-size: 20px;
}

.user-details p {
  margin: 0 0 15px 0;
  color: #666;
}

.points-badge {
  margin-top: 10px;
}

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

.habit-icon {
  margin-right: 15px;
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
}

.default-icon.checked-in {
  background: linear-gradient(135deg, #999 0%, #666 100%);
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
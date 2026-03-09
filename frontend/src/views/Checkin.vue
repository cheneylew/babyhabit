<template>
  <div class="checkin-container">
    <el-card class="checkin-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button type="text" @click="$router.push('/home')" class="back-btn">
              ← 返回主页
            </el-button>
            <h3>习惯打卡</h3>
          </div>
          <el-date-picker v-model="selectedDate" type="date" placeholder="选择日期" />
        </div>
      </template>
      
      <div v-if="habits.length > 0" class="habits-list">
        <el-card v-for="habit in habits" :key="habit.id" class="habit-item">
          <div class="habit-info">
            <div class="habit-icon">
              <img v-if="habit.icon" :src="habit.icon" alt="habit icon" />
              <div v-else class="default-icon">{{ habit.name.charAt(0) }}</div>
            </div>
            <div class="habit-details">
              <h4>{{ habit.name }}</h4>
              <p>{{ habit.description }}</p>
              <div class="habit-time">
                <span>{{ habit.checkin_time_start }} - {{ habit.checkin_time_end }}</span>
                <el-tag size="small" type="success">{{ habit.reward_points }} 积分</el-tag>
              </div>
            </div>
            <div class="checkin-action">
              <el-button 
                v-if="!isCheckedIn(habit.id)" 
                type="primary" 
                @click="checkin(habit)"
                :loading="checkingIn"
              >
                立即打卡
              </el-button>
              <el-button 
                v-else 
                type="info" 
                plain
                disabled
              >
                已打卡
              </el-button>
            </div>
          </div>
        </el-card>
      </div>
      <div v-else class="empty-state">
        <el-empty description="暂无习惯" />
      </div>
    </el-card>

    <el-card class="history-card">
      <template #header>
        <h3>打卡历史</h3>
      </template>
      <el-table :data="checkinRecords" style="width: 100%">
        <el-table-column prop="habit_id" label="习惯 ID" width="100" />
        <el-table-column prop="checkin_date" label="打卡日期" width="120" />
        <el-table-column prop="checkin_time" label="打卡时间" width="150" />
        <el-table-column prop="checkin_type" label="打卡类型" width="100">
          <template #default="scope">
            {{ scope.row.checkin_type === 1 ? '正常' : '补卡' }}
          </template>
        </el-table-column>
        <el-table-column prop="points_rewarded" label="获得积分" width="100" />
        <el-table-column prop="remark" label="备注" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../store/user'
import api from '../api'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const habits = ref([])
const checkinRecords = ref([])
const selectedDate = ref(dayjs().format('YYYY-MM-DD'))
const checkingIn = ref(false)

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
  const found = checkinRecords.value.some(record => {
    const recordDate = dayjs(record.checkin_date).format('YYYY-MM-DD')
    return record.habit_id === habitId && recordDate === today
  })
  console.log('habitId:', habitId, 'today:', today, 'records:', checkinRecords.value, 'found:', found)
  return found
}

const checkin = async (habit) => {
  checkingIn.value = true
  try {
    await api.post('/checkin', {
      habit_id: habit.id,
      checkin_date: selectedDate.value,
      checkin_type: 1 // 正常打卡
    })
    await loadCheckinRecords()
    // 重新加载用户信息以更新积分
    await userStore.getUserInfo()
    // 重新加载习惯列表以更新状态
    await loadHabits()
    // 显示成功提示
    ElMessage.success('打卡成功，积分已更新')
  } catch (error) {
    console.error('Checkin failed:', error)
    ElMessage.error(error.response?.data?.error || '打卡失败')
  } finally {
    checkingIn.value = false
  }
}

onMounted(async () => {
  await loadHabits()
  await loadCheckinRecords()
})
</script>

<style scoped>
.checkin-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.checkin-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.back-btn {
  font-size: 14px;
  padding: 5px 10px;
  color: #409EFF;
}

.back-btn:hover {
  color: #66b1ff;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
}

.habits-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 15px;
}

.habit-item {
  transition: all 0.3s ease;
}

.habit-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.habit-info {
  display: flex;
  align-items: flex-start;
}

.habit-icon {
  margin-right: 15px;
}

.habit-icon img {
  width: 60px;
  height: 60px;
  border-radius: 50%;
}

.default-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: bold;
}

.habit-details {
  flex: 1;
}

.habit-details h4 {
  margin: 0 0 5px 0;
  font-size: 16px;
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

.checkin-action {
  margin-left: 20px;
  display: flex;
  align-items: center;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}

.history-card {
  margin-bottom: 20px;
}

.el-table {
  margin-top: 20px;
}
</style>
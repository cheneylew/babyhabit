<template>
  <view class="checkin-container">
    <!-- 名言警句悬浮显示 -->
    <view v-if="randomQuote" class="quote-float">
      <view class="quote-content">
        <text class="quote-text">"{{ randomQuote.content }}"</text>
        <text v-if="randomQuote.author" class="quote-author">—— {{ randomQuote.author }}</text>
      </view>
    </view>
    
    <view class="checkin-card">
      <view class="card-header">
        <view class="header-left">
          <button type="text" @click="navigateToHome" class="back-btn">
            ← 返回主页
          </button>
          <h3>习惯打卡</h3>
        </view>
        <picker mode="date" :value="selectedDate" @change="onDateChange">
          <view class="date-picker">
            {{ selectedDate }}
          </view>
        </picker>
      </view>
      
      <view v-if="availableHabits.length > 0" class="habits-list">
        <view v-for="habit in availableHabits" :key="habit.id" class="habit-item">
          <view class="habit-info">
            <view class="habit-icon">
              <image v-if="habit.icon" :src="habit.icon" mode="aspectFill"></image>
              <view v-else class="default-icon">{{ habit.name.charAt(0) }}</view>
            </view>
            <view class="habit-details">
              <h4>{{ habit.name }}</h4>
              <p>{{ habit.description }}</p>
              <view class="habit-time">
                <span>{{ habit.checkin_time_start.substring(0, 5) }} - {{ habit.checkin_time_end.substring(0, 5) }}</span>
                <view class="reward-points">{{ habit.reward_points }} 积分</view>
              </view>
            </view>
            <view class="checkin-action">
              <button 
                v-if="!isCheckedIn(habit.id) && canCheckin(habit) && !isRolledBack(habit.id)" 
                type="primary" 
                @click="showCheckinConfirm(habit)"
                :loading="checkingIn"
              >
                {{ getCheckinButtonText(habit) }}
              </button>
              <button 
                v-else-if="isRolledBack(habit.id)" 
                type="danger" 
                plain
                disabled
              >
                已回退
              </button>
              <button 
                v-else-if="!isCheckedIn(habit.id)" 
                type="info" 
                plain
                disabled
              >
                {{ getCheckinButtonText(habit) }}
              </button>
              <button 
                v-else 
                type="info" 
                plain
                disabled
              >
                已打卡
              </button>
            </view>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text>暂无习惯</text>
      </view>
    </view>

    <view class="history-card">
      <view class="card-header">
        <h3>打卡历史</h3>
      </view>
      <view class="checkin-records">
        <view v-if="checkinRecords.length > 0" class="records-list">
          <view v-for="record in checkinRecords" :key="record.id" class="record-item">
            <view class="record-info">
              <view class="record-habit">{{ record.habit_name }}</view>
              <view class="record-date">{{ dayjs(record.checkin_date).format('YYYY-MM-DD') }} {{ dayjs(record.checkin_time).format('HH:mm:ss') }}</view>
              <view class="record-details">
                <view class="record-type">{{ record.checkin_type === 1 ? '正常' : '补卡' }}</view>
                <view class="record-status" :class="{
                  'rolled-back': record.is_rolled_back === 1,
                  'normal': record.status === 1 && record.is_rolled_back !== 1,
                  'invalid': record.status !== 1 && record.is_rolled_back !== 1
                }">
                  {{ record.is_rolled_back === 1 ? '已回退' : record.status === 1 ? '正常' : '无效' }}
                </view>
                <view class="record-points">+{{ record.points_rewarded }} 积分</view>
              </view>
              <view class="record-photo" v-if="record.photo_url">
                <image :src="record.photo_url" mode="aspectFill" @click="previewImage(record.photo_url)"></image>
              </view>
              <view class="record-rate" v-if="record.self_rate && record.self_rate > 0">
                <text style="color: #FF9900;">{{ '⭐'.repeat(record.self_rate) }}</text>
              </view>
              <view class="record-rate" v-else>
                <text style="color: #999;">-</text>
              </view>
              <view class="record-remark" v-if="record.remark">
                <text>备注: {{ record.remark }}</text>
              </view>
            </view>
          </view>
        </view>
        <view v-else class="empty-records">
          <text>暂无打卡记录</text>
        </view>
      </view>
    </view>

    <!-- 二次确认对话框 -->
    <view v-if="confirmDialogVisible" class="popup-overlay">
      <view class="dialog-content">
        <view class="dialog-header">
          <text>确认打卡</text>
        </view>
        <view v-if="currentHabit?.checkin_prompt" class="checkin-prompt">
          <view class="prompt-content">{{ currentHabit.checkin_prompt }}</view>
        </view>
        <view class="dialog-body">
          <text>确定要打卡习惯 <text class="habit-name">{{ currentHabit?.name }}</text> 吗？</text>
        </view>
        <view class="dialog-footer">
          <button type="default" @click="confirmDialogVisible = false">取消</button>
          <button type="primary" @click="proceedToCheckin" :loading="checkingIn">确认</button>
        </view>
      </view>
    </view>

    <!-- 选择照片对话框 -->
    <view v-if="photoDialogVisible" class="popup-overlay">
      <view class="dialog-content">
        <view class="dialog-header">
          <text>选择照片</text>
        </view>
        <view v-if="currentHabit?.checkin_prompt" class="checkin-prompt">
          <view class="prompt-content">{{ currentHabit.checkin_prompt }}</view>
        </view>
        <view class="photo-container">
          <view class="upload-area" @click="chooseImage" v-if="!photoData">
            <view class="upload-icon">📷</view>
            <text>点击选择照片</text>
          </view>
          <image v-if="photoData" :src="photoData" mode="aspectFill" class="photo-preview"></image>
        </view>
        <view class="dialog-footer">
          <button type="default" @click="photoDialogVisible = false">取消</button>
          <button type="default" @click="photoData = ''" v-if="photoData">重新选择</button>
          <button type="primary" @click="confirmCheckinWithPhoto" v-if="photoData" :loading="checkingIn">确认打卡</button>
        </view>
      </view>
    </view>

    <!-- 自我评分对话框 -->
    <view v-if="rateDialogVisible" class="popup-overlay">
      <view class="dialog-content">
        <view class="dialog-header">
          <text>自我评分</text>
        </view>
        <view class="rate-container">
          <text class="rate-title">请为本次打卡评分（1-10 分）：</text>
          <view class="stars-container">
            <text 
              v-for="i in 10" 
              :key="i" 
              class="star" 
              :class="{ 'active': selfRate >= i }"
              @click="selfRate = i"
            >
              ⭐
            </text>
          </view>
          <text class="rate-text">{{ selfRate }} 分</text>
        </view>
        <view class="dialog-footer">
          <button type="default" @click="rateDialogVisible = false">取消</button>
          <button type="primary" @click="submitSelfRate" :loading="submittingRate">提交</button>
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
const habits = ref([])
const checkinRecords = ref([])
const selectedDate = ref(dayjs().format('YYYY-MM-DD'))
const checkingIn = ref(false)
const confirmDialogVisible = ref(false)
const photoDialogVisible = ref(false)
const rateDialogVisible = ref(false)
const currentHabit = ref(null)
const photoData = ref('')
const selfRate = ref(0)
const submittingRate = ref(false)
const currentRecordId = ref(null)
const randomQuote = ref(null)

// 获取随机名言警句
const loadRandomQuote = async () => {
  try {
    const response = await api.get('/api/quote/random')
    if (response.data.quote) {
      randomQuote.value = response.data.quote
    }
  } catch (error) {
    console.error('Failed to load random quote:', error)
  }
}

// 导航方法
const navigateToHome = () => {
  uni.navigateBack()
}

const loadHabits = async () => {
  try {
    const response = await api.get('/api/habits')
    habits.value = response.data.habits || []
  } catch (error) {
    console.error('Failed to load habits:', error)
    habits.value = []
  }
}

const loadCheckinRecords = async () => {
  try {
    const startDate = dayjs().subtract(30, 'day').format('YYYY-MM-DD')
    const endDate = dayjs().format('YYYY-MM-DD')
    const response = await api.get('/api/checkin/records', {
      params: { 
        start_date: startDate, 
        end_date: endDate 
      }
    })
    checkinRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load checkin records:', error)
  }
}

const performCheckin = async () => {
  const habit = currentHabit.value
  
  // 再次验证日期
  if (!canCheckin(habit)) {
    uni.showToast({ title: '当前日期不允许打卡', icon: 'none' })
    return
  }
  
  checkingIn.value = true
  try {
    // 判断是正常打卡还是补卡
    const selected = dayjs(selectedDate.value)
    const today = dayjs()
    const isMakeup = !selected.isSame(today, 'day')
    
    // 使用 FormData 上传文件
    const formData = new FormData()
    formData.append('habit_id', String(habit.id))
    formData.append('checkin_type', String(isMakeup ? 2 : 1))
    formData.append('checkin_date', selectedDate.value)
    formData.append('remark', '')
    
    // 如果有照片，添加照片文件
    if (photoData.value) {
      // 将本地路径转换为 File 对象
      const file = await pathToFile(photoData.value, 'photo.jpg')
      formData.append('photo', file)
    }
    
    // 发送请求，不设置 Content-Type，让 axios 自动处理
    const response = await api.post('/api/checkin', formData)
    
    uni.showToast({ title: '打卡成功', icon: 'success' })
    photoDialogVisible.value = false
    confirmDialogVisible.value = false
    photoData.value = ''
    
    // 如果需要自我评分，显示评分对话框
    if (habit.allow_self_rate === 1 && response.data.record && response.data.record.id) {
      currentRecordId.value = response.data.record.id
      selfRate.value = 0
      rateDialogVisible.value = true
    }
    
    // 刷新数据
    await loadHabits()
    await loadCheckinRecords()
    await userStore.getUserInfo()
  } catch (error) {
    console.error('Checkin error:', error)
    // 先关闭对话框
    photoDialogVisible.value = false
    confirmDialogVisible.value = false
    photoData.value = ''
    // 优先显示后端返回的错误消息
    let errorMsg = '打卡失败'
    if (error.response?.data?.error) {
      errorMsg = error.response.data.error
    } else if (error.response?.data?.message) {
      errorMsg = error.response.data.message
    } else if (error.message) {
      errorMsg = error.message
    }
    // 短暂延迟后显示 toast，确保对话框已关闭
    setTimeout(() => {
      uni.showToast({ title: errorMsg, icon: 'none', duration: 3000 })
    }, 100)
  } finally {
    checkingIn.value = false
  }
}

// 提交自我评分
const submitSelfRate = async () => {
  if (selfRate.value < 1 || selfRate.value > 10) {
    uni.showToast({ title: '评分必须在 1-10 分之间', icon: 'none' })
    return
  }
  
  submittingRate.value = true
  try {
    await api.post('/api/checkin/submit-rate', {
      record_id: currentRecordId.value,
      self_rate: selfRate.value
    })
    
    uni.showToast({ title: '评分成功', icon: 'success' })
    rateDialogVisible.value = false
    
    // 刷新打卡记录
    await loadCheckinRecords()
  } catch (error) {
    console.error('Submit rate error:', error)
    uni.showToast({ title: error.response?.data?.error || '评分失败', icon: 'none' })
  } finally {
    submittingRate.value = false
  }
}

// 监听日期变化
const onDateChange = (e) => {
  selectedDate.value = dayjs(e.detail.value).format('YYYY-MM-DD')
  loadCheckinRecords()
}

// 过滤当天可打卡的习惯并排序
const availableHabits = computed(() => {
  const selected = dayjs(selectedDate.value)
  
  // 过滤出当天可打卡的习惯
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
      
      const selectedWeekday = selected.day() // 0-6, 0 表示周日
      
      if (!allowedDays.includes(selectedWeekday)) {
        return false
      }
    }
    
    // 其他检查（日期范围等）在 canCheckin 中处理
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

const isCheckedIn = (habitId) => {
  const today = dayjs().format('YYYY-MM-DD')
  const found = checkinRecords.value.some(record => {
    const recordDate = dayjs(record.checkin_date).format('YYYY-MM-DD')
    return record.habit_id === habitId && recordDate === today && record.is_rolled_back !== 1
  })
  return found
}

const isRolledBack = (habitId) => {
  const today = dayjs().format('YYYY-MM-DD')
  const found = checkinRecords.value.some(record => {
    const recordDate = dayjs(record.checkin_date).format('YYYY-MM-DD')
    return record.habit_id === habitId && recordDate === today && record.is_rolled_back === 1
  })
  return found
}

// 检查是否可以打卡
const canCheckin = (habit) => {
  const selected = dayjs(selectedDate.value)
  const today = dayjs()
  const yesterday = today.subtract(1, 'day')
  
  // 明天及之后不允许打卡
  if (selected.isAfter(today, 'day')) {
    return false
  }
  
  // 检查周期性习惯
  if (habit.schedule_type === 2) {
    // 周期性习惯：检查选中的日期是否在允许的星期内
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
    
    const selectedWeekday = selected.day() // 0-6, 0 表示周日
    
    if (!allowedDays.includes(selectedWeekday)) {
      return false
    }
  }
  
  // 今天允许打卡
  if (selected.isSame(today, 'day')) {
    return true
  }
  
  // 昨天只有允许补卡的习惯才能打卡
  if (selected.isSame(yesterday, 'day')) {
    return habit.allow_makeup === 1
  }
  
  // 前天及更早的日期不允许打卡
  return false
}

// 获取打卡按钮文本
const getCheckinButtonText = (habit) => {
  const selected = dayjs(selectedDate.value)
  const today = dayjs()
  const yesterday = today.subtract(1, 'day')
  
  if (selected.isAfter(today, 'day')) {
    return '不可打卡'
  }
  
  if (selected.isSame(yesterday, 'day')) {
    if (habit.allow_makeup === 1) {
      return '补卡'
    } else {
      return '不支持补卡'
    }
  }
  
  if (selected.isBefore(yesterday, 'day')) {
    return '不可打卡'
  }
  
  return '立即打卡'
}

// 显示打卡确认对话框
const showCheckinConfirm = (habit) => {
  currentHabit.value = habit
  
  // 每次点击打卡时重置照片和评分数据
  photoData.value = ''
  selfRate.value = 0
  
  // 如果需要拍照，直接打开选择照片对话框
  if (habit.require_photo === 1) {
    photoDialogVisible.value = true
  } else {
    // 不需要拍照才显示二次确认对话框
    confirmDialogVisible.value = true
  }
}

// 确认打卡，直接执行打卡
const proceedToCheckin = async () => {
  confirmDialogVisible.value = false
  await performCheckin()
}

// 选择照片
const chooseImage = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['original', 'compressed'],
    sourceType: ['camera', 'album'],
    success: (res) => {
      photoData.value = res.tempFilePaths[0]
    },
    fail: (err) => {
      console.error('选择照片失败:', err)
      uni.showToast({ title: '选择照片失败', icon: 'none' })
    }
  })
}

// 确认打卡并上传照片
const confirmCheckinWithPhoto = async () => {
  if (!photoData.value) {
    uni.showToast({ title: '请先选择照片', icon: 'none' })
    return
  }
  
  await performCheckin()
}



// 将本地路径转换为 File 对象
const pathToFile = async (path, filename) => {
  // H5环境下使用fetch获取图片数据
  // uni.chooseImage 在H5返回的是 blob:http:// 或 data:image/ 开头的临时路径
  if (path.startsWith('blob:') || path.startsWith('data:') || path.startsWith('http')) {
    try {
      const response = await fetch(path)
      const blob = await response.blob()
      return new File([blob], filename, { type: blob.type || 'image/jpeg' })
    } catch (err) {
      throw new Error('获取图片失败: ' + err.message)
    }
  }
  
  // 小程序环境下使用文件系统管理器
  // #ifndef H5
  return new Promise((resolve, reject) => {
    const fs = uni.getFileSystemManager()
    if (!fs) {
      reject(new Error('文件系统管理器不可用'))
      return
    }
    fs.readFile({
      filePath: path,
      success: (res) => {
        const file = new File([res.data], filename, { type: 'image/jpeg' })
        resolve(file)
      },
      fail: (err) => {
        reject(new Error('读取文件失败: ' + (err.message || err.errMsg)))
      }
    })
  })
  // #endif
  
  // #ifdef H5
  throw new Error('不支持的图片路径格式')
  // #endif
}



// 预览图片
const previewImage = (url) => {
  uni.previewImage({
    urls: [url],
    current: url
  })
}

onMounted(async () => {
  await loadHabits()
  await loadCheckinRecords()
  await loadRandomQuote()
})
</script>

<style scoped>
.checkin-container {
  padding: 10px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 名言警句悬浮样式 */
.quote-float {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
  max-width: 90%;
  width: fit-content;
}

.quote-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 12px 20px;
  border-radius: 20px;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.quote-text {
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  text-align: center;
  line-height: 1.4;
}

.quote-author {
  color: rgba(255, 255, 255, 0.8);
  font-size: 12px;
}

.checkin-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  margin-top: 70px;
}

.checkin-prompt {
  margin-bottom: 15px;
  padding: 10px;
  background-color: #f0f9ff;
  border-left: 4px solid #409eff;
  border-radius: 4px;
}

.prompt-content {
  color: #303133;
  font-size: 14px;
  line-height: 1.4;
}

/* 卡片头部样式 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  padding: 10px 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
  min-width: 0;
}

.header-left h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.back-btn {
  font-size: 14px;
  padding: 5px 10px;
  color: #409eff;
  background-color: transparent;
  border: none;
}

.date-picker {
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
}

/* 习惯列表样式 */
.habits-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 15px;
  margin-top: 20px;
}

.habit-item {
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.habit-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.habit-info {
  display: flex;
  align-items: flex-start;
  gap: 15px;
}

.habit-icon {
  flex-shrink: 0;
}

.habit-icon image {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
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
  flex-shrink: 0;
}

.habit-details {
  flex: 1;
  min-width: 0;
}

.habit-details h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  white-space: normal;
  line-height: 1.3;
}

.habit-details p {
  margin: 0 0 12px 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.habit-time {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
  flex-wrap: wrap;
  gap: 8px;
}

.habit-time span {
  color: #999;
  font-size: 13px;
  white-space: nowrap;
  flex-shrink: 0;
}

.reward-points {
  background-color: #f0f9eb;
  color: #67c23a;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 14px;
}

.checkin-action {
  margin-left: 10px;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

/* 历史记录样式 */
.history-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.records-list {
  margin-top: 15px;
}

.record-item {
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
  margin-bottom: 10px;
}

.record-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.record-habit {
  font-size: 16px;
  font-weight: 600;
}

.record-date {
  font-size: 14px;
  color: #666;
}

.record-details {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.record-type {
  background-color: #f0f9ff;
  color: #409eff;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.record-status {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.record-status.rolled-back {
  background-color: #fef0f0;
  color: #f56c6c;
}

.record-status.normal {
  background-color: #f0f9eb;
  color: #67c23a;
}

.record-status.invalid {
  background-color: #fafafa;
  color: #909399;
}

.record-points {
  background-color: #fdf6ec;
  color: #e6a23c;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.record-photo {
  margin-top: 8px;
}

.record-photo image {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  object-fit: cover;
}

.record-rate {
  margin-top: 8px;
  font-size: 14px;
}

.record-remark {
  margin-top: 8px;
  font-size: 14px;
  color: #666;
}

.empty-state, .empty-records {
  padding: 40px 0;
  text-align: center;
  color: #999;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .checkin-container {
    padding: 8px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .header-left {
    width: 100%;
    justify-content: space-between;
  }
  
  .header-left h3 {
    font-size: 16px;
    flex: 1;
  }
  
  .habits-list {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .habit-info {
    flex-direction: row;
    align-items: flex-start;
    gap: 12px;
  }
  
  .habit-icon {
    align-self: flex-start;
  }
  
  .habit-details {
    flex: 1;
    min-width: 0;
  }
  
  .checkin-action {
    margin-left: 10px;
    margin-top: 0;
    width: auto;
    flex-shrink: 0;
  }
  
  .checkin-action button {
    width: auto;
    padding: 4px 12px;
    font-size: 12px;
  }
  
  .habit-time {
    flex-direction: row;
    align-items: center;
    gap: 8px;
  }
  
  .habit-time span {
    flex-shrink: 0;
  }
}

@media (max-width: 480px) {
  .header-left {
    gap: 10px;
  }
  
  .back-btn {
    font-size: 12px;
    padding: 4px 8px;
  }
  
  .header-left h3 {
    font-size: 15px;
  }
  
  .habit-details h4 {
    font-size: 15px;
  }
  
  .habit-details p {
    font-size: 13px;
  }
  
  .habit-time span {
    font-size: 12px;
  }
  
  .checkin-action button {
    padding: 3px 10px;
    font-size: 11px;
  }
}

/* 照片预览样式 */
.photo-container {
  margin: 20px 0;
  text-align: center;
}

.upload-area {
  border: 2px dashed #d9d9d9;
  border-radius: 6px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.upload-area:hover {
  border-color: #409eff;
  background-color: #f6f9ff;
}

.upload-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.upload-area text {
  color: #909399;
  font-size: 14px;
}

.photo-preview {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
  border-radius: 4px;
}

/* 对话框样式 */
.popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.dialog-content {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  width: 80%;
  max-width: 400px;
}

.dialog-header {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
}

.dialog-body {
  margin-bottom: 20px;
  line-height: 1.5;
}

.habit-name {
  font-weight: 600;
  color: #409eff;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

/* 评分样式 */
.rate-container {
  margin: 20px 0;
  text-align: center;
}

.rate-title {
  display: block;
  margin-bottom: 20px;
  font-size: 16px;
}

.stars-container {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 15px;
}

.star {
  font-size: 24px;
  cursor: pointer;
  color: #999;
  transition: all 0.3s ease;
}

.star.active {
  color: #FF9900;
  transform: scale(1.2);
}

.rate-text {
  font-size: 16px;
  font-weight: 600;
  color: #FF9900;
}

button {
  border: none;
  border-radius: 4px;
  font-size: 14px;
  padding: 8px 16px;
  cursor: pointer;
}

button[type="primary"] {
  background-color: #409eff;
  color: #fff;
}

button[type="default"] {
  background-color: #f0f0f0;
  color: #333;
}

button[type="danger"] {
  background-color: #f56c6c;
  color: #fff;
}

button[type="info"] {
  background-color: #909399;
  color: #fff;
}

button[plain] {
  background-color: transparent;
  border: 1px solid;
}

button[plain][type="danger"] {
  color: #f56c6c;
  border-color: #f56c6c;
}

button[plain][type="info"] {
  color: #909399;
  border-color: #909399;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
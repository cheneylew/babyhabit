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
          <el-date-picker v-model="selectedDate" type="date" placeholder="选择日期" @change="onDateChange" />
        </div>
      </template>
      
      <!-- 名言警句嵌入显示 -->
      <div v-if="randomQuote" class="quote-embed">
        <div class="quote-content">
          <button class="quote-close" @click="closeQuote">×</button>
          <div class="quote-text-container">
            <span class="quote-text">
              "{{ randomQuote.content }}"
              <span v-if="randomQuote.meaning" class="meaning-toggle" @click="toggleMeaning(randomQuote.id)">
                <span class="meaning-icon">{{ showMeaning[randomQuote.id] ? '✕' : '!' }}</span>
              </span>
            </span>
            <span v-if="randomQuote.meaning && showMeaning[randomQuote.id]" class="quote-meaning">{{ randomQuote.meaning }}</span>
            <span v-if="randomQuote.author" class="quote-author">—— {{ randomQuote.author }}</span>
          </div>
          <div class="quote-buttons">
            <button class="quote-prev" @click="loadPreviousQuote" :disabled="currentQuoteIndex <= 0">上一句</button>
            <button class="quote-next" @click="loadRandomQuote">下一句</button>
          </div>
        </div>
      </div>
      
      <div v-if="availableHabits.length > 0" class="habits-list">
        <el-card v-for="habit in availableHabits" :key="habit.id" class="habit-item">
          <div class="habit-info">
            <div class="habit-icon">
              <img v-if="habit.icon" :src="habit.icon" alt="habit icon" />
              <div v-else class="default-icon">{{ habit.name.charAt(0) }}</div>
            </div>
            <div class="habit-details">
              <h4>{{ habit.name }}</h4>
              <p>{{ habit.description }}</p>
              <div class="habit-time">
                <span>{{ habit.checkin_time_start.substring(0, 5) }} - {{ habit.checkin_time_end.substring(0, 5) }}</span>
                <el-tag size="small" type="success">{{ habit.reward_points }} 积分</el-tag>
              </div>
            </div>
            <div class="checkin-action">
              <el-button 
                v-if="!isCheckedIn(habit.id) && canCheckin(habit) && !isRolledBack(habit.id)" 
                type="primary" 
                @click="showCheckinConfirm(habit)"
                :loading="checkingIn"
              >
                {{ getCheckinButtonText(habit) }}
              </el-button>
              <el-button 
                v-else-if="isRolledBack(habit.id)" 
                type="danger" 
                plain
                disabled
              >
                已回退
              </el-button>
              <el-button 
                v-else-if="!isCheckedIn(habit.id)" 
                type="info" 
                plain
                disabled
              >
                {{ getCheckinButtonText(habit) }}
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
        <el-table-column prop="habit_name" label="习惯名称" width="200" />
        <el-table-column prop="checkin_date" label="打卡日期" width="120">
          <template #default="scope">
            {{ dayjs(scope.row.checkin_date).format('YYYY-MM-DD') }}
          </template>
        </el-table-column>
        <el-table-column prop="checkin_time" label="打卡时间" width="150">
          <template #default="scope">
            {{ dayjs(scope.row.checkin_time).format('HH:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column prop="checkin_type" label="打卡类型" width="100">
          <template #default="scope">
            {{ scope.row.checkin_type === 1 ? '正常' : '补卡' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.is_rolled_back === 1" type="danger">已回退</el-tag>
            <el-tag v-else-if="scope.row.status === 1" type="success">正常</el-tag>
            <el-tag v-else type="info">无效</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="照片" width="100">
          <template #default="scope">
            <el-image 
              v-if="scope.row.photo_url" 
              :src="scope.row.photo_url" 
              :preview-src-list="[scope.row.photo_url]"
              fit="cover"
              style="width: 60px; height: 60px; cursor: pointer;"
            >
              <template #placeholder>
                <div style="width: 60px; height: 60px; background: #f5f7fa; display: flex; align-items: center; justify-content: center;">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <span v-else style="color: #999;">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="self_rate" label="自我评分" width="100">
          <template #default="scope">
            <span v-if="scope.row.self_rate && scope.row.self_rate > 0" style="color: #FF9900;">
              {{ '⭐'.repeat(scope.row.self_rate) }}
            </span>
            <span v-else style="color: #999;">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="points_rewarded" label="获得积分" width="100" />
        <el-table-column prop="remark" label="备注" />
      </el-table>
    </el-card>

    <!-- 二次确认对话框 -->
    <el-dialog v-model="confirmDialogVisible" title="确认打卡" width="80%" :fullscreen="false">
      <div v-if="currentHabit?.checkin_prompt" class="checkin-prompt">
        <el-alert 
          :title="currentHabit.checkin_prompt" 
          type="info" 
          :closable="false" 
          show-icon 
        />
      </div>
      <p style="margin-top: 15px;">确定要打卡习惯 <strong>{{ currentHabit?.name }}</strong> 吗？</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="confirmDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="proceedToCheckin" :loading="checkingIn">确认</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 选择照片对话框 -->
    <el-dialog v-model="photoDialogVisible" title="选择照片" width="80%">
      <div v-if="currentHabit?.checkin_prompt" class="checkin-prompt">
        <el-alert 
          :title="currentHabit.checkin_prompt" 
          type="info" 
          :closable="false" 
          show-icon 
        />
      </div>
      <div class="photo-container">
        <div class="upload-area" @click="triggerFileInput" v-if="!photoData">
          <el-icon class="upload-icon"><Picture /></el-icon>
          <p>点击选择照片</p>
        </div>
        <img v-if="photoData" :src="photoData" alt="选择的照片" class="photo-preview" />
        <input 
          ref="fileInputRef" 
          type="file" 
          accept="image/*" 
          capture="environment"
          style="display: none"
          @change="handleFileChange"
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="photoDialogVisible = false">取消</el-button>
          <el-button @click="selectPhoto" v-if="!photoData">选择照片</el-button>
          <el-button @click="retakePhoto" v-if="photoData">重新选择</el-button>
          <el-button type="primary" @click="confirmCheckinWithPhoto" v-if="photoData" :loading="checkingIn">确认打卡</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 自我评分对话框 -->
    <el-dialog v-model="rateDialogVisible" title="自我评分" width="80%">
      <div class="rate-container">
        <p class="rate-title">请为本次打卡评分（1-10 分）：</p>
        <div class="stars-container">
          <el-rate
            v-model="selfRate"
            :max="10"
            :colors="['#99A9BF', '#F7BA2A', '#FF9900']"
            :low-threshold="3"
            :high-threshold="7"
            show-text
            :texts="['1 分', '2 分', '3 分', '4 分', '5 分', '6 分', '7 分', '8 分', '9 分', '10 分']"
          />
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="rateDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitSelfRate" :loading="submittingRate">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick } from 'vue'
import { useUserStore } from '../store/user'
import api from '../api'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { Picture } from '@element-plus/icons-vue'

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
const fileInputRef = ref(null)
const selfRate = ref(0)
const submittingRate = ref(false)
const currentRecordId = ref(null)
const randomQuote = ref(null)
const quoteHistory = ref([])
const currentQuoteIndex = ref(-1)
const showMeaning = ref({})

const loadHabits = async () => {
  try {
    const response = await api.get('/habits')
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
    const response = await api.get('/checkin/records', {
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

// 获取随机名言警句
const loadRandomQuote = async () => {
  try {
    // 检查是否有下一句（在历史记录中）
    if (currentQuoteIndex.value < quoteHistory.value.length - 1) {
      // 有下一句，使用历史记录
      currentQuoteIndex.value++
      randomQuote.value = quoteHistory.value[currentQuoteIndex.value]
    } else {
      // 没有下一句，从API获取
      const response = await api.get('/quote/random')
      if (response.data.quote) {
        // 添加到历史记录
        quoteHistory.value.push(response.data.quote)
        // 限制历史记录长度，最多保存20条
        if (quoteHistory.value.length > 20) {
          quoteHistory.value.shift()
          // 调整索引
          if (currentQuoteIndex.value > 0) {
            currentQuoteIndex.value--
          }
        }
        currentQuoteIndex.value = quoteHistory.value.length - 1
        randomQuote.value = response.data.quote
      }
    }
  } catch (error) {
    console.error('Failed to load random quote:', error)
  }
}

// 关闭名言警句
const closeQuote = () => {
  randomQuote.value = null
  quoteHistory.value = []
  currentQuoteIndex.value = -1
  showMeaning.value = {}
}

const toggleMeaning = (quoteId) => {
  showMeaning.value = {
    ...showMeaning.value,
    [quoteId]: !showMeaning.value[quoteId]
  }
}

// 加载上一句名言
const loadPreviousQuote = () => {
  if (currentQuoteIndex.value > 0) {
    currentQuoteIndex.value--
    randomQuote.value = quoteHistory.value[currentQuoteIndex.value]
  }
}

// 监听日期变化
const onDateChange = () => {
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
  console.log('habitId:', habitId, 'today:', today, 'records:', checkinRecords.value, 'found:', found)
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

// 触发文件选择
const triggerFileInput = () => {
  if (fileInputRef.value) {
    fileInputRef.value.click()
  }
}

// 选择照片
const selectPhoto = () => {
  triggerFileInput()
}

// 处理文件选择
const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }
  
  // 读取文件
  const reader = new FileReader()
  reader.onload = (e) => {
    photoData.value = e.target.result
  }
  reader.onerror = () => {
    ElMessage.error('读取文件失败')
  }
  reader.readAsDataURL(file)
  
  // 清空 input，允许重复选择同一文件
  event.target.value = ''
}

// 重新选择照片
const retakePhoto = () => {
  photoData.value = ''
  triggerFileInput()
}

// 确认打卡并上传照片
const confirmCheckinWithPhoto = async () => {
  if (!photoData.value) {
    ElMessage.warning('请先选择照片')
    return
  }
  
  await performCheckin()
}

// 执行打卡
const performCheckin = async () => {
  const habit = currentHabit.value
  
  // 再次验证日期
  if (!canCheckin(habit)) {
    ElMessage.error('当前日期不允许打卡')
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
      // 将 Base64 转换为 File 对象
      const file = await base64ToFile(photoData.value, 'photo.jpg')
      formData.append('photo', file)
    }
    
    // 发送请求，不设置 Content-Type，让 axios 自动处理
    const response = await api.post('/checkin', formData)
    
    console.log('=== Checkin Response ===')
    console.log('Response data:', response.data)
    console.log('Record:', response.data.record)
    console.log('Record ID:', response.data.record?.id)
    console.log('========================')
    
    ElMessage.success('打卡成功')
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
    ElMessage.error(error.response?.data?.error || '打卡失败')
  } finally {
    checkingIn.value = false
  }
}

// 将 Base64 转换为 File 对象
const base64ToFile = (base64, filename) => {
  return new Promise((resolve, reject) => {
    fetch(base64)
      .then(res => res.blob())
      .then(blob => {
        const file = new File([blob], filename, { type: 'image/jpeg' })
        resolve(file)
      })
      .catch(reject)
  })
}

// 提交自我评分
const submitSelfRate = async () => {
  if (selfRate.value < 1 || selfRate.value > 10) {
    ElMessage.error('评分必须在 1-10 分之间')
    return
  }
  
  submittingRate.value = true
  try {
    await api.post('/checkin/submit-rate', {
      record_id: currentRecordId.value,
      self_rate: selfRate.value
    })
    
    ElMessage.success('评分成功')
    rateDialogVisible.value = false
    
    // 刷新打卡记录
    await loadCheckinRecords()
  } catch (error) {
    console.error('Submit rate error:', error)
    ElMessage.error(error.response?.data?.error || '评分失败')
  } finally {
    submittingRate.value = false
  }
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
  position: relative;
}

/* 名言警句嵌入样式 */
.quote-embed {
  width: 100%;
  margin-bottom: 20px;
}

.quote-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 12px 20px;
  border-radius: 20px;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  position: relative;
  width: 100%;
  box-sizing: border-box;
  max-width: 1200px;
  margin: 0 auto;
}

.quote-text-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  text-align: center;
  width: 100%;
}

.quote-buttons {
  display: flex;
  justify-content: space-between;
  width: 100%;
  margin-top: 4px;
}

.quote-close {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  border: none;
  color: #667eea;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  padding: 0;
  line-height: 1;
}

.quote-close:hover {
  background: #f0f0f0;
  transform: scale(1.1);
  transition: all 0.2s ease;
}

.quote-prev,
.quote-next {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  font-size: 10px;
  padding: 3px 8px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 60px;
  text-align: center;
}

.quote-prev:hover,
.quote-next:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.05);
}

.quote-prev:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quote-prev:disabled:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: none;
}

.quote-text {
  color: white;
  font-size: 16px;
  font-style: italic;
  text-align: center;
  line-height: 1.4;
  font-weight: 500;
  position: relative;
}

.meaning-toggle {
  display: inline-block;
  margin-left: 8px;
  cursor: pointer;
  vertical-align: middle;
}

.meaning-icon {
  display: inline-block;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.3);
  color: white;
  font-size: 12px;
  font-weight: bold;
  line-height: 20px;
  text-align: center;
  transition: all 0.3s ease;
}

.meaning-icon:hover {
  background-color: rgba(255, 255, 255, 0.5);
  transform: scale(1.1);
}

.quote-meaning {
  color: rgba(255, 255, 255, 0.9);
  font-size: 16px;
  margin-top: 6px;
  text-align: center;
  line-height: 1.3;
  font-weight: 400;
  max-width: 90%;
}

.quote-author {
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  margin-top: 6px;
  text-align: center;
  font-weight: 400;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .quote-float {
    top: 10px;
  }
  
  .quote-content {
    padding: 12px 18px;
  }
  
  .quote-text {
    font-size: 15px;
  }
  
  .quote-meaning {
    font-size: 12px;
  }
  
  .quote-author {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .quote-text {
    font-size: 14px;
  }
  
  .quote-meaning {
    font-size: 11px;
  }
  
  .quote-author {
    font-size: 12px;
  }
  
  .quote-content {
    padding: 14px 20px;
  }
}

.checkin-card {
  margin-bottom: 20px;
  border-radius: 12px;
  overflow: hidden;
}

.checkin-prompt {
  margin-bottom: 15px;
}

.checkin-prompt .el-alert {
  background-color: #f0f9ff;
}

/* 提高图片预览的 z-index，确保显示在表格之上 */
:deep(.el-image-viewer) {
  z-index: 2000 !important;
}

:deep(.el-image) {
  position: relative;
  z-index: 1000;
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
  color: #409EFF;
}

/* 习惯列表样式 */
.habits-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 15px;
  margin-top: 20px;
}

.habit-item {
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

.habit-icon img {
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

.checkin-action {
  margin-left: 10px;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

/* 打卡历史样式 */
.checkin-records {
  margin-top: 30px;
}

.checkin-records .el-card {
  border-radius: 12px;
  overflow: hidden;
}

.checkin-records .el-table {
  margin-top: 15px;
  border-radius: 8px;
  overflow: hidden;
}

.checkin-records .el-table th {
  background-color: #f5f7fa;
  font-weight: 600;
  font-size: 14px;
}

.checkin-records .el-table td {
  font-size: 14px;
  padding: 10px;
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
  
  .checkin-action .el-button {
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
  
  .checkin-records .el-table {
    font-size: 13px;
  }
  
  .checkin-records .el-table th,
  .checkin-records .el-table td {
    padding: 8px;
  }
  
  /* 隐藏表格中在手机端显示不下的列 */
  .checkin-records .el-table .el-table__column--width-fixed {
    display: none;
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
  
  .checkin-action .el-button {
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
  border-color: #409EFF;
  background-color: #f6f9ff;
}

.upload-icon {
  font-size: 48px;
  color: #909399;
  margin-bottom: 16px;
}

.upload-area p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.photo-preview {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
  border-radius: 4px;
}

/* 确认对话框样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
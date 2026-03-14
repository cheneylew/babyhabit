<template>
  <view class="edit-habit-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <h3>{{ editingHabit ? '编辑习惯' : '添加习惯' }}</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="edit-habit-content">
      <view class="form-container">
        <view class="form-item">
          <label>习惯名称</label>
          <input type="text" v-model="habitForm.name" placeholder="请输入习惯名称" />
        </view>
        <view class="form-item">
          <label>描述</label>
          <textarea v-model="habitForm.description" placeholder="请输入习惯描述"></textarea>
        </view>
        <view class="form-item">
          <label>图标</label>
          <input type="text" v-model="habitForm.icon" placeholder="请输入图标" />
        </view>
        <view class="form-item">
          <label>分类</label>
          <input type="text" v-model="habitForm.category" placeholder="请输入分类" />
        </view>
        <view class="form-item">
          <label>打卡类型</label>
          <view class="radio-group">
            <label><radio :checked="habitForm.schedule_type === '1'" @click="habitForm.schedule_type = '1'" /> 每日</label>
            <label><radio :checked="habitForm.schedule_type === '2'" @click="habitForm.schedule_type = '2'" /> 周期性</label>
          </view>
        </view>
        <view class="form-item" v-if="habitForm.schedule_type === '2'">
          <label>周期选择</label>
          <view class="checkbox-group">
            <label><checkbox v-model="habitForm.schedule_detail" :value="1" /> 周一</label>
            <label><checkbox v-model="habitForm.schedule_detail" :value="2" /> 周二</label>
            <label><checkbox v-model="habitForm.schedule_detail" :value="3" /> 周三</label>
            <label><checkbox v-model="habitForm.schedule_detail" :value="4" /> 周四</label>
            <label><checkbox v-model="habitForm.schedule_detail" :value="5" /> 周五</label>
            <label><checkbox v-model="habitForm.schedule_detail" :value="6" /> 周六</label>
            <label><checkbox v-model="habitForm.schedule_detail" :value="0" /> 周日</label>
          </view>
        </view>
        <view class="form-item">
          <label>打卡开始时间</label>
          <picker mode="time" :value="habitForm.checkin_time_start" @change="onStartTimeChange">
            <view class="picker-view">
              {{ habitForm.checkin_time_start || '选择开始时间' }}
            </view>
          </picker>
        </view>
        <view class="form-item">
          <label>打卡结束时间</label>
          <picker mode="time" :value="habitForm.checkin_time_end" @change="onEndTimeChange">
            <view class="picker-view">
              {{ habitForm.checkin_time_end || '选择结束时间' }}
            </view>
          </picker>
        </view>
        <view class="form-item">
          <label>奖励积分</label>
          <input type="number" v-model="habitForm.reward_points" placeholder="请输入奖励积分" />
        </view>
        <view class="form-item">
          <label>允许补卡</label>
          <view class="radio-group">
            <label><radio :checked="habitForm.allow_makeup === '1'" @click="habitForm.allow_makeup = '1'" /> 是</label>
            <label><radio :checked="habitForm.allow_makeup === '0'" @click="habitForm.allow_makeup = '0'" /> 否</label>
          </view>
        </view>
        <view class="form-item" v-if="habitForm.allow_makeup === '1'">
          <label>补卡天数</label>
          <input type="number" v-model="habitForm.makeup_days" placeholder="请输入补卡天数" />
        </view>
        <view class="form-item">
          <label>需要拍照</label>
          <view class="radio-group">
            <label><radio :checked="habitForm.require_photo === '1'" @click="habitForm.require_photo = '1'" /> 是</label>
            <label><radio :checked="habitForm.require_photo === '0'" @click="habitForm.require_photo = '0'" /> 否</label>
          </view>
        </view>
        <view class="form-item">
          <label>允许自我评分</label>
          <view class="radio-group">
            <label><radio :checked="habitForm.allow_self_rate === '1'" @click="habitForm.allow_self_rate = '1'" /> 是</label>
            <label><radio :checked="habitForm.allow_self_rate === '0'" @click="habitForm.allow_self_rate = '0'" /> 否</label>
          </view>
        </view>
        <view class="form-item">
          <label>打卡提示</label>
          <textarea v-model="habitForm.checkin_prompt" placeholder="打卡时显示给孩子的提示内容（可选）"></textarea>
        </view>
        <view class="form-item">
          <label>状态</label>
          <view class="radio-group">
            <label><radio :checked="habitForm.status === '1'" @click="habitForm.status = '1'" /> 启用</label>
            <label><radio :checked="habitForm.status === '0'" @click="habitForm.status = '0'" /> 禁用</label>
          </view>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="action-buttons">
        <button type="default" @click="navigateBack" class="cancel-button">取消</button>
        <button type="primary" @click="saveHabit" class="save-button">保存</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'

const userStore = useUserStore()

// 习惯相关
const editingHabit = ref(null)
const habitForm = ref({
  name: '',
  description: '',
  icon: '',
  category: '',
  schedule_type: '1',
  schedule_detail: [],
  checkin_time_start: '',
  checkin_time_end: '',
  reward_points: 0,
  allow_makeup: '0',
  makeup_days: 0,
  require_photo: '0',
  allow_self_rate: '0',
  checkin_prompt: '',
  status: '1'
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

// 时间选择事件处理
const onStartTimeChange = (e) => {
  habitForm.value.checkin_time_start = e.detail.value
}

const onEndTimeChange = (e) => {
  habitForm.value.checkin_time_end = e.detail.value
}

// 保存习惯
const saveHabit = async () => {
  try {
    // 确保数值类型正确
    const formData = { ...habitForm.value }
    formData.schedule_type = parseInt(formData.schedule_type)
    formData.reward_points = parseInt(formData.reward_points)
    formData.allow_makeup = parseInt(formData.allow_makeup)
    formData.makeup_days = parseInt(formData.makeup_days)
    formData.require_photo = parseInt(formData.require_photo)
    formData.allow_self_rate = parseInt(formData.allow_self_rate)
    formData.status = parseInt(formData.status)
    
    // 如果是周期性习惯，将数组转换为逗号分隔的字符串
    if (formData.schedule_type === 2 && Array.isArray(formData.schedule_detail)) {
      formData.schedule_detail = formData.schedule_detail.join(',')
    } else {
      // 每日习惯，schedule_detail 设为空字符串
      formData.schedule_detail = ''
    }
    
    if (editingHabit.value) {
      // 更新习惯
      await api.put(`/api/admin/habits/${editingHabit.value.id}`, formData)
    } else {
      // 添加习惯
      await api.post('/api/admin/habits', formData)
    }
    uni.showToast({ title: '保存成功', icon: 'success' })
    uni.navigateBack()
  } catch (error) {
    console.error('Failed to save habit:', error)
    uni.showToast({ title: '保存失败', icon: 'none' })
  }
}

onMounted(() => {
  const params = getRouteParams()
  if (params.habitId) {
    // 编辑模式，加载习惯信息
    const loadHabitInfo = async () => {
      try {
        const response = await api.get(`/api/admin/habits/${params.habitId}`)
        editingHabit.value = response.data.habit
        
        // 确保所有字段都有默认值，数值转换为字符串
        const habit = response.data.habit
        const requirePhoto = (habit.require_photo !== undefined && habit.require_photo !== null) ? habit.require_photo : 0
        const allowSelfRate = (habit.allow_self_rate !== undefined && habit.allow_self_rate !== null) ? habit.allow_self_rate : 0
        const allowMakeup = (habit.allow_makeup !== undefined && habit.allow_makeup !== null) ? habit.allow_makeup : 0
        const status = (habit.status !== undefined && habit.status !== null) ? habit.status : 1
        const scheduleType = habit.schedule_type || 1
        
        habitForm.value = {
          name: habit.name || '',
          description: habit.description || '',
          icon: habit.icon || '',
          category: habit.category || '',
          schedule_type: String(scheduleType),
          schedule_detail: [],
          checkin_time_start: habit.checkin_time_start || '',
          checkin_time_end: habit.checkin_time_end || '',
          reward_points: habit.reward_points || 0,
          allow_makeup: String(allowMakeup),
          makeup_days: habit.makeup_days || 0,
          require_photo: String(requirePhoto),
          allow_self_rate: String(allowSelfRate),
          checkin_prompt: habit.checkin_prompt || '',
          status: String(status)
        }
        
        // 如果是周期性习惯，将逗号分隔的字符串转换为数字数组
        if (habit.schedule_type === 2 && habit.schedule_detail) {
          habitForm.value.schedule_detail = habit.schedule_detail.split(',').map(item => parseInt(item.trim()))
        }
      } catch (error) {
        console.error('Failed to load habit info:', error)
        uni.showToast({ title: '加载失败', icon: 'none' })
      }
    }
    loadHabitInfo()
  }
})
</script>

<style scoped>
.edit-habit-container {
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

.edit-habit-content {
  padding-top: 20px;
}

.form-container {
  background-color: #f9f9f9;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
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

.form-item input,
.form-item textarea,
.picker-view {
  width: 100%;
  padding: 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  font-size: 14px;
  background-color: #fff;
}

.picker-view {
  color: #606266;
  cursor: pointer;
}

.form-item textarea {
  min-height: 100px;
  resize: vertical;
}

.radio-group {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  padding: 10px 0;
}

.checkbox-group {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
  padding: 10px 0;
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

.cancel-button:hover,
.save-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .edit-habit-container {
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
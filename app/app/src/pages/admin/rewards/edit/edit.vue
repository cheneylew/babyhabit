<template>
  <view class="edit-reward-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>{{ editingReward ? '编辑奖励' : '添加奖励' }}</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="edit-reward-content">
      <view class="form-container">
        <view class="form-item">
          <label>奖励名称</label>
          <input type="text" v-model="rewardForm.name" placeholder="请输入奖励名称" />
        </view>
        <view class="form-item">
          <label>描述</label>
          <textarea v-model="rewardForm.description" placeholder="请输入奖励描述"></textarea>
        </view>
        <view class="form-item">
          <label>图片</label>
          <input type="text" v-model="rewardForm.image" placeholder="请输入图片" />
        </view>
        <view class="form-item">
          <label>所需积分</label>
          <input type="number" v-model="rewardForm.points_required" placeholder="请输入所需积分" />
        </view>
        <view class="form-item">
          <label>库存</label>
          <input type="number" v-model="rewardForm.stock" placeholder="-1表示无限" />
        </view>
        <view class="form-item">
          <label>状态</label>
          <view class="radio-group">
            <label><radio v-model="rewardForm.status" :value="1" /> 启用</label>
            <label><radio v-model="rewardForm.status" :value="0" /> 禁用</label>
          </view>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="action-buttons">
        <button type="default" @click="navigateBack" class="cancel-button">取消</button>
        <button type="primary" @click="saveReward" class="save-button">保存</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'

const userStore = useUserStore()

// 奖励相关
const editingReward = ref(null)
const rewardForm = ref({
  name: '',
  description: '',
  image: '',
  points_required: 0,
  stock: -1,
  status: 1
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

// 保存奖励
const saveReward = async () => {
  try {
    // 确保数值类型正确
    const formData = { ...rewardForm.value }
    formData.points_required = parseInt(formData.points_required)
    formData.stock = parseInt(formData.stock)
    formData.status = parseInt(formData.status)
    
    // 如果 start_time 和 end_time 为空，则不传递
    if (!formData.start_time) delete formData.start_time
    if (!formData.end_time) delete formData.end_time
    
    if (editingReward.value) {
      // 更新奖励
      await api.put(`/api/admin/rewards/${editingReward.value.id}`, formData)
    } else {
      // 添加奖励
      await api.post('/api/admin/rewards', formData)
    }
    uni.showToast({ title: '保存成功', icon: 'success' })
    uni.navigateBack()
  } catch (error) {
    console.error('Failed to save reward:', error)
    uni.showToast({ title: '保存失败', icon: 'none' })
  }
}

onMounted(() => {
  const params = getRouteParams()
  if (params.rewardId) {
    // 编辑模式，加载奖励信息
    const loadRewardInfo = async () => {
      try {
        const response = await api.get(`/api/admin/rewards/${params.rewardId}`)
        editingReward.value = response.data.item
        rewardForm.value = { ...response.data.item }
      } catch (error) {
        console.error('Failed to load reward info:', error)
        uni.showToast({ title: '加载失败', icon: 'none' })
      }
    }
    loadRewardInfo()
  }
})
</script>

<style scoped>
.edit-reward-container {
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

.edit-reward-content {
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
.form-item textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  font-size: 14px;
  background-color: #fff;
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
  .edit-reward-container {
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
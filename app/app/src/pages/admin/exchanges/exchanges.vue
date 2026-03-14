<template>
  <view class="exchanges-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>兑换管理</h3>
      </view>
      <view class="header-right">
        <text class="user-name">{{ userStore.user?.name }}</text>
        <button type="primary" plain @click="logout" class="logout-button">
          🚪
        </button>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="exchanges-content">
      <view class="filter-section">
        <picker @change="handleStatusChange" :value="statusFilter" :range="statusOptions" range-key="label">
          <view class="picker">
            {{ statusOptions[statusFilter].label }}
          </view>
        </picker>
        <button type="default" @click="loadExchanges" class="refresh-button">
          刷新
        </button>
      </view>
      <view v-if="exchanges.length > 0" class="exchanges-list">
        <view v-for="exchange in exchanges" :key="exchange.id" class="exchange-item">
          <view class="exchange-header">
            <view class="exchange-child">{{ exchange.item?.description }}</view>
            <view class="exchange-status" :class="getExchangeStatusClass(exchange.status)">
              {{ getExchangeStatusText(exchange.status) }}
            </view>
          </view>
          <view class="exchange-details">
            <text class="detail-item">奖励: {{ exchange.item?.name }}</text>
            <text class="detail-item">积分: {{ exchange.points }}</text>
            <text class="detail-item">兑换时间: {{ exchange.exchange_time }}</text>
          </view>
          <view class="exchange-actions">
            <button 
              type="primary" 
              size="small" 
              @click="handleExchange(exchange)"
              :disabled="exchange.status !== 0"
              class="action-button primary"
            >
              {{ exchange.status === 0 ? '处理' : '已处理' }}
            </button>
            <button 
              type="danger" 
              size="small" 
              @click="cancelExchange(exchange)"
              :disabled="exchange.status !== 0"
              class="action-button danger"
            >
              取消
            </button>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">📦</text>
        <text>暂无兑换记录</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../store/user'
import api from '../../../api'

const userStore = useUserStore()

// 兑换记录相关
const exchanges = ref([])
const statusFilter = ref(0)
const statusOptions = [
  { value: 0, label: '全部' },
  { value: 1, label: '待处理' },
  { value: 2, label: '已完成' },
  { value: 3, label: '已取消' }
]

// 导航返回
const navigateBack = () => {
  uni.navigateBack()
}

// 加载兑换记录
const loadExchanges = async () => {
  try {
    const params = {}
    if (statusFilter.value > 0) {
      params.status = statusFilter.value - 1
    }
    const response = await api.get('/api/admin/exchanges', { params })
    exchanges.value = response.data.records || []
  } catch (error) {
    console.error('Failed to load exchanges:', error)
  }
}

// 处理状态变化
const handleStatusChange = (e) => {
  statusFilter.value = e.detail.value
  loadExchanges()
}

// 获取兑换状态类名
const getExchangeStatusClass = (status) => {
  switch (status) {
    case 0:
      return 'status-pending'
    case 1:
      return 'status-active'
    case 2:
      return 'status-canceled'
    default:
      return ''
  }
}

// 获取兑换状态文本
const getExchangeStatusText = (status) => {
  switch (status) {
    case 0:
      return '待处理'
    case 1:
      return '已完成'
    case 2:
      return '已取消'
    default:
      return '未知'
  }
}

// 处理兑换
const handleExchange = async (exchange) => {
  try {
    await uni.showModal({
      title: '处理兑换',
      content: '确定要处理该兑换请求吗？',
      success: async (res) => {
        if (res.confirm) {
          await api.put(`/api/admin/exchanges/${exchange.id}`, {
            status: 1
          })
          uni.showToast({ title: '处理成功', icon: 'success' })
          await loadExchanges()
        }
      }
    })
  } catch (error) {
    console.error('Failed to handle exchange:', error)
    uni.showToast({ title: '处理失败', icon: 'none' })
  }
}

// 取消兑换
const cancelExchange = async (exchange) => {
  try {
    await uni.showModal({
      title: '取消兑换',
      content: '确定要取消该兑换请求吗？',
      success: async (res) => {
        if (res.confirm) {
          await api.put(`/api/admin/exchanges/${exchange.id}`, {
            status: 2
          })
          uni.showToast({ title: '取消成功', icon: 'success' })
          await loadExchanges()
        }
      }
    })
  } catch (error) {
    console.error('Failed to cancel exchange:', error)
    uni.showToast({ title: '取消失败', icon: 'none' })
  }
}

// 登出
const logout = () => {
  userStore.logout()
  uni.redirectTo({ url: '/pages/login/login' })
}

onMounted(async () => {
  await loadExchanges()
})
</script>

<style scoped>
.exchanges-container {
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

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-name {
  font-size: 14px;
  color: #606266;
}

.logout-button {
  font-size: 16px;
  padding: 6px 10px;
}

.exchanges-content {
  padding-top: 20px;
}

.filter-section {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  align-items: center;
}

.picker {
  padding: 10px 15px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
  flex: 1;
  max-width: 200px;
}

.refresh-button {
  padding: 10px 20px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
  background-color: #f9f9f9;
  transition: all 0.2s ease;
}

.refresh-button:hover {
  background-color: #f0f0f0;
  transform: translateY(-1px);
}

.exchanges-list {
  margin-top: 20px;
}

.exchange-item {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 12px;
  margin-bottom: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.exchange-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.exchange-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.exchange-child {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.exchange-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 15px;
}

.detail-item {
  font-size: 14px;
  color: #666;
}

.exchange-actions {
  display: flex;
  gap: 10px;
}

.action-button {
  font-size: 12px;
  padding: 6px 16px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-button:hover {
  transform: translateY(-1px);
}

.action-button.primary {
  background-color: #409eff;
  color: #fff;
}

.action-button.danger {
  background-color: #f56c6c;
  color: #fff;
}

.exchange-status {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-active {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-pending {
  background-color: #ecf5ff;
  color: #409eff;
}

.status-canceled {
  background-color: #f5f7fa;
  color: #909399;
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

button {
  border: none;
  border-radius: 4px;
  font-size: 14px;
  padding: 8px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

button:hover {
  transform: translateY(-1px);
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

button[size="small"] {
  font-size: 12px;
  padding: 6px 12px;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .exchanges-container {
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
  
  .header-right {
    width: 100%;
    justify-content: space-between;
  }
  
  .filter-section {
    flex-direction: column;
    align-items: stretch;
  }
  
  .picker {
    max-width: none;
  }
  
  .exchange-actions {
    justify-content: space-between;
  }
  
  .action-button {
    flex: 1;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .exchange-item {
    padding: 15px;
  }
  
  .exchange-child {
    font-size: 16px;
  }
  
  .detail-item {
    font-size: 13px;
  }
  
  .action-button {
    font-size: 11px;
    padding: 5px 10px;
  }
}
</style>
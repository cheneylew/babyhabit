<template>
  <view class="rewards-container">
    <view class="rewards-card">
      <view class="card-header">
        <h3>奖励中心</h3>
        <view class="points-badge">
          <button type="primary" size="small">我的积分 {{ user?.points_balance || 0 }}</button>
        </view>
      </view>
      
      <view v-if="rewards.length > 0" class="rewards-grid">
        <view v-for="reward in rewards" :key="reward.id" class="reward-item">
          <view class="reward-image">
            <image v-if="reward.image" :src="reward.image" mode="aspectFill"></image>
            <view v-else class="default-image">{{ reward.name.charAt(0) }}</view>
          </view>
          <view class="reward-info">
            <h4>{{ reward.name }}</h4>
            <p>{{ reward.description }}</p>
            <view class="reward-price">
              <view class="price-tag">{{ reward.points_required }} 积分</view>
              <button 
                type="primary" 
                @click="exchange(reward)"
                :disabled="!canExchange(reward)"
              >
                {{ canExchange(reward) ? '立即兑换' : '积分不足' }}
              </button>
            </view>
            <view class="reward-stock" v-if="reward.stock > 0 && reward.stock !== -1">
              库存: {{ reward.stock }}
            </view>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text>暂无奖励物品</text>
      </view>
    </view>

    <view class="exchange-history-card">
      <view class="card-header">
        <h3>兑换历史</h3>
      </view>
      <view class="exchange-records">
        <view v-if="exchangeRecords.length > 0" class="records-list">
          <view v-for="record in exchangeRecords" :key="record.id" class="record-item">
            <view class="record-info">
              <view class="record-name">{{ record.item.name }}</view>
              <view class="record-details">
                <view class="record-quantity">数量: {{ record.quantity }}</view>
                <view class="record-points">消耗积分: {{ record.points }}</view>
                <view class="record-time">{{ record.exchange_time }}</view>
                <view class="record-status" :class="getStatusClass(record.status)">
                  {{ getStatusText(record.status) }}
                </view>
              </view>
            </view>
          </view>
        </view>
        <view v-else class="empty-records">
          <text>暂无兑换记录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../../store/user'
import api from '../../api'

const userStore = useUserStore()
const user = computed(() => userStore.user)

const rewards = ref([])
const exchangeRecords = ref([])

const loadRewards = async () => {
  try {
    const response = await api.get('/api/rewards')
    rewards.value = response.data.items
  } catch (error) {
    console.error('Failed to load rewards:', error)
  }
}

const loadExchangeRecords = async () => {
  try {
    const response = await api.get('/api/exchange/records')
    exchangeRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load exchange records:', error)
  }
}

const canExchange = (reward) => {
  return (user.value?.points_balance || 0) >= reward.points_required && 
         (reward.stock === -1 || reward.stock > 0)
}

const exchange = async (reward) => {
  try {
    await api.post('/api/exchange', {
      item_id: reward.id,
      quantity: 1
    })
    // 重新加载数据
    await loadRewards()
    await loadExchangeRecords()
    await userStore.getUserInfo()
    uni.showToast({ title: '兑换成功', icon: 'success' })
  } catch (error) {
    console.error('Exchange failed:', error)
    uni.showToast({ title: error.response?.data?.error || '兑换失败', icon: 'none' })
  }
}

const getStatusClass = (status) => {
  switch (status) {
    case 1: return 'status-success'
    case 2: return 'status-warning'
    case 3: return 'status-info'
    case 4: return 'status-success'
    default: return 'status-default'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 1: return '已完成'
    case 2: return '处理中'
    case 3: return '已发货'
    case 4: return '已收货'
    default: return '未知'
  }
}

onMounted(async () => {
  await loadRewards()
  await loadExchangeRecords()
})
</script>

<style scoped>
.rewards-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.rewards-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.points-badge {
  margin-left: auto;
}

.rewards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.reward-item {
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.reward-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.reward-image {
  text-align: center;
  margin-bottom: 15px;
}

.reward-image image {
  width: 150px;
  height: 150px;
  object-fit: cover;
  border-radius: 8px;
}

.default-image {
  width: 150px;
  height: 150px;
  border-radius: 8px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: bold;
  margin: 0 auto;
}

.reward-info {
  text-align: center;
}

.reward-info h4 {
  margin: 0 0 10px 0;
  font-size: 16px;
  font-weight: 600;
}

.reward-info p {
  margin: 0 0 15px 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
  height: 40px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.reward-price {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.price-tag {
  background-color: #fef0f0;
  color: #f56c6c;
  padding: 4px 12px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
}

.reward-stock {
  margin-top: 10px;
  font-size: 14px;
  color: #999;
}

.empty-state, .empty-records {
  padding: 60px 0;
  text-align: center;
  color: #999;
}

.exchange-history-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.records-list {
  margin-top: 20px;
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

.record-name {
  font-size: 16px;
  font-weight: 600;
}

.record-details {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  font-size: 14px;
  color: #666;
}

.record-status {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.status-success {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-warning {
  background-color: #fdf6ec;
  color: #e6a23c;
}

.status-info {
  background-color: #f0f9ff;
  color: #409eff;
}

.status-default {
  background-color: #fafafa;
  color: #909399;
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

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

button[size="small"] {
  font-size: 12px;
  padding: 6px 12px;
}
</style>
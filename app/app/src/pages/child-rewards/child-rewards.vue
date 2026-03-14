<template>
  <view class="rewards-container">
    <view class="header-card">
      <view class="user-info">
        <h3>积分商城</h3>
        <view class="points-tag">当前积分：{{ user?.points_balance || 0 }}</view>
      </view>
    </view>

    <view class="rewards-card">
      <view class="card-header">
        <h3>可兑换奖品</h3>
      </view>
      <view v-if="rewards.length > 0" class="rewards-list">
        <view v-for="reward in rewards" :key="reward.id" class="reward-item">
          <view class="reward-info">
            <view class="reward-image">
              <image v-if="reward.image" :src="reward.image" mode="aspectFill"></image>
              <view v-else class="default-image">🎁</view>
            </view>
            <view class="reward-details">
              <h4>{{ reward.name }}</h4>
              <p>{{ reward.description }}</p>
              <view class="reward-meta">
                <view class="price-tag">{{ reward.points_required }} 积分</view>
                <view class="stock" v-if="reward.stock !== -1">库存：{{ reward.stock }}</view>
                <view class="stock" v-else>无限库存</view>
              </view>
            </view>
            <view class="reward-action">
              <button 
                :class="['exchange-btn', { 'insufficient': user?.points_balance < reward.points_required, 'soldout': reward.stock === 0 }]"
                @click="exchange(reward)"
                :disabled="user?.points_balance < reward.points_required || reward.stock === 0"
              >
                {{ user?.points_balance < reward.points_required ? '积分不足' : (reward.stock === 0 ? '已兑完' : '立即兑换') }}
              </button>
            </view>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text>暂无奖品</text>
      </view>
    </view>

    <view class="records-card">
      <view class="card-header">
        <h3>兑换记录</h3>
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

    <!-- 兑换确认对话框 -->
    <view v-if="exchangeDialogVisible" class="dialog-overlay" @click="exchangeDialogVisible = false">
      <view class="dialog-content" @click.stop>
        <view class="dialog-header">
          <text>确认兑换</text>
        </view>
        <view class="dialog-body">
          <text>确定要兑换 <text class="reward-name">{{ currentReward?.name }}</text> 吗？</text>
          <text class="points-info">需要消耗：<text style="color: #f56c6c;">{{ currentReward?.points_required }} 积分</text></text>
          <text class="current-points">当前积分：<text style="color: #67c23a;">{{ user?.points_balance || 0 }} 积分</text></text>
        </view>
        <view class="dialog-footer">
          <button type="default" @click="exchangeDialogVisible = false">取消</button>
          <button type="primary" @click="confirmExchange" :loading="exchanging">确认兑换</button>
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
const exchangeDialogVisible = ref(false)
const currentReward = ref(null)
const exchanging = ref(false)

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

const exchange = (reward) => {
  currentReward.value = reward
  exchangeDialogVisible.value = true
}

const confirmExchange = async () => {
  if (!currentReward.value) return
  
  try {
    exchanging.value = true
    await api.post('/api/exchange', {
      item_id: currentReward.value.id,
      quantity: 1
    })
    uni.showToast({ title: '兑换成功，请等待管理员核销', icon: 'success' })
    // 重新加载数据
    await loadRewards()
    await loadExchangeRecords()
    await userStore.getUserInfo()
    exchangeDialogVisible.value = false
  } catch (error) {
    console.error('Exchange failed:', error)
    uni.showToast({ title: error.response?.data?.error || '兑换失败', icon: 'none' })
  } finally {
    exchanging.value = false
  }
}

const getStatusClass = (status) => {
  const classes = {
    1: 'status-success',    // 审批通过
    2: 'status-warning'     // 待审批
  }
  return classes[status] || 'status-info'
}

const getStatusText = (status) => {
  const texts = {
    1: '审批通过',
    2: '待审批'
  }
  return texts[status] || '未知'
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

.header-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.user-info h3 {
  margin: 0;
  flex: 1;
  min-width: 0;
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.points-tag {
  font-size: 14px;
  padding: 6px 12px;
  background-color: #f0f9eb;
  color: #67c23a;
  border-radius: 10px;
  flex-shrink: 0;
}

/* 奖励列表样式 */
.rewards-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  margin-bottom: 20px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.rewards-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 15px;
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
}

.reward-info {
  display: flex;
  align-items: flex-start;
  gap: 15px;
}

.reward-image {
  flex-shrink: 0;
  margin-right: 15px;
}

.reward-image image {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
}

.default-image {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
}

.reward-details {
  flex: 1;
  min-width: 0;
}

.reward-details h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
}

.reward-details p {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
}

.reward-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.price-tag {
  background-color: #fdf6ec;
  color: #e6a23c;
  padding: 4px 12px;
  border-radius: 10px;
  font-size: 14px;
}

.stock {
  color: #999;
  font-size: 13px;
}

.reward-action {
  margin-left: 20px;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.empty-state, .empty-records {
  padding: 40px 0;
  text-align: center;
  color: #999;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .user-info h3 {
    width: 100%;
    font-size: 16px;
  }
  
  .points-tag {
    font-size: 13px;
    padding: 5px 10px;
  }
  
  .rewards-list {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .reward-info {
    flex-direction: column;
    gap: 12px;
  }
  
  .reward-image {
    margin-right: 0;
    align-self: flex-start;
  }
  
  .reward-details {
    width: 100%;
  }
  
  .reward-action {
    margin-left: 0;
    margin-top: 10px;
    width: 100%;
  }
  
  .reward-action button {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .user-info h3 {
    font-size: 15px;
  }
  
  .points-tag {
    font-size: 12px;
    padding: 4px 8px;
  }
  
  .reward-details h4 {
    font-size: 15px;
  }
  
  .reward-details p {
    font-size: 13px;
  }
}

.records-card {
  margin-top: 20px;
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

/* 确认对话框样式 */
.dialog-overlay {
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
  width: 80vw;
  max-width: 400px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  animation: dialogFadeIn 0.2s ease;
}

@keyframes dialogFadeIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
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

.reward-name {
  font-weight: 600;
  color: #409eff;
}

.points-info {
  display: block;
  margin: 15px 0;
  font-size: 14px;
}

.current-points {
  display: block;
  margin: 10px 0 15px 0;
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
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

button:disabled {
  cursor: not-allowed;
}

/* 兑换按钮样式 */
.exchange-btn {
  background-color: #409eff;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  padding: 8px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.exchange-btn.insufficient,
.exchange-btn.soldout {
  background-color: #c0c4cc;
  color: #fff;
}

.exchange-btn:disabled {
  background-color: #c0c4cc;
  color: #fff;
  opacity: 1;
}
</style>
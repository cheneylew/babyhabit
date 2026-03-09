<template>
  <div class="rewards-container">
    <el-card class="rewards-card">
      <template #header>
        <div class="card-header">
          <h3>奖励中心</h3>
          <el-badge :value="user?.points_balance || 0" type="success" class="points-badge">
            <el-button type="primary" size="small">我的积分</el-button>
          </el-badge>
        </div>
      </template>
      
      <div v-if="rewards.length > 0" class="rewards-grid">
        <el-card v-for="reward in rewards" :key="reward.id" class="reward-item">
          <div class="reward-image">
            <img v-if="reward.image" :src="reward.image" alt="reward image" />
            <div v-else class="default-image">{{ reward.name.charAt(0) }}</div>
          </div>
          <div class="reward-info">
            <h4>{{ reward.name }}</h4>
            <p>{{ reward.description }}</p>
            <div class="reward-price">
              <el-tag type="danger" size="large">{{ reward.points_required }} 积分</el-tag>
              <el-button 
                type="primary" 
                @click="exchange(reward)"
                :disabled="!canExchange(reward)"
              >
                {{ canExchange(reward) ? '立即兑换' : '积分不足' }}
              </el-button>
            </div>
            <div class="reward-stock" v-if="reward.stock > 0 && reward.stock !== -1">
              库存: {{ reward.stock }}
            </div>
          </div>
        </el-card>
      </div>
      <div v-else class="empty-state">
        <el-empty description="暂无奖励物品" />
      </div>
    </el-card>

    <el-card class="exchange-history-card">
      <template #header>
        <h3>兑换历史</h3>
      </template>
      <el-table :data="exchangeRecords" style="width: 100%">
        <el-table-column prop="item.name" label="物品名称" width="200" />
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column prop="points" label="消耗积分" width="120" />
        <el-table-column prop="exchange_time" label="兑换时间" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../store/user'
import api from '../api'

const userStore = useUserStore()
const user = computed(() => userStore.user)

const rewards = ref([])
const exchangeRecords = ref([])

const loadRewards = async () => {
  try {
    const response = await api.get('/rewards')
    rewards.value = response.data.items
  } catch (error) {
    console.error('Failed to load rewards:', error)
  }
}

const loadExchangeRecords = async () => {
  try {
    const response = await api.get('/exchange/records')
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
    await api.post('/exchange', {
      item_id: reward.id,
      quantity: 1
    })
    // 重新加载数据
    await loadRewards()
    await loadExchangeRecords()
    await userStore.getUserInfo()
  } catch (error) {
    console.error('Exchange failed:', error)
  }
}

const getStatusType = (status) => {
  switch (status) {
    case 1: return 'success'
    case 2: return 'warning'
    case 3: return 'info'
    case 4: return 'success'
    default: return 'default'
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

.points-badge {
  margin-left: auto;
}

.rewards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.reward-item {
  transition: all 0.3s ease;
  cursor: pointer;
}

.reward-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.reward-image {
  text-align: center;
  margin-bottom: 15px;
}

.reward-image img {
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

.reward-stock {
  margin-top: 10px;
  font-size: 14px;
  color: #999;
}

.empty-state {
  padding: 60px 0;
  text-align: center;
}

.exchange-history-card {
  margin-bottom: 20px;
}

.el-table {
  margin-top: 20px;
}
</style>
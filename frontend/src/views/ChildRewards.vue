<template>
  <div class="rewards-container">
    <el-card class="header-card">
      <div class="user-info">
        <el-button type="text" @click="$router.push('/home')" class="back-btn">
          ← 返回主页
        </el-button>
        <h3>积分商城</h3>
        <el-tag type="success" class="points-tag">当前积分：{{ user?.points_balance || 0 }}</el-tag>
      </div>
    </el-card>

    <el-card class="rewards-card">
      <template #header>
        <h3>可兑换奖品</h3>
      </template>
      <div v-if="rewards.length > 0" class="rewards-list">
        <el-card v-for="reward in rewards" :key="reward.id" class="reward-item">
          <div class="reward-info">
            <div class="reward-image">
              <img v-if="reward.image" :src="reward.image" alt="reward" />
              <div v-else class="default-image">🎁</div>
            </div>
            <div class="reward-details">
              <h4>{{ reward.name }}</h4>
              <p>{{ reward.description }}</p>
              <div class="reward-meta">
                <el-tag type="warning" size="small">{{ reward.points_required }} 积分</el-tag>
                <span class="stock" v-if="reward.stock !== -1">库存：{{ reward.stock }}</span>
                <span class="stock" v-else>无限库存</span>
              </div>
            </div>
            <div class="reward-action">
              <el-button 
                type="primary" 
                @click="exchange(reward)"
                :disabled="user?.points_balance < reward.points_required || reward.stock === 0"
              >
                {{ user?.points_balance < reward.points_required ? '积分不足' : (reward.stock === 0 ? '已兑完' : '立即兑换') }}
              </el-button>
            </div>
          </div>
        </el-card>
      </div>
      <div v-else class="empty-state">
        <el-empty description="暂无奖品" />
      </div>
    </el-card>

    <el-card class="records-card">
      <template #header>
        <h3>兑换记录</h3>
      </template>
      <el-table :data="exchangeRecords" style="width: 100%">
        <el-table-column prop="item.name" label="奖品名称" width="150" />
        <el-table-column prop="quantity" label="数量" width="80" />
        <el-table-column prop="points" label="消耗积分" width="100" />
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

  <!-- 兑换确认对话框 -->
  <el-dialog v-model="exchangeDialogVisible" title="确认兑换" width="80%">
    <p>确定要兑换 <strong>{{ currentReward?.name }}</strong> 吗？</p>
    <p class="points-info">需要消耗：<span style="color: #f56c6c;">{{ currentReward?.points_required }} 积分</span></p>
    <p class="current-points">当前积分：<span style="color: #67c23a;">{{ user?.points_balance || 0 }} 积分</span></p>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="exchangeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmExchange" :loading="exchanging">确认兑换</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../store/user'
import api from '../api'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const user = computed(() => userStore.user)
const rewards = ref([])
const exchangeRecords = ref([])
const exchangeDialogVisible = ref(false)
const currentReward = ref(null)
const exchanging = ref(false)

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

const exchange = (reward) => {
  currentReward.value = reward
  exchangeDialogVisible.value = true
}

const confirmExchange = async () => {
  if (!currentReward.value) return
  
  try {
    exchanging.value = true
    await api.post('/exchange', {
      item_id: currentReward.value.id,
      quantity: 1
    })
    ElMessage.success('兑换成功，请等待管理员核销')
    // 重新加载数据
    await loadRewards()
    await loadExchangeRecords()
    await userStore.getUserInfo()
    exchangeDialogVisible.value = false
  } catch (error) {
    console.error('Exchange failed:', error)
    ElMessage.error(error.response?.data?.error || '兑换失败')
  } finally {
    exchanging.value = false
  }
}

const getStatusType = (status) => {
  const types = {
    1: 'success',    // 审批通过
    2: 'warning'     // 待审批
  }
  return types[status] || 'info'
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
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.back-btn {
  font-size: 14px;
  color: #409EFF;
  flex-shrink: 0;
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
  flex-shrink: 0;
}

/* 奖励列表样式 */
.rewards-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 15px;
}

.reward-item {
  transition: all 0.3s ease;
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

.reward-image img {
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

.empty-state {
  padding: 40px 0;
  text-align: center;
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
  
  .reward-action .el-button {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .back-btn {
    font-size: 12px;
  }
  
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
}

.el-table {
  margin-top: 15px;
}

/* 确认对话框样式 */
.points-info {
  margin: 15px 0;
  font-size: 14px;
}

.current-points {
  margin: 10px 0 15px 0;
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>

<template>
  <view class="exchange-container">
    <view class="exchange-card">
      <view class="card-header">
        <h3>兑换记录</h3>
      </view>
      <view class="exchange-records">
        <view v-if="exchangeRecords.length > 0" class="records-list">
          <view v-for="record in exchangeRecords" :key="record.id" class="record-item">
            <view class="record-info">
              <view class="record-id">ID: {{ record.id }}</view>
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
import { ref, onMounted } from 'vue'
import api from '../../api'

const exchangeRecords = ref([])

const loadExchangeRecords = async () => {
  try {
    const response = await api.get('/api/exchange/records')
    exchangeRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load exchange records:', error)
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
  await loadExchangeRecords()
})
</script>

<style scoped>
.exchange-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.exchange-card {
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

.record-id {
  font-size: 12px;
  color: #999;
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

.empty-records {
  padding: 60px 0;
  text-align: center;
  color: #999;
}
</style>
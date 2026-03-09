<template>
  <div class="exchange-container">
    <el-card class="exchange-card">
      <template #header>
        <h3>兑换记录</h3>
      </template>
      <el-table :data="exchangeRecords" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
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
import { ref, onMounted } from 'vue'
import api from '../api'

const exchangeRecords = ref([])

const loadExchangeRecords = async () => {
  try {
    const response = await api.get('/exchange/records')
    exchangeRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load exchange records:', error)
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
}

.el-table {
  margin-top: 20px;
}
</style>
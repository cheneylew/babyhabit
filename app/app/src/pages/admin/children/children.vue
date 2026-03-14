<template>
  <view class="children-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-right">
        <text class="user-name">{{ userStore.user?.name }}</text>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="children-content">
      <button type="primary" @click="navigateToAddChild" class="add-button">
        <text class="add-icon">+</text>
        添加小孩
      </button>
      <view v-if="children.length > 0" class="children-list">
        <view v-for="child in children" :key="child.id" class="child-item">
          <view class="child-info">
            <view class="child-header">
              <view class="child-name">{{ child.name }}</view>
              <view class="child-status" :class="child.status === 1 ? 'status-active' : 'status-inactive'">
                {{ child.status === 1 ? '正常' : '禁用' }}
              </view>
            </view>
            <view class="child-details">
              <text class="detail-item">用户名: {{ child.username }}</text>
              <text class="detail-item">ID: {{ child.id }}</text>
            </view>
          </view>
          <view class="child-actions">
            <button type="primary" size="small" @click="navigateToEditChild(child)" class="action-button primary">编辑</button>
            <button type="danger" size="small" @click="deleteChild(child.id)" class="action-button danger">删除</button>
            <button type="warning" size="small" @click="navigateToAssignHabit(child)" class="action-button warning">分配习惯</button>
            <button type="info" size="small" @click="navigateToAssignedHabits(child)" class="action-button info">查看习惯</button>
            <button type="success" size="small" @click="navigateToCheckinRecords(child)" class="action-button success">查看打卡</button>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text class="empty-icon">👶</text>
        <text>暂无小孩账号</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../store/user'
import api from '../../../api'

const userStore = useUserStore()

// 小孩相关
const children = ref([])

// 导航返回
const navigateBack = () => {
  uni.navigateBack()
}

// 导航到添加小孩页面
const navigateToAddChild = () => {
  uni.navigateTo({ url: '/pages/admin/children/edit/edit' })
}

// 导航到编辑小孩页面
const navigateToEditChild = (child) => {
  uni.navigateTo({ url: `/pages/admin/children/edit/edit?childId=${child.id}` })
}

// 导航到分配习惯页面
const navigateToAssignHabit = (child) => {
  uni.navigateTo({ url: `/pages/admin/children/assign-habit/assign-habit?childId=${child.id}` })
}

// 导航到已分配习惯页面
const navigateToAssignedHabits = (child) => {
  uni.navigateTo({ url: `/pages/admin/children/assigned-habits/assigned-habits?childId=${child.id}` })
}

// 导航到打卡记录页面
const navigateToCheckinRecords = (child) => {
  uni.navigateTo({ url: `/pages/admin/children/checkin-records/checkin-records?childId=${child.id}` })
}

// 加载小孩列表
const loadChildren = async () => {
  try {
    const response = await api.get('/api/admin/children')
    children.value = response.data.children || []
  } catch (error) {
    console.error('Failed to load children:', error)
  }
}

// 删除小孩
const deleteChild = async (id) => {
  try {
    await uni.showModal({
      title: '删除小孩',
      content: '确定要删除这个小孩账号吗？',
      success: async (res) => {
        if (res.confirm) {
          await api.delete(`/api/admin/children/${id}`)
          await loadChildren()
          uni.showToast({ title: '删除成功', icon: 'success' })
        }
      }
    })
  } catch (error) {
    console.error('Failed to delete child:', error)
    uni.showToast({ title: '删除失败', icon: 'none' })
  }
}

// 登出
const logout = () => {
  userStore.logout()
  uni.redirectTo({ url: '/pages/login/login' })
}

onMounted(async () => {
  await loadChildren()
})
</script>

<style scoped>
.children-container {
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

.children-content {
  padding-top: 20px;
}

.add-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #409eff 0%, #667eea 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transition: all 0.3s ease;
}

.add-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.4);
}

.add-icon {
  font-size: 20px;
  font-weight: bold;
}

.children-list {
  margin-top: 20px;
}

.child-item {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 12px;
  margin-bottom: 15px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.child-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.child-info {
  flex: 1;
  min-width: 0;
}

.child-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.child-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.child-details {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.detail-item {
  font-size: 14px;
  color: #666;
}

.child-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
}

.action-button {
  font-size: 12px;
  padding: 6px 12px;
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

.action-button.warning {
  background-color: #e6a23c;
  color: #fff;
}

.action-button.info {
  background-color: #909399;
  color: #fff;
}

.action-button.success {
  background-color: #67c23a;
  color: #fff;
}

.child-status {
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

.status-inactive {
  background-color: #fef0f0;
  color: #f56c6c;
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

/* 对话框样式 */
.dialog-content {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  width: 80%;
  max-width: 500px;
}

.dialog-header {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
}

.dialog-body {
  margin-bottom: 20px;
}

.form-item {
  margin-bottom: 15px;
}

.form-item label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.form-item input, .form-item textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
}

.form-item textarea {
  min-height: 80px;
  resize: vertical;
}

.radio-group {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.picker {
  padding: 10px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

/* 已分配习惯和打卡记录列表样式 */
.assigned-habits-list, .checkin-records-list {
  max-height: 400px;
  overflow-y: auto;
}

.assigned-habit-item, .checkin-record-item {
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.assigned-habit-name, .checkin-record-habit {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.assigned-habit-points, .checkin-record-points {
  font-size: 14px;
  color: #67c23a;
  font-weight: 500;
}

.assigned-habit-time, .checkin-record-date, .checkin-record-type {
  font-size: 14px;
  color: #666;
}

.checkin-record-status {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  align-self: flex-start;
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

button[type="warning"] {
  background-color: #e6a23c;
  color: #fff;
}

button[type="info"] {
  background-color: #909399;
  color: #fff;
}

button[type="success"] {
  background-color: #67c23a;
  color: #fff;
}

button[plain] {
  background-color: transparent;
  border: 1px solid;
}

button[plain][type="primary"] {
  color: #409eff;
  border-color: #409eff;
}

button[size="small"] {
  font-size: 12px;
  padding: 6px 12px;
  align-self: flex-start;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .children-container {
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
  
  .child-item {
    flex-direction: column;
    align-items: stretch;
  }
  
  .child-actions {
    flex-direction: row;
    flex-wrap: wrap;
    margin-top: 15px;
  }
  
  .action-button {
    flex: 1;
    min-width: 80px;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .add-button {
    padding: 10px 16px;
    font-size: 14px;
  }
  
  .child-item {
    padding: 15px;
  }
  
  .child-name {
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
<template>
  <view class="edit-child-container">
    <!-- 顶部导航栏 -->
    <view class="admin-header">
      <view class="header-left">
        <button type="text" @click="navigateBack" class="back-btn">
          ← 返回
        </button>
        <h3>{{ editingChild ? '编辑小孩' : '添加小孩' }}</h3>
      </view>
    </view>

    <!-- 内容区域 -->
    <view class="edit-child-content">
      <view class="form-container">
        <view class="form-item">
          <label>用户名</label>
          <input type="text" v-model="childForm.username" :disabled="editingChild" placeholder="请输入用户名" />
        </view>
        <view class="form-item" v-if="!editingChild">
          <label>密码</label>
          <input type="password" v-model="childForm.password" placeholder="请输入密码" />
        </view>
        <view class="form-item">
          <label>姓名</label>
          <input type="text" v-model="childForm.name" placeholder="请输入姓名" />
        </view>
        <view class="form-item">
          <label>手机号</label>
          <input type="text" v-model="childForm.phone" placeholder="请输入手机号" />
        </view>
        <view class="form-item">
          <label>邮箱</label>
          <input type="email" v-model="childForm.email" placeholder="请输入邮箱" />
        </view>
        <view class="form-item">
          <label>状态</label>
          <view class="radio-group">
            <label><radio v-model="childForm.status" :value="1" /> 正常</label>
            <label><radio v-model="childForm.status" :value="0" /> 禁用</label>
          </view>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="action-buttons">
        <button type="default" @click="navigateBack" class="cancel-button">取消</button>
        <button type="primary" @click="saveChild" class="save-button">保存</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../../../store/user'
import api from '../../../../api'

const userStore = useUserStore()

// 小孩相关
const editingChild = ref(null)
const childForm = ref({
  username: '',
  password: '',
  name: '',
  phone: '',
  email: '',
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

// 保存小孩
const saveChild = async () => {
  try {
    // 如果 phone 和 email 为空，则不传递
    const formData = { ...childForm.value }
    if (!formData.phone) delete formData.phone
    if (!formData.email) delete formData.email
    
    if (editingChild.value) {
      // 更新小孩信息
      await api.put(`/api/admin/children/${editingChild.value.id}`, formData)
    } else {
      // 添加小孩
      await api.post('/api/admin/children', formData)
    }
    uni.showToast({ title: '保存成功', icon: 'success' })
    uni.navigateBack()
  } catch (error) {
    console.error('Failed to save child:', error)
    uni.showToast({ title: '保存失败', icon: 'none' })
  }
}

onMounted(() => {
  const params = getRouteParams()
  if (params.childId) {
    // 编辑模式，加载小孩信息
    const loadChildInfo = async () => {
      try {
        const response = await api.get(`/api/admin/children/${params.childId}`)
        editingChild.value = response.data.child
        childForm.value = { ...response.data.child }
      } catch (error) {
        console.error('Failed to load child info:', error)
        uni.showToast({ title: '加载失败', icon: 'none' })
      }
    }
    loadChildInfo()
  }
})
</script>

<style scoped>
.edit-child-container {
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

.edit-child-content {
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

.form-item input {
  width: 100%;
  padding: 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  font-size: 14px;
  background-color: #fff;
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
  .edit-child-container {
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
<template>
  <div class="admin-container">
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <h3>管理员控制台</h3>
          <div class="user-info">
            <span>{{ userStore.user?.name }}</span>
            <el-button type="primary" plain @click="logout">登出</el-button>
          </div>
        </div>
      </template>
      
      <el-tabs v-model="activeTab">
        <!-- 小孩账号管理 -->
        <el-tab-pane label="小孩管理" name="children">
          <el-button type="primary" @click="childDialogVisible = true">添加小孩</el-button>
          <el-table :data="children" style="width: 100%" class="children-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="name" label="姓名" width="120" />
            <el-table-column prop="phone" label="手机号" width="120" />
            <el-table-column prop="email" label="邮箱" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                  {{ scope.row.status === 1 ? '正常' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="250">
              <template #default="scope">
                <el-button type="primary" size="small" @click="editChild(scope.row)">编辑</el-button>
                <el-button type="danger" size="small" @click="deleteChild(scope.row.id)">删除</el-button>
                <el-button type="warning" size="small" @click="assignHabit(scope.row)">分配习惯</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 习惯管理 -->
        <el-tab-pane label="习惯管理" name="habits">
          <el-button type="primary" @click="dialogVisible = true">添加习惯</el-button>
          <el-table :data="habits" style="width: 100%" class="habits-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="习惯名称" width="150" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="reward_points" label="奖励积分" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                  {{ scope.row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button type="primary" size="small" @click="editHabit(scope.row)">编辑</el-button>
                <el-button type="danger" size="small" @click="deleteHabit(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 奖励管理 -->
        <el-tab-pane label="奖励管理" name="rewards">
          <el-button type="primary" @click="rewardDialogVisible = true">添加奖励</el-button>
          <el-table :data="rewards" style="width: 100%" class="rewards-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="奖励名称" width="150" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="points_required" label="所需积分" width="100" />
            <el-table-column prop="stock" label="库存" width="100">
              <template #default="scope">
                {{ scope.row.stock === -1 ? '无限' : scope.row.stock }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                  {{ scope.row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button type="primary" size="small" @click="editReward(scope.row)">编辑</el-button>
                <el-button type="danger" size="small" @click="deleteReward(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 兑换管理 -->
        <el-tab-pane label="兑换管理" name="exchanges">
          <el-table :data="exchangeRecords" style="width: 100%" class="exchanges-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="user_id" label="用户ID" width="100" />
            <el-table-column prop="item.name" label="物品名称" width="150" />
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
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button type="primary" size="small" @click="updateExchangeStatus(scope.row)">更新状态</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 小孩编辑对话框 -->
    <el-dialog v-model="childDialogVisible" :title="editingChild ? '编辑小孩' : '添加小孩'">
      <el-form :model="childForm" :rules="childRules" ref="childFormRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="childForm.username" :disabled="editingChild" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!editingChild">
          <el-input type="password" v-model="childForm.password" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="childForm.name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="childForm.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="childForm.email" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="childForm.status">
            <el-radio label="1">正常</el-radio>
            <el-radio label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="childDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveChild">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 习惯编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="editingHabit ? '编辑习惯' : '添加习惯'">
      <el-form :model="habitForm" :rules="habitRules" ref="habitFormRef" label-width="100px">
        <el-form-item label="习惯名称" prop="name">
          <el-input v-model="habitForm.name" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input type="textarea" v-model="habitForm.description" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="habitForm.icon" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="habitForm.category" />
        </el-form-item>
        <el-form-item label="打卡类型" prop="schedule_type">
          <el-radio-group v-model="habitForm.schedule_type">
            <el-radio label="1">每日</el-radio>
            <el-radio label="2">周期性</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="打卡开始时间" prop="checkin_time_start">
          <el-time-picker v-model="habitForm.checkin_time_start" format="HH:mm" value-format="HH:mm:ss" placeholder="选择开始时间" />
        </el-form-item>
        <el-form-item label="打卡结束时间" prop="checkin_time_end">
          <el-time-picker v-model="habitForm.checkin_time_end" format="HH:mm" value-format="HH:mm:ss" placeholder="选择结束时间" />
        </el-form-item>
        <el-form-item label="奖励积分" prop="reward_points">
          <el-input type="number" v-model="habitForm.reward_points" />
        </el-form-item>
        <el-form-item label="允许补卡" prop="allow_makeup">
          <el-radio-group v-model="habitForm.allow_makeup">
            <el-radio label="1">是</el-radio>
            <el-radio label="0">否</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="habitForm.status">
            <el-radio label="1">启用</el-radio>
            <el-radio label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveHabit">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 奖励编辑对话框 -->
    <el-dialog v-model="rewardDialogVisible" :title="editingReward ? '编辑奖励' : '添加奖励'">
      <el-form :model="rewardForm" :rules="rewardRules" ref="rewardFormRef" label-width="100px">
        <el-form-item label="奖励名称" prop="name">
          <el-input v-model="rewardForm.name" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input type="textarea" v-model="rewardForm.description" />
        </el-form-item>
        <el-form-item label="图片" prop="image">
          <el-input v-model="rewardForm.image" />
        </el-form-item>
        <el-form-item label="所需积分" prop="points_required">
          <el-input type="number" v-model="rewardForm.points_required" />
        </el-form-item>
        <el-form-item label="库存" prop="stock">
          <el-input type="number" v-model="rewardForm.stock" placeholder="-1表示无限" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="rewardForm.status">
            <el-radio label="1">启用</el-radio>
            <el-radio label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="rewardDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveReward">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 兑换状态更新对话框 -->
    <el-dialog v-model="statusDialogVisible" title="更新兑换状态">
      <el-form :model="statusForm" :rules="statusRules" ref="statusFormRef" label-width="100px">
        <el-form-item label="状态" prop="status">
          <el-select v-model="statusForm.status" placeholder="选择状态">
            <el-option label="已完成" value="1" />
            <el-option label="处理中" value="2" />
            <el-option label="已发货" value="3" />
            <el-option label="已收货" value="4" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="statusDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveExchangeStatus">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分配习惯对话框 -->
    <el-dialog v-model="assignHabitDialogVisible" title="分配习惯">
      <el-form :model="assignHabitForm" :rules="assignHabitRules" ref="assignHabitFormRef" label-width="100px">
        <el-form-item label="选择习惯" prop="habit_id">
          <el-select v-model="assignHabitForm.habit_id" placeholder="选择习惯" style="width: 100%">
            <el-option
              v-for="habit in habits"
              :key="habit.id"
              :label="habit.name"
              :value="habit.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="assignHabitDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveAssignHabit">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import api from '../api'

const router = useRouter()

const userStore = useUserStore()
const activeTab = ref('children')

// 小孩相关
const children = ref([])
const childDialogVisible = ref(false)
const editingChild = ref(null)
const childForm = ref({
  username: '',
  password: '',
  name: '',
  phone: '',
  email: '',
  status: 1
})
const childRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}
const childFormRef = ref(null)

// 习惯相关
const habits = ref([])
const dialogVisible = ref(false)
const editingHabit = ref(null)
const habitForm = ref({
  name: '',
  description: '',
  icon: '',
  category: '',
  schedule_type: 1,
  schedule_detail: '',
  checkin_time_start: '',
  checkin_time_end: '',
  reward_points: 0,
  allow_makeup: 0,
  makeup_days: 0,
  status: 1
})
const habitRules = {
  name: [{ required: true, message: '请输入习惯名称', trigger: 'blur' }],
  schedule_type: [{ required: true, message: '请选择打卡类型', trigger: 'change' }],
  reward_points: [{ required: true, message: '请输入奖励积分', trigger: 'blur' }]
}
const habitFormRef = ref(null)

// 奖励相关
const rewards = ref([])
const rewardDialogVisible = ref(false)
const editingReward = ref(null)
const rewardForm = ref({
  name: '',
  description: '',
  image: '',
  points_required: 0,
  stock: -1,
  status: 1
})
const rewardRules = {
  name: [{ required: true, message: '请输入奖励名称', trigger: 'blur' }],
  points_required: [{ required: true, message: '请输入所需积分', trigger: 'blur' }]
}
const rewardFormRef = ref(null)

// 兑换相关
const exchangeRecords = ref([])
const statusDialogVisible = ref(false)
const currentExchange = ref(null)
const statusForm = ref({
  status: 1
})
const statusRules = {
  status: [{ required: true, message: '请选择状态', trigger: 'blur' }]
}
const statusFormRef = ref(null)

// 分配习惯相关
const assignHabitDialogVisible = ref(false)
const currentChild = ref(null)
const assignHabitForm = ref({
  habit_id: null
})
const assignHabitRules = {
  habit_id: [{ required: true, message: '请选择习惯', trigger: 'change' }]
}
const assignHabitFormRef = ref(null)

// 加载小孩列表
const loadChildren = async () => {
  try {
    const response = await api.get('/user/children')
    children.value = response.data.children
  } catch (error) {
    console.error('Failed to load children:', error)
  }
}

// 保存小孩
const saveChild = async () => {
  if (!childFormRef.value.validate()) return
  
  try {
    // 如果 phone 和 email 为空，则不传递
    const formData = { ...childForm.value }
    if (!formData.phone) delete formData.phone
    if (!formData.email) delete formData.email
    
    if (editingChild.value) {
      // 更新小孩
      await api.put(`/admin/children/${editingChild.value.id}`, formData)
    } else {
      // 添加小孩
      await api.post('/admin/children', formData)
    }
    childDialogVisible.value = false
    await loadChildren()
  } catch (error) {
    console.error('Failed to save child:', error)
  }
}

// 编辑小孩
const editChild = (child) => {
  editingChild.value = child
  childForm.value = { ...child }
  childDialogVisible.value = true
}

// 删除小孩
const deleteChild = async (id) => {
  try {
    await api.delete(`/admin/children/${id}`)
    await loadChildren()
  } catch (error) {
    console.error('Failed to delete child:', error)
  }
}

// 加载习惯列表
const loadHabits = async () => {
  try {
    const response = await api.get('/admin/habits')
    habits.value = response.data.habits
  } catch (error) {
    console.error('Failed to load habits:', error)
  }
}

// 加载奖励列表
const loadRewards = async () => {
  try {
    const response = await api.get('/admin/rewards')
    rewards.value = response.data.items
  } catch (error) {
    console.error('Failed to load rewards:', error)
  }
}

// 加载兑换记录
const loadExchangeRecords = async () => {
  try {
    const response = await api.get('/exchange/records')
    exchangeRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load exchange records:', error)
  }
}

// 保存习惯
const saveHabit = async () => {
  if (!habitFormRef.value.validate()) return
  
  try {
    // 确保数值类型正确
    const formData = { ...habitForm.value }
    formData.schedule_type = parseInt(formData.schedule_type)
    formData.reward_points = parseInt(formData.reward_points)
    formData.allow_makeup = parseInt(formData.allow_makeup)
    formData.makeup_days = parseInt(formData.makeup_days)
    formData.status = parseInt(formData.status)
    
    if (editingHabit.value) {
      // 更新习惯
      await api.put(`/admin/habits/${editingHabit.value.id}`, formData)
    } else {
      // 添加习惯
      await api.post('/admin/habits', formData)
    }
    dialogVisible.value = false
    await loadHabits()
  } catch (error) {
    console.error('Failed to save habit:', error)
  }
}

// 编辑习惯
const editHabit = (habit) => {
  editingHabit.value = habit
  habitForm.value = {
    ...habit,
    schedule_type: String(habit.schedule_type),
    allow_makeup: String(habit.allow_makeup),
    status: String(habit.status)
  }
  dialogVisible.value = true
}

// 删除习惯
const deleteHabit = async (id) => {
  try {
    await api.delete(`/admin/habits/${id}`)
    await loadHabits()
  } catch (error) {
    console.error('Failed to delete habit:', error)
  }
}

// 保存奖励
const saveReward = async () => {
  if (!rewardFormRef.value.validate()) return
  
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
      await api.put(`/admin/rewards/${editingReward.value.id}`, formData)
    } else {
      // 添加奖励
      await api.post('/admin/rewards', formData)
    }
    rewardDialogVisible.value = false
    await loadRewards()
  } catch (error) {
    console.error('Failed to save reward:', error)
  }
}

// 编辑奖励
const editReward = (reward) => {
  editingReward.value = reward
  rewardForm.value = { ...reward }
  rewardDialogVisible.value = true
}

// 删除奖励
const deleteReward = async (id) => {
  try {
    await api.delete(`/admin/rewards/${id}`)
    await loadRewards()
  } catch (error) {
    console.error('Failed to delete reward:', error)
  }
}

// 更新兑换状态
const updateExchangeStatus = (exchange) => {
  currentExchange.value = exchange
  statusForm.value = { status: exchange.status }
  statusDialogVisible.value = true
}

// 保存兑换状态
const saveExchangeStatus = async () => {
  if (!statusFormRef.value.validate()) return
  
  try {
    await api.put(`/admin/exchange/${currentExchange.value.id}/status`, { status: statusForm.value.status })
    statusDialogVisible.value = false
    await loadExchangeRecords()
  } catch (error) {
    console.error('Failed to update exchange status:', error)
  }
}

// 分配习惯
const assignHabit = (child) => {
  currentChild.value = child
  assignHabitForm.value = { habit_id: null }
  assignHabitDialogVisible.value = true
}

// 保存分配习惯
const saveAssignHabit = async () => {
  if (!assignHabitFormRef.value.validate()) return
  
  try {
    await api.post('/admin/habits/assign', {
      habit_id: assignHabitForm.value.habit_id,
      child_id: currentChild.value.id
    })
    assignHabitDialogVisible.value = false
    // 显示成功提示
    alert('习惯分配成功')
  } catch (error) {
    console.error('Failed to assign habit:', error)
    alert('习惯分配失败：' + (error.response?.data?.error || error.message))
  }
}

// 获取状态类型
const getStatusType = (status) => {
  switch (status) {
    case 1: return 'success'
    case 2: return 'warning'
    case 3: return 'info'
    case 4: return 'success'
    default: return 'default'
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 1: return '已完成'
    case 2: return '处理中'
    case 3: return '已发货'
    case 4: return '已收货'
    default: return '未知'
  }
}

// 登出
const logout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(async () => {
  await loadChildren()
  await loadHabits()
  await loadRewards()
  await loadExchangeRecords()
})
</script>

<style scoped>
.admin-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.admin-card {
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

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info span {
  color: #606266;
  font-size: 14px;
}

.habits-table,
.rewards-table,
.exchanges-table {
  margin-top: 20px;
}

.dialog-footer {
  text-align: right;
}
</style>
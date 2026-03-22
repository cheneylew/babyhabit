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
            <el-table-column label="操作" width="420">
              <template #default="scope">
                <el-button type="primary" size="small" @click="editChild(scope.row)">编辑</el-button>
                <el-button type="danger" size="small" @click="deleteChild(scope.row.id)">删除</el-button>
                <el-button type="warning" size="small" @click="assignHabit(scope.row)">分配习惯</el-button>
                <el-button type="info" size="small" @click="viewAssignedHabits(scope.row)">查看习惯</el-button>
                <el-button type="success" size="small" @click="viewCheckinRecords(scope.row)">查看打卡记录</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 习惯管理 -->
        <el-tab-pane label="习惯管理" name="habits">
          <el-button type="primary" @click="addHabit">添加习惯</el-button>
          <el-table :data="habits" style="width: 100%" class="habits-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="习惯名称" width="150" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="reward_points" label="奖励积分" width="100" />
            <el-table-column prop="schedule_type" label="打卡类型" width="100">
              <template #default="scope">
                {{ scope.row.schedule_type === 1 ? '每日' : '周期性' }}
              </template>
            </el-table-column>
            <el-table-column prop="schedule_detail" label="周期选择" width="150">
              <template #default="scope">
                {{ formatScheduleDetail(scope.row.schedule_detail) }}
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
            <el-table-column label="操作" width="220">
              <template #default="scope">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="approveExchange(scope.row)"
                  :disabled="scope.row.status === 1"
                >
                  通过
                </el-button>
                <el-button 
                  type="danger" 
                  size="small" 
                  @click="deleteExchangeRecord(scope.row.id)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 名言警句管理 -->
        <el-tab-pane label="名言警句" name="quotes">
          <div class="quotes-toolbar">
            <el-button type="primary" @click="addQuote">添加名言</el-button>
            <el-button type="success" @click="batchImportDialogVisible = true">批量导入</el-button>
            <el-button type="danger" :disabled="selectedQuotes.length === 0" @click="batchDeleteQuotes">批量删除</el-button>
          </div>
          <el-table :data="quotes" style="width: 100%" class="quotes-table" @selection-change="handleQuoteSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="content" label="内容" show-overflow-tooltip />
            <el-table-column prop="meaning" label="意思" show-overflow-tooltip />
            <el-table-column prop="author" label="作者" width="150" />
            <el-table-column prop="create_time" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.create_time) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button type="primary" size="small" @click="editQuote(scope.row)">编辑</el-button>
                <el-button type="danger" size="small" @click="deleteQuote(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="quotePage"
            v-model:page-size="quotePageSize"
            :total="quoteTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @size-change="loadQuotes"
            @current-change="loadQuotes"
            class="pagination"
          />
        </el-tab-pane>

        <!-- 词汇管理 -->
        <el-tab-pane label="词汇管理" name="vocabulary">
          <div class="vocabulary-toolbar" style="display: flex; justify-content: space-between; align-items: center;">
            <div>
              <el-button type="primary" @click="addVocabulary">添加词汇</el-button>
              <el-button type="success" @click="openVocabularyBatchImport">批量导入</el-button>
              <el-button type="danger" :disabled="selectedVocabularies.length === 0" @click="batchDeleteVocabularies">批量删除</el-button>
            </div>
            <div style="display: flex; gap: 10px;">
              <el-button type="warning" @click="detectIncompleteVocabularies">检测</el-button>
              <el-button type="info" @click="loadVocabularies">刷新</el-button>
            </div>
          </div>
          <el-table :data="vocabularies" style="width: 100%" class="vocabulary-table" @selection-change="handleVocabularySelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="english" label="英文" width="150" />
            <el-table-column prop="chinese" label="中文" width="150" />
            <el-table-column prop="phonetic" label="音标" width="150" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="scope">
                {{ scope.row.type === 'word' ? '单词' : '句子' }}
              </template>
            </el-table-column>
            <el-table-column label="教材" width="150">
              <template #default="scope">
                {{ getBookName(scope.row.book_id) }}
              </template>
            </el-table-column>
            <el-table-column prop="category" label="分类" width="100" />
            <el-table-column prop="create_time" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.create_time) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button type="primary" size="small" @click="editVocabulary(scope.row)">编辑</el-button>
                <el-button type="danger" size="small" @click="deleteVocabulary(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="vocabularyPage"
            v-model:page-size="vocabularyPageSize"
            :total="vocabularyTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @size-change="loadVocabularies"
            @current-change="loadVocabularies"
            class="pagination"
          />
        </el-tab-pane>

        <!-- 教材管理 -->
        <el-tab-pane label="教材管理" name="books">
          <div class="books-toolbar">
            <el-button type="primary" @click="addBook">添加教材</el-button>
          </div>
          <el-table :data="books" style="width: 100%" class="books-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="教材名称" />
            <el-table-column prop="create_time" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.create_time) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button type="primary" size="small" @click="editBook(scope.row)">编辑</el-button>
                <el-button type="danger" size="small" @click="deleteBook(scope.row.id)">删除</el-button>
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
            <el-radio value="1">正常</el-radio>
            <el-radio value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="每日单词数量" prop="daily_word_limit">
          <el-input-number v-model="childForm.daily_word_limit" :min="1" :max="100" :step="1" />
          <span style="margin-left: 10px; font-size: 12px; color: #909399;">默认5个</span>
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
            <el-radio value="1">每日</el-radio>
            <el-radio value="2">周期性</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="habitForm.schedule_type === '2'" label="周期选择" prop="schedule_detail">
          <el-checkbox-group v-model="habitForm.schedule_detail">
            <el-checkbox :label="1">周一</el-checkbox>
            <el-checkbox :label="2">周二</el-checkbox>
            <el-checkbox :label="3">周三</el-checkbox>
            <el-checkbox :label="4">周四</el-checkbox>
            <el-checkbox :label="5">周五</el-checkbox>
            <el-checkbox :label="6">周六</el-checkbox>
            <el-checkbox :label="0">周日</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item v-if="habitForm.schedule_type === '2'" label="打卡开始时间" prop="checkin_time_start">
          <el-time-picker v-model="habitForm.checkin_time_start" format="HH:mm" value-format="HH:mm:ss" placeholder="选择开始时间" />
        </el-form-item>
        <el-form-item v-if="habitForm.schedule_type === '2'" label="打卡结束时间" prop="checkin_time_end">
          <el-time-picker v-model="habitForm.checkin_time_end" format="HH:mm" value-format="HH:mm:ss" placeholder="选择结束时间" />
        </el-form-item>
        <el-form-item v-if="habitForm.schedule_type === '1'" label="打卡开始时间" prop="checkin_time_start">
          <el-time-picker v-model="habitForm.checkin_time_start" format="HH:mm" value-format="HH:mm:ss" placeholder="选择开始时间" />
        </el-form-item>
        <el-form-item v-if="habitForm.schedule_type === '1'" label="打卡结束时间" prop="checkin_time_end">
          <el-time-picker v-model="habitForm.checkin_time_end" format="HH:mm" value-format="HH:mm:ss" placeholder="选择结束时间" />
        </el-form-item>
        <el-form-item label="奖励积分" prop="reward_points">
          <el-input type="number" v-model="habitForm.reward_points" />
        </el-form-item>
        <el-form-item label="允许补卡" prop="allow_makeup">
          <el-radio-group v-model="habitForm.allow_makeup">
            <el-radio value="1">是</el-radio>
            <el-radio value="0">否</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="需要拍照" prop="require_photo">
          <el-radio-group v-model="habitForm.require_photo">
            <el-radio value="1">是</el-radio>
            <el-radio value="0">否</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="允许自我评分" prop="allow_self_rate">
          <el-radio-group v-model="habitForm.allow_self_rate">
            <el-radio value="1">是</el-radio>
            <el-radio value="0">否</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="打卡提示" prop="checkin_prompt">
          <el-input 
            v-model="habitForm.checkin_prompt" 
            type="textarea" 
            :rows="2"
            placeholder="打卡时显示给孩子的提示内容（可选）"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="habitForm.status">
            <el-radio value="1">启用</el-radio>
            <el-radio value="0">禁用</el-radio>
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
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
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
            <el-option label="待审批" :value="2" />
            <el-option label="审批通过" :value="1" />
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
    <el-dialog v-model="assignHabitDialogVisible" title="分配习惯" width="600px">
      <el-form :model="assignHabitForm" :rules="assignHabitRules" ref="assignHabitFormRef" label-width="100px">
        <el-form-item label="选择习惯" prop="habit_ids">
          <el-checkbox-group v-model="assignHabitForm.habit_ids">
            <el-checkbox 
              v-for="habit in habits" 
              :key="habit.id" 
              :label="habit.id"
            >
              {{ habit.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="assignHabitDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveAssignHabit">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 查看已分配习惯对话框 -->
    <el-dialog v-model="viewHabitsDialogVisible" title="已分配习惯" width="800px">
      <el-table :data="assignedHabits" style="width: 100%">
        <el-table-column prop="habit.name" label="习惯名称" width="150" />
        <el-table-column prop="habit.description" label="描述" />
        <el-table-column prop="habit.reward_points" label="奖励积分" width="100" />
        <el-table-column prop="assign_time" label="分配时间" width="180" />
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button type="danger" size="small" @click="deleteAssignedHabit(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="viewHabitsDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 查看打卡记录对话框 -->
    <el-dialog v-model="checkinRecordsDialogVisible" title="打卡记录" width="900px">
      <el-table :data="childCheckinRecords" style="width: 100%">
        <el-table-column prop="habit_name" label="习惯名称" width="150" />
        <el-table-column prop="checkin_date" label="打卡日期" width="120">
          <template #default="scope">
            {{ dayjs(scope.row.checkin_date).format('YYYY-MM-DD') }}
          </template>
        </el-table-column>
        <el-table-column prop="checkin_time" label="打卡时间" width="150">
          <template #default="scope">
            {{ dayjs(scope.row.checkin_time).format('HH:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column prop="checkin_type" label="打卡类型" width="100">
          <template #default="scope">
            {{ scope.row.checkin_type === 1 ? '正常' : '补卡' }}
          </template>
        </el-table-column>
        <el-table-column label="照片" width="100">
          <template #default="scope">
            <el-image 
              v-if="scope.row.photo_url" 
              :src="scope.row.photo_url" 
              :preview-src-list="[scope.row.photo_url]"
              fit="cover"
              style="width: 60px; height: 60px; cursor: pointer;"
            >
              <template #placeholder>
                <div style="width: 60px; height: 60px; background: #f5f7fa; display: flex; align-items: center; justify-content: center;">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <span v-else style="color: #999;">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="self_rate" label="自我评分" width="100">
          <template #default="scope">
            <span v-if="scope.row.self_rate && scope.row.self_rate > 0" style="color: #FF9900;">
              {{ '⭐'.repeat(scope.row.self_rate) }}
            </span>
            <span v-else style="color: #999;">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="points_rewarded" label="获得积分" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '正常' : '已回退' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button 
              type="danger" 
              size="small" 
              @click="rollbackCheckin(scope.row)"
              :disabled="scope.row.status !== 1"
            >
              回退
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="checkinRecordsDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 名言警句编辑对话框 -->
    <el-dialog v-model="quoteDialogVisible" :title="editingQuote ? '编辑名言' : '添加名言'" width="600px">
      <el-form :model="quoteForm" :rules="quoteRules" ref="quoteFormRef" label-width="80px">
        <el-form-item label="内容" prop="content">
          <el-input type="textarea" :rows="4" v-model="quoteForm.content" placeholder="请输入名言内容" />
        </el-form-item>
        <el-form-item label="意思" prop="meaning">
          <el-input type="textarea" :rows="3" v-model="quoteForm.meaning" placeholder="请输入名言的意思，通俗易懂的解释，方便小学生理解（可选）" />
        </el-form-item>
        <el-form-item label="作者" prop="author">
          <el-input v-model="quoteForm.author" placeholder="请输入作者（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="quoteDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveQuote">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 批量导入对话框 -->
    <el-dialog v-model="batchImportDialogVisible" title="批量导入名言" width="700px">
      <el-form label-width="100px">
        <el-form-item label="导入格式">
          <p class="import-hint">每行一条名言，格式：内容|意思|作者（意思和作者可选）</p>
          <p class="import-hint">示例：</p>
          <p class="import-example">学而时习之，不亦说乎|学习后经常复习，不是很愉快吗？|孔子</p>
          <p class="import-example">知识就是力量|知识能够帮助我们解决问题，是一种强大的力量|培根</p>
          <p class="import-example">千里之行，始于足下|再长的路，都是从第一步开始的</p>
        </el-form-item>
        <el-form-item label="名言内容">
          <el-input type="textarea" :rows="10" v-model="batchImportContent" placeholder="请输入要导入的名言，每行一条" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="batchImportDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveBatchImport">导入</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 词汇编辑对话框 -->
    <el-dialog v-model="vocabularyDialogVisible" :title="editingVocabulary ? '编辑词汇' : '添加词汇'" width="600px">
      <el-form :model="vocabularyForm" :rules="vocabularyRules" ref="vocabularyFormRef" label-width="80px">
        <el-form-item label="英文" prop="english">
          <el-input v-model="vocabularyForm.english" placeholder="请输入英文单词或句子" />
        </el-form-item>
        <el-form-item label="中文" prop="chinese">
          <el-input v-model="vocabularyForm.chinese" placeholder="请输入中文翻译" />
        </el-form-item>
        <el-form-item label="音标" prop="phonetic">
          <el-input v-model="vocabularyForm.phonetic" placeholder="请输入音标（可选）" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="vocabularyForm.type" placeholder="请选择类型">
            <el-option label="单词" value="word" />
            <el-option label="句子" value="sentence" />
          </el-select>
        </el-form-item>
        <el-form-item label="教材" prop="book_id">
          <el-select-v2
            v-model="vocabularyForm.book_id"
            :options="bookOptions.map(b => ({ label: b.name, value: b.id }))"
            placeholder="请选择教材"
            filterable
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="vocabularyForm.category" placeholder="请输入分类（如：水果）" />
        </el-form-item>
        <el-form-item label="例句" prop="example_sentence">
          <el-input type="textarea" :rows="3" v-model="vocabularyForm.example_sentence" placeholder="请输入例句（可选）" />
        </el-form-item>
        <el-form-item label="音频URL" prop="audio_url">
          <el-input v-model="vocabularyForm.audio_url" placeholder="请输入音频URL（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="vocabularyDialogVisible = false" :disabled="vocabularyLoading">取消</el-button>
          <el-button type="primary" @click="saveVocabulary" :loading="vocabularyLoading">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 词汇批量导入对话框 -->
    <el-dialog v-model="vocabularyBatchImportDialogVisible" title="批量导入词汇" width="700px">
      <el-form label-width="100px">
        <el-form-item label="教材" required>
          <el-select-v2
            v-model="batchImportBookId"
            :options="bookOptions.map(b => ({ label: b.name, value: b.id }))"
            placeholder="请选择教材"
            filterable
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="导入格式">
          <p class="import-hint">每行一条词汇，直接输入英文内容即可</p>
          <p class="import-hint">系统会自动判断类型：包含空格、长度较长或以标点结尾的视为句子，否则视为单词</p>
          <p class="import-hint">示例：</p>
          <p class="import-example">Gesture（单词）</p>
          <p class="import-example">Hello, how are you?（句子）</p>
        </el-form-item>
        <el-form-item label="词汇内容">
          <el-input type="textarea" :rows="10" v-model="vocabularyBatchImportContent" placeholder="请输入要导入的词汇，每行一条" :disabled="vocabularyBatchImportLoading" />
        </el-form-item>
        <el-form-item label="批量备注">
          <el-input type="textarea" :rows="3" v-model="batchImportRemark" placeholder="请输入批量备注，将应用于所有导入的词汇（可选）" :disabled="vocabularyBatchImportLoading" />
        </el-form-item>
        <el-form-item v-if="vocabularyBatchImportLoading">
          <el-progress :percentage="vocabularyImportProgress" :format="formatProgress" />
          <p class="import-progress-text">{{ vocabularyImportCurrentWord }} - {{ vocabularyImportStatus }}</p>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelVocabularyBatchImport" :disabled="vocabularyBatchImportLoading">取消</el-button>
          <el-button type="primary" @click="saveVocabularyBatchImport" :loading="vocabularyBatchImportLoading">导入</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 词汇检测对话框 -->
    <el-dialog v-model="vocabularyDetectDialogVisible" title="词汇检测" width="800px">
      <div v-if="!vocabularyDetectLoading && incompleteVocabularies.length > 0">
        <p class="detect-hint">检测到以下词汇缺少AI生成的信息（中文、音标、例句或发音）：</p>
        <el-table :data="incompleteVocabularies" style="width: 100%" class="detect-table">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="english" label="英文" width="200" />
          <el-table-column prop="chinese" label="中文" width="150" show-overflow-tooltip />
          <el-table-column prop="phonetic" label="音标" width="150" show-overflow-tooltip />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column label="状态" width="150">
            <template #default="scope">
              <div>
                <span v-if="!scope.row.chinese" class="status-missing">缺少中文</span>
                <span v-if="!scope.row.phonetic" class="status-missing">缺少音标</span>
                <span v-if="!scope.row.example_sentence" class="status-missing">缺少例句</span>
                <span v-if="!scope.row.audio_url" class="status-missing">缺少发音</span>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-else-if="!vocabularyDetectLoading && incompleteVocabularies.length === 0">
        <el-empty description="未检测到缺少AI生成信息的词汇" />
      </div>
      <div v-else>
        <el-progress :percentage="vocabularyDetectProgress" :format="formatProgress" />
        <p class="import-progress-text">{{ vocabularyDetectCurrentWord }} - {{ vocabularyDetectStatus }}</p>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="vocabularyDetectDialogVisible = false" :disabled="vocabularyDetectLoading">关闭</el-button>
          <el-button v-if="!vocabularyDetectLoading && incompleteVocabularies.length > 0" type="primary" @click="regenerateVocabularies">重新生成</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 教材编辑对话框 -->
    <el-dialog v-model="bookDialogVisible" :title="editingBook ? '编辑教材' : '添加教材'" width="500px">
      <el-form :model="bookForm" :rules="bookRules" ref="bookFormRef" label-width="100px">
        <el-form-item label="教材名称" prop="name">
          <el-input v-model="bookForm.name" placeholder="请输入教材名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="bookDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveBook" :loading="bookLoading">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import api from '../api'
import dayjs from 'dayjs'
import { Picture } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const userStore = useUserStore()
const activeTab = ref('children')

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 格式化周期选择
const formatScheduleDetail = (detail) => {
  if (!detail) return '-'
  
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const days = detail.split(',')
  const formattedDays = days.map(day => {
    const dayNum = parseInt(day)
    if (dayNum >= 0 && dayNum <= 6) {
      return weekdays[dayNum]
    }
    return day
  })
  
  return formattedDays.join('、')
}

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
  status: 1,
  daily_word_limit: 5
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
  schedule_type: '1',
  schedule_detail: [],
  checkin_time_start: '',
  checkin_time_end: '',
  reward_points: 0,
  allow_makeup: '0',
  makeup_days: 0,
  require_photo: '0',
  allow_self_rate: '0',
  checkin_prompt: '',
  status: '1'
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

// 查看已分配习惯相关
const viewHabitsDialogVisible = ref(false)
const assignedHabits = ref([])
const currentChild = ref(null)

// 分配习惯相关
const assignHabitDialogVisible = ref(false)
const assignHabitForm = ref({
  habit_ids: []
})
const assignHabitRules = {
  habit_ids: [{ required: true, message: '请选择习惯', trigger: 'change', type: 'array', min: 1 }]
}
const assignHabitFormRef = ref(null)

// 查看打卡记录相关
const checkinRecordsDialogVisible = ref(false)
const childCheckinRecords = ref([])

// 名言警句相关
const quotes = ref([])
const quotePage = ref(1)
const quotePageSize = ref(10)
const quoteTotal = ref(0)
const quoteDialogVisible = ref(false)
const editingQuote = ref(null)
const quoteForm = ref({
  content: '',
  author: ''
})
const quoteRules = {
  content: [{ required: true, message: '请输入名言内容', trigger: 'blur' }]
}
const quoteFormRef = ref(null)
const selectedQuotes = ref([])
const batchImportDialogVisible = ref(false)
const batchImportContent = ref('')

// 词汇管理相关
const vocabularies = ref([])
const vocabularyPage = ref(1)
const vocabularyPageSize = ref(10)
const vocabularyTotal = ref(0)
const vocabularyDialogVisible = ref(false)
const vocabularyBatchImportDialogVisible = ref(false)
const editingVocabulary = ref(null)
const vocabularyForm = ref({
  english: '',
  chinese: '',
  phonetic: '',
  type: 'word',
  book_id: '',
  category: '',
  example_sentence: '',
  audio_url: ''
})
const vocabularyRules = {
  english: [{ required: true, message: '请输入英文', trigger: 'blur' }],
  chinese: [{ required: false, message: '请输入中文', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  book_id: [{ required: true, message: '请输入或选择教材', trigger: 'change' }]
}

// 教材选项
const bookOptions = ref([])
const batchImportBookId = ref('')
const vocabularyFormRef = ref(null)
const selectedVocabularies = ref([])
const vocabularyBatchImportContent = ref('')
const vocabularyLoading = ref(false)
const vocabularyBatchImportLoading = ref(false)
const vocabularyImportProgress = ref(0)
const vocabularyImportCurrentWord = ref('')
const vocabularyImportStatus = ref('')
const vocabularyImportAbort = ref(false)
const batchImportRemark = ref('')

// 词汇检测相关
const vocabularyDetectDialogVisible = ref(false)
const incompleteVocabularies = ref([])
const vocabularyDetectLoading = ref(false)
const vocabularyDetectProgress = ref(0)
const vocabularyDetectCurrentWord = ref('')
const vocabularyDetectStatus = ref('')
const vocabularyDetectAbort = ref(false)

// 教材管理相关
const books = ref([])
const bookDialogVisible = ref(false)
const editingBook = ref(null)
const bookForm = ref({
  name: ''
})
const bookRules = {
  name: [{ required: true, message: '请输入教材名称', trigger: 'blur' }]
}
const bookFormRef = ref(null)
const bookLoading = ref(false)

// 加载名言警句列表
const loadQuotes = async () => {
  try {
    const response = await api.get('/admin/quotes', {
      params: {
        page: quotePage.value,
        page_size: quotePageSize.value
      }
    })
    quotes.value = response.data.quotes
    quoteTotal.value = response.data.total
  } catch (error) {
    console.error('Failed to load quotes:', error)
    ElMessage.error('加载名言警句失败')
  }
}

// 添加名言
const addQuote = () => {
  editingQuote.value = null
  quoteForm.value = {
    content: '',
    meaning: '',
    author: ''
  }
  quoteDialogVisible.value = true
}

// 编辑名言
const editQuote = (quote) => {
  editingQuote.value = quote
  quoteForm.value = {
    content: quote.content,
    meaning: quote.meaning || '',
    author: quote.author || ''
  }
  quoteDialogVisible.value = true
}

// 保存名言
const saveQuote = async () => {
  if (!quoteFormRef.value.validate()) return

  try {
    if (!editingQuote.value) {
      // 添加名言需要二次确认
      await ElMessageBox.confirm(
        '确定要添加这条名言吗？',
        '添加确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info',
        }
      )
    }

    if (editingQuote.value) {
      // 更新名言
      await api.put(`/admin/quotes/${editingQuote.value.id}`, quoteForm.value)
      ElMessage.success('更新成功')
    } else {
      // 添加名言
      await api.post('/admin/quotes', quoteForm.value)
      ElMessage.success('添加成功')
    }
    quoteDialogVisible.value = false
    await loadQuotes()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to save quote:', error)
      ElMessage.error('操作失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 删除名言
const deleteQuote = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这条名言吗？此操作不可恢复。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await api.delete(`/admin/quotes/${id}`)
    ElMessage.success('删除成功')
    await loadQuotes()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete quote:', error)
      ElMessage.error('删除失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 处理名言选择变化
const handleQuoteSelectionChange = (selection) => {
  selectedQuotes.value = selection
}

// 批量删除名言
const batchDeleteQuotes = async () => {
  if (selectedQuotes.value.length === 0) {
    ElMessage.warning('请先选择要删除的名言')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedQuotes.value.length} 条名言吗？此操作不可恢复。`,
      '批量删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const ids = selectedQuotes.value.map(q => q.id)
    await api.delete('/admin/quotes/batch', { data: { ids } })
    ElMessage.success('批量删除成功')
    selectedQuotes.value = []
    await loadQuotes()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to batch delete quotes:', error)
      ElMessage.error('批量删除失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 批量导入名言
const saveBatchImport = async () => {
  if (!batchImportContent.value.trim()) {
    ElMessage.warning('请输入要导入的名言内容')
    return
  }

  try {
    await ElMessageBox.confirm(
      '确定要导入这些名言吗？',
      '导入确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info',
      }
    )

    // 解析导入内容
    const lines = batchImportContent.value.trim().split('\n')
    const quotes = []

    for (const line of lines) {
      if (!line.trim()) continue

      const parts = line.split('|')
      const content = parts[0].trim()
      const meaning = parts[1] ? parts[1].trim() : ''
      const author = parts[2] ? parts[2].trim() : ''

      if (content) {
        quotes.push({ content, meaning, author })
      }
    }

    if (quotes.length === 0) {
      ElMessage.warning('没有有效的名言内容')
      return
    }

    await api.post('/admin/quotes/batch', { quotes })
    ElMessage.success(`成功导入 ${quotes.length} 条名言`)
    batchImportDialogVisible.value = false
    batchImportContent.value = ''
    await loadQuotes()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to batch import quotes:', error)
      ElMessage.error('批量导入失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 加载词汇列表
const loadVocabularies = async () => {
  try {
    const response = await api.get('/admin/vocabulary', {
      params: {
        page: vocabularyPage.value,
        page_size: vocabularyPageSize.value
      }
    })
    vocabularies.value = response.data.vocabularies
    vocabularyTotal.value = response.data.total
  } catch (error) {
    console.error('Failed to load vocabularies:', error)
    ElMessage.error('加载词汇失败')
  }
}

// 加载教材选项
const loadBookOptions = async () => {
  try {
    const response = await api.get('/vocabulary/options/book')
    bookOptions.value = response.data.books || []
  } catch (error) {
    console.error('Failed to load book options:', error)
  }
}

// 加载教材列表
const loadBooks = async () => {
  try {
    const response = await api.get('/admin/books')
    books.value = response.data.books || []
  } catch (error) {
    console.error('Failed to load books:', error)
    ElMessage.error('加载教材列表失败')
  }
}

// 根据教材ID获取教材名称
const getBookName = (bookId) => {
  if (!bookId) return '-'
  const book = books.value.find(b => b.id === bookId)
  return book ? book.name : '-'
}

// 添加教材
const addBook = () => {
  editingBook.value = null
  bookForm.value = {
    name: ''
  }
  bookDialogVisible.value = true
}

// 编辑教材
const editBook = (book) => {
  editingBook.value = book
  bookForm.value = {
    name: book.name
  }
  bookDialogVisible.value = true
}

// 保存教材
const saveBook = async () => {
  if (!bookFormRef.value.validate()) return

  bookLoading.value = true
  try {
    if (editingBook.value) {
      await api.put(`/admin/books/${editingBook.value.id}`, bookForm.value)
      ElMessage.success('更新成功')
    } else {
      await api.post('/admin/books', bookForm.value)
      ElMessage.success('添加成功')
    }
    bookDialogVisible.value = false
    await loadBooks()
    await loadBookOptions()
  } catch (error) {
    console.error('Failed to save book:', error)
    ElMessage.error('保存失败')
  } finally {
    bookLoading.value = false
  }
}

// 删除教材
const deleteBook = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个教材吗？',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    await api.delete(`/admin/books/${id}`)
    ElMessage.success('删除成功')
    await loadBooks()
    await loadBookOptions()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete book:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 添加词汇
const addVocabulary = async () => {
  editingVocabulary.value = null
  vocabularyForm.value = {
    english: '',
    chinese: '',
    phonetic: '',
    type: 'word',
    book_id: '',
    category: '',
    example_sentence: '',
    audio_url: ''
  }
  await loadBookOptions()
  vocabularyDialogVisible.value = true
}

// 编辑词汇
const editVocabulary = async (vocabulary) => {
  editingVocabulary.value = vocabulary
  vocabularyForm.value = {
    english: vocabulary.english || '',
    chinese: vocabulary.chinese || '',
    phonetic: vocabulary.phonetic || '',
    type: vocabulary.type || 'word',
    book_id: vocabulary.book_id || '',
    category: vocabulary.category || '',
    example_sentence: vocabulary.example_sentence || '',
    audio_url: vocabulary.audio_url || ''
  }
  await loadBookOptions()
  vocabularyDialogVisible.value = true
}

// 保存词汇
const saveVocabulary = async () => {
  if (!vocabularyFormRef.value.validate()) return

  try {
    if (!editingVocabulary.value) {
      // 添加词汇需要二次确认
      await ElMessageBox.confirm(
        '确定要添加这个词汇吗？',
        '添加确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info',
        }
      )
    }

    vocabularyLoading.value = true

    if (editingVocabulary.value) {
      // 更新词汇
      await api.put(`/admin/vocabulary/${editingVocabulary.value.id}`, vocabularyForm.value)
      ElMessage.success('更新成功')
    } else {
      // 添加词汇
      await api.post('/admin/vocabulary', vocabularyForm.value)
      ElMessage.success('添加成功')
    }
    vocabularyDialogVisible.value = false
    await loadVocabularies()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to save vocabulary:', error)
      ElMessage.error('操作失败：' + (error.response?.data?.error || error.message))
    }
  } finally {
    vocabularyLoading.value = false
  }
}

// 删除词汇
const deleteVocabulary = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个词汇吗？此操作不可恢复。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await api.delete(`/admin/vocabulary/${id}`)
    ElMessage.success('删除成功')
    await loadVocabularies()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete vocabulary:', error)
      ElMessage.error('删除失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 处理词汇选择变化
const handleVocabularySelectionChange = (selection) => {
  selectedVocabularies.value = selection
}

// 批量删除词汇
const batchDeleteVocabularies = async () => {
  if (selectedVocabularies.value.length === 0) {
    ElMessage.warning('请先选择要删除的词汇')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedVocabularies.value.length} 个词汇吗？此操作不可恢复。`,
      '批量删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const ids = selectedVocabularies.value.map(v => v.id)
    await api.delete('/admin/vocabulary/batch', { data: { ids } })
    ElMessage.success('批量删除成功')
    selectedVocabularies.value = []
    await loadVocabularies()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to batch delete vocabularies:', error)
      ElMessage.error('批量删除失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 格式化进度显示
const formatProgress = (percentage) => {
  return `${percentage}%`
}

// 打开词汇批量导入对话框
const openVocabularyBatchImport = async () => {
  batchImportBookId.value = ''
  vocabularyBatchImportContent.value = ''
  batchImportRemark.value = ''
  await loadBookOptions()
  vocabularyBatchImportDialogVisible.value = true
}

// 取消词汇批量导入
const cancelVocabularyBatchImport = () => {
  vocabularyImportAbort.value = true
  vocabularyBatchImportDialogVisible.value = false
  vocabularyBatchImportLoading.value = false
  vocabularyImportProgress.value = 0
  vocabularyImportCurrentWord.value = ''
  vocabularyImportStatus.value = ''
  batchImportRemark.value = ''
}

// 检测缺少AI生成信息的词汇
const detectIncompleteVocabularies = async () => {
  try {
    vocabularyDetectLoading.value = true
    vocabularyDetectProgress.value = 0
    vocabularyDetectCurrentWord.value = ''
    vocabularyDetectStatus.value = '正在检测...'
    
    const response = await api.get('/admin/vocabulary/incomplete')
    incompleteVocabularies.value = response.data.vocabularies || []
    
    vocabularyDetectLoading.value = false
    vocabularyDetectDialogVisible.value = true
  } catch (error) {
    console.error('Failed to detect incomplete vocabularies:', error)
    ElMessage.error('检测失败：' + (error.response?.data?.error || error.message))
    vocabularyDetectLoading.value = false
  }
}

// 重新生成词汇的AI信息
const regenerateVocabularies = async () => {
  if (incompleteVocabularies.value.length === 0) {
    ElMessage.warning('没有需要重新生成的词汇')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要重新生成 ${incompleteVocabularies.value.length} 个词汇的AI信息吗？`,
      '重新生成确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info',
      }
    )
    
    vocabularyDetectLoading.value = true
    vocabularyDetectProgress.value = 0
    vocabularyDetectCurrentWord.value = ''
    vocabularyDetectStatus.value = '正在重新生成...'
    
    // 构建词汇ID列表
    const vocabularyIds = incompleteVocabularies.value.map(v => v.id)
    
    // 发送重新生成请求
    const response = await fetch(api.defaults.baseURL + '/admin/vocabulary/regenerate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ ids: vocabularyIds })
    })
    
    if (!response.ok) {
      throw new Error('重新生成失败')
    }
    
    // 处理流式响应
    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      
      const chunk = decoder.decode(value, { stream: true })
      const lines = chunk.split('\n')
      
      for (const line of lines) {
        if (!line.trim()) continue
        
        try {
          const data = JSON.parse(line)
          if (data.status === 'progress') {
            // 更新进度
            const progress = Math.round((data.current / data.total) * 100)
            vocabularyDetectProgress.value = progress
            vocabularyDetectCurrentWord.value = data.word
            vocabularyDetectStatus.value = `正在重新生成: ${data.word} (${data.current}/${data.total})`
          } else if (data.status === 'completed') {
            // 重新生成完成
            vocabularyDetectProgress.value = 100
            vocabularyDetectStatus.value = '重新生成完成'
            
            ElMessage.success(`重新生成完成：成功 ${data.success} 个词汇`)
            
            // 关闭对话框并重置状态
            setTimeout(() => {
              vocabularyDetectDialogVisible.value = false
              vocabularyDetectLoading.value = false
              vocabularyDetectProgress.value = 0
              vocabularyDetectCurrentWord.value = ''
              vocabularyDetectStatus.value = ''
              // 重新加载词汇列表
              loadVocabularies()
            }, 1000)
          }
        } catch (e) {
          console.error('Error parsing JSON:', e)
        }
      }
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to regenerate vocabularies:', error)
      ElMessage.error('重新生成失败：' + (error.response?.data?.error || error.message))
    }
    vocabularyDetectLoading.value = false
  }
}

// 批量导入词汇
const saveVocabularyBatchImport = async () => {
  // 验证教材
  if (!batchImportBookId.value) {
    ElMessage.warning('请选择教材')
    return
  }
  if (!vocabularyBatchImportContent.value.trim()) {
    ElMessage.warning('请输入要导入的词汇内容')
    return
  }

  try {
    await ElMessageBox.confirm(
      '确定要导入这些词汇吗？',
      '导入确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info',
      }
    )

    // 解析导入内容
    const lines = vocabularyBatchImportContent.value.trim().split('\n')
    const vocabularies = []

    for (const line of lines) {
      if (!line.trim()) continue

      const english = line.trim()
      // 自动判断类型：包含空格、长度较长或以标点结尾的视为句子
      let type = 'word'
      if (english.includes(' ') || english.length > 10 || /[.!?]$/.test(english)) {
        type = 'sentence'
      }

      if (english) {
        vocabularies.push({ english, type })
      }
    }

    if (vocabularies.length === 0) {
      ElMessage.warning('没有有效的词汇内容')
      return
    }

    vocabularyBatchImportLoading.value = true
    vocabularyImportProgress.value = 0
    vocabularyImportCurrentWord.value = ''
    vocabularyImportStatus.value = ''
    vocabularyImportAbort.value = false

    try {
      // 调用批量导入API（使用Fetch API处理流式响应）
      const response = await fetch('/api/admin/vocabulary/batch', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          vocabularies: vocabularies,
          book_id: batchImportBookId.value,
          remark: batchImportRemark.value
        })
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        if (vocabularyImportAbort.value) {
          reader.cancel()
          break
        }

        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        
        // 处理每一行JSON
        const lines = buffer.split('\n')
        buffer = lines.pop() // 保留不完整的最后一行

        for (const line of lines) {
          if (!line.trim()) continue

          try {
            const data = JSON.parse(line)
            if (data.status === 'processing') {
              // 更新进度
              const progress = Math.round((data.current / data.total) * 100)
              vocabularyImportProgress.value = progress
              vocabularyImportCurrentWord.value = data.word
              vocabularyImportStatus.value = `正在导入: ${data.word} (${data.current}/${data.total})`
            } else if (data.status === 'completed') {
              // 导入完成
              vocabularyImportProgress.value = 100
              vocabularyImportStatus.value = '导入完成'
              
              ElMessage.success(`导入完成：成功 ${data.success} 个词汇，重复 ${data.duplicates} 个词汇`)
              
              // 关闭对话框并重置状态
              setTimeout(() => {
                vocabularyBatchImportDialogVisible.value = false
                vocabularyBatchImportContent.value = ''
                batchImportBookId.value = ''
                vocabularyBatchImportLoading.value = false
                vocabularyImportProgress.value = 0
                vocabularyImportCurrentWord.value = ''
                vocabularyImportStatus.value = ''
                // 重新加载词汇列表
                loadVocabularies()
              }, 1000)
            }
          } catch (e) {
            console.error('Error parsing JSON:', e)
          }
        }
      }
    } catch (error) {
      console.error('Failed to batch import vocabulary:', error)
      ElMessage.error('批量导入失败：' + error.message)
      vocabularyBatchImportLoading.value = false
      vocabularyImportProgress.value = 0
      vocabularyImportCurrentWord.value = ''
      vocabularyImportStatus.value = ''
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to batch import vocabulary:', error)
      ElMessage.error('批量导入失败：' + (error.response?.data?.error || error.message))
      vocabularyBatchImportLoading.value = false
      vocabularyImportProgress.value = 0
      vocabularyImportCurrentWord.value = ''
      vocabularyImportStatus.value = ''
    }
  }
}

// 加载小孩列表
const loadChildren = async () => {
  try {
    const response = await api.get('/user/children')
    const childList = response.data.children
    
    // 为每个小孩获取每日单词数量偏好设置
    for (const child of childList) {
      try {
        const prefResponse = await api.get(`/user/preference?key=daily_word_limit&user_id=${child.id}`)
        if (prefResponse.data.preference) {
          child.daily_word_limit = parseInt(prefResponse.data.preference.preference_value) || 5
        } else {
          child.daily_word_limit = 5
        }
      } catch (error) {
        console.error(`Failed to get daily word limit for child ${child.id}:`, error)
        child.daily_word_limit = 5
      }
    }
    
    children.value = childList
  } catch (error) {
    console.error('Failed to load children:', error)
  }
}

// 保存小孩
const saveChild = async () => {
  if (!childFormRef.value.validate()) return
  
  try {
    if (!editingChild.value) {
      // 添加小孩需要二次确认
      await ElMessageBox.confirm(
        '确定要添加这个小孩吗？',
        '添加确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info',
        }
      )
    }
    
    // 如果 phone 和 email 为空，则不传递
    const formData = { ...childForm.value }
    if (!formData.phone) delete formData.phone
    if (!formData.email) delete formData.email
    // 将状态值转换为整数类型
    formData.status = parseInt(formData.status)
    // 提取每日单词数量设置
    const dailyWordLimit = formData.daily_word_limit
    delete formData.daily_word_limit
    
    let childId
    if (editingChild.value) {
      // 更新小孩
      await api.put(`/admin/children/${editingChild.value.id}`, formData)
      ElMessage.success('更新成功')
      childId = editingChild.value.id
    } else {
      // 添加小孩
      const response = await api.post('/admin/children', formData)
      ElMessage.success('添加成功')
      childId = response.data.child.id
    }
    
    // 保存每日单词数量到用户偏好设置
    try {
      await api.post('/user/preference', {
        key: 'daily_word_limit',
        value: dailyWordLimit.toString(),
        user_id: childId
      })
      console.log('每日单词数量保存成功')
    } catch (error) {
      console.error('保存每日单词数量失败:', error)
      ElMessage.error('保存每日单词数量失败：' + (error.response?.data?.error || error.message))
    }
    
    childDialogVisible.value = false
    await loadChildren()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to save child:', error)
      ElMessage.error('操作失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 编辑小孩
const editChild = async (child) => {
  editingChild.value = child
  
  // 获取每日单词数量偏好设置
  let dailyWordLimit = 5
  try {
    const prefResponse = await api.get(`/user/preference?key=daily_word_limit&user_id=${child.id}`)
    if (prefResponse.data.preference) {
      dailyWordLimit = parseInt(prefResponse.data.preference.preference_value) || 5
    }
  } catch (error) {
    console.error(`Failed to get daily word limit for child ${child.id}:`, error)
  }
  
  childForm.value = { 
    ...child,
    status: child.status.toString(), // 确保状态值是字符串类型
    daily_word_limit: dailyWordLimit
  }
  childDialogVisible.value = true
}

// 删除小孩
const deleteChild = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个小孩吗？此操作不可恢复。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await api.delete(`/admin/children/${id}`)
    ElMessage.success('删除成功')
    await loadChildren()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete child:', error)
      ElMessage.error('删除失败：' + (error.response?.data?.error || error.message))
    }
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
    const response = await api.get('/admin/exchanges')
    exchangeRecords.value = response.data.records
  } catch (error) {
    console.error('Failed to load exchange records:', error)
  }
}

// 保存习惯
const saveHabit = async () => {
  if (!habitFormRef.value.validate()) return
  
  try {
    if (!editingHabit.value) {
      // 添加习惯需要二次确认
      await ElMessageBox.confirm(
        '确定要添加这个习惯吗？',
        '添加确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info',
        }
      )
    }
    
    // 确保数值类型正确
    const formData = { ...habitForm.value }
    formData.schedule_type = parseInt(formData.schedule_type)
    formData.reward_points = parseInt(formData.reward_points)
    formData.allow_makeup = parseInt(formData.allow_makeup)
    formData.makeup_days = parseInt(formData.makeup_days)
    formData.require_photo = parseInt(formData.require_photo)
    formData.allow_self_rate = parseInt(formData.allow_self_rate)
    formData.status = parseInt(formData.status)
    
    // 调试日志
    console.log('=== saveHabit Debug ===')
    console.log('checkin_prompt:', formData.checkin_prompt)
    console.log('formData:', formData)
    console.log('======================')
    
    // 如果是周期性习惯，将数组转换为逗号分隔的字符串
    if (formData.schedule_type === 2 && Array.isArray(formData.schedule_detail)) {
      formData.schedule_detail = formData.schedule_detail.join(',')
    } else {
      // 每日习惯，schedule_detail 设为空字符串
      formData.schedule_detail = ''
    }
    
    if (editingHabit.value) {
      // 更新习惯
      await api.put(`/admin/habits/${editingHabit.value.id}`, formData)
      ElMessage.success('更新成功')
    } else {
      // 添加习惯
      await api.post('/admin/habits', formData)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    await loadHabits()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to save habit:', error)
      ElMessage.error('操作失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 添加习惯
const addHabit = () => {
  editingHabit.value = null
  habitForm.value = {
    name: '',
    description: '',
    icon: '',
    category: '',
    schedule_type: '1',
    schedule_detail: [],
    checkin_time_start: '',
    checkin_time_end: '',
    reward_points: 0,
    allow_makeup: '0',
    makeup_days: 0,
    require_photo: '0',
    allow_self_rate: '0',
    checkin_prompt: '',
    status: '1'
  }
  dialogVisible.value = true
}

// 编辑习惯
const editHabit = (habit) => {
  editingHabit.value = habit
  
  // 确保所有字段都有默认值
  const requirePhoto = (habit.require_photo !== undefined && habit.require_photo !== null) ? habit.require_photo : 0
  const allowSelfRate = (habit.allow_self_rate !== undefined && habit.allow_self_rate !== null) ? habit.allow_self_rate : 0
  const allowMakeup = (habit.allow_makeup !== undefined && habit.allow_makeup !== null) ? habit.allow_makeup : 0
  const status = (habit.status !== undefined && habit.status !== null) ? habit.status : 1
  const scheduleType = habit.schedule_type || 1
  
  habitForm.value = {
    name: habit.name || '',
    description: habit.description || '',
    icon: habit.icon || '',
    category: habit.category || '',
    schedule_type: String(scheduleType),
    schedule_detail: [],
    checkin_time_start: habit.checkin_time_start || '',
    checkin_time_end: habit.checkin_time_end || '',
    reward_points: habit.reward_points || 0,
    allow_makeup: String(allowMakeup),
    makeup_days: habit.makeup_days || 0,
    require_photo: String(requirePhoto),
    allow_self_rate: String(allowSelfRate),
    checkin_prompt: habit.checkin_prompt || '',
    status: String(status)
  }
  
  // 如果是周期性习惯，将逗号分隔的字符串转换为数字数组
  if (habit.schedule_type === 2 && habit.schedule_detail) {
    habitForm.value.schedule_detail = habit.schedule_detail.split(',').map(item => parseInt(item.trim()))
  }
  
  dialogVisible.value = true
}

// 删除习惯
const deleteHabit = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个习惯吗？此操作不可恢复。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await api.delete(`/admin/habits/${id}`)
    ElMessage.success('删除成功')
    await loadHabits()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete habit:', error)
      ElMessage.error('删除失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 保存奖励
const saveReward = async () => {
  if (!rewardFormRef.value.validate()) return
  
  try {
    if (!editingReward.value) {
      // 添加奖励需要二次确认
      await ElMessageBox.confirm(
        '确定要添加这个奖励吗？',
        '添加确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info',
        }
      )
    }
    
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
      ElMessage.success('更新成功')
    } else {
      // 添加奖励
      await api.post('/admin/rewards', formData)
      ElMessage.success('添加成功')
    }
    rewardDialogVisible.value = false
    await loadRewards()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to save reward:', error)
      ElMessage.error('操作失败：' + (error.response?.data?.error || error.message))
    }
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
    await ElMessageBox.confirm(
      '确定要删除这个奖励吗？此操作不可恢复。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await api.delete(`/admin/rewards/${id}`)
    ElMessage.success('删除成功')
    await loadRewards()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete reward:', error)
      ElMessage.error('删除失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 确认审批通过
const approveExchange = async (exchange) => {
  try {
    await api.put(`/admin/exchange/${exchange.id}/status`, { status: 1 })
    ElMessage.success('审批通过成功')
    await loadExchangeRecords()
  } catch (error) {
    console.error('Failed to approve exchange:', error)
    ElMessage.error('审批通过失败')
  }
}

// 删除兑换记录
const deleteExchangeRecord = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个兑换记录吗？此操作不可恢复。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await api.delete(`/admin/exchanges/${id}`)
    ElMessage.success('删除成功')
    await loadExchangeRecords()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete exchange record:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 分配习惯
const assignHabit = async (child) => {
  currentChild.value = child
  
  // 加载已分配的习惯
  try {
    const response = await api.get('/admin/assigned-habits', {
      params: { child_id: child.id }
    })
    const assignedHabitIds = response.data.assignments.map(a => a.habit_id)
    assignHabitForm.value = { habit_ids: assignedHabitIds }
  } catch (error) {
    console.error('Failed to load assigned habits:', error)
    assignHabitForm.value = { habit_ids: [] }
  }
  
  assignHabitDialogVisible.value = true
}

// 保存分配习惯
const saveAssignHabit = async () => {
  if (!assignHabitFormRef.value.validate()) return
  
  try {
    // 调用批量分配接口
    const habitIds = assignHabitForm.value.habit_ids
    const childId = currentChild.value.id
    
    await api.post('/admin/habits/batch-assign', {
      habit_ids: habitIds,
      child_id: childId
    })
    
    assignHabitDialogVisible.value = false
    // 显示成功提示
    ElMessage.success('习惯分配成功')
  } catch (error) {
    console.error('Failed to assign habit:', error)
    ElMessage.error('习惯分配失败：' + (error.response?.data?.error || error.message))
  }
}

// 查看已分配习惯
const viewAssignedHabits = async (child) => {
  currentChild.value = child
  try {
    const response = await api.get('/admin/assigned-habits', {
      params: { child_id: child.id }
    })
    assignedHabits.value = response.data.assignments
    viewHabitsDialogVisible.value = true
  } catch (error) {
    console.error('Failed to load assigned habits:', error)
    alert('加载已分配习惯失败')
  }
}

// 删除已分配习惯
const deleteAssignedHabit = async (assignment) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个已分配的习惯吗？',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await api.delete(`/admin/habit-assignments/${assignment.id}`)
    ElMessage.success('删除成功')
    // 重新加载列表
    await viewAssignedHabits(currentChild.value)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete assigned habit:', error)
      ElMessage.error('删除失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 查看打卡记录
const viewCheckinRecords = async (child) => {
  currentChild.value = child
  try {
    const startDate = dayjs().subtract(30, 'day').format('YYYY-MM-DD')
    const endDate = dayjs().format('YYYY-MM-DD')
    const response = await api.get('/admin/child/checkin-records', {
      params: {
        child_id: child.id,
        start_date: startDate,
        end_date: endDate
      }
    })
    childCheckinRecords.value = response.data.records
    checkinRecordsDialogVisible.value = true
  } catch (error) {
    console.error('Failed to load checkin records:', error)
    ElMessage.error('加载打卡记录失败：' + (error.response?.data?.error || error.message))
  }
}

// 回退打卡
const rollbackCheckin = async (record) => {
  try {
    await ElMessageBox.confirm(
      `确定要回退这次打卡吗？回退后将扣除 ${record.points_rewarded} 积分。`,
      '回退打卡',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await api.post('/admin/checkin/rollback', {
      checkin_id: record.id,
      reason: '管理员手动回退'
    })
    
    ElMessage.success('回退成功')
    // 重新加载列表
    await viewCheckinRecords(currentChild.value)
  } catch (error) {
    if (error === 'cancel') return
    console.error('Failed to rollback checkin:', error)
    ElMessage.error('回退失败：' + (error.response?.data?.error || error.message))
  }
}

// 获取状态类型
const getStatusType = (status) => {
  switch (status) {
    case 1: return 'success'  // 审批通过
    case 2: return 'warning'  // 待审批
    default: return 'default'
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 1: return '审批通过'
    case 2: return '待审批'
    default: return '未知'
  }
}

// 登出
const logout = () => {
  userStore.logout()
  router.push('/login')
}
// 页面加载时初始化数据
onMounted(async () => {
  await loadChildren()
  await loadHabits()
  await loadRewards()
  await loadExchangeRecords()
  await loadQuotes()
  await loadVocabularies()
  await loadBooks()
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
.exchanges-table,
.quotes-table {
  margin-top: 20px;
}

.dialog-footer {
  text-align: right;
}

.quotes-toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.import-hint {
  margin: 5px 0;
  color: #606266;
  font-size: 14px;
}

.import-example {
  margin: 3px 0;
  color: #909399;
  font-size: 13px;
  font-family: monospace;
  background-color: #f5f7fa;
  padding: 2px 8px;
  border-radius: 4px;
}

.books-toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.books-table {
  margin-top: 20px;
}
</style>
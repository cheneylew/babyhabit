<template>
  <div class="vocabulary-container">
    <el-card class="vocabulary-header">
      <template #header>
        <div class="header-content">
          <div class="header-left">
            <el-button type="primary" @click="goHome" size="small">
              <el-icon><ArrowLeft /></el-icon>
              返回主页
            </el-button>
            <h2>单词记忆</h2>
             <el-button type="info" @click="showLearningHistory" size="small">
              学习历史
            </el-button>
          </div>
          <div class="header-stats">
           
            <el-statistic :value="`${todayLearnedNewWords}/${todayPlan.newWords || 10}`" title="今日新单词" />
            <el-statistic :value="`${todayReviewedWords}/${todayPlan.reviewWords || (todayReviewedWords > 0 ? todayReviewedWords : 0)}`" title="今日复习" />
          </div>
        </div>
      </template>
      <div class="book-selection">
        <h3>选择教材</h3>
        <el-select
          v-model="selectedBookIds"
          multiple
          filterable
          placeholder="请选择要学习的教材"
          style="width: 100%"
          @change="handleBookSelectionChange"
        >
          <el-option
            v-for="book in bookOptions"
            :key="book.id"
            :label="book.name"
            :value="book.id"
          />
        </el-select>
      </div>
      <div class="plan-summary">
        <p>今天需要学习 {{ todayPlan.newWords }} 个新单词（已学习 {{ todayLearnedNewWords }} 个），复习 {{ todayPlan.reviewWords || (todayReviewedWords > 0 ? todayReviewedWords : 0) }} 个单词（已复习 {{ todayReviewedWords }} 个）</p>
        <div style="display: flex; gap: 10px; justify-content: center;">
          <el-button type="primary" size="large" @click="startLearning" :disabled="isLearning">
            {{ isLearning ? '学习中...' : '开始学习' }}
          </el-button>
          <el-button type="success" size="large" @click="startDictation">
            默写单词
          </el-button>
        </div>
      </div>
    </el-card>

    <el-dialog
      v-model="historyDialogVisible"
      title="学习历史"
      width="90%"
      :close-on-click-modal="false"
    >
      <div v-if="learningHistory.length > 0" class="history-list">
        <el-table :data="learningHistory" style="width: 100%">
          <el-table-column prop="english" label="英文" min-width="100">
            <template #default="scope">
              <div class="truncate-text">{{ scope.row.english }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="chinese" label="中文" min-width="150">
            <template #default="scope">
              <div class="truncate-text">{{ scope.row.chinese }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="learnDate" label="学习日期" width="180" />
        </el-table>
      </div>
      <div v-else class="empty-history">
        <el-empty description="暂无学习历史" />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="historyDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      v-model="learningDialogVisible"
      title="单词学习"
      width="90%"
      :close-on-click-modal="false"
      @close="closeLearning"
    >
      <div v-if="currentWord" class="word-card">
        <div class="word-info" v-show="currentStage === 'recognition' || showWordInfoInMemory">
          <!-- 新单词/复习单词标签 -->
          <div class="word-type-tag" v-if="currentWord.is_new">
            <el-tag type="success" size="small">新单词</el-tag>
          </div>
          <div class="word-type-tag" v-else>
            <el-tag type="warning" size="small">复习单词</el-tag>
          </div>
          
          <h3 class="word-english" v-show="learningMode !== 'english'">{{ currentWord.english }}</h3>
          <div style="display: flex; align-items: center; gap: 10px;" v-show="learningMode !== 'english'">
            <p class="word-phonetic" style="margin: 0;">{{ formatPhonetic(currentWord.phonetic) }}</p>
            <el-button type="info" size="small" @click="playAudio" :loading="isPlaying" style="margin: 0;">
              播放
            </el-button>
          </div>
          <p class="word-chinese" v-show="learningMode !== 'chinese'">{{ currentWord.chinese }}</p>
          <div class="word-audio">
            <el-button size="small" :type="learningMode === 'normal' ? 'primary' : ''" @click="learningMode = 'normal'">正常</el-button>
            <el-button size="small" :type="learningMode === 'english' ? 'primary' : ''" @click="learningMode = 'english'">默英文</el-button>
            <el-button size="small" :type="learningMode === 'chinese' ? 'primary' : ''" @click="learningMode = 'chinese'">默中文</el-button>
          </div>
          <div class="word-example" v-if="currentWord && currentWord.example_sentence && learningMode === 'normal'">
            <div v-if="currentWord && typeof currentWord.example_sentence === 'string'">
              <div v-if="currentWord && isJsonString(currentWord.example_sentence)">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                  <h4 style="margin: 0;">例句：</h4>
                  <el-button type="text" size="small" @click="showAllExamples = !showAllExamples" v-if="currentWord && currentWord.example_sentence && parseJson(currentWord.example_sentence).length > 1">
                    {{ showAllExamples ? '收起' : '更多' }}
                  </el-button>
                </div>
                <div v-for="(example, index) in parseJson(currentWord.example_sentence)" :key="index" class="example-item" v-show="index === 0 || showAllExamples">
                  <div class="example-english-container">
                    <span v-html="processSentenceWords(highlightWordInSentence(example.english, currentWord.english))"></span>
                    <el-button type="primary" size="small" @click="playExampleAudio(example.english)">播放</el-button>
                  </div>
                  <p class="example-chinese">{{ example.chinese }}</p>
                </div>
              </div>
              <div v-else>
                <h4 style="margin: 0 0 10px 0;">例句：</h4>
                <p v-html="currentWord && highlightWordInSentence(currentWord.example_sentence, currentWord.english)"></p>
              </div>
            </div>
            <div v-else-if="currentWord && Array.isArray(currentWord.example_sentence)">
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                <h4 style="margin: 0;">例句：</h4>
                <el-button type="text" size="small" @click="showAllExamples = !showAllExamples" v-if="currentWord && currentWord.example_sentence.length > 1">
                  {{ showAllExamples ? '收起' : '更多' }}
                </el-button>
              </div>
              <div v-for="(example, index) in currentWord.example_sentence" :key="index" class="example-item" v-show="index === 0 || showAllExamples">
                <div class="example-english-container">
                  <span v-html="processSentenceWords(highlightWordInSentence(example.english, currentWord.english))"></span>
                  <el-button type="primary" size="small" @click="playExampleAudio(example.english)">播放</el-button>
                </div>
                <p class="example-chinese">{{ example.chinese }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="learning-stage" v-if="currentStage === 'memory'">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px;">
            <h4>记忆检测</h4>
            <el-button type="info" size="small" @click="showWordInfoInMemory = !showWordInfoInMemory">
              {{ showWordInfoInMemory ? '隐藏' : '回看' }}
            </el-button>
          </div>
          <div v-if="currentWord && currentCheckType === 'chineseToEnglish'" class="check-section">
            <p class="check-prompt">请根据中文选择正确的英文单词：</p>
            <p class="check-question">{{ currentWord.chinese }}</p>
            <el-radio-group v-model="userAnswer" class="check-options">
              <el-radio v-for="(option, index) in answerOptions" :key="index" :label="option">{{ option }}</el-radio>
            </el-radio-group>
          </div>
          <div v-else-if="currentWord && currentCheckType === 'englishToChinese'" class="check-section">
            <p class="check-prompt">请根据英文选择正确的中文意思：</p>
            <p class="check-question">{{ currentWord.english }}</p>
            <el-radio-group v-model="userAnswer" class="check-options">
              <el-radio v-for="(option, index) in answerOptions" :key="index" :label="option">{{ option }}</el-radio>
            </el-radio-group>
          </div>
          <div v-else-if="currentCheckType === 'listening'" class="check-section">
            <p class="check-prompt">请听发音选择正确的单词：</p>
            <el-button type="info" @click="playAudio" class="listen-button">
              <el-icon><i-ep-audio /></el-icon>
              再听一遍
            </el-button>
            <el-radio-group v-model="userAnswer" class="check-options">
              <el-radio v-for="(option, index) in answerOptions" :key="index" :label="option">{{ option }}</el-radio>
            </el-radio-group>
          </div>
        </div>

        <div class="learning-feedback" v-if="showFeedback">
          <el-alert
            :title="isCorrect ? '正确！' : '错误！'"
            :type="isCorrect ? 'success' : 'error'"
            show-icon
          />
          <p v-if="!isCorrect" class="correct-answer">正确答案：{{ correctAnswer }}</p>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button v-if="currentStage === 'recognition'" type="warning" @click="skipWord">
            跳过
          </el-button>
          <el-button @click="closeLearning">退出学习</el-button>
          <el-button v-if="currentStage === 'recognition'" type="primary" @click="nextStage">
            开始记忆检测
          </el-button>
          <el-button v-else-if="!showFeedback" type="primary" @click="submitAnswer">
            提交答案
          </el-button>
          <el-button v-else-if="isCorrect" type="primary" @click="nextWord">
            下一个
          </el-button>
          <el-button v-else type="warning" @click="nextWord">
            重新尝试
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-card class="vocabulary-stats">
      <template #header>
        <h3>学习统计</h3>
      </template>
      <div class="stats-grid">
        <el-statistic :value="totalWords" title="累计学习" />
        <el-statistic :value="todayLearnedWords" title="今日学习" />
        <el-statistic :value="masteredWords" title="已掌握" />
        <el-statistic :value="learningStreak" title="连续学习" />
        <el-statistic :value="accuracyRate" title="正确率" :formatter="(value) => value + '%'" />
      </div>
    </el-card>

    <!-- 单词意思对话框 -->
    <el-dialog
      v-model="wordMeaningDialogVisible"
      title="单词意思"
      width="80%"
      :close-on-click-modal="false"
    >
      <div v-loading="isLoading">
        <h3>{{ currentWordForMeaning }}</h3>
        <p>{{ currentWordMeaning }}</p>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="playExampleAudio(currentWordForMeaning)">
            <el-icon><i-ep-audio /></el-icon>
            播放发音
          </el-button>
          <el-button type="primary" @click="closeWordMeaningDialog">我知道了</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 默写单词对话框 -->
    <el-dialog
      v-model="dictationDialogVisible"
      title="默写单词"
      width="90%"
      :close-on-click-modal="false"
      @close="closeDictation"
    >
      <div v-if="currentDictationWord" class="word-card">
        <div class="word-info">
          <!-- 新单词/复习单词标签 -->
          <div class="word-type-tag" v-if="currentDictationWord.is_new">
            <el-tag type="success" size="small">新单词</el-tag>
          </div>
          <div class="word-type-tag" v-else>
            <el-tag type="warning" size="small">复习单词</el-tag>
          </div>
          
          <h3 class="word-english" v-show="dictationMode === 'normal' || dictationMode === 'english' || (dictationMode !== 'chinese' && dictationCheckType !== 'englishToChinese')">{{ currentDictationWord.english }}</h3>
          <div style="display: flex; align-items: center; gap: 10px;" v-show="dictationMode === 'normal' || dictationMode === 'english' || (dictationMode !== 'chinese' && dictationCheckType !== 'englishToChinese')">
            <p class="word-phonetic" style="margin: 0;">{{ formatPhonetic(currentDictationWord.phonetic) }}</p>
            <el-button type="info" size="small" @click="playDictationAudio" :loading="isPlaying" style="margin: 0;">
              播放
            </el-button>
          </div>
          <p class="word-chinese" v-show="dictationMode === 'normal' || dictationMode === 'chinese' || (dictationMode !== 'english' && dictationCheckType !== 'chineseToEnglish')">{{ currentDictationWord.chinese }}</p>
          <div class="word-example" v-if="currentDictationWord && currentDictationWord.example_sentence && dictationMode === 'normal'">
            <div v-if="currentDictationWord && typeof currentDictationWord.example_sentence === 'string'">
              <div v-if="currentDictationWord && isJsonString(currentDictationWord.example_sentence)">
                <h4>例句 <el-button type="text" size="small" @click="showAllExamples = !showAllExamples" v-if="currentDictationWord && currentDictationWord.example_sentence && parseJson(currentDictationWord.example_sentence).length > 1">{{ showAllExamples ? '收起' : '更多' }}</el-button></h4>
                <div v-for="(example, index) in parseJson(currentDictationWord.example_sentence)" :key="index" class="example-item" v-show="index === 0 || showAllExamples">
                  <div class="example-english-container">
                    <span v-html="processSentenceWords(highlightWordInSentence(example.english, currentDictationWord.english))"></span>
                    <el-button type="primary" size="small" @click="playExampleAudio(example.english)">播放</el-button>
                  </div>
                  <p class="example-chinese">{{ example.chinese }}</p>
                </div>
              </div>
              <div v-else>
                <h4>例句</h4>
                <p v-html="currentDictationWord && highlightWordInSentence(currentDictationWord.example_sentence, currentDictationWord.english)"></p>
              </div>
            </div>
            <div v-else-if="currentDictationWord && Array.isArray(currentDictationWord.example_sentence)">
              <h4>例句 <el-button type="text" size="small" @click="showAllExamples = !showAllExamples" v-if="currentDictationWord && currentDictationWord.example_sentence.length > 1">{{ showAllExamples ? '收起' : '更多' }}</el-button></h4>
              <div v-for="(example, index) in currentDictationWord.example_sentence" :key="index" class="example-item" v-show="index === 0 || showAllExamples">
                <div class="example-english-container">
                  <span v-html="processSentenceWords(highlightWordInSentence(example.english, currentDictationWord.english))"></span>
                  <el-button type="primary" size="small" @click="playExampleAudio(example.english)">播放</el-button>
                </div>
                <p class="example-chinese">{{ example.chinese }}</p>
              </div>
            </div>
          </div>
          <div class="word-audio">
            <el-button size="small" :type="dictationMode === 'normal' ? 'primary' : ''" @click="dictationMode = 'normal'">看答案</el-button>
            <el-button size="small" :type="dictationMode === 'english' ? 'primary' : ''" @click="dictationMode = 'english'">只看英文</el-button>
            <el-button size="small" :type="dictationMode === 'chinese' ? 'primary' : ''" @click="dictationMode = 'chinese'">只看中文</el-button>
          </div>
        </div>



        <div class="learning-feedback" v-if="dictationShowFeedback">
          <el-alert
            :title="dictationIsCorrect ? '正确！' : '错误！'"
            :type="dictationIsCorrect ? 'success' : 'error'"
            show-icon
          />
          <p v-if="!dictationIsCorrect" class="correct-answer">正确答案：{{ dictationCorrectAnswer }}</p>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeDictation">退出默写</el-button>
          <el-button v-if="!dictationShowFeedback" type="danger" @click="dictationFailed">
            还不太会
          </el-button>
          <el-button v-if="!dictationShowFeedback" type="success" @click="submitDictationAnswer">
            默写成功
          </el-button>
          <el-button v-else type="primary" @click="nextDictationWord">
            下一个
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import { ElMessage, ElMessageBox, ElLoading, ElIcon } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import api from '../api'

const router = useRouter()
const userStore = useUserStore()
const user = computed(() => userStore.user)

// 学习计划
const todayPlan = ref({
  newWords: 0,
  reviewWords: 0
})
const totalWords = ref(0)
const masteredWords = ref(0)
const learningStreak = ref(0)
const accuracyRate = ref('0.0')
const todayLearnedWords = ref(0)
const todayLearnedNewWords = ref(0)
const todayReviewedWords = ref(0)

// 教材选择
const bookOptions = ref([])
const selectedBookIds = ref([])
const isLoadingBooks = ref(false)

// 学习状态
const isLearning = ref(false)
const learningDialogVisible = ref(false)
const currentWord = ref(null)
const currentStage = ref('recognition') // recognition, memory
const currentCheckType = ref('chineseToEnglish') // chineseToEnglish, englishToChinese, listening
const userAnswer = ref('')
const answerOptions = ref([])
const correctAnswer = ref('')
const showFeedback = ref(false)
const isCorrect = ref(false)
const isPlaying = ref(false)

// 学习模式：normal (正常), english (默英文), chinese (默中文)
const learningMode = ref('normal')

// 控制例句显示
const showAllExamples = ref(false)

// 控制记忆检测时是否显示单词信息
const showWordInfoInMemory = ref(false)

// 加载教材选项
const loadBookOptions = async () => {
  try {
    isLoadingBooks.value = true
    const response = await api.get('/vocabulary/options/book')
    bookOptions.value = response.data.books || []
  } catch (error) {
    console.error('Failed to load book options:', error)
  } finally {
    isLoadingBooks.value = false
  }
}

// 加载用户教材偏好
const loadUserBookPreferences = async () => {
  try {
    const response = await api.get('/user/preference', {
      params: { key: 'selected_books' }
    })
    if (response.data.preference && response.data.preference.preference_value) {
      selectedBookIds.value = JSON.parse(response.data.preference.preference_value)
    }
  } catch (error) {
    console.error('Failed to load user book preferences:', error)
    // 降级到localStorage
    const savedPreferences = localStorage.getItem('userBookPreferences')
    if (savedPreferences) {
      selectedBookIds.value = JSON.parse(savedPreferences)
    }
  }
}

// 保存用户教材偏好
const saveUserBookPreferences = async () => {
  try {
    await api.post('/user/preference', {
      key: 'selected_books',
      value: JSON.stringify(selectedBookIds.value)
    })
    // 同时保存到localStorage作为备份
    localStorage.setItem('userBookPreferences', JSON.stringify(selectedBookIds.value))
  } catch (error) {
    console.error('Failed to save user book preferences:', error)
    // 降级到localStorage
    localStorage.setItem('userBookPreferences', JSON.stringify(selectedBookIds.value))
  }
}

// 处理教材选择变化
const handleBookSelectionChange = async () => {
  await saveUserBookPreferences()
  await loadTodayPlan()
}

// 加载今日学习计划
const loadTodayPlan = async () => {
  try {
    const response = await api.get('/vocabulary/plan', {
      params: {
        book_ids: selectedBookIds.value.join(',')
      }
    })
    todayPlan.value = response.data.plan
    totalWords.value = response.data.stats.totalWords
    masteredWords.value = response.data.stats.masteredWords
    learningStreak.value = response.data.stats.learningStreak
    accuracyRate.value = response.data.stats.accuracyRate
    todayLearnedWords.value = response.data.stats.todayLearnedWords || 0
    todayLearnedNewWords.value = response.data.stats.todayLearnedNewWords || 0
    todayReviewedWords.value = response.data.stats.todayReviewedWords || 0
  } catch (error) {
    console.error('Failed to load learning plan:', error)
  }
}



// 加载下一个单词
const loadNextWord = async () => {
  if (!currentWords.value || currentWords.value.length === 0) {
    ElMessageBox.alert(
      '恭喜你！今天的单词学习任务已经完成，继续保持！',
      '学习完成',
      {
        confirmButtonText: '确定',
        type: 'success'
      }
    )
    closeLearning()
    await loadTodayPlan()
    return
  }
  
  currentWord.value = currentWords.value.shift()
  currentStage.value = 'recognition'
  showFeedback.value = false
  userAnswer.value = ''
  showAllExamples.value = false
  showWordInfoInMemory.value = false
  
  // 生成选项
  await generateAnswerOptions()
  
  // 自动播放发音三遍
  setTimeout(() => {
    playAudio()
    setTimeout(() => {
      // playAudio()
      setTimeout(() => {
        // playAudio()
      }, 2000)
    }, 2000)
  }, 500)
}

// 生成答案选项
const generateAnswerOptions = async () => {
  try {
    const response = await api.get('/vocabulary/options', {
      params: { wordId: currentWord.value.id, type: currentCheckType.value }
    })
    answerOptions.value = response.data.options
    
    // 设置正确答案
    if (currentCheckType.value === 'chineseToEnglish' || currentCheckType.value === 'listening') {
      correctAnswer.value = currentWord.value.english
    } else if (currentCheckType.value === 'englishToChinese') {
      correctAnswer.value = currentWord.value.chinese
    }
  } catch (error) {
    console.error('Failed to generate options:', error)
  }
}

// 进入记忆阶段
const nextStage = () => {
  currentStage.value = 'memory'
  showWordInfoInMemory.value = false
  // 随机选择检测类型
  const checkTypes = ['chineseToEnglish', 'englishToChinese', 'listening']
  currentCheckType.value = checkTypes[Math.floor(Math.random() * checkTypes.length)]
  generateAnswerOptions()
}

// 提交答案
const submitAnswer = async () => {
  if (!userAnswer.value) {
    ElMessage.warning('请选择答案')
    return
  }
  
  isCorrect.value = userAnswer.value === correctAnswer.value
  
  // 只有当记忆检测通过后，才记录学习结果
  if (isCorrect.value) {
    try {
      await api.post('/vocabulary/record', {
        wordId: currentWord.value.id,
        isCorrect: isCorrect.value,
        checkType: currentCheckType.value
      })
      // 正确时直接进入下一个单词
      await nextWord()
    } catch (error) {
      console.error('Failed to record learning result:', error)
      // 出错时显示反馈
      showFeedback.value = true
    }
  } else {
    // 错误时显示反馈
    showFeedback.value = true
  }
}

// 下一个单词
// 存储当前学习的单词列表
const currentWords = ref([])

// 学习历史相关
const historyDialogVisible = ref(false)
const learningHistory = ref([])

const startLearning = async () => {
  try {
    const response = await api.get('/vocabulary/start', {
      params: {
        book_ids: selectedBookIds.value.join(',')
      }
    })
    if (!response.data.words || response.data.words.length === 0) {
      ElMessageBox.alert(
        '恭喜你！今天的单词学习任务已经完成，继续保持！',
        '学习完成',
        {
          confirmButtonText: '确定',
          type: 'success'
        }
      )
      closeLearning()
      return false
    }
    
    // 存储当前学习的单词列表
    currentWords.value = [...response.data.words]
    
    // 恢复到正常模式
    learningMode.value = 'normal'
    isLearning.value = true
    learningDialogVisible.value = true
    currentStage.value = 'recognition'
    currentCheckType.value = 'chineseToEnglish'
    showFeedback.value = false
    
    // 开始学习第一个单词
    await loadNextWord()
    return true
  } catch (error) {
    console.error('Failed to start learning:', error)
    ElMessage.error('开始学习失败：' + (error.response?.data?.error || error.message))
    return false
  }
}

const nextWord = async () => {
  // 只有当记忆检测通过后，才继续加载下一个单词
  if (isCorrect.value) {
    // 恢复到正常模式
    learningMode.value = 'normal'
    // 继续加载下一个单词
    await loadNextWord()
  } else {
    // 记忆检测失败，重新开始当前单词的记忆检测
    currentStage.value = 'memory'
    // 随机选择检测类型
    const checkTypes = ['chineseToEnglish', 'englishToChinese', 'listening']
    currentCheckType.value = checkTypes[Math.floor(Math.random() * checkTypes.length)]
    generateAnswerOptions()
    showFeedback.value = false
    userAnswer.value = ''
  }
}

// 跳过单词（标记为掌握）
const skipWord = async () => {
  try {
    // 二次确认对话框
    await ElMessageBox.confirm(
      '确定要跳过这个单词吗？跳过将标记为已掌握，后续不会再重复学习。',
      '跳过单词',
      {
        confirmButtonText: '确定跳过',
        cancelButtonText: '取消',
        type: 'warning',
        distinguishCancelAndClose: true
      }
    )

    // 标记单词为掌握
    await api.post('/vocabulary/record', {
      wordId: currentWord.value.id,
      isCorrect: true,
      checkType: 'skip',
      mastered: true
    })
    
    // 从当前单词列表中移除刚刚跳过的单词
    currentWords.value = currentWords.value.filter(word => word.id !== currentWord.value.id)
    
    // 如果是新单词，获取新的单词来补充
    if (currentWord.value.is_new) {
      // 获取新的单词列表
      const response = await api.get('/vocabulary/start', {
        params: {
          book_ids: selectedBookIds.value.join(',')
        }
      })
      
      if (response.data.words && response.data.words.length > 0) {
        // 过滤掉已经学习过的单词
        const newWords = response.data.words.filter(word => 
          !currentWords.value.some(w => w.id === word.id) && 
          word.id !== currentWord.value.id
        )
        
        // 将新单词添加到当前列表
        currentWords.value = [...currentWords.value, ...newWords]
      }
    }
    
    // 加载下一个单词
    await loadNextWord()
  } catch (error) {
    // 如果是用户取消，不显示错误信息
    if (error !== 'cancel') {
      console.error('Failed to skip word:', error)
      ElMessage.error('跳过单词失败：' + (error.response?.data?.error || error.message))
    }
  }
}

// 播放发音
const playAudio = () => {
  if (currentWord.value?.audio_url) {
    const audio = new Audio(currentWord.value.audio_url)
    audio.play()
    isPlaying.value = true
    audio.onended = () => {
      isPlaying.value = false
    }
  }
}

// 播放例句音频
const playExampleAudio = async (sentence) => {
  console.log('playExampleAudio called with sentence:', sentence)
  try {
    // 直接调用后端API，让后端处理文件存在性检查和生成
    console.log('Calling API to generate or get audio...')
    const response = await api.post('/vocabulary/generate-sentence-audio', {
      sentence: sentence
    })
    console.log('API response:', response)
    
    if (response.data.audio_url) {
      // 播放生成的音频
      console.log('Playing audio...')
      const audio = new Audio(response.data.audio_url)
      audio.play()
      console.log('Audio played successfully')
    } else {
      console.error('API returned no audio_url')
      ElMessage.error('生成音频失败：没有返回音频URL')
    }
  } catch (apiError) {
    console.error('API request failed:', apiError)
    ElMessage.error('调用API失败：' + (apiError.message || '未知错误'))
  }
}

// 检查音频文件是否存在
const checkAudioFileExists = async (url) => {
  try {
    const response = await fetch(url, { method: 'HEAD' })
    return response.ok
  } catch (error) {
    return false
  }
}

// 计算SHA-256哈希值
const calculateMD5 = async (text) => {
  const encoder = new TextEncoder()
  const data = encoder.encode(text)
  const hashBuffer = await crypto.subtle.digest('SHA-256', data)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
}

// 检查字符串是否为JSON格式
const isJsonString = (str) => {
  try {
    JSON.parse(str)
    return true
  } catch (error) {
    return false
  }
}

// 解析JSON字符串
const parseJson = (str) => {
  if (!str) {
    return []
  }
  try {
    return JSON.parse(str)
  } catch (error) {
    return []
  }
}

// 处理音标显示
const formatPhonetic = (phonetic) => {
  if (!phonetic) return ''
  try {
    const phoneticObj = JSON.parse(phonetic)
    if (phoneticObj.uk && phoneticObj.us) {
      return `英式: ${phoneticObj.uk} / 美式: ${phoneticObj.us}`
    } else if (phoneticObj.uk) {
      return `英式: ${phoneticObj.uk}`
    } else if (phoneticObj.us) {
      return `美式: ${phoneticObj.us}`
    }
  } catch (error) {
    // 如果不是JSON格式，直接返回
    return phonetic
  }
  return phonetic
}

// 高亮例句中的单词
const highlightWordInSentence = (sentence, word) => {
  if (!sentence || !word) return sentence
  console.log('highlightWordInSentence called with:', { sentence, word })
  const regex = new RegExp(`\\b${word}\\b`, 'gi')
  // 使用内联样式确保高亮效果
  const result = sentence.replace(regex, '<span class="highlight-word" style="background-color: #ff4d4f; color: white; padding: 0 4px; border-radius: 4px; font-weight: bold; display: inline;">$&</span>')
  console.log('highlightWordInSentence result:', result)
  return result
}

// 为例句中的所有单词添加点击事件
const processSentenceWords = (sentence) => {
  if (!sentence) return sentence
  
  // 处理已经高亮的单词，为其添加点击事件
  let processedSentence = sentence.replace(/(<span class="highlight-word"[^>]*>)([^<]+)(<\/span>)/g, (match, openTag, word, closeTag) => {
    // 在highlight-word的span标签上直接添加点击事件属性
    const newOpenTag = openTag.replace('class="highlight-word"', 'class="highlight-word clickable-word"').replace('style="', 'data-word="' + word + '" style="cursor: pointer;')
    return `${newOpenTag}${word}${closeTag}`
  })
  
  // 对非HTML标签部分的单词添加点击事件
  // 使用正则表达式，避免对HTML标签的误处理
  // 匹配HTML标签和纯文本部分
  return processedSentence.replace(/(<[^>]+>)|(\b\w+\b)/g, (match, tag, word) => {
    // 如果是HTML标签，直接返回
    if (tag) {
      return tag
    }
    // 如果是单词，并且不是已经处理过的（不在highlight-word标签内）
    if (word && !match.includes('clickable-word') && !match.includes('highlight-word')) {
      return `<span class="clickable-word" data-word="${word}" style="cursor: pointer; text-decoration: underline;">${word}</span>`
    }
    return match
  })
}

// 加载状态
const isLoading = ref(false)

// 单词意思对话框
const wordMeaningDialogVisible = ref(false)
const currentWordMeaning = ref('')
const currentWordForMeaning = ref('')

// 默写单词相关状态
const dictationDialogVisible = ref(false)
const currentDictationWord = ref(null)
const dictationMode = ref('normal') // normal, english, chinese
const dictationCheckType = ref('chineseToEnglish') // chineseToEnglish, englishToChinese
const dictationUserAnswer = ref('')
const dictationCorrectAnswer = ref('')
const dictationShowFeedback = ref(false)
const dictationIsCorrect = ref(false)
const dictationWords = ref([])

// 查询单词的中文意思
const getWordMeaning = async (word) => {
  let loadingInstance = null
  try {
    console.log('查询单词意思:', word)
    // 显示全局loading动画
    loadingInstance = ElLoading.service({
      lock: true,
      text: '查询单词意思中...',
      background: 'rgba(0, 0, 0, 0.7)',
    })
    // 调用后端API查询单词意思
    const response = await api.post('/vocabulary/get-word-meaning', {
      word: word
    })
    if (response.data.meaning) {
      // 显示单词意思对话框
      currentWordForMeaning.value = word
      currentWordMeaning.value = response.data.meaning
      wordMeaningDialogVisible.value = true
      // 自动播放单词发音
      playExampleAudio(word)
    } else {
      ElMessage.error('查询单词意思失败')
    }
  } catch (error) {
    console.error('查询单词意思失败:', error)
    ElMessage.error('查询单词意思失败')
  } finally {
    // 关闭loading动画
    if (loadingInstance) {
      loadingInstance.close()
    }
  }
}

// 关闭单词意思对话框
const closeWordMeaningDialog = () => {
  wordMeaningDialogVisible.value = false
  currentWordMeaning.value = ''
  currentWordForMeaning.value = ''
}

// 关闭学习
const closeLearning = async () => {
  learningDialogVisible.value = false
  isLearning.value = false
  currentWord.value = null
  // 重新加载学习计划数据
  await loadTodayPlan()
}

// 显示学习历史
const showLearningHistory = async () => {
  try {
    const response = await api.get('/vocabulary/history')
    learningHistory.value = response.data.words.map(word => ({
      ...word,
      learnDate: new Date(word.create_time).toLocaleString()
    }))
    historyDialogVisible.value = true
  } catch (error) {
    console.error('Failed to load learning history:', error)
    ElMessage.error('加载学习历史失败：' + (error.response?.data?.error || error.message))
  }
}

// 返回主页
const goHome = () => {
  router.push('/')
}

// 事件处理函数
const handleWordClick = (event) => {
  // 防止事件冒泡
  event.stopPropagation()
  const target = event.target
  if (target.classList.contains('clickable-word')) {
    const word = target.getAttribute('data-word')
    if (word) {
      getWordMeaning(word)
    }
  }
}

// 开始默写单词
const startDictation = async () => {
  try {
    // 先获取需要默写的单词
    const response = await api.get('/vocabulary/dictation', {
      params: {
        book_ids: selectedBookIds.value.join(',')
      }
    })
    if (!response.data.words || response.data.words.length === 0) {
      ElMessage.info('今天没有需要默写的单词')
      return
    }
    
    // 弹出选择对话框
    ElMessageBox.confirm(
      '请选择默写模式',
      '默写单词',
      {
        confirmButtonText: '根据中文默英文',
        cancelButtonText: '根据英文默中文',
        type: 'info',
        distinguishCancelAndClose: true
      }
    ).then(() => {
      // 用户选择根据中文默英文
      startDictationWithMode(response.data.words, 'chineseToEnglish')
    }).catch((action) => {
      // 用户选择根据英文默中文
      if (action === 'cancel') {
        startDictationWithMode(response.data.words, 'englishToChinese')
      }
    })
  } catch (error) {
    console.error('Failed to start dictation:', error)
    ElMessage.error('开始默写失败')
  }
}

// 开始默写单词（指定模式）
const startDictationWithMode = async (words, checkType) => {
  // 根据选择的默写模式设置默认显示模式
  if (checkType === 'chineseToEnglish') {
    // 根据中文默英文，默认只看中文
    dictationMode.value = 'chinese'
  } else if (checkType === 'englishToChinese') {
    // 根据英文默中文，默认只看英文
    dictationMode.value = 'english'
  } else {
    // 默认正常模式
    dictationMode.value = 'normal'
  }
  
  dictationDialogVisible.value = true
  dictationShowFeedback.value = false
  dictationWords.value = words
  
  // 开始默第一个单词
  await loadNextDictationWordWithMode(checkType)
}

// 加载下一个默写单词（指定模式）
const loadNextDictationWordWithMode = async (checkType) => {
  if (!dictationWords.value || dictationWords.value.length === 0) {
    ElMessage.info('没有需要默写的单词')
    closeDictation()
    return
  }
  
  currentDictationWord.value = dictationWords.value.shift()
  dictationShowFeedback.value = false
  dictationUserAnswer.value = ''
  
  // 使用指定的检测类型
  dictationCheckType.value = checkType
  
  // 设置正确答案
  if (dictationCheckType.value === 'chineseToEnglish') {
    dictationCorrectAnswer.value = currentDictationWord.value.english
  } else if (dictationCheckType.value === 'englishToChinese') {
    dictationCorrectAnswer.value = currentDictationWord.value.chinese
  }
}

// 加载下一个默写单词
const loadNextDictationWord = async () => {
  if (!dictationWords.value || dictationWords.value.length === 0) {
    ElMessage.info('没有需要默写的单词')
    closeDictation()
    return
  }
  
  currentDictationWord.value = dictationWords.value.shift()
  dictationShowFeedback.value = false
  dictationUserAnswer.value = ''
  
  // 随机选择检测类型
  const checkTypes = ['chineseToEnglish', 'englishToChinese']
  dictationCheckType.value = checkTypes[Math.floor(Math.random() * checkTypes.length)]
  
  // 设置正确答案
  if (dictationCheckType.value === 'chineseToEnglish') {
    dictationCorrectAnswer.value = currentDictationWord.value.english
  } else if (dictationCheckType.value === 'englishToChinese') {
    dictationCorrectAnswer.value = currentDictationWord.value.chinese
  }
}

// 播放默单词发音
const playDictationAudio = () => {
  if (currentDictationWord.value?.audio_url) {
    const audio = new Audio(currentDictationWord.value.audio_url)
    audio.play()
    isPlaying.value = true
    audio.onended = () => {
      isPlaying.value = false
    }
  }
}

// 提交默答案
const submitDictationAnswer = async () => {
  // 用户点击默写成功，直接记录成功结果
  dictationIsCorrect.value = true
  dictationShowFeedback.value = true
  
  // 记录默写结果
  try {
    await api.post('/vocabulary/dictation/record', {
      wordId: currentDictationWord.value.id,
      isCorrect: true,
      checkType: dictationCheckType.value
    })
  } catch (error) {
    console.error('Failed to record dictation result:', error)
  }
}

// 默写失败
const dictationFailed = async () => {
  dictationIsCorrect.value = false
  dictationShowFeedback.value = true
  // 显示完整的单词信息，相当于点击了"正常"按钮
  dictationMode.value = 'normal'
  
  // 记录默写失败的单词到艾宾浩斯不会的单词逻辑
  try {
    await api.post('/vocabulary/dictation/record', {
      wordId: currentDictationWord.value.id,
      isCorrect: false,
      checkType: dictationCheckType.value
    })
  } catch (error) {
    console.error('Failed to record dictation result:', error)
  }
}

// 下一个默单词
const nextDictationWord = async () => {
  // 根据检测类型设置相应的显示模式
  if (dictationCheckType.value === 'chineseToEnglish') {
    // 根据中文默英文，应该只看中文
    dictationMode.value = 'chinese'
  } else {
    // 根据英文默中文，应该只看英文
    dictationMode.value = 'english'
  }
  // 继续加载下一个单词，使用当前的检测类型
  await loadNextDictationWordWithMode(dictationCheckType.value)
}

// 关闭默
const closeDictation = () => {
  dictationDialogVisible.value = false
  currentDictationWord.value = null
  dictationWords.value = []
}

onMounted(async () => {
  if (!user.value) {
    await userStore.getUserInfo()
  }
  await loadBookOptions()
  await loadUserBookPreferences()
  await loadTodayPlan()
  
  // 添加事件监听器，使用事件委托处理单词点击事件
  document.addEventListener('click', handleWordClick)
})

onUnmounted(() => {
  // 移除事件监听器，避免重复绑定
  document.removeEventListener('click', handleWordClick)
})
</script>

<style scoped>
.vocabulary-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.vocabulary-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.header-content h2 {
  margin: 0;
  color: #667eea;
}

.header-stats {
  display: flex;
  gap: 30px;
  align-items: center;
}

.truncate-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.word-example .highlight-word {
  background-color: #ff4d4f !important;
  color: white !important;
  padding: 0 4px !important;
  border-radius: 4px !important;
  font-weight: bold !important;
  display: inline !important;
}

.example-english-container {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 5px;
}

.book-selection {
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
  margin-bottom: 15px;
}

.book-selection h3 {
  margin: 0 0 15px 0;
  color: #666;
  font-size: 16px;
}

.plan-summary {
  text-align: center;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
  margin-top: 15px;
}

.plan-summary p {
  font-size: 16px;
  margin-bottom: 20px;
  color: #666;
}

.word-card {
  padding: 10px;
}

.word-info {
  text-align: center;
  margin-bottom: 30px;
  position: relative;
}

.word-type-tag {
  position: absolute;
  top: -10px;
  right: 10px;
}

.word-english {
  font-size: 36px;
  font-weight: bold;
  color: #667eea;
  margin: 0 0 10px 0;
}

.word-phonetic {
  font-size: 18px;
  color: #999;
  margin: 0 0 15px 0;
}

.word-chinese {
  font-size: 24px;
  color: #333;
  margin: 0 0 20px 0;
}

.word-audio {
  margin: 20px 0;
}

.word-example {
  margin-top: 20px;
  padding: 15px;
  background: #f5f5f5;
  border-radius: 8px;
  text-align: left;
}

.word-example h4 {
  margin: 0 0 10px 0;
  color: #666;
}

.learning-stage {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eaeaea;
}

.learning-stage h4 {
  margin: 0 0 20px 0;
  color: #666;
}

.check-section {
  text-align: left;
}

.check-prompt {
  font-size: 16px;
  color: #666;
  margin-bottom: 10px;
}

.check-question {
  font-size: 20px;
  font-weight: bold;
  color: #333;
  margin-bottom: 20px;
}

.check-section {
  text-align: left;
  max-width: 100%;
  overflow: hidden;
}

.check-options {
  margin-top: 20px;
  max-width: 100%;
}

.check-options .el-radio {
  display: block;
  margin-bottom: 10px;
  font-size: 16px;
  word-wrap: break-word;
  overflow-wrap: break-word;
  max-width: 100%;
}

.check-options .el-radio__input {
  float: left;
  margin-right: 8px;
}

.check-options .el-radio__label {
  word-wrap: break-word;
  overflow-wrap: break-word;
  white-space: normal !important;
  display: block;
  margin-left: 30px;
  max-width: 100%;
}

/* 确保对话框内容区域有适当的内边距 */
.word-card {
  padding: 10px;
  max-width: 100%;
  box-sizing: border-box;
  overflow: hidden;
}

/* 确保对话框内容不超出范围 */
.el-dialog__body {
  padding: 0 !important;
  max-width: 100%;
  overflow: hidden;
}

/* 强制选项文字换行 */
.learning-stage .el-radio-group.check-options {
  max-width: 100% !important;
  width: 100% !important;
}

.learning-stage .el-radio-group.check-options .el-radio {
  display: block !important;
  position: relative !important;
  padding-left: 5px !important;
  margin-bottom: 15px !important;
  line-height: 1.4 !important;
  max-width: 100% !important;
  width: 100% !important;
  overflow: hidden !important;
  word-wrap: break-word !important;
  overflow-wrap: break-word !important;
}

.learning-stage .el-radio-group.check-options .el-radio__input {
  position: absolute !important;
  left: 0 !important;
  top: 2px !important;
  margin: 0 !important;
}

.learning-stage .el-radio-group.check-options .el-radio__label {
  display: block !important;
  word-break: break-all !important;
  white-space: normal !important;
  line-height: 1.4 !important;
  max-width: 100% !important;
  width: 100% !important;
  box-sizing: border-box !important;
  overflow: visible !important;
}

/* 确保整个对话框内容不超出范围 */
.el-dialog {
  max-width: 90vw !important;
  box-sizing: border-box !important;
}

.el-dialog__body {
  padding: 0 !important;
  max-width: 100% !important;
  width: 100% !important;
  overflow: hidden !important;
  box-sizing: border-box !important;
}

.word-card {
  padding: 10px !important;
  max-width: 100% !important;
  width: 100% !important;
  box-sizing: border-box !important;
  overflow: hidden !important;
}

.learning-stage {
  max-width: 100% !important;
  width: 100% !important;
  overflow: hidden !important;
  box-sizing: border-box !important;
}

.check-section {
  max-width: 100% !important;
  width: 100% !important;
  overflow: hidden !important;
  box-sizing: border-box !important;
}

.listen-button {
  margin-bottom: 20px;
}

.word-card {
  padding: 10px !important;
  max-width: 100% !important;
  box-sizing: border-box !important;
  overflow: hidden !important;
}

.learning-stage {
  max-width: 100% !important;
  overflow: hidden !important;
  box-sizing: border-box !important;
}

.check-section {
  max-width: 100% !important;
  overflow: hidden !important;
  box-sizing: border-box !important;
}

.check-options {
  max-width: 100% !important;
  overflow: hidden !important;
  box-sizing: border-box !important;
}

.listen-button {
  margin-bottom: 20px;
}

.learning-feedback {
  margin-top: 20px;
  padding: 15px;
  background: #f0f9ff;
  border-radius: 8px;
}

.correct-answer {
  margin-top: 10px;
  font-weight: bold;
  color: #666;
}

.vocabulary-stats {
  margin-top: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.el-statistic {
  text-align: center;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
}

.el-statistic__value {
  font-size: 24px;
  font-weight: bold;
  color: #667eea;
}

.el-statistic__label {
  font-size: 14px;
  color: #666;
  margin-top: 10px;
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-stats {
    width: 100%;
    justify-content: space-between;
  }
  
  .word-english {
    font-size: 28px;
  }
  
  .word-chinese {
    font-size: 20px;
  }
}
</style>
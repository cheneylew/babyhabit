<template>
  <div class="vocabulary-container">
    <el-card class="vocabulary-header">
      <template #header>
        <div class="header-content">
          <h2>艾宾浩斯单词记忆</h2>
          <div class="header-stats">
            <el-statistic :value="todayPlan.newWords" title="今日新单词" />
            <el-statistic :value="todayPlan.reviewWords" title="今日复习" />
          </div>
        </div>
      </template>
      <div class="plan-summary">
        <p>今天需要学习 {{ todayPlan.newWords }} 个新单词，复习 {{ todayPlan.reviewWords }} 个单词</p>
        <el-button type="primary" size="large" @click="startLearning" :disabled="isLearning">
          {{ isLearning ? '学习中...' : '开始学习' }}
        </el-button>
      </div>
    </el-card>

    <el-dialog
      v-model="learningDialogVisible"
      title="单词学习"
      width="80%"
      :close-on-click-modal="false"
    >
      <div v-if="currentWord" class="word-card">
        <div class="word-info">
          <h3 class="word-english">{{ currentWord.english }}</h3>
          <p class="word-phonetic">{{ currentWord.phonetic }}</p>
          <p class="word-chinese">{{ currentWord.chinese }}</p>
          <div class="word-audio">
            <el-button type="info" @click="playAudio" :loading="isPlaying">
              <el-icon><i-ep-audio /></el-icon>
              播放发音
            </el-button>
          </div>
          <div class="word-example" v-if="currentWord.example_sentence">
            <h4>例句：</h4>
            <p>{{ currentWord.example_sentence }}</p>
          </div>
        </div>

        <div class="learning-stage" v-if="currentStage === 'memory'">
          <h4>记忆检测</h4>
          <div v-if="currentCheckType === 'chineseToEnglish'" class="check-section">
            <p class="check-prompt">请根据中文选择正确的英文单词：</p>
            <p class="check-question">{{ currentWord.chinese }}</p>
            <el-radio-group v-model="userAnswer" class="check-options">
              <el-radio v-for="(option, index) in answerOptions" :key="index" :label="option">{{ option }}</el-radio>
            </el-radio-group>
          </div>
          <div v-else-if="currentCheckType === 'englishToChinese'" class="check-section">
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
          <el-button @click="closeLearning">退出学习</el-button>
          <el-button v-if="currentStage === 'recognition'" type="primary" @click="nextStage">
            开始记忆检测
          </el-button>
          <el-button v-else-if="!showFeedback" type="primary" @click="submitAnswer">
            提交答案
          </el-button>
          <el-button v-else type="primary" @click="nextWord">
            下一个
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
        <el-statistic :value="masteredWords" title="已掌握" />
        <el-statistic :value="learningStreak" title="连续学习" />
        <el-statistic :value="accuracyRate + '%'" title="正确率" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import api from '../api'

const router = useRouter()
const userStore = useUserStore()
const user = computed(() => userStore.user)

// 学习计划
const todayPlan = ref({ newWords: 0, reviewWords: 0 })
const totalWords = ref(0)
const masteredWords = ref(0)
const learningStreak = ref(0)
const accuracyRate = ref(0)

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

// 加载今日学习计划
const loadTodayPlan = async () => {
  try {
    const response = await api.get('/vocabulary/plan')
    todayPlan.value = response.data.plan
    totalWords.value = response.data.stats.totalWords
    masteredWords.value = response.data.stats.masteredWords
    learningStreak.value = response.data.stats.learningStreak
    accuracyRate.value = response.data.stats.accuracyRate
  } catch (error) {
    console.error('Failed to load learning plan:', error)
  }
}

// 开始学习
const startLearning = async () => {
  try {
    const response = await api.get('/vocabulary/start')
    if (response.data.words.length === 0) {
      ElMessage.info('今天没有需要学习的单词')
      return
    }
    
    isLearning.value = true
    learningDialogVisible.value = true
    currentStage.value = 'recognition'
    currentCheckType.value = 'chineseToEnglish'
    showFeedback.value = false
    
    // 开始学习第一个单词
    await loadNextWord(response.data.words)
  } catch (error) {
    console.error('Failed to start learning:', error)
    isLearning.value = false
  }
}

// 加载下一个单词
const loadNextWord = async (words) => {
  if (words.length === 0) {
    ElMessage.success('今日学习完成！')
    closeLearning()
    await loadTodayPlan()
    return
  }
  
  currentWord.value = words.shift()
  currentStage.value = 'recognition'
  showFeedback.value = false
  userAnswer.value = ''
  
  // 生成选项
  await generateAnswerOptions()
}

// 生成答案选项
const generateAnswerOptions = async () => {
  try {
    const response = await api.get('/vocabulary/options', {
      params: { wordId: currentWord.value.id, type: currentCheckType.value }
    })
    answerOptions.value = response.data.options
  } catch (error) {
    console.error('Failed to generate options:', error)
  }
}

// 进入记忆阶段
const nextStage = () => {
  currentStage.value = 'memory'
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
  showFeedback.value = true
  
  // 记录学习结果
  try {
    await api.post('/vocabulary/record', {
      wordId: currentWord.value.id,
      isCorrect: isCorrect.value,
      checkType: currentCheckType.value
    })
  } catch (error) {
    console.error('Failed to record learning result:', error)
  }
}

// 下一个单词
const nextWord = async () => {
  // 继续加载下一个单词
  await startLearning()
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

// 关闭学习
const closeLearning = () => {
  learningDialogVisible.value = false
  isLearning.value = false
  currentWord.value = null
}

onMounted(async () => {
  if (!user.value) {
    await userStore.getUserInfo()
  }
  await loadTodayPlan()
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

.header-content h2 {
  margin: 0;
  color: #667eea;
}

.header-stats {
  display: flex;
  gap: 30px;
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
  padding: 20px;
}

.word-info {
  text-align: center;
  margin-bottom: 30px;
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

.check-options {
  margin-top: 20px;
}

.check-options .el-radio {
  display: block;
  margin-bottom: 10px;
  font-size: 16px;
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
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
      @close="closeLearning"
    >
      <div v-if="currentWord" class="word-card">
        <div class="word-info">
          <h3 class="word-english" v-show="learningMode !== 'chinese'">{{ currentWord.english }}</h3>
          <div style="display: flex; align-items: center; gap: 10px;" v-show="learningMode !== 'chinese'">
            <p class="word-phonetic" style="margin: 0;">{{ formatPhonetic(currentWord.phonetic) }}</p>
            <el-button type="info" size="small" @click="playAudio" :loading="isPlaying" style="margin: 0;">
              播放
            </el-button>
          </div>
          <p class="word-chinese" v-show="learningMode !== 'english'">{{ currentWord.chinese }}</p>
          <div class="word-audio">
            <el-button size="small" :type="learningMode === 'normal' ? 'primary' : ''" @click="learningMode = 'normal'">正常</el-button>
            <el-button size="small" :type="learningMode === 'english' ? 'primary' : ''" @click="learningMode = 'english'">默英文</el-button>
            <el-button size="small" :type="learningMode === 'chinese' ? 'primary' : ''" @click="learningMode = 'chinese'">默中文</el-button>
          </div>
          <div class="word-example" v-if="currentWord.example_sentence && learningMode === 'normal'">
            <h4>例句：</h4>
            <div v-if="typeof currentWord.example_sentence === 'string'">
              <div v-if="isJsonString(currentWord.example_sentence)">
                <div v-for="(example, index) in parseJson(currentWord.example_sentence)" :key="index" class="example-item">
                  <div class="example-english-container">
                    <span v-html="processSentenceWords(highlightWordInSentence(example.english, currentWord.english))"></span>
                    <el-button type="primary" size="small" @click="playExampleAudio(example.english)">播放</el-button>
                  </div>
                  <p class="example-chinese">{{ example.chinese }}</p>
                </div>
              </div>
              <div v-else>
                <p v-html="highlightWordInSentence(currentWord.example_sentence, currentWord.english)"></p>
              </div>
            </div>
            <div v-else-if="Array.isArray(currentWord.example_sentence)">
              <div v-for="(example, index) in currentWord.example_sentence" :key="index" class="example-item">
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
      width="50%"
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
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
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
const todayLearnedWords = ref(0)

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

// 加载今日学习计划
const loadTodayPlan = async () => {
  try {
    const response = await api.get('/vocabulary/plan')
    todayPlan.value = response.data.plan
    totalWords.value = response.data.stats.totalWords
    masteredWords.value = response.data.stats.masteredWords
    learningStreak.value = response.data.stats.learningStreak
    accuracyRate.value = response.data.stats.accuracyRate
    todayLearnedWords.value = response.data.stats.todayLearnedWords || 0
  } catch (error) {
    console.error('Failed to load learning plan:', error)
  }
}

// 开始学习
const startLearning = async () => {
  try {
    const response = await api.get('/vocabulary/start')
    if (!response.data.words || response.data.words.length === 0) {
      ElMessage.info('今天没有需要学习的单词')
      return
    }
    
    // 恢复到正常模式
    learningMode.value = 'normal'
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
  if (!words || words.length === 0) {
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
  
  // 只有当记忆检测通过后，才记录学习结果
  if (isCorrect.value) {
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
}

// 下一个单词
const nextWord = async () => {
  // 只有当记忆检测通过后，才继续加载下一个单词
  if (isCorrect.value) {
    // 恢复到正常模式
    learningMode.value = 'normal'
    // 继续加载下一个单词
    await startLearning()
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
const closeLearning = () => {
  learningDialogVisible.value = false
  isLearning.value = false
  currentWord.value = null
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

onMounted(async () => {
  if (!user.value) {
    await userStore.getUserInfo()
  }
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

.header-content h2 {
  margin: 0;
  color: #667eea;
}

.header-stats {
  display: flex;
  gap: 30px;
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
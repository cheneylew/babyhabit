# 英语单词背诵系统 PRD v1.1

## 1. 功能概述
英语单词和句型背诵，使用艾宾浩斯记忆法。每天背诵10~15分钟，利用碎片时间。首页新增一个按钮，进入艾宾浩斯单词记忆法

## 2. 核心设计原则
- **简单实用**：功能聚焦，避免过度设计
- **循序渐进**：从认识到应用，分阶段学习
- **及时反馈**：每次学习都有明确的进度和结果
- **趣味激励**：积分和成就系统，保持学习动力

## 3. 记忆方法（三阶段）

### 3.1 认识阶段
- 看英文单词 + 中文意思 + 音标
- 点击播放发音，跟读2-3遍
- 查看例句理解用法

### 3.2 记忆阶段
- **中英互译**：看中文选英文 / 看英文选中文
- **拼写练习**：看中文写英文单词
- **听力识别**：听发音选单词

### 3.3 巩固阶段
- **快速复习**：随机抽取已学单词快速过一遍
- **错题重练**：针对之前做错的单词重点练习

## 4. 艾宾浩斯复习计划

### 4.1 复习间隔
- 第1次：学习当天
- 第2次：1天后
- 第3次：3天后
- 第4次：7天后
- 第5次：15天后
- 第6次：30天后

### 4.2 单词状态
- **新单词**：初次学习
- **复习中**：按计划在复习周期内
- **已掌握**：连续6次复习正确
- **需加强**：复习时错误，回到复习队列

## 5. 每日学习计划

### 5.1 学习内容
- **新单词**：每天5个（避免认知负荷过重）
- **复习单词**：根据艾宾浩斯计划自动安排
- **总时长**：控制在10-15分钟

### 5.2 学习流程
1. 查看今日学习计划
2. 学习新单词（认识阶段）
3. 复习待复习单词（记忆阶段检测）
4. 查看学习结果和统计

## 6. 前台小孩端功能

### 6.1 今日学习页面
- 显示今日计划：新单词X个 + 复习单词Y个
- 开始学习按钮
- 学习进度条

### 6.2 学习界面
- 单词卡片展示（英文、音标、中文、例句）
- 发音播放按钮
- 三种检测模式切换
- 正确/错误即时反馈
- 下一题按钮

### 6.3 学习统计
- 累计学习天数
- 已掌握单词数
- 今日学习情况
- 学习日历（打卡记录）

### 6.4 个人中心
- 我的积分
- 我的成就
- 学习设置

## 7. 后台管理员端功能

### 7.1 单词管理
- 添加单词（英文、中文、音标、例句、分类）
- **智能生成**：使用字节语音大模型自动生成单词的英式发音：豆包语音合成 Doubao-TTS2.0
- **智能例句**：使用字节大模型自动生成符合语境的例句：豆包大模型Doubao-Seed-2.0-pro
- 批量导入（Excel模板）
- 编辑/删除单词
- 按分类筛选（年级、课本、类型）

### 7.2 学习记录查看
- 查看每个孩子的学习统计
- 单词掌握情况
- 学习活跃度

### 7.3 系统设置
- 每日新单词数量设置
- 复习提醒时间设置

## 8. 数据库设计

### 8.1 单词表
```sql
CREATE TABLE ab_vocabulary (
  id INT PRIMARY KEY AUTO_INCREMENT,
  english VARCHAR(255) NOT NULL,
  chinese VARCHAR(255) NOT NULL,
  phonetic VARCHAR(100),
  audio_url VARCHAR(255), -- 字节语音大模型生成的英式发音
  example_sentence TEXT, -- 字节大模型生成的例句
  type ENUM('word', 'sentence') DEFAULT 'word',
  grade VARCHAR(50),
  textbook VARCHAR(100),
  category VARCHAR(100),
  create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 8.2 学习记录表
```sql
CREATE TABLE ab_learning_record (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  vocabulary_id INT NOT NULL,
  status ENUM('new', 'learning', 'mastered', 'reviewing') DEFAULT 'new',
  review_stage INT DEFAULT 0, -- 0-6对应6次复习
  next_review_date DATE,
  correct_count INT DEFAULT 0,
  wrong_count INT DEFAULT 0,
  create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY (user_id, vocabulary_id)
);
```

### 8.3 学习打卡表
```sql
CREATE TABLE ab_study_checkin (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  checkin_date DATE NOT NULL,
  new_words_count INT DEFAULT 0,
  review_words_count INT DEFAULT 0,
  correct_count INT DEFAULT 0,
  total_count INT DEFAULT 0,
  duration_minutes INT,
  create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY (user_id, checkin_date)
);
```

## 9. 核心算法

### 9.1 计算下次复习时间
```go
func calculateNextReview(stage int, lastReview time.Time) time.Time {
    intervals := []int{1, 1, 3, 7, 15, 30} // 天数
    if stage >= len(intervals) {
        stage = len(intervals) - 1
    }
    return lastReview.AddDate(0, 0, intervals[stage])
}
```

### 9.2 生成每日学习计划
```go
func generateDailyPlan(userID int) ([]int, []int) {
    today := time.Now().Format("2006-01-02")
    
    // 获取今天需要复习的单词
    reviewWords := getDueReviewWords(userID, today)
    
    // 获取新单词（如果今天还没学够5个）
    learnedToday := getTodayLearnedCount(userID)
    newWordsCount := 5 - learnedToday
    if newWordsCount < 0 {
        newWordsCount = 0
    }
    newWords := getNewWords(userID, newWordsCount)
    
    return newWords, reviewWords
}
```

### 9.3 更新单词状态
```go
func updateWordStatus(recordID int, isCorrect bool) {
    record := getLearningRecord(recordID)
    
    if isCorrect {
        record.CorrectCount++
        record.ReviewStage++
        
        if record.ReviewStage >= 6 {
            record.Status = "mastered"
        } else {
            record.Status = "reviewing"
            record.NextReviewDate = calculateNextReview(record.ReviewStage, time.Now())
        }
    } else {
        record.WrongCount++
        record.ReviewStage = 0 // 错误后回到第一阶段
        record.Status = "reviewing"
        record.NextReviewDate = time.Now().AddDate(0, 0, 1) // 第二天复习
    }
    
    saveLearningRecord(record)
}
```

## 10. 激励机制

### 10.1 积分规则
- 学习1个新单词：+5分
- 复习1个单词正确：+3分
- 连续学习7天：+20分
- 掌握1个单词（6次复习完成）：+10分

### 10.2 成就系统
- **初出茅庐**：累计学习10个单词
- **坚持不懈**：连续学习7天
- **小有成就**：掌握50个单词
- **单词达人**：掌握100个单词
- **学习之星**：连续学习30天

## 11. 界面设计要点

### 11.1 学习界面
- 简洁清晰的单词卡片
- 大按钮，方便点击
- 即时反馈（正确绿色，错误红色）
- 进度条显示剩余单词数

### 11.2 统计界面
- 直观的数字展示
- 简单的柱状图或折线图
- 日历形式的打卡记录

## 12. 技术实现要点

### 12.1 前端
- 单词卡片滑动切换
- 本地缓存学习进度（防止意外退出）
- 音频预加载

### 12.2 后端
- 每日凌晨生成学习计划（定时任务）
- 学习记录批量提交（减少请求次数）
- **智能语音集成**：调用字节语音大模型API生成单词的英式发音
- **智能例句生成**：调用字节大模型API为单词生成符合语境的例句
- **音频文件管理**：存储和管理生成的语音文件，确保快速加载

## 13. 后续可扩展功能（v2.0）
- 单词拼写游戏
- 学习排行榜
- 家长端查看功能
- 智能推荐薄弱单词

---

**设计原则总结**：
1. **先完成核心功能**，再考虑扩展
2. **保持界面简洁**，减少认知负担
3. **数据驱动**，通过实际使用情况迭代优化

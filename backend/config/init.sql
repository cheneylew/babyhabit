-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS babyhabit CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE babyhabit;

-- 用户表
CREATE TABLE IF NOT EXISTS user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(11),
    email VARCHAR(100),
    avatar VARCHAR(255),
    register_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login_time DATETIME,
    last_activity_time DATETIME,
    user_type TINYINT NOT NULL COMMENT '1:父母，2:小孩',
    status TINYINT DEFAULT 1 COMMENT '1:正常，0:禁用',
    parent_id BIGINT COMMENT '小孩账号必填，父母账号为NULL',
    login_fail_count TINYINT DEFAULT 0,
    locked_until DATETIME,
    points_balance INT DEFAULT 0,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_phone (phone) WHERE phone IS NOT NULL
);

-- 积分记录表
CREATE TABLE IF NOT EXISTS points_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    points SMALLINT NOT NULL,
    source VARCHAR(50) NOT NULL,
    related_type VARCHAR(20),
    related_id BIGINT,
    expire_time DATETIME,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id)
);

-- 习惯表
CREATE TABLE IF NOT EXISTS habit (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    category VARCHAR(50),
    schedule_type TINYINT NOT NULL COMMENT '1:每日，2:周期性',
    schedule_detail VARCHAR(255),
    checkin_time_start TIME,
    checkin_time_end TIME,
    reward_points SMALLINT DEFAULT 0,
    allow_makeup TINYINT DEFAULT 0 COMMENT '1:是，0:否',
    makeup_days TINYINT DEFAULT 0,
    creator_id BIGINT NOT NULL,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status TINYINT DEFAULT 1 COMMENT '1:启用，0:禁用',
    FOREIGN KEY (creator_id) REFERENCES user(id)
);

-- 习惯分配表
CREATE TABLE IF NOT EXISTS habit_assignment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    habit_id INT NOT NULL,
    child_id INT NOT NULL,
    assign_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TINYINT DEFAULT 1 COMMENT '1:有效，0:已取消'
);

-- 打卡记录表
CREATE TABLE IF NOT EXISTS checkin_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    habit_id BIGINT NOT NULL,
    assignment_id BIGINT,
    checkin_date DATE NOT NULL,
    checkin_time DATETIME NOT NULL,
    checkin_type TINYINT DEFAULT 1 COMMENT '1:正常打卡，2:补卡',
    remark VARCHAR(255),
    status TINYINT DEFAULT 1 COMMENT '1:成功，0:失败',
    points_rewarded SMALLINT DEFAULT 0,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (habit_id) REFERENCES habit(id),
    UNIQUE KEY unique_checkin (user_id, habit_id, checkin_date)
);

-- 奖励物品表
CREATE TABLE IF NOT EXISTS reward_item (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image VARCHAR(255),
    category VARCHAR(50),
    points_required INT NOT NULL,
    stock INT DEFAULT -1 COMMENT '-1表示不限量',
    exchange_limit INT DEFAULT 0 COMMENT '0表示不限',
    user_exchanged INT DEFAULT 0,
    start_time DATETIME,
    end_time DATETIME,
    creator_id BIGINT NOT NULL,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status TINYINT DEFAULT 1 COMMENT '1:上架，0:下架',
    FOREIGN KEY (creator_id) REFERENCES user(id)
);

-- 兑换记录表
CREATE TABLE IF NOT EXISTS exchange_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    points INT NOT NULL,
    quantity INT DEFAULT 1,
    exchange_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    delivery_info JSON,
    status TINYINT DEFAULT 2 COMMENT '1:已完成，2:处理中，3:已发货，4:已完成收货',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (item_id) REFERENCES reward_item(id)
);

-- 连续打卡表
CREATE TABLE IF NOT EXISTS streak_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    habit_id BIGINT NOT NULL,
    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,
    last_checkin_date DATE,
    streak_start_date DATE,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (habit_id) REFERENCES habit(id),
    UNIQUE KEY unique_streak (user_id, habit_id)
);

-- 成就表
CREATE TABLE IF NOT EXISTS achievement (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    condition_type VARCHAR(50) NOT NULL,
    condition_value INT NOT NULL,
    reward_points INT DEFAULT 0,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户成就表
CREATE TABLE IF NOT EXISTS user_achievement (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    achievement_id BIGINT NOT NULL,
    unlocked TINYINT DEFAULT 0 COMMENT '1:已解锁，0:未解锁',
    unlocked_time DATETIME,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (achievement_id) REFERENCES achievement(id),
    UNIQUE KEY unique_user_achievement (user_id, achievement_id)
);

-- 通知消息表
CREATE TABLE IF NOT EXISTS notification (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    related_type VARCHAR(20),
    related_id BIGINT,
    is_read TINYINT DEFAULT 0 COMMENT '1:已读，0:未读',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id)
);

-- 习惯提醒表
CREATE TABLE IF NOT EXISTS habit_reminder (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    habit_id BIGINT NOT NULL,
    reminder_time TIME NOT NULL,
    is_enabled TINYINT DEFAULT 1 COMMENT '1:启用，0:禁用',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (habit_id) REFERENCES habit(id)
);

-- 习惯模板表
CREATE TABLE IF NOT EXISTS habit_template (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    category VARCHAR(50),
    default_schedule_type TINYINT DEFAULT 1,
    default_schedule_detail VARCHAR(255),
    default_reward_points SMALLINT DEFAULT 0,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TINYINT DEFAULT 1 COMMENT '1:启用，0:禁用'
);

-- 创建索引
CREATE INDEX idx_user_parent_id ON user(parent_id);
CREATE INDEX idx_points_record_user_id ON points_record(user_id);
CREATE INDEX idx_habit_creator_id ON habit(creator_id);
CREATE INDEX idx_habit_assignment_child_id ON habit_assignment(child_id, status);
CREATE INDEX idx_checkin_record_user_date ON checkin_record(user_id, checkin_date);
CREATE INDEX idx_exchange_record_user_id ON exchange_record(user_id);
CREATE INDEX idx_notification_user_read ON notification(user_id, is_read);
CREATE INDEX idx_habit_reminder_user_id ON habit_reminder(user_id);

-- 插入默认数据
-- 插入默认成就
INSERT INTO achievement (code, name, description, condition_type, condition_value, reward_points) VALUES
('first_checkin', '初次打卡', '完成第一次打卡', 'total_checkins', 1, 10),
('streak_7', '坚持一周', '连续打卡7天', 'streak_days', 7, 50),
('streak_30', '月度达人', '连续打卡30天', 'streak_days', 30, 200),
('checkin_100', '百日成就', '累计打卡100次', 'total_checkins', 100, 500),
('habit_10', '习惯养成', '创建10个习惯', 'habit_count', 10, 100),
('exchange_5', '兑换达人', '累计兑换5次', 'exchange_count', 5, 150);

-- 插入默认习惯模板
INSERT INTO habit_template (name, description, category, default_schedule_type, default_reward_points) VALUES
('每日刷牙', '每天早晚刷牙，保持口腔健康', 'health', 1, 5),
('洗脸洗手', '每天洗脸洗手，保持清洁', 'health', 1, 5),
('做作业', '每天完成家庭作业', 'study', 1, 10),
('阅读', '每天阅读30分钟', 'study', 1, 10),
('运动', '每天运动30分钟', 'sports', 1, 10),
('整理房间', '每天整理自己的房间', 'other', 1, 8),
('早睡早起', '每天按时睡觉起床', 'health', 1, 15),
('喝水', '每天喝足够的水', 'health', 1, 3);
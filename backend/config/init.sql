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
    parent_id BIGINT COMMENT '小孩账号必填，父母账号为 NULL',
    login_fail_count TINYINT DEFAULT 0,
    locked_until DATETIME,
    points_balance INT DEFAULT 0,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    require_photo TINYINT DEFAULT 0 COMMENT '1:需要拍照，0:不需要',
    allow_self_rate TINYINT DEFAULT 0 COMMENT '是否允许自我评分：0-否，1-是',
    checkin_prompt VARCHAR(500) COMMENT '打卡提示内容',
    creator_id BIGINT NOT NULL,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status TINYINT DEFAULT 1 COMMENT '1:启用，0:禁用',
    FOREIGN KEY (creator_id) REFERENCES user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 习惯分配表
CREATE TABLE IF NOT EXISTS habit_assignment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    habit_id INT NOT NULL,
    child_id INT NOT NULL,
    assign_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TINYINT DEFAULT 1 COMMENT '1:有效，0:已取消'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 打卡记录表
CREATE TABLE IF NOT EXISTS checkin_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    habit_id BIGINT NOT NULL,
    assignment_id BIGINT,
    checkin_date DATE NOT NULL,
    checkin_time DATETIME NOT NULL,
    checkin_type TINYINT DEFAULT 1 COMMENT '1:正常打卡，2:补卡',
    photo_url VARCHAR(1000) COMMENT '打卡照片文件路径',
    self_rate INT COMMENT '自我评分：1-10 分',
    remark VARCHAR(255),
    status TINYINT DEFAULT 1 COMMENT '1:成功，0:失败',
    is_rolled_back TINYINT DEFAULT 0 COMMENT '是否已回退：0-否，1-是',
    rollback_time DATETIME COMMENT '回退时间',
    rollback_reason VARCHAR(500) COMMENT '回退原因',
    points_rewarded SMALLINT DEFAULT 0,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (habit_id) REFERENCES habit(id),
    UNIQUE KEY unique_checkin (user_id, habit_id, checkin_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 奖励物品表
CREATE TABLE IF NOT EXISTS reward_item (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image VARCHAR(255),
    category VARCHAR(50),
    points_required INT NOT NULL,
    stock INT DEFAULT -1 COMMENT '-1 表示不限量',
    exchange_limit INT DEFAULT 0 COMMENT '0 表示不限',
    user_exchanged INT DEFAULT 0,
    start_time DATETIME,
    end_time DATETIME,
    creator_id BIGINT NOT NULL,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status TINYINT DEFAULT 1 COMMENT '1:上架，0:下架',
    FOREIGN KEY (creator_id) REFERENCES user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
('streak_7', '坚持一周', '连续打卡 7 天', 'streak_days', 7, 50),
('streak_30', '月度达人', '连续打卡 30 天', 'streak_days', 30, 200),
('checkin_100', '百日成就', '累计打卡 100 次', 'total_checkins', 100, 500),
('habit_10', '习惯养成', '创建 10 个习惯', 'habit_count', 10, 100),
('exchange_5', '兑换达人', '累计兑换 5 次', 'exchange_count', 5, 150);

-- 插入默认习惯模板
INSERT INTO habit_template (name, description, category, default_schedule_type, default_reward_points) VALUES
('每日刷牙', '每天早晚刷牙，保持口腔健康', 'health', 1, 5),
('洗脸洗手', '每天洗脸洗手，保持清洁', 'health', 1, 5),
('做作业', '每天完成家庭作业', 'study', 1, 10),
('阅读', '每天阅读 30 分钟', 'study', 1, 10),
('运动', '每天运动 30 分钟', 'sports', 1, 10),
('整理房间', '每天整理自己的房间', 'other', 1, 8),
('早睡早起', '每天按时睡觉起床', 'health', 1, 15),
('喝水', '每天喝足够的水', 'health', 1, 3);

-- 名言警句表
CREATE TABLE IF NOT EXISTS quote (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    meaning TEXT,
    author VARCHAR(100),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 名言警句默认数据（补充了每条名言的“meaning”释义）
INSERT INTO quote (content, meaning, author) VALUES
('“好好学习，天天向上。”', '强调持续学习的重要性，每天进步一点点就能取得长足发展。', '毛泽东'),
('“世上无难事，只要肯登攀。”', '只要有决心和毅力，任何困难都能克服。', '毛泽东'),
('“一万年太久，只争朝夕。”', '时间宝贵，要珍惜当下，抓紧每一天去奋斗。', '毛泽东'),
('“虚心使人进步，骄傲使人落后。”', '保持谦逊才能不断学习进步，自满只会停滞不前。', '毛泽东'),
('“星星之火，可以燎原。”', '微小的力量只要持续积累，也能发展成巨大影响。', '毛泽东'),
('“学习的敌人是自己的满足，要认真学习一点东西，必须从不自满开始。”', '学习需要保持空杯心态，自满会阻碍进步。', '毛泽东'),
('“饭可以一日不吃，觉可以一日不睡，书不可以一日不读。”', '读书学习比吃饭睡觉更重要，应养成每日阅读的习惯。', '毛泽东'),
('“读书是学习，使用也是学习，而且是更重要的学习。”', '学以致用，实践是检验和深化知识的关键环节。', '毛泽东'),
('“自信人生二百年，会当水击三千里。”', '保持自信，敢于挑战，人生才能拥有无限可能。', '毛泽东'),
('“不管风吹浪打，胜似闲庭信步。”', '面对困难要沉着冷静，保持从容心态。', '毛泽东'),
('“不到长城非好汉。”', '不达到目标决不罢休，体现坚韧不拔的精神。', '毛泽东'),
('“下定决心，不怕牺牲，排除万难，去争取胜利。”', '坚定信念、勇往直前才能最终取得成功。', '毛泽东'),
('“自己动手，丰衣足食。”', '强调自立自强，依靠自身努力创造美好生活。', '毛泽东'),
('“团结就是力量。”', '集体协作能产生巨大能量，团结一致才能战胜困难。', '毛泽东'),
('“知识就是力量。”', '掌握知识能够改变命运，增强个人与社会竞争力。', '培根'),
('“业精于勤，荒于嬉；行成于思，毁于随。”', '学业因勤奋而精，因玩乐而废；做事要独立思考，不可随波逐流。', '韩愈'),
('“学而不思则罔，思而不学则殆。”', '只学习不思考会迷惑，只思考不学习会陷入危险。', '孔子'),
('“温故而知新，可以为师矣。”', '复习旧知识能有新体会，就具备传授他人的能力。', '孔子'),
('“三人行，必有我师焉。”', '身边的人都有值得学习之处，要保持谦逊好学。', '孔子'),
('“宝剑锋从磨砺出，梅花香自苦寒来。”', '只有经历艰苦磨炼，才能取得卓越成就。', '警世贤文'),
('“书山有路勤为径，学海无涯苦作舟。”', '求学之路没有捷径，唯有勤奋与刻苦才能前行。', '韩愈'),
('“少壮不努力，老大徒伤悲。”', '年轻时不努力，年老只能后悔，强调把握青春。', '长歌行'),
('“一寸光阴一寸金，寸金难买寸光阴。”', '时间极其宝贵，失去就无法用金钱换回。', '增广贤文'),
('“成功是一种习惯，而不是一种结果。”', '成功源于日常良好习惯的积累，而非偶然。', '亚里士多德'),
('“成功不是最终的，失败不是致命的，重要的是你从失败中吸取的经验。”', '持续学习与改进比一时的成败更重要。', '温斯顿·丘吉尔'),
('“成功的秘诀在于保持耐心和持续的努力。”', '长期坚持与耐心是通向成功的核心要素。', '爱迪生'),
('“失败是成功之母。”', '失败提供经验教训，是迈向成功的必经之路。', '中国谚语'),
('“不积跬步，无以至千里；不积小流，无以成江海。”', '伟大成就源于每一小步、每一滴水的持续积累。', '荀子'),
('“天行健，君子以自强不息。”', '天地运行刚健有力，君子应不断自我强大。', '周易'),
('“路漫漫其修远兮，吾将上下而求索。”', '追求真理与理想的道路虽长，也要百折不挠探索。', '屈原'),
('“有志者，事竟成，破釜沉舟，百二秦关终属楚；苦心人，天不负，卧薪尝胆，三千越甲可吞吴。”', '坚定志向并付出艰辛，终能完成看似不可能的目标。', '蒲松龄'),
('“锲而舍之，朽木不折；锲而不舍，金石可镂。”', '持之以恒才能攻克最难的事情。', '荀子'),
('“千里之行，始于足下。”', '任何远大目标都要从当下第一步开始。', '老子'),
('“读万卷书，行万里路。”', '理论结合实践，广泛阅读与亲身体验同样重要。', '刘彝'),
('“学如逆水行舟，不进则退。”', '学习需要持续前进，稍有懈怠就会退步。', '增广贤文'),
('“勤能补拙是良训，一分辛苦一分才。”', '勤奋能够弥补天资不足，付出与收获成正比。', '华罗庚'),
('“天才就是百分之九十九的汗水加百分之一的灵感。”', '所谓天才主要靠后天努力，灵感只是点睛之笔。', '爱迪生'),
('“不经一番寒彻骨，怎得梅花扑鼻香。”', '不经历严酷考验，就无法获得美好成果。', '黄檗'),
('“只要功夫深，铁杵磨成针。”', '只要有恒心与毅力，再难的事也能完成。', '中国谚语'),
('“水滴石穿，绳锯木断。”', '持续微小的力量最终能完成看似不可能的任务。', '班固'),
('“千里之堤，溃于蚁穴。”', '小问题若不及时处理，可能引发大灾难，防微杜渐。', '韩非子'),
('“勿以恶小而为之，勿以善小而不为。”', '再小的坏事也不做，再小的好事也要做，积善成德。', '刘备'),
('“己所不欲，勿施于人。”', '自己不愿意的事不要强加给别人，体现同理心。', '孔子'),
('“三人行，必有我师焉。择其善者而从之，其不善者而改之。”', '向他人优点学习，以他人缺点自警并改正。', '孔子'),
('“知之者不如好之者，好之者不如乐之者。”', '学习最好以兴趣为驱动，乐在其中才能持久。', '孔子'),
('“学而时习之，不亦说乎？”', '经常复习实践所学，会获得愉悦与成就感。', '孔子'),
('“敏而好学，不耻下问。”', '聪敏又勤学，不以向地位低的人请教为耻。', '孔子'),
('“学而不厌，诲人不倦。”', '自己学习不满足，教导他人不疲倦，终身学习并分享。', '孔子'),
('“知之为知之，不知为不知，是知也。”', '诚实面对自己的知识边界，才是真正的智慧。', '孔子'),
('“人无远虑，必有近忧。”', '缺乏长远规划，眼前就会出现忧患。', '孔子'),
('“工欲善其事，必先利其器。”', '要做好工作，先准备好合适工具与方法。', '孔子'),
('“欲速则不达。”', '一味求快反而达不到目的，应循序渐进。', '孔子'),
('“逝者如斯夫，不舍昼夜。”', '时间流逝不停，要珍惜每分每秒。', '孔子'),
('“岁寒，然后知松柏之后凋也。”', '严峻环境才能检验出真正的坚韧品质。', '孔子'),
('“志不强者智不达。”', '志向不坚定，智慧也难以充分发挥。', '墨子'),
('“长风破浪会有时，直挂云帆济沧海。”', '相信总有一天会乘风破浪，实现远大抱负。', '李白'),
('“天生我材必有用，千金散尽还复来。”', '自信天生我才必有用，财富失去还能再创造。', '李白'),
('“会当凌绝顶，一览众山小。”', '勇于攀登最高峰，才能拥有最广阔的视野。', '杜甫'),
('“不畏浮云遮望眼，自缘身在最高层。”', '站位高远就能排除干扰，看清事物本质。', '王安石'),
('“山重水复疑无路，柳暗花明又一村。”', '困境中常含转机，坚持就能迎来新天地。', '陆游'),
('“纸上得来终觉浅，绝知此事要躬行。”', '书本知识需亲身实践才能深刻理解和掌握。', '陆游'),
('“春蚕到死丝方尽，蜡炬成灰泪始干。”', '无私奉献到底，燃烧自己照亮他人。', '李商隐'),
('“落红不是无情物，化作春泥更护花。”', '退居幕后也要继续贡献力量，滋养后来者。', '龚自珍'),
('“千磨万击还坚劲，任尔东西南北风。”', '历经无数考验仍坚韧不拔，从容面对挑战。', '郑燮'),
('“粉骨碎身浑不怕，要留清白在人间。”', '即使是最脆弱的人，也能在时间的长河中保持健康。', '李清照');

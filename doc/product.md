## 1. 产品概述
### 1.1 核心目标
打造一个以习惯养成为核心的儿童成长应用，通过游戏化机制培养孩子的良好习惯，助力孩子全面发展。

### 1.2 产品定位
一款面向家庭的习惯管理工具，通过父母引导、孩子参与的方式，建立正向激励机制，帮助孩子养成良好习惯。

## 2. 功能需求

### 2.1 用户管理模块
- **注册登录**：支持手机号/邮箱注册，密码登录
- **用户类型**：分为父母账号和小孩账号
- **账号状态**：正常、禁用
- **个人信息管理**：支持修改基本信息

### 2.2 习惯管理模块
- **习惯创建**：父母账号可创建、编辑、删除习惯
- **习惯设置**：包含习惯名称、描述、打卡时间、奖励积分
- **习惯分配**：父母可将习惯分配给指定小孩
- **习惯类型**：支持每日习惯、周期性习惯（如每周一、三、五）

### 2.3 打卡模块
- **打卡机制**：小孩账号每日登录后可进行打卡
- **打卡验证**：支持时间验证（在指定时间范围内可打卡）
- **打卡记录**：系统自动记录打卡时间、状态
- **打卡提醒**：支持设置打卡提醒通知

### 2.4 奖励系统模块
- **积分管理**：打卡成功后自动发放积分
- **积分记录**：详细记录积分获取和使用情况
- **奖励物品管理**：父母可添加、编辑、删除奖励物品，设置兑换积分
- **积分兑换**：小孩可使用积分兑换奖励物品
- **兑换记录**：系统记录兑换历史

### 2.5 数据展示模块
- **日历视图**：以日历形式展示习惯完成情况
- **统计报表**：展示习惯完成率、连续打卡天数等数据
- **个人中心**：展示积分余额、已完成习惯数量、已兑换物品等信息
- **数据导出**：支持导出习惯完成情况和积分记录

## 3. 角色与权限

### 3.1 小孩账号
- **核心权限**：登录、查看分配的习惯、进行打卡、查看积分、兑换奖励
- **操作限制**：不能创建/修改习惯，不能管理奖励物品
- **界面布局**：
  - 底部导航：日历打卡、个人中心
  - 日历打卡页：展示当日可打卡项目、已打卡项目、历史打卡记录
  - 个人中心页：展示积分余额、已完成习惯数量、已兑换物品记录

### 3.2 父母账号
- **核心权限**：管理小孩账号、创建/管理习惯、设置奖励机制、查看统计数据
- **操作功能**：
  - 小孩账号管理：添加、编辑、删除、禁用小孩账号
  - 习惯管理：创建、编辑、删除习惯，分配给小孩
  - 奖励管理：添加、编辑、删除奖励物品，设置兑换积分
  - 数据查看：查看所有小孩的习惯完成情况和积分记录

## 4. 数据模型

### 4.1 用户表（user）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 用户ID（主键，自增） |
| username | VARCHAR(50) | 账号（唯一索引） |
| password | VARCHAR(100) | 密码（bcrypt加密存储） |
| name | VARCHAR(50) | 姓名 |
| phone | VARCHAR(11) | 手机号（中国手机号11位，唯一索引） |
| email | VARCHAR(100) | 邮箱 |
| avatar | VARCHAR(255) | 头像URL |
| register_time | DATETIME | 注册时间 |
| last_login_time | DATETIME | 最后登录时间 |
| last_activity_time | DATETIME | 最后活动时间 |
| user_type | TINYINT | 用户类型（1:父母，2:小孩） |
| status | TINYINT | 状态（1:正常，0:禁用） |
| parent_id | BIGINT | 所属父母ID（小孩账号必填，父母账号为NULL） |
| login_fail_count | TINYINT | 登录失败次数（超过5次锁定） |
| locked_until | DATETIME | 账号锁定截止时间 |
| points_balance | INT | 当前积分余额（冗余字段，用于快速查询） |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |

### 4.2 积分记录表（points_record）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 记录ID（主键，自增） |
| user_id | BIGINT | 用户ID（索引） |
| points | SMALLINT | 积分变动值（正数为增加，负数为减少） |
| source | VARCHAR(50) | 积分来源（checkin:打卡获得、exchange:兑换物品、admin:管理员调整、bonus:奖励） |
| related_type | VARCHAR(20) | 关联类型（checkin_record、exchange_record等） |
| related_id | BIGINT | 关联ID（如打卡记录ID、兑换记录ID） |
| expire_time | DATETIME | 积分过期时间（NULL表示永不过期） |
| create_time | DATETIME | 记录时间 |

**设计说明**：
- 移除balance字段，积分余额通过汇总用户所有积分记录计算（或从用户表的points_balance冗余字段获取）
- 添加expire_time支持积分过期功能

### 4.3 习惯表（habit）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 习惯ID（主键，自增） |
| name | VARCHAR(100) | 习惯名称 |
| description | TEXT | 习惯描述 |
| icon | VARCHAR(255) | 习惯图标/图片URL |
| category | VARCHAR(50) | 习惯分类（health:卫生、study:学习、sports:运动、other:其他） |
| schedule_type | TINYINT | 打卡类型（1:每日，2:周期性） |
| schedule_detail | VARCHAR(255) | 打卡时间/周期详情（JSON格式，如{"weekdays":[1,3,5]}表示周一三五） |
| checkin_time_start | TIME | 每日打卡开始时间 |
| checkin_time_end | TIME | 每日打卡结束时间 |
| reward_points | SMALLINT | 完成奖励积分 |
| allow_makeup | TINYINT | 是否允许补卡（1:是，0:否） |
| makeup_days | TINYINT | 允许补卡天数（0表示不允许补卡） |
| creator_id | BIGINT | 创建者ID（父母账号） |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |
| status | TINYINT | 状态（1:启用，0:禁用） |

**schedule_detail格式说明**：
- schedule_type=1（每日）：schedule_detail可为空或设置为{"time":"08:00"}
- schedule_type=2（周期性）：schedule_detail为JSON，如{"weekdays":[1,3,5],"time":"08:00"}

### 4.4 习惯分配表（habit_assignment）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | INT | 分配ID（主键） |
| habit_id | INT | 习惯ID |
| child_id | INT | 小孩用户ID |
| assign_time | DATETIME | 分配时间 |
| status | TINYINT | 状态（1:有效，0:已取消） |

### 4.5 打卡记录表（checkin_record）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 记录ID（主键，自增） |
| user_id | BIGINT | 用户ID（索引） |
| habit_id | BIGINT | 习惯ID（索引） |
| assignment_id | BIGINT | 习惯分配ID |
| checkin_date | DATE | 打卡日期（索引，用于按日期查询） |
| checkin_time | DATETIME | 打卡时间 |
| checkin_type | TINYINT | 打卡类型（1:正常打卡，2:补卡） |
| remark | VARCHAR(255) | 打卡备注（用户可添加简短说明） |
| status | TINYINT | 状态（1:成功，0:失败） |
| points_rewarded | SMALLINT | 获得的积分 |
| create_time | DATETIME | 创建时间 |

**唯一索引**：(user_id, habit_id, checkin_date) - 确保同一用户同一天同一习惯只能打卡一次

### 4.6 奖励物品表（reward_item）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 物品ID（主键，自增） |
| name | VARCHAR(100) | 物品名称 |
| description | TEXT | 物品描述 |
| image | VARCHAR(255) | 物品图片URL |
| category | VARCHAR(50) | 物品分类（physical:实物、virtual:虚拟） |
| points_required | INT | 兑换所需积分 |
| stock | INT | 库存数量（-1表示不限量） |
| exchange_limit | INT | 每人限兑数量（0表示不限） |
| user_exchanged | INT | 已兑换数量（冗余字段） |
| start_time | DATETIME | 上架开始时间 |
| end_time | DATETIME | 上架结束时间 |
| creator_id | BIGINT | 创建者ID（父母账号） |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |
| status | TINYINT | 状态（1:上架，0:下架） |

### 4.7 兑换记录表（exchange_record）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 记录ID（主键，自增） |
| user_id | BIGINT | 用户ID（索引） |
| item_id | BIGINT | 奖励物品ID（索引） |
| points | INT | 兑换积分 |
| quantity | INT | 兑换数量 |
| exchange_time | DATETIME | 兑换时间 |
| delivery_info | JSON | 收货信息（实物填写收货地址、电话等） |
| status | TINYINT | 状态（1:已完成，2:处理中，3:已发货，4:已完成收货） |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |

### 4.8 连续打卡表（streak_record）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 记录ID（主键，自增） |
| user_id | BIGINT | 用户ID（唯一索引） |
| habit_id | BIGINT | 习惯ID |
| current_streak | INT | 当前连续打卡天数 |
| longest_streak | INT | 最长连续打卡天数 |
| last_checkin_date | DATE | 最后打卡日期 |
| streak_start_date | DATE | 当前连续打卡开始日期 |
| update_time | DATETIME | 更新时间 |

### 4.9 成就表（achievement）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 成就ID（主键，自增） |
| code | VARCHAR(50) | 成就代码（唯一） |
| name | VARCHAR(100) | 成就名称 |
| description | TEXT | 成就描述 |
| icon | VARCHAR(255) | 成就图标 |
| condition_type | VARCHAR(50) | 触发条件类型（streak_days:连续天数、total_checkins:累计打卡次数等） |
| condition_value | INT | 触发条件值 |
| reward_points | INT | 成就奖励积分 |
| create_time | DATETIME | 创建时间 |

### 4.10 用户成就表（user_achievement）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 记录ID（主键，自增） |
| user_id | BIGINT | 用户ID（索引） |
| achievement_id | BIGINT | 成就ID |
| unlocked | TINYINT | 是否解锁（1:已解锁，0:未解锁） |
| unlocked_time | DATETIME | 解锁时间 |
| create_time | DATETIME | 创建时间 |

### 4.11 通知消息表（notification）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 通知ID（主键，自增） |
| user_id | BIGINT | 接收用户ID（索引） |
| title | VARCHAR(100) | 通知标题 |
| content | TEXT | 通知内容 |
| type | VARCHAR(50) | 通知类型（checkin_reminder:打卡提醒、points_change:积分变动、achievement:成就解锁、exchange:兑换通知） |
| related_type | VARCHAR(20) | 关联类型 |
| related_id | BIGINT | 关联ID |
| is_read | TINYINT | 是否已读（1:已读，0:未读） |
| create_time | DATETIME | 创建时间 |

### 4.12 习惯提醒表（habit_reminder）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 提醒ID（主键，自增） |
| user_id | BIGINT | 用户ID（索引） |
| habit_id | BIGINT | 习惯ID |
| reminder_time | TIME | 提醒时间 |
| is_enabled | TINYINT | 是否启用（1:启用，0:禁用） |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |

### 4.13 习惯模板表（habit_template）
| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| id | BIGINT | 模板ID（主键，自增） |
| name | VARCHAR(100) | 模板名称 |
| description | TEXT | 模板描述 |
| icon | VARCHAR(255) | 模板图标 |
| category | VARCHAR(50) | 习惯分类 |
| default_schedule_type | TINYINT | 默认打卡类型 |
| default_schedule_detail | VARCHAR(255) | 默认打卡时间/周期 |
| default_reward_points | SMALLINT | 默认奖励积分 |
| create_time | DATETIME | 创建时间 |
| status | TINYINT | 状态（1:启用，0:禁用） |

## 5. 技术架构

### 5.1 前端
- **框架**：Vue.js
- **状态管理**：Vuex
- **路由**：Vue Router
- **UI组件**：Element UI
- **网络请求**：Axios
- **本地存储**：localStorage（缓存用户信息）
- **推送通知**：集成第三方推送服务

### 5.2 后端
- **语言**：Golang
- **框架**：Gin
- **数据库**：MySQL
- **缓存**：Redis（用于会话管理和热点数据）
- **认证**：JWT
- **API设计**：RESTful风格
- **日志**：结构化日志
- **结构**：
  - **控制器**：处理HTTP请求，调用服务层
  - **服务层**：业务逻辑处理，调用数据访问层
  - **数据访问层**：与数据库交互，执行CRUD操作

### 5.3 数据库
- **类型**：MySQL 8.0+
- **字符集**：UTF-8mb4
- **存储引擎**：InnoDB
- **索引策略**：针对常用查询建立合适的索引
- **备份策略**：定期全量备份，增量备份

## 6. 非功能需求

### 6.1 性能要求
- **响应时间**：页面加载时间 < 2秒，API响应时间 < 500ms
- **并发处理**：支持1000+并发用户
- **数据处理**：每日可处理10万+打卡记录

### 6.2 安全要求
- **数据加密**：密码使用bcrypt加密存储，敏感数据传输使用HTTPS
- **权限控制**：严格的角色权限管理，防止越权操作
- **SQL注入防护**：使用参数化查询，防止SQL注入
- **XSS防护**：对用户输入进行过滤和转义

### 6.3 可用性要求
- **系统可用性**：99.9%
- **容灾备份**：定期数据备份，支持快速恢复
- **错误处理**：完善的错误处理机制，提供友好的错误提示

### 6.4 可扩展性
- **模块化设计**：采用模块化架构，便于功能扩展
- **API设计**：预留扩展接口，支持未来功能集成
- **数据库设计**：考虑数据量增长，设计合理的表结构和索引

## 7. 运营需求

### 7.1 数据统计
- **习惯完成率**：统计每个习惯的完成情况
- **用户活跃度**：统计每日/周/月活跃用户数
- **积分发放与使用**：统计积分的发放和使用情况
- **兑换热度**：统计奖励物品的兑换热度

### 7.2 运营工具
- **批量操作**：支持批量导入习惯、批量添加奖励物品
- **消息推送**：支持向特定用户群发送通知
- **数据导出**：支持导出各类统计数据

## 8. 项目实施计划

### 8.1 开发阶段
1. **需求分析与设计**：1周
2. **数据库设计与搭建**：1周
3. **后端API开发**：2周
4. **前端界面开发**：2周
5. **系统集成与测试**：1周

### 8.2 上线计划
1. **内部测试**：1周
2. **beta测试**：2周
3. **正式上线**：持续迭代

## 9. 风险评估

### 9.1 潜在风险
- **用户粘性**：如何保持孩子的持续使用兴趣
- **数据安全**：确保用户数据，尤其是儿童信息的安全
- **技术挑战**：如何实现高效的打卡验证和数据统计

### 9.2 应对策略
- **游戏化设计**：增加成就系统、等级系统等游戏化元素
- **隐私保护**：严格遵守数据隐私法规，不收集不必要的信息
- **技术优化**：采用缓存、异步处理等技术提高系统性能

## 10. 附录

### 10.1 术语定义
- **习惯**：需要孩子定期完成的任务，如刷牙、洗脸等
- **打卡**：孩子完成习惯后在应用中进行的确认操作
- **积分**：完成习惯后获得的虚拟奖励，可用于兑换实物或虚拟奖励
- **奖励物品**：可通过积分兑换的物品，由父母设置

### 10.2 参考资料
- 儿童行为心理学相关研究
- 类似应用的成功案例分析
- 教育专家的建议和意见

---

## 11. 业务逻辑说明

### 11.1 打卡业务规则

#### 11.1.1 正常打卡
- 用户登录后，可对当日已分配的习惯进行打卡
- 打卡时间必须在习惯设定的 `checkin_time_start` 和 `checkin_time_end` 时间范围内
- 同一用户同一天同一习惯只能打卡一次（唯一索引约束）
- 打卡成功后，自动发放积分到用户账户，并记录积分变动

#### 11.1.2 补卡机制
- 若习惯设置允许补卡（`allow_makeup=1`），用户可在规定天数内补卡
- 补卡天数由 `makeup_days` 字段控制
- 补卡类型标记为 `checkin_type=2`
- 补卡是否奖励积分可由父母配置

#### 11.1.3 连续打卡计算
- 每日凌晨校验前一天的连续打卡情况
- 若前一天已打卡且当天未断签，`current_streak` +1
- 若前一天未打卡，重置 `current_streak` 为0
- 更新 `longest_streak` 为 max(current_streak, longest_streak)

### 11.2 积分业务规则

#### 11.2.1 积分获取
- 打卡成功：获得习惯设定的 `reward_points` 积分
- 成就解锁：获得成就设定的 `reward_points` 积分
- 父母奖励：父母手动调整用户积分

#### 11.2.2 积分使用
- 兑换奖励物品：扣除相应积分
- 积分不足时无法兑换

#### 11.2.3 积分过期
- 若积分设置了 `expire_time`，到期后积分失效
- 可通过定时任务每日凌晨检查并清理过期积分

### 11.3 兑换业务规则

#### 11.3.1 兑换限制
- 库存不足（stock <= 0 且 stock != -1）时无法兑换
- 每人限兑数量超过时无法兑换
- 积分不足时无法兑换

#### 11.3.2 兑换流程
1. 用户选择奖励物品并确认兑换
2. 验证库存、兑换限制、积分余额
3. 扣除积分，记录积分变动
4. 创建兑换记录，状态为"处理中"
5. 父母处理兑换请求（发货/发放虚拟奖励）
6. 更新兑换状态

### 11.4 成就系统规则

#### 11.4.1 成就类型
| 条件类型 | 说明 | 示例 |
|---------|------|------|
| streak_days | 连续打卡天数 | 连续7天打卡 |
| total_checkins | 累计打卡次数 | 累计打卡100次 |
| habit_count | 习惯数量 | 拥有10个习惯 |
| exchange_count | 兑换次数 | 累计兑换5次 |

#### 11.4.2 成就触发
- 每次打卡后检查相关成就
- 每次积分变动后检查相关成就
- 成就解锁后发送通知，发放奖励积分

### 11.5 账号安全规则

#### 11.5.1 登录失败处理
- 登录失败时，`login_fail_count` +1
- 登录失败次数超过5次，锁定账号15分钟
- 锁定期间 `locked_until` 设置为解锁时间
- 登录成功后重置 `login_fail_count`

#### 11.5.2 密码安全
- 密码必须使用 bcrypt 加密存储
- 密码长度不少于8位

### 11.6 数据索引策略

| 表名 | 索引字段 | 索引类型 | 用途 |
|-----|---------|---------|------|
| user | username | UNIQUE | 账号登录 |
| user | phone | UNIQUE | 手机号登录 |
| user | parent_id | INDEX | 查询小孩所属父母 |
| points_record | user_id, create_time | INDEX | 查询用户积分记录 |
| habit | creator_id | INDEX | 查询父母创建的习惯 |
| habit_assignment | child_id, status | INDEX | 查询小孩的有效习惯 |
| checkin_record | user_id, checkin_date | INDEX | 查询用户打卡记录 |
| checkin_record | (user_id, habit_id, checkin_date) | UNIQUE | 防止重复打卡 |
| exchange_record | user_id | INDEX | 查询用户兑换记录 |
| notification | user_id, is_read | INDEX | 查询未读通知 |

---

## 12. 数据库配置信息

**MySQL连接信息**：
- 主机：47.91.151.207
- 端口：3306
- 用户名：root
- 密码：cnldj1988
- 数据库名：babyhabit
- 字符集：utf8mb4
- 存储引擎：InnoDB

**Redis 连接信息**：
- 主机：47.91.151.207
- 端口：6379
- 用户名：
- 密码：cnldj1988
- 数据库索引：0
# Lottery Backend 项目学习指南

> 基于 Go-Zero 微服务框架的抽奖系统后端项目

## 📖 目录

1. [项目概述](#1-项目概述)
2. [技术架构](#2-技术架构)
3. [核心功能模块](#3-核心功能模块)
4. [环境搭建](#4-环境搭建)
5. [项目结构](#5-项目结构)
6. [业务流程](#6-业务流程)
7. [数据库设计](#7-数据库设计)
8. [API接口](#8-api接口)
9. [学习路线](#9-学习路线)

---

## 1. 项目概述

### 1.1 项目介绍

Lottery Backend 是一个功能完整的微服务后端系统，主要实现在线抽奖平台的核心功能。项目采用 Go-Zero 框架开发，包含用户系统、抽奖系统、签到系统、商城系统等多个业务模块。

### 1.2 核心特性

- ✅ **微服务架构**: API + RPC 双层设计，服务解耦
- ✅ **完整监控体系**: Jaeger链路追踪 + Prometheus监控 + Grafana可视化
- ✅ **日志收集**: ELK Stack (Elasticsearch + Logstash + Kibana)
- ✅ **消息队列**: Kafka + Asynq 实现异步任务和定时任务
- ✅ **缓存策略**: Redis 缓存 + Go-Zero内置缓存
- ✅ **数据持久化**: MySQL 8.0
- ✅ **第三方集成**: 微信小程序登录、OSS文件存储
- ✅ **容器化部署**: Docker + Docker Compose + Kubernetes

### 1.3 业务模块

| 模块 | 功能说明 | 端口 |
|------|---------|------|
| **lottery** | 抽奖核心功能（发起、参与、开奖、打卡任务） | API:1005, RPC:2005 |
| **usercenter** | 用户中心（注册、登录、认证、用户信息） | API:2006, RPC:2004 |
| **checkin** | 签到系统（每日签到、积分、任务） | API:2009, RPC:- |
| **shop** | 商城系统（商品、订单、心愿兑换） | API:-, RPC:- |
| **vote** | 投票系统（投票配置、投票记录） | API:-, RPC:- |
| **comment** | 评论系统（评论、点赞） | API:-, RPC:- |
| **notice** | 通知系统（消息订阅、推送） | API:-, RPC:- |
| **upload** | 文件上传（图片、视频上传到OSS） | API:-, RPC:- |
| **mqueue** | 消息队列（定时任务、异步任务） | Scheduler, Job |

---

## 2. 技术架构

### 2.1 整体架构图

```
┌────────────────────────────────────────────────┐
│              前端 (小程序/H5)                   │
└───────────────────┬────────────────────────────┘
                    │ HTTPS
                    ▼
┌────────────────────────────────────────────────┐
│         Nginx Gateway (端口: 8888)             │
│         - 反向代理                              │
│         - CORS处理                              │
│         - SSL终止                               │
└───────────────────┬────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
        ▼           ▼           ▼
    ┌─────┐    ┌─────┐    ┌─────┐
    │ API │    │ API │    │ API │
    │层   │    │层   │    │层   │
    └──┬──┘    └──┬──┘    └──┬──┘
       │ gRPC     │ gRPC     │ gRPC
       ▼          ▼          ▼
    ┌─────┐    ┌─────┐    ┌─────┐
    │ RPC │    │ RPC │    │ RPC │
    │层   │    │层   │    │层   │
    └──┬──┘    └──┬──┘    └──┬──┘
       │          │          │
       └──────────┼──────────┘
                  │
    ┌─────────────┼─────────────┐
    │             │             │
    ▼             ▼             ▼
┌────────┐  ┌────────┐  ┌────────┐
│ MySQL  │  │ Redis  │  │ Kafka  │
└────────┘  └────────┘  └────────┘
```

### 2.2 服务分层

**三层架构模式**:

```
┌─────────────────────────────────────┐
│          API Layer (HTTP)           │
│  - 接收HTTP请求                      │
│  - 参数验证                          │
│  - JWT鉴权                          │
│  - 调用RPC服务                       │
└─────────────────┬───────────────────┘
                  │ gRPC
┌─────────────────▼───────────────────┐
│          RPC Layer (gRPC)           │
│  - 处理业务逻辑                      │
│  - 数据库操作                        │
│  - 缓存处理                          │
│  - 事务管理                          │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│         Model Layer (数据层)         │
│  - 数据库CRUD                        │
│  - 缓存操作                          │
│  - 数据模型定义                      │
└─────────────────────────────────────┘
```

### 2.3 技术栈

**后端框架**:
- Go 1.23
- Go-Zero 1.5.3 (微服务框架)
- gRPC 1.64.0 (RPC通信)

**数据存储**:
- MySQL 8.0.28 (主数据库)
- Redis 6.2.5 (缓存)

**消息队列**:
- Kafka (消息队列)
- Asynq (异步任务队列)

**监控与追踪**:
- Jaeger (链路追踪)
- Prometheus (监控指标)
- Grafana (可视化)
- Elasticsearch + Kibana (日志系统)

**其他**:
- Docker & Docker Compose (容器化)
- Nginx (网关)
- JWT (认证)

---

## 3. 核心功能模块

### 3.1 抽奖模块 (Lottery)

**核心功能**:

1. **发起抽奖**
   - 创建抽奖活动
   - 配置奖品（支持多个等级）
   - 设置开奖方式（按时间/按人数/即抽即中）
   - 可选打卡任务（增加中奖概率）

2. **参与抽奖**
   - 用户参与抽奖
   - 检查参与资格
   - 完成打卡任务提升中奖率

3. **自动开奖**
   - 按时间定时开奖
   - 按人数自动开奖
   - 中奖概率算法（基础概率 + 打卡倍数）
   - 奖品等级分配

4. **打卡任务系统**
   - 体验小程序（15秒）
   - 浏览公众号文章（15秒）
   - 浏览图片（6秒）
   - 浏览视频号（15秒）
   - 增加中奖倍数（1-10倍）

**开奖类型**:

| 类型 | 说明 | 触发条件 |
|------|------|----------|
| 按时间开奖 | 到达指定时间自动开奖 | announce_time <= 当前时间 |
| 按人数开奖 | 参与人数达到设定值自动开奖 | 参与人数 >= join_number |
| 即抽即中 | 用户参与后立即开奖 | 用户点击参与 |

### 3.2 用户中心 (Usercenter)

**核心功能**:
- 微信小程序登录
- 用户信息管理
- JWT Token生成与验证
- 用户权限管理

### 3.3 签到模块 (Checkin)

**核心功能**:
- 每日签到
- 积分累积
- 任务系统
- 任务奖励领取

### 3.4 商城模块 (Shop)

**核心功能**:
- 商品管理
- 订单管理
- 心愿商品兑换
- 积分商城

### 3.5 消息队列 (Mqueue)

**核心功能**:
- 定时任务调度 (Scheduler)
- 异步任务处理 (Job)
- 抽奖定时开奖任务
- 签到任务处理

---

## 4. 环境搭建

### 4.1 前置要求

```bash
# 必需软件
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.23+
- Git
```

### 4.2 快速启动

**第一步：启动基础环境**

```bash
# 启动 MySQL、Redis、Kafka、Jaeger、Prometheus等
docker-compose -f docker-compose-env.yml up -d
```

**第二步：初始化数据库**

```bash
# 连接MySQL (密码: PXDN93VRKUm8TeE7)
mysql -h 127.0.0.1 -P 33069 -u root -p

# 导入SQL文件
source deploy/sql/lottery.sql
source deploy/sql/looklook_usercenter.sql
source deploy/sql/checkin.sql
# ... 其他SQL文件
```

**第三步：启动应用服务**

```bash
# 方式1: 使用docker-compose启动
docker-compose up -d

# 方式2: 本地开发模式 (支持热重载)
go install github.com/cortesi/modd/cmd/modd@latest
modd
```

### 4.3 服务端口

| 服务 | 端口 | 访问地址 |
|------|------|----------|
| Nginx Gateway | 8888 | http://localhost:8888 |
| Lottery API | 1005 | http://localhost:1005 |
| Usercenter API | 2006 | http://localhost:2006 |
| Jaeger UI | 16686 | http://localhost:16686 |
| Grafana | 3001 | http://localhost:3001 |
| Kibana | 5601 | http://localhost:5601 |
| Asynqmon | 8980 | http://localhost:8980 |
| Prometheus | 9091 | http://localhost:9091 |
| MySQL | 33069 | localhost:33069 |
| Redis | 36379 | localhost:36379 |

---

## 5. 项目结构

### 5.1 目录结构

```
lottery-backend/
├── app/                          # 应用服务目录
│   ├── lottery/                 # 抽奖服务
│   │   ├── cmd/
│   │   │   ├── api/            # HTTP API服务
│   │   │   │   ├── desc/       # API定义文件
│   │   │   │   ├── etc/        # 配置文件
│   │   │   │   ├── internal/   # 内部实现
│   │   │   │   └── lottery.go  # 主入口
│   │   │   └── rpc/            # gRPC服务
│   │   │       ├── etc/        # 配置文件
│   │   │       ├── internal/   # 内部实现
│   │   │       ├── pb/         # Protobuf定义
│   │   │       └── lottery.go  # 主入口
│   │   └── model/              # 数据模型
│   ├── usercenter/             # 用户中心服务
│   ├── checkin/                # 签到服务
│   ├── shop/                   # 商城服务
│   ├── vote/                   # 投票服务
│   ├── comment/                # 评论服务
│   ├── notice/                 # 通知服务
│   ├── upload/                 # 上传服务
│   └── mqueue/                 # 消息队列服务
│       ├── cmd/job/            # 任务消费者
│       └── cmd/scheduler/      # 任务调度器
├── common/                      # 公共代码
│   ├── constants/              # 常量定义
│   ├── ctxdata/               # 上下文数据
│   ├── globalkey/             # 全局键
│   ├── middleware/            # 中间件
│   ├── result/                # 响应封装
│   ├── tool/                  # 工具函数
│   ├── uniqueid/              # 唯一ID生成
│   └── xerr/                  # 错误定义
├── deploy/                     # 部署配置
│   ├── sql/                   # 数据库脚本
│   ├── nginx/                 # Nginx配置
│   ├── prometheus/            # Prometheus配置
│   ├── filebeat/              # Filebeat配置
│   └── goctl/                 # 代码生成模板
├── doc/                        # 文档
│   ├── chinese/               # 中文文档
│   └── english/               # 英文文档
├── data/                       # 数据目录（Docker挂载）
├── docker-compose.yml          # 应用服务编排
├── docker-compose-env.yml      # 基础环境编排
├── modd.conf                   # 热重载配置
├── go.mod                      # Go模块依赖
└── README.md                   # 项目说明
```

### 5.2 单个服务结构

```
app/[service]/cmd/api/
├── desc/                       # API定义
│   ├── [module]/
│   │   └── [module].api       # 模块API定义
│   └── main.api               # 主API定义
├── etc/
│   └── [service].yaml         # 配置文件
├── internal/
│   ├── config/
│   │   └── config.go          # 配置结构体
│   ├── handler/               # HTTP处理器
│   │   └── [module]/
│   │       └── [handler].go
│   ├── logic/                 # 业务逻辑
│   │   └── [module]/
│   │       └── [logic].go
│   ├── middleware/            # 中间件
│   ├── svc/
│   │   └── serviceContext.go # 服务上下文
│   └── types/
│       └── types.go           # 类型定义
└── [service].go               # 主入口
```

---

## 6. 业务流程

### 6.1 用户登录流程

```
1. 用户在小程序端获取code
   ↓
2. 前端调用 /usercenter/v1/user/wxMiniAuth
   ↓
3. 后端使用code换取openid (调用微信API)
   ↓
4. 查询或创建用户记录
   ↓
5. 生成JWT Token
   ↓
6. 返回token和用户信息给前端
   ↓
7. 前端后续请求携带token
```

### 6.2 发起抽奖流程

```
1. 用户填写抽奖信息（名称、奖品、开奖设置等）
   ↓
2. 前端调用 POST /lottery/v1/lottery/createLottery
   ↓
3. API层验证JWT token，获取userId
   ↓
4. API层验证参数（使用validator）
   ↓
5. 调用RPC服务 AddLottery
   ↓
6. RPC层开启数据库事务
   ├─ 插入抽奖记录 (lottery表)
   ├─ 插入奖品记录 (prize表，可能多条)
   └─ 插入打卡任务 (clock_task表，如果有)
   ↓
7. 提交事务
   ↓
8. 返回抽奖ID给前端
```

### 6.3 参与抽奖流程

```
1. 用户点击参与抽奖
   ↓
2. 前端调用 POST /lottery/v1/lottery/participation
   ↓
3. 验证用户登录状态
   ↓
4. 检查该用户是否已参与
   ↓
5. 检查抽奖是否已开奖
   ↓
6. 插入参与记录 (lottery_participation表)
   ↓
7. 如果是"按人数开奖"，检查是否达到人数
   ↓
8. 返回参与成功
```

### 6.4 自动开奖流程

```
定时任务 (每分钟执行)
   ↓
1. 查询需要开奖的抽奖
   ├─ 按时间: announce_time <= now() && is_announced = 0
   └─ 按人数: 参与人数 >= join_number && is_announced = 0
   ↓
2. 对每个抽奖执行开奖逻辑
   ├─ 获取所有参与者
   ├─ 获取所有奖品
   ├─ 计算每个用户的中奖概率
   │   └─ 基础概率1 + 打卡任务增加的倍数
   ├─ 按奖品等级分配中奖者
   │   └─ 使用加权随机算法选择
   ├─ 更新中奖记录
   └─ 更新抽奖状态为"已开奖"
   ↓
3. 发送中奖通知（可选）
```

### 6.5 打卡任务流程

```
1. 用户完成指定任务（浏览小程序/文章/视频等）
   ↓
2. 前端计时达到要求时长
   ↓
3. 调用 POST /lottery/v1/lottery/createClockTaskRecord
   ↓
4. 记录打卡完成记录 (clock_task_record表)
   ↓
5. 记录增加的中奖倍数
   ↓
6. 开奖时这个倍数会影响用户的中奖概率
```

---

## 7. 数据库设计

### 7.1 核心表结构

**lottery 表 (抽奖表)**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 主键 |
| user_id | int | 发起用户ID |
| name | varchar(50) | 抽奖名称 |
| thumb | varchar(255) | 封面图 |
| publish_time | datetime | 发布时间 |
| join_number | int | 自动开奖人数 |
| introduce | varchar(255) | 抽奖说明 |
| award_deadline | datetime | 领奖截止时间 |
| is_selected | tinyint | 是否精选 |
| announce_type | tinyint | 开奖类型：1时间 2人数 3即抽即中 |
| announce_time | datetime | 开奖时间 |
| is_announced | tinyint | 是否已开奖 |
| sponsor_id | int | 赞助商ID |
| is_clocked | tinyint | 是否开启打卡任务 |
| clock_task_id | int | 打卡任务ID |

**prize 表 (奖品表)**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 主键 |
| lottery_id | int | 抽奖ID |
| type | tinyint | 奖品类型：1奖品 2优惠券 3兑换码 等 |
| name | varchar(24) | 奖品名称 |
| level | int | 奖品等级（几等奖） |
| thumb | varchar(255) | 奖品图片 |
| count | int | 奖品数量 |
| grant_type | tinyint | 发放方式 |

**lottery_participation 表 (参与记录表)**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| lottery_id | int | 抽奖ID |
| user_id | int | 用户ID |
| is_won | tinyint | 是否中奖 |
| prize_id | bigint | 中奖奖品ID |

**clock_task 表 (打卡任务表)**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| lottery_id | bigint | 抽奖ID |
| type | tinyint | 任务类型：1小程序 2文章 3图片 4视频 |
| seconds | int | 需要完成的秒数 |
| chance_type | tinyint | 增加概率类型：1随机 2指定 |
| increase_multiple | int | 增加倍数 |

**clock_task_record 表 (打卡记录表)**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| lottery_id | bigint | 抽奖ID |
| user_id | bigint | 用户ID |
| clock_task_id | bigint | 任务ID |
| increase_multiple | bigint | 本次增加的倍数 |

### 7.2 表关系图

```
lottery (抽奖)
    ├─→ prize (奖品) [1:N]
    ├─→ lottery_participation (参与记录) [1:N]
    └─→ clock_task (打卡任务) [1:1]
            └─→ clock_task_record (打卡记录) [1:N]

user (用户)
    ├─→ lottery (发起的抽奖) [1:N]
    └─→ lottery_participation (参与的抽奖) [1:N]
```

---

## 8. API接口

### 8.1 抽奖相关接口

**公开接口（无需登录）**

```http
# 获取抽奖列表
POST /lottery/v1/lottery/lotteryList
Content-Type: application/json

{
  "page": 1,
  "size": 10,
  "isSelected": 0
}

# 查看抽奖详情
POST /lottery/v1/lottery/lotteryDetail
Content-Type: application/json

{
  "id": 1
}

# 查看参与者列表
POST /lottery/v1/lottery/participations
Content-Type: application/json

{
  "lotteryId": 1,
  "pageIndex": 1,
  "pageSize": 20
}

# 查看中奖名单
POST /lottery/v1/lottery/getLotteryWinnersList
Content-Type: application/json

{
  "lotteryId": 1
}
```

**需要登录的接口**

```http
# 发起抽奖
POST /lottery/v1/lottery/createLottery
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "新年抽奖",
  "thumb": "https://...",
  "introduce": "活动说明",
  "announceType": 1,
  "announceTime": 1704067200,
  "awardDeadline": 1704672000,
  "prizes": [
    {
      "type": 1,
      "name": "一等奖",
      "level": 1,
      "count": 1,
      "thumb": "https://...",
      "grantType": 1
    }
  ],
  "isClocked": 1,
  "clockTask": {
    "type": 1,
    "seconds": 15,
    "chanceType": 2,
    "increaseMultiple": 5
  }
}

# 参与抽奖
POST /lottery/v1/lottery/participation
Authorization: Bearer {token}
Content-Type: application/json

{
  "lotteryId": 1
}

# 完成打卡任务
POST /lottery/v1/lottery/createClockTaskRecord
Authorization: Bearer {token}
Content-Type: application/json

{
  "lotteryId": 1
}

# 获取我的中奖列表
POST /lottery/v1/lottery/getLotteryWinList
Authorization: Bearer {token}
Content-Type: application/json

{
  "lastId": 0,
  "size": 10
}
```

### 8.2 用户相关接口

```http
# 微信小程序登录
POST /usercenter/v1/user/wxMiniAuth
Content-Type: application/json

{
  "code": "wx_code",
  "iv": "...",
  "encryptedData": "..."
}

# 获取用户信息
GET /usercenter/v1/user/detail
Authorization: Bearer {token}
```

### 8.3 签到相关接口

```http
# 每日签到
POST /checkin/v1/checkin/signin
Authorization: Bearer {token}

# 获取签到记录
GET /checkin/v1/checkin/records
Authorization: Bearer {token}
```

---

## 9. 学习路线

### 9.1 第一阶段：环境搭建 (1-2天)

**目标**: 成功运行项目

- [ ] 安装Docker和Docker Compose
- [ ] 启动基础环境 (MySQL、Redis等)
- [ ] 导入数据库
- [ ] 启动应用服务
- [ ] 访问各个监控面板

**验证**:
- 访问 http://localhost:8888 正常
- Jaeger、Grafana等监控系统可访问
- 数据库连接正常

### 9.2 第二阶段：理解架构 (3-5天)

**目标**: 理解项目架构和技术选型

**学习内容**:
1. **Go-Zero框架**
   - 阅读官方文档
   - 理解API定义语法
   - 了解RPC调用机制

2. **微服务架构**
   - API层 vs RPC层的职责
   - 服务间如何通信
   - 为什么要分层

3. **目录结构**
   - app目录下各服务的作用
   - common目录的公共代码
   - model层的数据访问

**实践**:
- 追踪一次完整的HTTP请求
- 查看Jaeger中的调用链路
- 阅读一个简单接口的完整代码

### 9.3 第三阶段：核心业务 (1-2周)

**目标**: 掌握抽奖系统的核心业务逻辑

**学习路径**:

1. **抽奖模块**
   - 阅读数据库表结构
   - 理解发起抽奖流程
   - 理解参与抽奖流程
   - 理解开奖算法
   - 理解打卡任务系统

2. **用户模块**
   - 微信登录流程
   - JWT认证机制
   - 用户信息管理

3. **定时任务**
   - Asynq使用
   - 定时开奖任务
   - 任务调度机制

**实践**:
- 自己发起一个抽奖
- 参与抽奖并完成打卡
- 观察自动开奖过程
- 修改中奖概率算法

### 9.4 第四阶段：深入源码 (2-3周)

**目标**: 深入理解代码实现细节

**学习重点**:

1. **数据库操作**
   - Model层的实现
   - 缓存策略
   - 事务处理

2. **中间件**
   - CORS中间件
   - JWT认证中间件
   - 日志中间件

3. **错误处理**
   - 自定义错误码
   - 错误包装与传递
   - 统一错误响应

4. **监控与追踪**
   - Jaeger集成
   - Prometheus指标
   - 日志收集

**实践**:
- 添加新的API接口
- 实现新的业务功能
- 优化现有代码
- 编写单元测试

### 9.5 第五阶段：扩展开发 (按需)

**目标**: 独立开发新功能

**实践项目**:
1. 添加新的抽奖类型
2. 实现抽奖分享功能
3. 添加抽奖统计报表
4. 实现抽奖推荐算法
5. 优化开奖性能

---

## 10. 常见问题

### 10.1 环境问题

**Q: Docker容器启动失败？**
A: 检查端口占用，确保相关端口未被占用

**Q: 数据库连接失败？**
A: 确认MySQL容器正常运行，检查密码配置

**Q: Redis连接失败？**
A: 确认Redis容器正常运行，检查密码和端口

### 10.2 功能问题

**Q: 登录失败？**
A: 确认微信AppID和Secret配置正确

**Q: 文件上传失败？**
A: 检查OSS配置，确保AK/SK正确

**Q: 抽奖不开奖？**
A: 检查mqueue-scheduler服务是否运行

### 10.3 开发问题

**Q: 如何添加新接口？**
A: 
1. 在desc目录下的.api文件定义接口
2. 使用goctl生成代码
3. 实现logic层业务逻辑

**Q: 如何调试代码？**
A: 
1. 查看日志输出
2. 使用Jaeger查看调用链路
3. 使用IDE的debug功能

---

## 11. 参考资源

### 11.1 官方文档

- [Go-Zero官方文档](https://go-zero.dev/)
- [Go-Zero GitHub](https://github.com/zeromicro/go-zero)
- [gRPC官方文档](https://grpc.io/docs/)

### 11.2 项目文档

- [开发环境搭建](chinese/01-开发环境搭建.md)
- [Nginx网关](chinese/02-nginx网关.md)
- [鉴权服务](chinese/03-鉴权服务.md)
- [用户服务](chinese/04-用户服务.md)
- [链路追踪](chinese/12-链路追踪.md)
- [服务监控](chinese/13-服务监控.md)

### 11.3 技术博客

- Go-Zero微服务实战
- 微服务架构最佳实践
- 分布式系统设计

---

## 12. 总结

Lottery Backend 是一个功能完整、架构清晰的微服务项目，涵盖了：

✅ **完整的业务场景** - 从用户登录到抽奖全流程
✅ **现代化架构** - 微服务 + 容器化 + 云原生
✅ **工程化实践** - 监控、日志、追踪一应俱全
✅ **可扩展性** - 易于添加新功能和服务
✅ **生产级代码** - 错误处理、事务、缓存等完善

通过学习这个项目，你可以掌握：
- Go微服务开发
- 分布式系统设计
- 数据库设计与优化
- 消息队列应用
- 容器化部署
- 微服务监控

**祝你学习顺利！** 🚀

---

## 附录

### A. 环境变量配置

```yaml
# MySQL
MYSQL_ROOT_PASSWORD: PXDN93VRKUm8TeE7
MYSQL_PORT: 33069

# Redis
REDIS_PASSWORD: G62m50oigInC30sf
REDIS_PORT: 36379

# JWT Secret
JWT_ACCESS_SECRET: ae0536f9-6450-4606-8e13-5a19ed505da0
```

### B. 常用命令

```bash
# 查看所有容器状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service-name]

# 重启服务
docker-compose restart [service-name]

# 进入容器
docker exec -it [container-name] sh

# 查看数据库
mysql -h 127.0.0.1 -P 33069 -u root -p

# 查看Redis
redis-cli -h 127.0.0.1 -p 36379 -a G62m50oigInC30sf
```

### C. 项目贡献者

- 抽奖服务：王中阳
- 其他贡献者：见项目README

---

**文档版本**: v1.0  
**最后更新**: 2025-01-05  
**维护者**: lottery-backend team


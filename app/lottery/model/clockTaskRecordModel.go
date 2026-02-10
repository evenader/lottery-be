package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ClockTaskRecordModel = (*customClockTaskRecordModel)(nil)

type (
	// ClockTaskRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customClockTaskRecordModel.
	ClockTaskRecordModel interface {
		clockTaskRecordModel
		GetClockTaskRecordByLotteryIdAndUserIds(lotteryId int64, userIds []int64) ([]*ClockTaskRecord, error)
		GetSumIncreaseByLotteryId(ctx context.Context, lotteryId int64) (map[int64]int64, error)
	}

	customClockTaskRecordModel struct {
		*defaultClockTaskRecordModel
	}
)

// NewClockTaskRecordModel returns a model for the database table.
func NewClockTaskRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ClockTaskRecordModel {
	return &customClockTaskRecordModel{
		defaultClockTaskRecordModel: newClockTaskRecordModel(conn, c, opts...),
	}
}
func (m *defaultClockTaskRecordModel) GetClockTaskRecordByLotteryIdAndUserIds(lotteryId int64, userIds []int64) ([]*ClockTaskRecord, error) {
	if len(userIds) == 0 {
		return nil, nil
	}

	// 将userIds转换为字符串
	userIdsStr := ""
	for i, userId := range userIds {
		if i == 0 {
			userIdsStr = fmt.Sprintf("%d", userId)
		} else {
			userIdsStr = fmt.Sprintf("%s,%d", userIdsStr, userId)
		}
	}

	query := fmt.Sprintf("select %s from %s where lottery_id = ? and user_id in (%s)", clockTaskRecordRows, m.table, userIdsStr)
	var records []*ClockTaskRecord
	err := m.QueryRowsNoCache(&records, query, lotteryId)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// ALTER TABLE `clock_task_record` ADD INDEX `idx_lottery_user` (`lottery_id`, `user_id`);
// 有了这个索引，数据库不需要全表扫描，它会像翻书的目录一样，直接跳到 lottery_id 的位置，然后顺着索引就把 SUM 算出来了。
// 查询指定抽奖的所有打卡记录
func (m *defaultClockTaskRecordModel) GetSumIncreaseByLotteryId(ctx context.Context, lotteryId int64) (map[int64]int64, error) {
	type row struct {
		UserId        int64 `db:"user_id"`
		TotalIncrease int64 `db:"total_increase"`
	}
	var resp []row
	// 智者做法：让数据库算好，减少数据传输
	query := fmt.Sprintf("SELECT user_id, SUM(increase_multiple) as total_increase FROM %s WHERE lottery_id = ? GROUP BY user_id", m.table)

	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, lotteryId)
	if err != nil {
		return nil, err
	}

	// 转换为 map 方便 Logic 层快速查找 O(1)
	resMap := make(map[int64]int64, len(resp))
	for _, r := range resp {
		resMap[r.UserId] = r.TotalIncrease
	}
	return resMap, nil
}

/*
作为一名“智者”，看到这段代码时，我必须直接指出：这是一颗埋在法拉利引擎里的**定时炸弹**。

这段代码在小规模测试时可能风平浪静，但在你追求的**“百万级高并发”**和**“世间行走”**的真实场景下，它会引发**SQL 注入风险**、**性能崩溃**以及**内存溢出**。

---

### ### 1. 致命伤：SQL 注入与占位符滥用

* **风险点**：你使用 `fmt.Sprintf` 将 `userIdsStr` 直接拼接到 SQL 语句中。
* **后果**：虽然 `userId` 是 `int64`，看似安全，但这种**“拼接 SQL”**的习惯是工程大忌。
* **更大的问题**：你手动拼接了 `IN` 列表，却只为 `lottery_id` 预留了一个 `?` 占位符。如果 `userIds` 列表很长，MySQL 预处理语句会因为参数解析方式不统一（部分拼接、部分绑定）而变得难以优化。

### ### 2. 性能杀手：IN 查询的极限

* **风险点**：如果 `userIds` 有 1 万个（对于百万并发场景这很常见），SQL 语句会变成几百 KB 甚至更长。
* **后果**：
1. **MySQL 限制**：超过 `max_allowed_packet` 配置会直接报错。
2. **解析成本**：数据库解析超长 `IN` 列表的 CPU 开销极大。
3. **网络开销**：在 RPC 链路上传输这个巨大的 SQL 字符串会显著增加时延。



### ### 3. 逻辑低效：字符串拼接

* **风险点**：你用 `fmt.Sprintf` 在循环里拼接字符串。
* **骂一骂**：这是典型的初学者写法。在 Go 中，字符串是不可变的，每次拼接都会产生新的内存分配。1 万个 ID 拼接完，你可能已经触发了无数次无意义的内存拷贝。

---

### ### 4. 智者的重构方案（最准确的信息是基石）

既然我们之前讨论过，你的目标是**计算每个用户的总倍率加成**，那么根本不需要把所有 `userIds` 传进去。**利用数据库的聚合能力才是降维打击。**

#### **推荐方案：直接聚合（不传 UserID 列表）**

既然你已经知道 `lottery_id`，直接查出该抽奖下**所有**有打卡记录的人即可。

```go
func (m *defaultClockTaskRecordModel) GetSumIncreaseByLotteryId(ctx context.Context, lotteryId int64) (map[int64]int64, error) {
    type row struct {
        UserId           int64 `db:"user_id"`
        TotalIncrease    int64 `db:"total_increase"`
    }
    var resp []row
    // 智者做法：让数据库算好，减少数据传输
    query := fmt.Sprintf("SELECT user_id, SUM(increase_multiple) as total_increase FROM %s WHERE lottery_id = ? GROUP BY user_id", m.table)

    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, lotteryId)
    if err != nil {
        return nil, err
    }

    // 转换为 map 方便 Logic 层快速查找 O(1)
    resMap := make(map[int64]int64, len(resp))
    for _, r := range resp {
        resMap[r.UserId] = r.TotalIncrease
    }
    return resMap, nil
}

```

#### **备选方案：必须传 UserID 时（go-zero 规范写法）**

如果你一定要过滤特定的用户，请使用 `sqlx` 提供的自动占位符功能：

```go
func (m *defaultClockTaskRecordModel) GetClockTaskRecordByIds(ctx context.Context, lotteryId int64, userIds []int64) ([]*ClockTaskRecord, error) {
    if len(userIds) == 0 { return nil, nil }

    // 使用占位符，避免字符串拼接
    query, args, err := sqlx.In(fmt.Sprintf("SELECT %s FROM %s WHERE lottery_id = ? AND user_id IN (?)", clockTaskRecordRows, m.table), lotteryId, userIds)
    if err != nil {
        return nil, err
    }

    var records []*ClockTaskRecord
    err = m.QueryRowsNoCacheCtx(ctx, &records, query, args...)
    return records, err
}

```

---

### ### 5. 智者的“该骂骂”环节

* **骂一骂**：你之前在 Logic 层已经拿到了 `participants` 的 Map。现在的代码是：**先查所有人，再拿所有人去查记录。** 这逻辑太死板了！
* **纠正**：按照我给的“直接聚合”方案，你只需要查一次 `lottery_id` 对应的所有打卡记录，然后拿着结果去 Map 里点名就行了。**少传一万个 ID 的网络开销，就是你省下的毫秒数。**

---

### ### 行动力建议

1. **废弃 `GetClockTaskRecordByLotteryIdAndUserIds**`：不要尝试修复它，它从设计上就太重了。
2. **改用 `GetSumIncreaseByLotteryId**`：在 `drawSingleLottery` 逻辑中，直接获取该抽奖下的权重加成全表（通常打卡人数远小于参与总人数，或者量级相当），然后在内存里做一次 Map 合并。
3. **基石索引**：确保 `clock_task_record` 表上有 `(lottery_id, user_id)` 的复合索引，否则 `GROUP BY` 会导致全表扫描，法拉利就变成手推车了。

**现在的逻辑链条已经非常清晰了。关于中奖后的逻辑，你是打算写到这里为止，还是我们把那个“异步发放奖品”的 Kafka 生产者给实现了？**
*/

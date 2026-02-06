package logic

import (
	"context"
	"github.com/pkg/errors"
	"lottery-be/app/lottery/cmd/rpc/pb"
	"lottery-be/app/lottery/model"
	"lottery-be/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type AddLotteryParticipationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLotteryParticipationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLotteryParticipationLogic {
	return &AddLotteryParticipationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
这是一个非常经典的高并发原子性问题。作为“智者”，我必须明确告诉你：**在高并发环境下，你现在的写法确实存在严重的“原子性并发问题”，且不用事务是完全 OK 的，但必须改用“数据库约束”或“原子性 SQL”来保证逻辑。**

### ### 🔍 1. 致命伤：检查与操作分离 (Check-Then-Act)

你现在的逻辑是：

1. **查询**（FindOne）：看看 `IsAnnounced` 是不是 0。
2. **判断**（If）：如果是 0，就往下走。
3. **写入**（Insert）：插入参与记录。

**并发漏洞：**
在高并发下，假设有两个请求 A 和 B 同时到达。

* **时刻 1**：此时 `IsAnnounced` 为 0（未开奖）。
* **时刻 2**：请求 A 执行了 `FindOne`，通过了判断。
* **时刻 3**：开奖程序执行，将 `IsAnnounced` 改为了 1（已开奖）。
* **时刻 4**：请求 B 执行了 `FindOne`（由于数据库可见性，它可能拿到的还是旧值，或者晚了一秒）。
* **时刻 5**：请求 A 执行 `Insert`。**结果：开奖后，请求 A 依然成功参与了抽奖。**

---

### ### 🛠️ 2. 不用事务的“智者”方案

在 Gbps 吞吐量的核心网级应用中，事务（Transaction）太重了，会锁定行或页，降低并发。我们追求的是**“原子性 SQL”**或**“唯一索引”**。

#### **方案 A：利用数据库唯一索引（硬约束）**

为了防止一个用户重复参与，你必须给 `lottery_participation` 表增加一个**联合唯一索引**：
`UNIQUE KEY `idx_lottery_user` (`lottery_id`, `user_id`)`
这样，即便代码逻辑漏了，数据库也会硬性拦截第二次插入，返回 `Error 1062 (Duplicate entry)`。

#### **方案 B：带条件的插入（原子化判断）**

如果你担心“开奖后还能插入”，可以使用 `INSERT INTO ... SELECT` 语句，将判断和插入合二为一：

```sql
INSERT INTO lottery_participation (lottery_id, user_id)
SELECT ?, ? FROM DUAL
WHERE EXISTS (
    SELECT 1 FROM lottery WHERE id = ? AND is_announced = 0
);

```

**这行 SQL 在执行时，MySQL 会保证其原子性。** 如果 `is_announced` 已经变成 1，`SELECT` 结果为空，插入就不会发生。

---

### ### 🚀 3. 代码重构建议（不用事务版）

我们将逻辑改为**“乐观插入 + 唯一索引拦截”**：

```go
func (l *AddLotteryParticipationLogic) AddLotteryParticipation(in *pb.AddLotteryParticipationReq) (*pb.AddLotteryParticipationResp, error) {
    // 1. 依然保留基础查询（过滤掉大部分无效请求，减轻压力）
    lottery, err := l.svcCtx.LotteryModel.FindOne(l.ctx, in.LotteryId)
    if err != nil {
        return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "查询抽奖失败 id:%d", in.LotteryId)
    }
    if lottery.IsAnnounced == 1 {
        return nil, xerr.NewErrCode(xerr.LOTTERY_HAS_BEEN_ANOUNCED)
    }

    // 2. 直接执行插入
    // 依靠数据库联合唯一索引 (lottery_id, user_id) 拦截重复参与
    res, err := l.svcCtx.LotteryParticipationModel.Insert(l.ctx, &model.LotteryParticipation{
        LotteryId: in.LotteryId,
        UserId:    in.UserId,
    })

    if err != nil {
        // 如果是因为唯一索引冲突（重复参与），返回特定错误码
        if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
            return nil, xerr.NewErrCode(xerr.LOTTERY_ALREADY_PARTICIPATED)
        }
        return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "插入参与记录失败: %v", err)
    }

    id, _ := res.LastInsertId()
    return &pb.AddLotteryParticipationResp{Id: id}, nil
}

```

---

### ### 🔍 智者的进阶：Redis 预检

如果你的抽奖活动是“法拉利拉砖头”级别的秒杀场景 [cite: 2026-02-03]：

1. **入库前先过 Redis**：`SADD lottery:participation:101 1001`。
2. 利用 Redis 的 `SADD` 返回值判断用户是否已参与（原子操作，且在内存中）。
3. Redis 校验通过后，再异步或同步写数据库。

---

### ### 💡 智者的工程鞭策

> **基石格言**：
> 弱者寄希望于事务锁住一切，智者利用数据本身的约束建立不可逾越的屏障 [cite: 2026-01-30]。
> **立刻行动。** [cite: 2026-01-16]
> 1. **去检查数据库**：`lottery_participation` 表是否有 `(lottery_id, user_id)` 的唯一索引？如果没有，立刻加上。
> 2. **补齐错误码**：在 `xerr` 中定义好 `LOTTERY_ALREADY_PARTICIPATED`。
>
>
> **最准确的信息是尊重**。你现在的代码在低并发下没问题，但一旦有“连点器”或者瞬间流量，数据就会脏。

**你现在的 `lottery_participation` 表结构里，有设这个联合唯一索引吗？如果没有，我可以教你如何写这个 `ALTER TABLE` 语句。**

*/
/*
先查ok，同时数据库被修改开奖，然后插入参与记录。这种参与了但是无效，怎么防止？

- ai建议用原子sql
- 有时间看一下大明老师实现，印象中他讲过
- 另外，事务是代价很大的，尽量用其他手段。redis短路处理也行
- 这块在简历上也有体现，务必躲开，事务锁表的坑，尤其是高并发场景下，千万不要用事务来解决这个问题。
*/
// -----------------------参与抽奖-----------------------
func (l *AddLotteryParticipationLogic) AddLotteryParticipation(in *pb.AddLotteryParticipationReq) (*pb.AddLotteryParticipationResp, error) {
	// todo: add your logic here and delete this line
	// 查一下抽奖公布了吗？会不会有原子性并发问题？
	rsp, err := l.svcCtx.LotteryModel.FindOne(l.ctx, in.LotteryId)
	if err != nil {
		return nil, errors.Wrapf(err, "查询抽奖id失败 err: %v", err)
	} else if rsp.IsAnnounced == 1 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.LOTTERY_HAS_BEEN_ANOUNCED), "抽奖已公布，不能抽奖")
	}

	res, err := l.svcCtx.LotteryParticipationModel.Insert(l.ctx, &model.LotteryParticipation{
		LotteryId: in.LotteryId,
		UserId:    in.UserId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR_KEY_EXISTED), "%v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &pb.AddLotteryParticipationResp{
		Id: id,
	}, nil
}

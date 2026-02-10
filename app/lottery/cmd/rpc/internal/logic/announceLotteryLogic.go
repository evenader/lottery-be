package logic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"lottery-be/app/lottery/cmd/rpc/internal/fastpicker"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
	"lottery-be/app/lottery/cmd/rpc/pb"
	"lottery-be/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

/*
按照人数开奖的逻辑暂时不按照 抽奖服务增加抽奖时判断人数来处理。事件不够。
先按照原始轮询处理
*/
type AnnounceLotteryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnnounceLotteryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnnounceLotteryLogic {
	logic := &AnnounceLotteryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
	logic.initStrategies()

	return logic
}

func (l *AnnounceLotteryLogic) initStrategies() {
	l.svcCtx.RegisterStrategy(1, &timeLotteryStrategy{})
	l.svcCtx.RegisterStrategy(2, &memberCountLotteryStrategy{})
}

type Winner struct {
	LotteryId int64
	UserId    int64
	PrizeId   int64
}
type timeLotteryStrategy struct{}

func (t *timeLotteryStrategy) Run(ctx context.Context, svcCtx *svc.ServiceContext, drawSingleLottery func(lotteryId int64) error) error {
	// 扫描获取满足条件的lotteryId列表，针对每个id调用drawSingleLottery
	ids, err := svcCtx.LotteryModel.SearchTimeOutIds(ctx, time.Now(), 1)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("search timeout lottery fail"),
			"timeLotteryStrategy search timeout lottery err:%v", err)
	}
	if len(ids) == 0 {
		// 没有符合条件的开奖，退出，报成功？
		return nil
	}

	// 去参与者表查询
	// 获取所有 lottery_id 为 ids列表的参与记录的用户id列表
	// 策略不关心落库，落库交给logic做。同时不把多个id作为事务处理，防止单点失败整体失败
	for _, id := range ids {
		err := drawSingleLottery(id)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("draw single lottery fail"),
				"timeLotteryStrategy draw single lottery err:%v, lotteryId:%d", err, id)
		}
	}
	return nil
}

type memberCountLotteryStrategy struct{}

func (m *memberCountLotteryStrategy) Run(ctx context.Context, svcCtx *svc.ServiceContext, drawSingleLottery func(lotteryId int64) error) error {
	return nil
}

// 开奖一般逻辑，与开奖策略无关，所以抽离在logic层面，避免策略模式重复代码（如果几类模式代码不一样，就直接写在策略里）
/*
一般流程 针对单个id的开奖流程：
*/
/*
ai
单个算法 计算 + 事务落库:
ai
### 1. 为什么“跨抽奖”的大批量落库是危险的？
如果你想把所有 lottery_id 的中奖人堆在一起落库，你会面临以下三个“法拉利翻车”现场：

事务锁定范围过大： 一个事务如果同时修改 100 个抽奖的状态并插入 10,000 条记录，MySQL 会长时间占用大量的行锁。如果这时候有个用户想查询其中一个抽奖的状态，或者想在开奖瞬间加入，数据库会产生严重的阻塞。

“全全或全无”的代价： 如果第 99 个抽奖由于数据异常报错了，整个大事务会回滚。这意味着前 98 个本来已经算好的开奖结果全部作废。这在分布式系统中是不可接受的。

内存压力： 你需要把所有中奖信息（UserID, PrizeID）存在内存里。100 个抽奖 × 1000 人 = 10 万个结构体。这会造成频繁的 GC（垃圾回收），反而拖慢了你的算法性能。

### 2. 智者的权衡：两级批量策略
我们要做的不是“放弃批量”，而是**“在正确的维度批量”**。

第一级：单抽奖内批量（必须做） 在 drawSingleLottery 内部，针对这 1000 个中奖人，我们使用 INSERT INTO ... VALUES (),()...。这已经把 IO 损耗降低了 99%。

第二级：跨抽奖并发执行（高性能基石） 既然你是单实例，不要串行开奖。利用 Go 的并发特性，同时开 5-10 个抽奖。

### 3. 修正后的高性能代码结构
我们把“落库”留在 drawSingleLottery，但用并发去抵消事务的延迟。
*/
func (l *AnnounceLotteryLogic) drawSingleLottery(lotteryId int64) error {
	// 返回什么呢？
	prizes, err := l.svcCtx.PrizeModel.FindByLotteryIdWithOrder(l.ctx, lotteryId)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("search prize by lottery id fail"),
			"timeLotteryStrategy search prize by lottery id err:%v, lotteryId:%d", err, lotteryId)
	}

	participation, err := l.svcCtx.LotteryParticipationModel.GetParticipationUserIdsByLotteryId(l.ctx, lotteryId)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("search prize by lottery id fail"),
			"timeLotteryStrategy search prize by lottery id err:%v, lotteryId:%d", err, lotteryId)
	}
	if len(participation) == 0 {
		return l.markAsEmpty(lotteryId) // 没有参与者 标记处理
	}

	lottery, err := l.svcCtx.LotteryModel.FindOne(l.ctx, lotteryId)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("search lottery by id fail"), "err:")
	}

	weights := make(map[int64]int64, len(participation))
	for _, uid := range participation {
		weights[uid] = 1
	}
	if lottery.IsClocked == 1 { // 需要打卡 如果 uids 极多（比如 10w+），建议分批查询或使用连接查询
		clockMap, err := l.svcCtx.ClockTaskRecordModel.GetSumIncreaseByLotteryId(l.ctx, lotteryId)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("search clock task record by lottery id fail"),
				"timeLotteryStrategy search clock task record by lottery id err:%v, lotteryId:%d", err, lotteryId)
		}
		if clockMap != nil {
			for uid, increase := range clockMap {
				weights[uid] += increase
			}
		}
	}

	// 开始抽奖
	finalUids := make([]int64, 0, len(weights))
	finalRatios := make([]int64, 0, len(weights))
	for uid, ratio := range weights {
		finalUids = append(finalUids, uid)
		finalRatios = append(finalRatios, ratio)
	}
	// 抽完奖品有剩余？
	picker := fastpicker.NewFastPicker(finalRatios)
	defer picker.Recycle()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	winners := make([]*Winner, 0)

	for _, p := range prizes {
		for i := 0; i < int(p.Count); i++ {
			if len(finalUids) == 0 {
				break
			}
			idx := picker.PickAndRemove(rng)
			if idx == -1 {
				break
			}
			// 记录结果
			winners = append(winners, &Winner{
				UserId:    finalUids[idx],
				PrizeId:   p.Id,
				LotteryId: lotteryId,
			})
			// 物理移除业务 UID，保持与 Picker 索引同步
			finalUids = append(finalUids[:idx], finalUids[idx+1:]...)
		}
	}

	return l.saveWinner2DB(lotteryId, winners)
}

// 是不是所有的rpc接口都得支持超时重试那一堆东西
// 单实例不考虑分布式锁
// 单个实例不考虑扩展出处理中这种状态位（扫描时 置标志位 防止其他节点扫到）
func (l *AnnounceLotteryLogic) AnnounceLottery(in *pb.AnnounceLotteryReq) (*pb.AnnounceLotteryResp, error) {
	// 这里应该用枚举的IsValid
	if in.AnnounceType != 1 && in.AnnounceType != 2 {
		return nil, errors.Wrapf(xerr.NewErrMsg("invalid announce type"), "AnnounceLotteryLogic invalid announce type:%d", in.AnnounceType)
	}

	strategy := l.svcCtx.GetLotteryStrategy(in.AnnounceType)
	if err := strategy.Run(l.ctx, l.svcCtx, l.drawSingleLottery); err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("announce lottery fail"), "AnnounceLotteryLogic announce lottery err:%v, announceType:%d", err, in.AnnounceType)
	}

	return &pb.AnnounceLotteryResp{}, nil
}
func (l *AnnounceLotteryLogic) saveWinner2DB(lotteryId int64, winners []*Winner) error {
	// 1. 开启唯一的事务入口
	return l.svcCtx.LotteryModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {

		// A. 状态抢占（CAS）：这是整场演出的“发令枪”
		res, err := session.ExecCtx(ctx,
			"UPDATE lottery SET is_announced = 1 WHERE id = ? AND is_announced = 0",
			lotteryId)
		if err != nil {
			return err
		}
		aff, _ := res.RowsAffected()
		if aff == 0 {
			return nil // 已经被抢占，优雅退出
		}

		// B. 只有抢占成功，才执行落库逻辑，并将 session 传下去
		if len(winners) > 0 {
			return l.batchUpdateParticipation(ctx, session, lotteryId, winners)
		}
		return nil
	})
}

// 注意：这个函数不再自开事务，而是接收 session
func (l *AnnounceLotteryLogic) batchUpdateParticipation(ctx context.Context, session sqlx.Session, lotteryId int64, winners []*Winner) error {
	for _, w := range winners {
		// 1. 扣减奖品库存（新增：防止超发）
		resPrize, err := session.ExecCtx(ctx,
			"UPDATE prize SET count = count - 1 WHERE id = ? AND count > 0",
			w.PrizeId)
		if err != nil {
			return err
		}
		pAff, _ := resPrize.RowsAffected()
		if pAff == 0 {
			return errors.New("prize stock insufficient")
		}

		// 2. 标记中奖（带 RowsAffected 检查）
		updateSql := "UPDATE lottery_participation SET is_winner = 1, prize_id = ? WHERE lottery_id = ? AND user_id = ? AND is_winner = 0"
		resPart, err := session.ExecCtx(ctx, updateSql, w.PrizeId, lotteryId, w.UserId)
		if err != nil {
			return err
		}

		aff, _ := resPart.RowsAffected()
		if aff != 1 {
			return fmt.Errorf("unexpected affected rows for user %d", w.UserId)
		}
	}
	return nil
}
func (l *AnnounceLotteryLogic) markAsEmpty(lotteryId int64) error {
	// 即使没人参加，也要通过事务或带条件的 UPDATE 抢占状态
	// 防止在标记为空的同时，又有极端的并发写入导致逻辑冲突
	res, err := l.svcCtx.LotteryModel.UpdateStatusToAnnounced(l.ctx, lotteryId)
	if err != nil {
		return errors.Wrapf(err, "markAsEmpty fail, lotteryId: %d", lotteryId)
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		// 已经被抢占了，或者该抽奖本就不存在，逻辑上视为成功或跳过
		l.Logger.Infof("lottery %d already announced or not exists when marking empty", lotteryId)
		return nil
	}

	l.Logger.Infof("lottery %d announced with zero participants", lotteryId)
	return nil
}

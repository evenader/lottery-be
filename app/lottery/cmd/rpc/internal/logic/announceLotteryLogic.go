package logic

import (
	"context"
	"time"

	"lottery-be/app/lottery/cmd/rpc/internal/svc"
	"lottery-be/app/lottery/cmd/rpc/pb"
	"lottery-be/app/lottery/model"
	"lottery-be/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

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

type timeLotteryStrategy struct{}

func (t *timeLotteryStrategy) Run(ctx context.Context, svcCtx *svc.ServiceContext, drawSingleLottery func(lotteryId int64) ([]int64, error)) error {
	// 扫描获取满足条件的lotteryId列表，针对每个id调用drawSingleLottery
	ids, err := svcCtx.LotteryModel.SearchTimeOutIds(ctx, time.Now(), 1)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("search timeout lottery fail"),
			"timeLotteryStrategy search timeout lottery err:%v", err)
	}

	// 去参与者表查询
	// 获取所有 lottery_id 为 ids列表的参与记录的用户id列表
	var AllWinners [][]int64
	for _, lotId := range ids {
		winner, err := drawSingleLottery(lotId)
		// 这里需要copy吗？
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("draw single lottery fail"),
				"timeLotteryStrategy draw single lottery err:%v, lotteryId:%d", err, lotId)
		}
		AllWinners = append(AllWinners, winner)
	}

	var prizes [][]*model.Prize
	for _, lotId := range ids {
		prize, err := svcCtx.PrizeModel.FindByLotteryId(ctx, lotId)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("search prize by lottery id fail"),
				"timeLotteryStrategy search prize by lottery id err:%v, lotteryId:%d", err, lotId)
		}
		prizes = append(prizes, prize)
	}

	var participations [][]int64
	participations = make([][]int64, 0)
	for _, lotId := range ids {
		participation, err := svcCtx.LotteryParticipationModel.GetParticipationUserIdsByLotteryId(ctx, lotId)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("search prize by lottery id fail"),
				"timeLotteryStrategy search prize by lottery id err:%v, lotteryId:%d", err, lotId)
		}
		participations = append(participations, participation)
	}

}

type memberCountLotteryStrategy struct{}

func (m *memberCountLotteryStrategy) Run(ctx context.Context, svcCtx *svc.ServiceContext, drawSingleLottery func(lotteryId int64) ([]int64, error)) error {

}

// 开奖一般逻辑，与开奖策略无关，所以抽离在logic层面，避免策略模式重复代码（如果几类模式代码不一样，就直接写在策略里）
/*
一般流程 针对单个id的开奖流程：
*/
func (l *AnnounceLotteryLogic) drawSingleLottery(lotteryId int64) (int64, error) {
	// 返回什么呢？

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

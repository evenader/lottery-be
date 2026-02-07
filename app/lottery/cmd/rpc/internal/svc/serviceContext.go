package svc

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"lottery-be/app/lottery/cmd/rpc/internal/config"
	"lottery-be/app/lottery/model"
	"lottery-be/app/usercenter/cmd/rpc/usercenter"
)

type ServiceContext struct {
	Config               config.Config
	RedisClient          *redis.Redis
	LotteryModel         model.LotteryModel
	PrizeModel           model.PrizeModel
	ClockTaskModel       model.ClockTaskModel
	ClockTaskRecordModel model.ClockTaskRecordModel
	UserCenterRpc        usercenter.Usercenter
	//NoticeRpc                 notice.Notice
	LotteryParticipationModel model.LotteryParticipationModel
	AnnounceStrategies        map[int64]LotteryStrategy
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		LotteryModel:  model.NewLotteryModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		PrizeModel:    model.NewPrizeModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UserCenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpcConf)),
		//NoticeRpc:                 notice.NewNotice(zrpc.MustNewClient(c.NoticeRpcConf)),
		LotteryParticipationModel: model.NewLotteryParticipationModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		ClockTaskModel:            model.NewClockTaskModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		ClockTaskRecordModel:      model.NewClockTaskRecordModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		//RedisClient:               redis.New(c.RedisConf),
	}
}

type LotteryStrategy interface {
	// 策略本身不持有对象，策略是无状态的，每次调用都传入上下文对象
	// 策略模式自身仅仅负责策略触发，至于获取开奖概率，生成获奖用户，持久化数据库属于通用业务流程交给logic处理，所以这部分能力由logic注入
	Run(ctx context.Context, serviceContext *ServiceContext, drawSingleLottery func(lotteryId int64) ([]int64, error)) error
}

func (s *ServiceContext) RegisterStrategy(announceType int64, strategy LotteryStrategy) {
	if s.AnnounceStrategies == nil {
		s.AnnounceStrategies = make(map[int64]LotteryStrategy)
	}
	s.AnnounceStrategies[announceType] = strategy
}

func (s *ServiceContext) GetLotteryStrategy(announceType int64) LotteryStrategy {
	return s.AnnounceStrategies[announceType]
}

package logic

import (
	"context"
	"lottery-be/app/mqueue/cmd/scheduler/internal/svc"
)

type MqueueScheduler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronScheduler(ctx context.Context, svcCtx *svc.ServiceContext) *MqueueScheduler {
	return &MqueueScheduler{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MqueueScheduler) Register() {
	// looklook
	// l.settleRecordScheduler()

	// lottery
	l.LotteryDrawScheduler()
	l.WishCheckinScheduler()
}

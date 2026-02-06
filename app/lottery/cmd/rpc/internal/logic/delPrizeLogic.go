package logic

import (
	"context"
	"lottery-be/app/lottery/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type DelPrizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelPrizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelPrizeLogic {
	return &DelPrizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelPrizeLogic) DelPrize(in *pb.DelPrizeReq) (*pb.DelPrizeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DelPrizeResp{}, nil
}

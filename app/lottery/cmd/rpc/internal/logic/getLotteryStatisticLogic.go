package logic

import (
	"context"
	"lottery-be/app/lottery/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type GetLotteryStatisticLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLotteryStatisticLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLotteryStatisticLogic {
	return &GetLotteryStatisticLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLotteryStatisticLogic) GetLotteryStatistic(in *pb.GetLotteryStatisticReq) (*pb.GetLotteryStatisticResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetLotteryStatisticResp{}, nil
}

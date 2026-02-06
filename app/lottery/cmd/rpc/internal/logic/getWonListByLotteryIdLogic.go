package logic

import (
	"context"
	"lottery-be/app/lottery/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type GetWonListByLotteryIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWonListByLotteryIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWonListByLotteryIdLogic {
	return &GetWonListByLotteryIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWonListByLotteryIdLogic) GetWonListByLotteryId(in *pb.GetWonListByLotteryIdReq) (*pb.GetWonListByLotteryIdResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetWonListByLotteryIdResp{}, nil
}

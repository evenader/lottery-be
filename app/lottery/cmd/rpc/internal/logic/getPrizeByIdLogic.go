package logic

import (
	"context"
	"lottery-be/app/lottery/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type GetPrizeByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPrizeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrizeByIdLogic {
	return &GetPrizeByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPrizeByIdLogic) GetPrizeById(in *pb.GetPrizeByIdReq) (*pb.GetPrizeByIdResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetPrizeByIdResp{}, nil
}

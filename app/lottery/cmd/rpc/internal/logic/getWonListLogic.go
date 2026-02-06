package logic

import (
	"context"
	"lottery-be/app/lottery/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type GetWonListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWonListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWonListLogic {
	return &GetWonListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWonListLogic) GetWonList(in *pb.GetWonListReq) (*pb.GetWonListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetWonListResp{}, nil
}

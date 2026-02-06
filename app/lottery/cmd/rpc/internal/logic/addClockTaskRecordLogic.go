package logic

import (
	"context"
	"lottery-be/app/lottery/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type AddClockTaskRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddClockTaskRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddClockTaskRecordLogic {
	return &AddClockTaskRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------完成打卡任务-----------------------
func (l *AddClockTaskRecordLogic) AddClockTaskRecord(in *pb.AddClockTaskRecordReq) (*pb.AddClockTaskRecordResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddClockTaskRecordResp{}, nil
}

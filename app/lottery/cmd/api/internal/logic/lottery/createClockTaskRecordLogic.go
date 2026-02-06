// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"context"

	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateClockTaskRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 完成打卡任务
func NewCreateClockTaskRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateClockTaskRecordLogic {
	return &CreateClockTaskRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateClockTaskRecordLogic) CreateClockTaskRecord(req *types.CreateClockTaskRecordReq) (resp *types.CreateClockTaskRecordResp, err error) {
	// todo: add your logic here and delete this line

	return
}

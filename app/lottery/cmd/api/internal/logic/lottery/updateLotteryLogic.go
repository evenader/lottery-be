// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"context"

	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLotteryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布抽奖
func NewUpdateLotteryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLotteryLogic {
	return &UpdateLotteryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 传入一个id 数据库侧仅仅更新发布时间 此api 忽略 发布抽奖
func (l *UpdateLotteryLogic) UpdateLottery(req *types.UpdateLotteryReq) (resp *types.UpdateLotteryResp, err error) {
	// todo: add your logic here and delete this line

	return
}

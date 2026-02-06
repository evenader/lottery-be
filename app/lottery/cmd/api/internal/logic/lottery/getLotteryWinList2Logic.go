// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"context"

	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLotteryWinList2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前抽奖中奖者名单
func NewGetLotteryWinList2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLotteryWinList2Logic {
	return &GetLotteryWinList2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLotteryWinList2Logic) GetLotteryWinList2(req *types.GetLotteryWinList2Req) (resp *types.GetLotteryWinList2Resp, err error) {
	// todo: add your logic here and delete this line

	return
}

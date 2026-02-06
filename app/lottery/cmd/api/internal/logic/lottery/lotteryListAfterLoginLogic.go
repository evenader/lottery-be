// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"context"

	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LotteryListAfterLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录后获取抽奖列表
func NewLotteryListAfterLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LotteryListAfterLoginLogic {
	return &LotteryListAfterLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LotteryListAfterLoginLogic) LotteryListAfterLogin(req *types.LotteryListReq) (resp *types.LotteryListResp, err error) {
	// todo: add your logic here and delete this line

	return
}

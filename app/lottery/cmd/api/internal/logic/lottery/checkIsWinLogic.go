// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"context"

	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckIsWinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 判断当前用户当前抽奖是否中奖
func NewCheckIsWinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIsWinLogic {
	return &CheckIsWinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckIsWinLogic) CheckIsWin(req *types.CheckIsWinReq) (resp *types.CheckIsWinResp, err error) {
	// todo: add your logic here and delete this line

	return
}

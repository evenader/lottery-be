// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"lottery-be/app/usercenter/cmd/api/internal/svc"
	"lottery-be/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置user为admin
func NewSetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAdminLogic {
	return &SetAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetAdminLogic) SetAdmin(req *types.SetAdminReq) (resp *types.SetAdminResp, err error) {
	// todo: add your logic here and delete this line

	return
}

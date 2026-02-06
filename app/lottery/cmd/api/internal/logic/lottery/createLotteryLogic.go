// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"context"
	"github.com/pkg/errors"
	"lottery-be/app/lottery/cmd/api/internal/converter"
	"lottery-be/app/lottery/cmd/rpc/lottery"
	"lottery-be/common/ctxdata"
	"lottery-be/common/xerr"

	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLotteryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发起抽奖
func NewCreateLotteryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLotteryLogic {
	return &CreateLotteryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 用户创建抽奖
func (l *CreateLotteryLogic) CreateLottery(req *types.CreateLotteryReq) (resp *types.CreateLotteryResp, err error) {
	if req == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "CreateLotteryReq is nil")
	}
	addLotteryReq := converter.BuildAddLotteryReq(req, ctxdata.GetUidFromCtx(l.ctx))
	rsp, err := l.svcCtx.LotteryRpc.AddLottery(l.ctx, &lottery.AddLotteryReq{})
	if err != nil {
		// 日志不打印信元
		return nil, errors.Wrapf(err, "rpc AddLottery Failed addLotteryReq on user: %+v", addLotteryReq.UserId)
	}

	return &types.CreateLotteryResp{
		Id: rsp.Id,
	}, nil
}

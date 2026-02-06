// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"context"
	"github.com/pkg/errors"
	"lottery-be/app/lottery/cmd/rpc/lottery"
	"lottery-be/common/ctxdata"

	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLotteryParticipationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 参与抽奖
func NewAddLotteryParticipationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLotteryParticipationLogic {
	return &AddLotteryParticipationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 为什么没有信息 因为参与抽奖的接口 只需要传入抽奖id 其他信息都可以从上下文中获取？姑且放在这里
func (l *AddLotteryParticipationLogic) AddLotteryParticipation(req *types.AddLotteryParticipationReq) (resp *types.AddLotteryParticipationResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUidFromCtx(l.ctx)
	addLottertPartReq := &lottery.AddLotteryParticipationReq{
		LotteryId: req.LotteryId,
		UserId:    uid,
	}
	_, err = l.svcCtx.LotteryRpc.AddLotteryParticipation(l.ctx, addLottertPartReq)
	if err != nil {
		return nil, errors.Wrapf(err, "rpc AddLotteryParticipation Failed addLottertPartReq on user: %+v", addLottertPartReq.UserId)
	}

	return
}

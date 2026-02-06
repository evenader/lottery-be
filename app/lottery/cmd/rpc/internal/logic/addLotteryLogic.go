package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"lottery-be/app/lottery/cmd/rpc/pb"
	"lottery-be/app/lottery/model"
	"lottery-be/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"lottery-be/app/lottery/cmd/rpc/internal/svc"
)

type AddLotteryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLotteryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLotteryLogic {
	return &AddLotteryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------抽奖表-----------------------
func (l *AddLotteryLogic) AddLottery(in *pb.AddLotteryReq) (*pb.AddLotteryResp, error) {
	var lotteryId int64
	err := l.svcCtx.LotteryModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 插入一次抽奖、若干奖品、一次定时任务
		lottery := &model.Lottery{
			UserId:        in.UserId,
			Name:          in.Name,
			AwardDeadline: time.Unix(in.AwardDeadline, 0),
			Introduce:     in.Introduce,
			JoinNumber:    in.JoinNumber,
			AnnounceType:  in.AnnounceType,
			AnnounceTime:  time.Unix(in.AnnounceTime, 0),
			Thumb:         in.Thumb,
			IsSelected:    0,
			IsAnnounced:   0,
			SponsorId:     in.SponsorId,
			IsClocked:     in.IsClocked,
		}
		if in.PublishType == 1 {
			lottery.PublishTime.Time = time.Now()
			lottery.PublishTime.Valid = true
		}

		inRes, err := l.svcCtx.LotteryModel.Insert(ctx, lottery, model.WithSession(session))
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERTLOTTERY_ERROR), "Lottery Database Exception lottery : %+v , err: %v", lottery, err)
		}
		lotteryId, _ = inRes.LastInsertId()
		// 插入奖品
		for _, prize := range in.Prizes {
			modelPrize := &model.Prize{}
			copier.Copy(modelPrize, prize)
			modelPrize.LotteryId = lotteryId
			_, err := l.svcCtx.PrizeModel.Insert(l.ctx, modelPrize, model.WithSession(session))
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERTPRIZE_ERROR), "Prize Database Exception prize : %+v , err: %v", modelPrize, err)
			}
		}

		// 插入定时任务
		if in.ClockTask != nil {
			clockTask := &model.ClockTask{
				LotteryId:        lotteryId,
				Seconds:          in.ClockTask.Seconds,
				Type:             in.ClockTask.Type,
				AppletType:       in.ClockTask.AppletType,
				PageLink:         in.ClockTask.PageLink,
				AppId:            in.ClockTask.AppId,
				PagePath:         in.ClockTask.PagePath,
				Image:            in.ClockTask.Image,
				VideoAccountId:   in.ClockTask.VideoAccountId,
				VideoId:          in.ClockTask.VideoId,
				ArticleLink:      in.ClockTask.ArticleLink,
				Copywriting:      in.ClockTask.Copywriting,
				ChanceType:       in.ClockTask.ChanceType,
				IncreaseMultiple: in.ClockTask.IncreaseMultiple,
			}

			taskRes, err := l.svcCtx.ClockTaskModel.Insert(l.ctx, clockTask, model.WithSession(session))
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERTCLOCKTASK_ERROR), "ClockTask Database Exception clockTask : %+v , err: %v", clockTask, err)
			}
			clockTaskId, _ := taskRes.LastInsertId()
			// 更新抽奖表的 clockTaskId 字段
			if err := l.svcCtx.LotteryModel.UpdateClockTaskIdOnLottery(l.ctx, lotteryId, clockTaskId, model.WithSession(session)); err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Lottery Database Exception lottery : %+v , err: %v", lottery, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddLotteryResp{Id: lotteryId}, nil
}

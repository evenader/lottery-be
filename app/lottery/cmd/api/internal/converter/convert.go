// Copyright (c) 2026 evenader. All rights reserved.

package converter

import (
	"github.com/jinzhu/copier"
	"lottery-be/app/lottery/cmd/api/internal/types"
	"lottery-be/app/lottery/cmd/rpc/lottery"
	"lottery-be/app/lottery/cmd/rpc/pb"
)

func BuildAddLotteryReq(req *types.CreateLotteryReq, userId int64) *lottery.AddLotteryReq {
	if req == nil {
		return nil
	}
	// copier 有反射开销
	return &lottery.AddLotteryReq{
		UserId:        userId,
		Name:          req.Name,
		Thumb:         req.Thumb,
		AnnounceType:  req.AnnounceType,
		AnnounceTime:  req.AnnounceTime,
		JoinNumber:    req.JoinNumber,
		Introduce:     req.Introduce,
		AwardDeadline: req.AwardDeadline,
		Prizes:        convertCsrPrizesToPb(req.Prizes),
		SponsorId:     req.SponsorId,
		IsClocked:     req.IsClocked,
		ClockTask:     convertCsrClockTaskToPb(req),
		PublishType:   req.PublishType,
	}
}

func convertCsrPrizesToPb(csrPrizes []*types.CreatePrize) []*lottery.Prize {
	if csrPrizes == nil {
		return nil
	}
	pbPrizes := make([]*lottery.Prize, 0, len(csrPrizes))
	for _, csrPrize := range csrPrizes {
		pbPrize := lottery.Prize{}
		err := copier.Copy(&pbPrize, csrPrize)
		if err != nil {
			return nil
		}
		pbPrizes = append(pbPrizes, &pbPrize)
	}
	return pbPrizes
}

// 转换打卡任务
// 性能考虑，忽略copier
func convertCsrClockTaskToPb(req *types.CreateLotteryReq) *lottery.ClockTask {
	// 这里用了魔鬼数字，需要自定义枚举
	//
	if req.IsClocked == 1 || req.ClockTask == nil {
		return nil
	}
	return &pb.ClockTask{
		Type:             req.ClockTask.Type,
		Seconds:          req.ClockTask.Seconds,
		AppletType:       req.ClockTask.AppletType,
		PageLink:         req.ClockTask.PageLink,
		AppId:            req.ClockTask.AppId,
		PagePath:         req.ClockTask.PagePath,
		Image:            req.ClockTask.Image,
		VideoAccountId:   req.ClockTask.VideoAccountId,
		VideoId:          req.ClockTask.VideoId,
		ArticleLink:      req.ClockTask.ArticleLink,
		Copywriting:      req.ClockTask.Copywriting,
		ChanceType:       req.ClockTask.ChanceType,
		IncreaseMultiple: req.ClockTask.IncreaseMultiple}
}

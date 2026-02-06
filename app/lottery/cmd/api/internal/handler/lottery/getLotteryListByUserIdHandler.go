// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package lottery

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lottery-be/app/lottery/cmd/api/internal/logic/lottery"
	"lottery-be/app/lottery/cmd/api/internal/svc"
	"lottery-be/app/lottery/cmd/api/internal/types"
)

// 获取当前用户全部/发起/中奖的抽奖列表
func GetLotteryListByUserIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLotteryListByUserIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := lottery.NewGetLotteryListByUserIdLogic(r.Context(), svcCtx)
		resp, err := l.GetLotteryListByUserId(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

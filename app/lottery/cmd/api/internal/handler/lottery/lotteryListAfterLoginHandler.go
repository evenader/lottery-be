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

// 登录后获取抽奖列表
func LotteryListAfterLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LotteryListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := lottery.NewLotteryListAfterLoginLogic(r.Context(), svcCtx)
		resp, err := l.LotteryListAfterLogin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

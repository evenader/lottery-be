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

// CheckIsParticipated
func CheckIsParticipatedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckIsParticipatedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := lottery.NewCheckIsParticipatedLogic(r.Context(), svcCtx)
		resp, err := l.CheckIsParticipated(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

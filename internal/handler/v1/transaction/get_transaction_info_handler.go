// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package transaction

import (
	"net/http"

	"github.com/KyoUK4n/etherscan/internal/logic/v1/transaction"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTransactionInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTransactionInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := transaction.NewGetTransactionInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetTransactionInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

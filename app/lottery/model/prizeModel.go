package model

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"lottery-be/common/xerr"
)

var _ PrizeModel = (*customPrizeModel)(nil)

type (
	// PrizeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPrizeModel.
	PrizeModel interface {
		prizeModel
		FindByLotteryId(ctx context.Context, lotteryId int64) ([]*Prize, error)
	}

	customPrizeModel struct {
		*defaultPrizeModel
	}
)

// NewPrizeModel returns a model for the database table.
func NewPrizeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PrizeModel {
	return &customPrizeModel{
		defaultPrizeModel: newPrizeModel(conn, c, opts...),
	}
}

func (m *defaultPrizeModel) FindByLotteryId(ctx context.Context, lotteryId int64) ([]*Prize, error) {
	var resp []*Prize
	query := fmt.Sprintf("SELECT * FROM %s WHERE lottery_id = ?", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, lotteryId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_FIND_PRIZES_BYLOTTERYID_ERROR), "QueryRowsNoCacheCtx, &resp:%v, query:%v, lotteryId:%v, error: %v", &resp, query, lotteryId, err)
	}
	return resp, nil
}

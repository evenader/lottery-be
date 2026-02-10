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
		FindByLotteryIdWithOrder(ctx context.Context, lotteryId int64) ([]*Prize, error)
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

func (m *defaultPrizeModel) FindByLotteryIdWithOrder(ctx context.Context, lotteryId int64) ([]*Prize, error) {
	var resp []*Prize
	query := fmt.Sprintf("SELECT id, lottery_id, level, count,is_clocked FROM %s WHERE lottery_id = ? ORDER BY level ASC, id ASC", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, lotteryId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_FIND_PRIZES_BYLOTTERYID_ERROR), "QueryRowsNoCacheCtx,  query:%v, lotteryId:%v, error: %v", query, lotteryId, err)
	}
	return resp, nil
}

package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LotteryModel = (*customLotteryModel)(nil)

type (
	// LotteryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLotteryModel.
	LotteryModel interface {
		lotteryModel
		UpdateClockTaskIdOnLottery(ctx context.Context, id int64, clockTaskId int64, opts ...Option) error
	}

	customLotteryModel struct {
		*defaultLotteryModel
	}
)

// NewLotteryModel returns a model for the database table.
func NewLotteryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LotteryModel {
	return &customLotteryModel{
		defaultLotteryModel: newLotteryModel(conn, c, opts...),
	}
}

/*
1、 动态拼接字符串有性能瓶颈在高并发环境下
2、替代：
  - 使用预编译语句，减少字符串拼接的开销
  - strings.Builder
*/
func (m *defaultLotteryModel) UpdateClockTaskIdOnLottery(ctx context.Context, id int64, clockTaskId int64, opts ...Option) error {
	// 1. 解析 Options (保持原样，为了支持事务)
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	// 2. 准备缓存 Key
	lotteryLotteryIdKey := fmt.Sprintf("%s%v", cacheLotteryLotteryIdPrefix, id)

	// 3. 执行局部更新
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		// 【核心改动】：SQL 只针对单个字段
		query := fmt.Sprintf("update %s set `clock_task_id` = ? where `id` = ?", m.table)

		// 逻辑分流：判断是否在事务 Session 中
		if o.Session != nil {
			return o.Session.ExecCtx(ctx, query, clockTaskId, id)
		}
		return conn.ExecCtx(ctx, query, clockTaskId, id)
	}, lotteryLotteryIdKey)

	return err
}

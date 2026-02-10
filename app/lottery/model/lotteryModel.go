package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"lottery-be/common/xerr"
	"time"
)

var _ LotteryModel = (*customLotteryModel)(nil)

type (
	// LotteryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLotteryModel.
	LotteryModel interface {
		lotteryModel
		UpdateClockTaskIdOnLottery(ctx context.Context, id int64, clockTaskId int64, opts ...Option) error
		SearchTimeOutIds(ctx context.Context, currentTime time.Time, announceType int64) ([]int64, error)
		UpdateStatusToAnnounced(ctx context.Context, id int64) (sql.Result, error)
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

func (m *defaultLotteryModel) SearchTimeOutIds(ctx context.Context, currentTime time.Time, announceType int64) ([]int64, error) {
	var resp []int64
	// 1. 只查 id，保持与 []int64 匹配
	// 2. 将 announce_type 改为动态占位符
	query := fmt.Sprintf("SELECT id FROM %s WHERE announce_type = ? AND is_announced = 0 AND del_state = 0 AND announce_time <= ?", m.table)

	// 注意：参数顺序必须与 SQL 中的 ? 顺序严格一致
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, announceType, currentTime)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.GETLOTTERY_BYLESSTHAN_CURRENTTIME_ERROR),
			"SearchTimeOutIds fail, announceType:%v, currentTime:%v, error: %v", announceType, currentTime, err)
	}
	return resp, nil
}

func (m *defaultLotteryModel) UpdateStatusToAnnounced(ctx context.Context, id int64) (sql.Result, error) {
	query := fmt.Sprintf("update %s set is_announced = 1 where id = ? and is_announced = 0", m.table)

	// go-zero 的 model 接口直接提供了 ExecCtx 方法
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query, id)
	})
}

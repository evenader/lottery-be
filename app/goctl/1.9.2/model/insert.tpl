{{/* 这是老版本的 Insert 逻辑，留作基石参考
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
    ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
    ret,err:=m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return ret,err
}
*/}}


func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}, opts ...Option) (sql.Result, error) {
	// 1. 解析 Options
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	{{if .withCache}}
	// 2. 带缓存模式：利用 m.ExecCtx 保证数据写入后缓存被清理
	{{.keys}}
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)

		// 逻辑分流：如果 Options 传入了 Session，则在事务内执行
		if o.Session != nil {
			return o.Session.ExecCtx(ctx, query, {{.expressionValues}})
		}
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}})
	{{else}}
	// 3. 无缓存模式：直接判断 Session
	query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
	if o.Session != nil {
		return o.Session.ExecCtx(ctx, query, {{.expressionValues}})
	}
	ret, err := m.conn.ExecCtx(ctx, query, {{.expressionValues}})
	{{end}}

	return ret, err
}
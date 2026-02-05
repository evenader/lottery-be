{{/* 老版本 Delete 逻辑备份
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		return conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		_,err:=m.conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}
*/}}

func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}, opts ...Option) error {
	// 1. 解析 Options
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	{{if .withCache}}
	{{if .containsIndexCache}}
	// 2. 如果存在唯一索引缓存，必须先查出完整数据，否则无法构建所有索引的 Redis Key
	data, err := m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err != nil {
		return err
	}
	{{end}}

	{{.keys}}
	_, err {{if .containsIndexCache}}={{else}}:={{end}} m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)

		// 逻辑分流：判断是否在事务 Session 中执行
		if o.Session != nil {
			return o.Session.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
		}
		return conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}})
	{{else}}
	// 3. 无缓存模式
	query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
	if o.Session != nil {
		_, err := o.Session.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
		return err
	}
	_, err := m.conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
	{{end}}
	return err
}
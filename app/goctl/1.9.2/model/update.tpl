{{/* 老版本 Update 逻辑备份
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    _, {{if .containsIndexCache}}err{{else}}err:{{end}}= m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return err
}
*/}}

func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}, opts ...Option) error {
	// 1. 解析 Options
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	{{if .withCache}}
	{{if .containsIndexCache}}
	// 2. 如果存在唯一索引缓存，需要先查出旧数据以清理旧索引缓存
	data, err := m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}})
	if err != nil {
		return err
	}
	{{end}}

	{{.keys}}
	_, {{if .containsIndexCache}}err{{else}}err :{{end}}= m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)

		// 逻辑分流：判断是否在事务 Session 中
		if o.Session != nil {
			return o.Session.ExecCtx(ctx, query, {{.expressionValues}})
		}
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}})
	{{else}}
	// 3. 无缓存模式
	query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
	if o.Session != nil {
		_, err := o.Session.ExecCtx(ctx, query, {{.expressionValues}})
		return err
	}
	_, err := m.conn.ExecCtx(ctx, query, {{.expressionValues}})
	{{end}}
	return err
}
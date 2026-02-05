type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}sqlc.CachedConn{{else}}conn sqlx.SqlConn{{end}}
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)


package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"lottery-be/app/usercenter/cmd/rpc/internal/config"
	"lottery-be/app/usercenter/model"
)

type ServiceContext struct {
	Config        config.Config
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config:        c,
		UserModel:     model.NewUserModel(sqlConn, c.Cache),
		UserAuthModel: model.NewUserAuthModel(sqlConn, c.Cache),
	}
}

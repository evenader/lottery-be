package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"lottery-be/app/usercenter/model"
	"lottery-be/common/utils"
	"lottery-be/common/xerr"

	"lottery-be/app/usercenter/cmd/rpc/internal/svc"
	"lottery-be/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const NickNameDftLength = 8

var ErrUserAlreadyRegisterError = xerr.NewErrMsg("user has been registered")

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "mobile:%s,err:%v", in.Mobile, err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "Register user exists mobile:%s,err:%v", in.Mobile, err)
	}

	nickname := in.Nickname
	if len(nickname) == 0 {
		nickname = utils.Krand(NickNameDftLength, utils.KC_RAND_KIND_ALL) // 待优化 这里会有碰撞风险
	}
	password, err := utils.BcryptWithString(in.Password) // 建议换成 bcrypt
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "Bcrypt password hash failed, err:%v", err)
	}

	//trans是重资源，非必要计算放在外面
	var userId int64
	if err := l.svcCtx.UserModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		user := &model.User{
			Mobile:   in.Mobile,
			Nickname: nickname,
			Password: password,
		}

		insertRes, err := l.svcCtx.UserModel.Insert(ctx, user, model.WithSession(session))
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user Insert err:%v,user:%+v", err, user)
		}

		lastId, err := insertRes.LastInsertId() // todo 优化
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user insertResult.LastInsertId err:%v,user:%+v", err, user)
		}
		userId = lastId
		userAuth := &model.UserAuth{
			UserId:   lastId,
			AuthKey:  in.AuthKey,
			AuthType: in.AuthType}

		if _, err = l.svcCtx.UserAuthModel.Insert(ctx, userAuth, model.WithSession(session)); err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user_auth Insert err:%v,userAuth:%v", err, userAuth)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "db trans exec fail err:%v", err)
	}
	l.Logger.Info("-----> DEBUG: Current userId to generate token: %d \n", userId)
	tokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenRsp, err := tokenLogic.GenerateToken(&pb.GenerateTokenReq{
		Id: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", userId)
	}
	// 🌟 核心问题：在这里打印一下，看 tokenRsp 到底有没有值
	l.Logger.Info("-----> TOKEN DEBUG: %s ", tokenRsp.AccessToken)
	return &pb.RegisterResp{
		AccessToken:  tokenRsp.AccessToken,
		AccessExpire: tokenRsp.AccessExpire,
		RefreshAfter: tokenRsp.RefreshAfter,
	}, nil
}

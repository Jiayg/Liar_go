package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/jiayg/liar/internal/model"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/internal/service/internal/do"
	"github.com/jiayg/liar/library/libUtils"
	"github.com/jiayg/liar/library/liberr"
	"github.com/mssola/user_agent"
)

type ILoginLog interface {
	Invoke(ctx context.Context, data *model.LoginLogParams)
	Update(ctx context.Context, id uint64, ip string) (err error)
}

type loginLogImpl struct {
	Pool *grpool.Pool
}

var (
	loginLogService = loginLogImpl{
		Pool: grpool.New(100),
	}
)

func LoginLog() ILoginLog {
	return ILoginLog(&loginLogService)
}

func (s *loginLogImpl) Invoke(ctx context.Context, data *model.LoginLogParams) {
	s.Pool.Add(
		ctx,
		func(ctx context.Context) {
			//写入日志数据
			ua := user_agent.New(data.UserAgent)
			browser, _ := ua.Browser()
			loginData := &do.SysLoginLog{
				LoginName:     data.Username,
				Ipaddr:        data.Ip,
				LoginLocation: libUtils.GetCityByIp(data.Ip),
				Browser:       browser,
				Os:            ua.OS(),
				Status:        data.Status,
				Msg:           data.Msg,
				LoginTime:     gtime.Now(),
				Module:        data.Module,
			}
			_, err := dao.SysLoginLog.Ctx(ctx).Insert(loginData)
			if err != nil {
				g.Log().Error(ctx, err)
			}
		},
	)
}

func (s *loginLogImpl) Update(ctx context.Context, id uint64, ip string) (err error) {
	g.Try(func() {
		_, err = dao.SysUser.Ctx(ctx).WherePri(id).Update(g.Map{
			dao.SysUser.Columns().LastLoginIp:   ip,
			dao.SysUser.Columns().LastLoginTime: gtime.Now(),
		})
		liberr.ErrIsNil(ctx, err, "更新用户登录信息失败")
	})
	return
}

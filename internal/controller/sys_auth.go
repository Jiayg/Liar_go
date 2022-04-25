package controller

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gmode"
	"github.com/jiayg/liar/apiv1"
	"github.com/jiayg/liar/internal/model"
	"github.com/jiayg/liar/internal/service"
	"github.com/jiayg/liar/library/libUtils"
)

var Auth = authController{}

type authController struct {
}

// 验证码
func (c *authController) Captcha(ctx context.Context, req *apiv1.CaptchaReq) (res *apiv1.CaptchaRes, err error) {
	var (
		idKeyC, base64stringC string
	)
	idKeyC, base64stringC, err = service.Captcha().GetVerifyImgString(ctx)
	res = &apiv1.CaptchaRes{
		Key: idKeyC,
		Img: base64stringC,
	}
	return
}

// 登录
func (c *authController) Login(ctx context.Context, req *apiv1.LoginReq) (res *apiv1.LoginRes, err error) {
	var (
		user *model.LoginUserRes
		// token       string
		// permissions []string
		// menuList    []*model.UserMenus
	)
	//判断验证码是否正确
	debug := gmode.IsDevelop()
	if !debug {
		if !service.Captcha().VerifyString(req.VerifyKey, req.VerifyCode) {
			err = gerror.New("验证码输入错误")
			return
		}
	}

	ip := libUtils.GetClientIp(ctx)
	userAgent := libUtils.GetUserAgent(ctx)
	user, err = service.User().GetUserByLogin(ctx, req)
	if err != nil {
		// 保存登录失败的日志信息
		service.LoginLog().Invoke(ctx, &model.LoginLogParams{
			Status:    0,
			Username:  req.Username,
			Ip:        ip,
			UserAgent: userAgent,
			Msg:       err.Error(),
			Module:    "系统后台",
		})
		return
	}
	err = service.LoginLog().Update(ctx, user.Id, ip)
	if err != nil {
		return
	}
	// 报存登录成功的日志信息
	service.LoginLog().Invoke(ctx, &model.LoginLogParams{
		Status:    1,
		Username:  req.Username,
		Ip:        ip,
		UserAgent: userAgent,
		Msg:       "登录成功",
		Module:    "系统后台",
	})
	// key := gconv.String(user.Id) + "-" + gmd5.MustEncryptString(user.UserName) + gmd5.MustEncryptString(user.UserPassword)
	// if g.Cfg().MustGet(ctx, "gfToken.multiLogin").Bool() {
	// 	key = gconv.String(user.Id) + "-" + gmd5.MustEncryptString(user.UserName) + gmd5.MustEncryptString(user.UserPassword+ip+userAgent)
	// }
	// user.UserPassword = ""
	// token, err = service.GfToken().GenerateToken(ctx, key, user)
	// if err != nil {
	// 	return
	// }
	// //获取用户菜单数据
	// menuList, permissions, err = service.User().GetUserRoles(ctx, user.Id)
	// if err != nil {
	// 	return
	// }
	// res = &apiv1.LoginRes{
	// 	UserInfo:    user,
	// 	Token:       token,
	// 	MenuList:    menuList,
	// 	Permissions: permissions,
	// }
	return
}

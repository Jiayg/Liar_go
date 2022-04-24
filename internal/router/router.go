package router

import (
	"github.com/jiayg/liar/internal/controller"
	"github.com/jiayg/liar/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

func BindController(group *ghttp.RouterGroup) {
	group.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Middleware(middleware.Middleware().MiddlewareCORS)
		sysRouter(group)
		demoRouter(group)
		commRouter(group)
	})

}

// 后台路由
func sysRouter(group *ghttp.RouterGroup) {
	group.Group("/system", func(group *ghttp.RouterGroup) {
		// 系统初始化
		// group.Bind(
		// 	controller.DbInit,
		// )
		//登录验证拦截
		// service.GfToken().Middleware(group)
		//context拦截器
		// group.Middleware(service.Middleware().Ctx, service.Middleware().Auth)
		group.Bind(
			controller.User,
			controller.Role,
		)
	})
}

// 绑定测试路由
func demoRouter(group *ghttp.RouterGroup) {
	group.Group("/demo", func(group *ghttp.RouterGroup) {
		group.Bind(
			controller.Demo,
		)
	})
}

// 绑定公共路由
func commRouter(group *ghttp.RouterGroup) {
	group.Group("/pub", func(group *ghttp.RouterGroup) {
		group.Group("/captcha", func(group *ghttp.RouterGroup) {
			group.Bind(
				controller.Captcha,
			)
		})
	})
}

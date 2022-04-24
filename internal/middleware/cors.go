package middleware

import "github.com/gogf/gf/v2/net/ghttp"

type middlewareImpl struct{}

var middlewareService = middlewareImpl{}

type IMiddleware interface {
	MiddlewareCORS(r *ghttp.Request)
}

func Middleware() IMiddleware {
	return IMiddleware(&middlewareImpl{})
}

func (s *middlewareImpl) MiddlewareCORS(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	// you can set options
	//corsOptions.AllowDomain = []string{"goframe.org", "baidu.com"}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

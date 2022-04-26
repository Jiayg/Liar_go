package service

import (
	"context"
	"sync"

	"github.com/gogf/gf/os/gctx"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jiayg/liar/internal/consts"
	"github.com/jiayg/liar/internal/model"
	"github.com/jiayg/liar/library/liberr"
	"github.com/tiger1103/gfast-token/gftoken"
)

type IGfToken interface {
	GenerateToken(ctx context.Context, key string, data interface{}) (keys string, err error)
	Middleware(group *ghttp.RouterGroup) error
	ParseToken(r *ghttp.Request) (*gftoken.CustomClaims, error)
	IsLogin(r *ghttp.Request) (b bool, failed *gftoken.AuthFailed)
}

type gfTokenImpl struct {
	*gftoken.GfToken
}

var gT = gfTokenImpl{
	GfToken: gftoken.NewGfToken(),
}

func NewGfToken(options *model.TokenOptions) IGfToken {
	var fun gftoken.OptionFunc
	if options.CacheModel == consts.CacheModelRedis {
		fun = gftoken.WithGRedis()
	} else {
		fun = gftoken.WithGCache()
	}
	gT.GfToken = gftoken.NewGfToken(
		gftoken.WithCacheKey(options.CacheKey),
		gftoken.WithTimeout(options.Timeout),
		gftoken.WithMaxRefresh(options.MaxRefresh),
		gftoken.WithMultiLogin(options.MultiLogin),
		gftoken.WithExcludePaths(options.ExcludePaths),
		fun,
	)
	return IGfToken(&gT)
}

type gft struct {
	options *model.TokenOptions
	gT      IGfToken
	lock    *sync.Mutex
}

var gftService = &gft{
	options: nil,
	gT:      nil,
	lock:    &sync.Mutex{},
}

func GfToken() IGfToken {
	if gftService.gT == nil {
		gftService.lock.Lock()
		defer gftService.lock.Unlock()
		if gftService.gT == nil {
			ctx := gctx.New()
			err := g.Cfg().MustGet(ctx, "gfToken").Struct(&gftService.options)
			liberr.ErrIsNil(ctx, err)
			gftService.gT = NewGfToken(gftService.options)
		}
	}
	return gftService.gT
}

package controller

import (
	"context"

	"github.com/jiayg/liar/apiv1"
	"github.com/jiayg/liar/internal/service"
)

var User = userController{}

type userController struct {
	BaseController
}

// 添加用户
func (c *userController) Add(ctx context.Context, req *apiv1.UserAddReq) (res *apiv1.UserAddRes, err error) {
	err = service.User().Add(ctx, req)
	return
}

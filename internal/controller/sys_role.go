package controller

import (
	"context"

	"github.com/jiayg/liar/apiv1"
	"github.com/jiayg/liar/internal/service"
)

var Role = roleController{}

type roleController struct {
	BaseController
}

// 角色分页列表
func (c *roleController) GetPageList(ctx context.Context, req *apiv1.RoleListReq) (res *apiv1.RoleListRes, err error) {
	res, err = service.Role().GetPageList(ctx, req)
	return
}

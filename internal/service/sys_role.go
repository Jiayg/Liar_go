package service

import (
	"context"

	"github.com/jiayg/liar/apiv1"
	"github.com/jiayg/liar/internal/consts"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/library/liberr"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type IRole interface {
	// GetRoleList(ctx context.Context) (list []*entity.SysRole, err error)
	GetRoleListSearch(ctx context.Context, req *apiv1.RoleListReq) (res *apiv1.RoleListRes, err error)
	// AddRole(ctx context.Context, req *apiv1.RoleAddReq) (err error)
	// Get(ctx context.Context, id uint) (res *entity.SysRole, err error)
	// GetFilteredNamedPolicy(ctx context.Context, id uint) (gpSlice []int, err error)
	// EditRole(ctx context.Context, req *apiv1.RoleEditReq) error
	// DeleteByIds(ctx context.Context, ids []int64) (err error)
}

type roleImpl struct {
}

var roleService = roleImpl{}

func Role() IRole {
	return IRole(&roleService)
}

func (s *roleImpl) GetRoleListSearch(ctx context.Context, req *apiv1.RoleListReq) (res *apiv1.RoleListRes, err error) {
	res = new(apiv1.RoleListRes)
	g.Try(func() {
		model := dao.SysRole.Ctx(ctx)
		if req.RoleName != "" {
			model = model.Where("name like ?", "%"+req.RoleName+"%")
		}
		if req.Status != "" {
			model = model.Where("status", gconv.Int(req.Status))
		}
		res.Total, err = model.Count()
		liberr.ErrIsNil(ctx, err, "获取角色数据失败")
		if req.PageNum == 0 {
			req.PageNum = 1
		}
		res.CurrentPage = req.PageNum
		if req.PageSize == 0 {
			req.PageSize = consts.PageSize
		}
		err = model.Page(res.CurrentPage, req.PageSize).Order("id asc").Scan(&res.List)
		liberr.ErrIsNil(ctx, err, "获取数据失败")
	})
	return
}

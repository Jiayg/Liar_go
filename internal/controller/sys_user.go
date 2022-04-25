package controller

import (
	"context"

	"github.com/jiayg/liar/apiv1"
	"github.com/jiayg/liar/internal/model/entity"
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

// 删除用户
func (c *userController) Delete(ctx context.Context, req *apiv1.UserDeleteReq) (res *apiv1.UserDeleteRes, err error) {
	err = service.User().Delete(ctx, req.Ids)
	return
}

// 修改用户
func (c *userController) Update(ctx context.Context, req *apiv1.UserUpdateReq) (res *apiv1.UserUpdateRes, err error) {
	err = service.User().Update(ctx, req)
	return
}

// 获取用户
func (c *userController) Get(ctx context.Context, req *apiv1.UserGetReq) (res *apiv1.UserGetRes, err error) {
	res, err = service.User().Get(ctx, req.Id)
	return
}

// TODO 需要重构，多角色数据一起查询
// 获取用户列表
func (c *userController) List(ctx context.Context, req *apiv1.UserSearchReq) (res *apiv1.UserSearchRes, err error) {
	var (
		total int
		list  []*entity.SysUser
	)
	res = new(apiv1.UserSearchRes)
	total, list, err = service.User().List(ctx, req)
	if err != nil || total == 0 {
		return
	}
	res.Total = total
	res.UserList, err = service.User().GetUsersRoleDept(ctx, list)
	return
}

// 修改用户状态
func (c *userController) ChangeStatus(ctx context.Context, req *apiv1.UserStatusReq) (res *apiv1.UserStatusRes, err error) {
	err = service.User().ChangeStatus(ctx, req)
	return
}

// 重置密码
func (c *userController) ResetPwd(ctx context.Context, req *apiv1.UserResetPwdReq) (res *apiv1.UserResetPwdRes, err error) {
	err = service.User().ResetPwd(ctx, req)
	return
}

package service

import (
	"context"
	"fmt"

	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/jiayg/liar/apiv1"
	"github.com/jiayg/liar/internal/model/entity"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/internal/service/internal/do"
	"github.com/jiayg/liar/library/libUtils"
	"github.com/jiayg/liar/library/liberr"
)

type IUser interface {
	Add(ctx context.Context, req *apiv1.UserAddReq) (err error)
}

type userImpl struct{}

var (
	notCheckAuthAdminIds *gset.Set
	userService          = userImpl{}
)

func User() IUser {
	return IUser(&userService)
}

// 检查用户名和手机号是否存在
func (s *userImpl) userNameOrMobileExists(ctx context.Context, userName, mobile string, id ...int64) error {
	user := (*entity.SysUser)(nil)
	err := g.Try(func() {
		model := dao.SysUser.Ctx(ctx)
		if len(id) > 0 {
			model = model.Where(dao.SysUser.Columns().Id+"!=", id)
		}
		model = model.Where(fmt.Sprintf("%s='%s' OR %s='%s'",
			dao.SysUser.Columns().UserName,
			userName,
			dao.SysUser.Columns().Mobile,
			mobile))
		err := model.Limit(1).Scan(&user)
		liberr.ErrIsNil(ctx, err, "获取用户信息失败")
		if user == nil {
			return
		}
		if user.UserName == userName {
			liberr.ErrIsNil(ctx, gerror.New("用户名已存在"))
		}
		if user.Mobile == mobile {
			liberr.ErrIsNil(ctx, gerror.New("手机号已存在"))
		}
	})
	return err
}

// 添加用户角色
func (s *userImpl) addUserRole(ctx context.Context, roleIds []int64, userId int64) (err error) {
	err = g.Try(func() {
		enforcer, e := CasbinEnforcer(ctx)
		liberr.ErrIsNil(ctx, e)
		for _, v := range roleIds {
			_, e = enforcer.AddGroupingPolicy(gconv.String(userId), gconv.String(v))
			liberr.ErrIsNil(ctx, e)
		}
	})
	return
}

// 添加用户
func (s *userImpl) Add(ctx context.Context, req *apiv1.UserAddReq) (err error) {
	err = s.userNameOrMobileExists(ctx, req.UserName, req.Mobile)
	if err != nil {
		return
	}
	req.UserSalt = grand.S(10)
	req.Password = libUtils.EncryptPassword(req.Password, req.UserSalt)
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			userId, e := dao.SysUser.Ctx(ctx).TX(tx).InsertAndGetId(do.SysUser{
				UserName:     req.UserName,
				Mobile:       req.Mobile,
				UserNickname: req.NickName,
				UserPassword: req.Password,
				UserSalt:     req.UserSalt,
				UserStatus:   req.Status,
				UserEmail:    req.Email,
				Sex:          req.Sex,
				DeptId:       req.DeptId,
				Remark:       req.Remark,
				IsAdmin:      req.IsAdmin,
			})
			liberr.ErrIsNil(ctx, e, "添加用户失败")
			e = s.addUserRole(ctx, req.RoleIds, userId)
			liberr.ErrIsNil(ctx, e, "设置用户权限失败")
			e = s.AddUserPost(ctx, tx, req.PostIds, userId)
			liberr.ErrIsNil(ctx, e)
		})
		return err
	})
	return
}

// 添加用户岗位
func (s *userImpl) AddUserPost(ctx context.Context, tx *gdb.TX, postIds []int64, userId int64) (err error) {
	err = g.Try(func() {
		// 删除旧岗位信息
		_, err = dao.SysUserPost.Ctx(ctx).TX(tx).Where(dao.SysUserPost.Columns().UserId, userId).Delete()
		liberr.ErrIsNil(ctx, err, "设置用户岗位失败")
		if len(postIds) == 0 {
			return
		}
		// 添加岗位信息
		data := g.List{}
		for _, v := range postIds {
			data = append(data, g.Map{
				dao.SysUserPost.Columns().UserId: userId,
				dao.SysUserPost.Columns().PostId: v,
			})
		}
		_, err = dao.SysUserPost.Ctx(ctx).TX(tx).Data(data).Insert()
		liberr.ErrIsNil(ctx, err, "设置用户岗位失败")
	})
	return
}

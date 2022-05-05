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
	"github.com/jiayg/liar/internal/consts"
	"github.com/jiayg/liar/internal/model"
	"github.com/jiayg/liar/internal/model/entity"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/internal/service/internal/do"
	"github.com/jiayg/liar/library/libUtils"
	"github.com/jiayg/liar/library/liberr"
)

type IUser interface {
	GetUserByLogin(ctx context.Context, req *apiv1.LoginReq) (user *model.LoginUserRes, err error)
	Add(ctx context.Context, req *apiv1.UserAddReq) (err error)
	Delete(ctx context.Context, ids []int) (err error)
	Update(ctx context.Context, req *apiv1.UserUpdateReq) (err error)
	Get(ctx context.Context, id uint64) (res *apiv1.UserGetRes, err error)
	GetPageList(ctx context.Context, req *apiv1.UserSearchReq) (total int, list []*entity.SysUser, err error)
	GetUserRolesDepts(ctx context.Context, list []*entity.SysUser) (users []*model.SysUserRoleDeptRes, err error)
	GetUserMenus(ctx context.Context, userId uint64) (menuList []*model.UserMenus, permissions []string, err error)
	ChangeStatus(ctx context.Context, req *apiv1.UserStatusReq) (err error)
	ResetPwd(ctx context.Context, req *apiv1.UserResetPwdReq) (err error)
	NotCheckAuthUserIds(ctx context.Context) *gset.Set
}

type userImpl struct{}

var (
	notCheckAuthUserIds *gset.Set
	userService         = userImpl{}
)

func User() IUser {
	return IUser(&userService)
}

func (s *userImpl) NotCheckAuthUserIds(ctx context.Context) *gset.Set {
	ids := g.Cfg().MustGet(ctx, "system.notCheckAuthAdminIds")
	if !g.IsNil(ids) {
		notCheckAuthUserIds = gset.NewFrom(ids)
	}
	return notCheckAuthUserIds
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

// 根据id获取用户信息
func (s *userImpl) getUserById(ctx context.Context, id uint64, withPwd ...bool) (user *entity.SysUser, err error) {
	err = g.Try(func() {
		// 用户信息
		if len(withPwd) > 0 && withPwd[0] {
			err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, id).Scan(&user)
		} else {
			err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, id).
				FieldsEx(dao.SysUser.Columns().UserPassword, dao.SysUser.Columns().UserSalt).Scan(&user)
		}
		liberr.ErrIsNil(ctx, err, "获取用户数据失败")
	})
	return
}

// 根据用户账户密码返回登录信息
func (s *userImpl) GetUserByLogin(ctx context.Context, req *apiv1.LoginReq) (user *model.LoginUserRes, err error) {
	err = g.Try(func() {
		user = &model.LoginUserRes{}
		err = dao.SysUser.Ctx(ctx).Fields(user).Where(dao.SysUser.Columns().UserName, req.Username).Scan(user)
		liberr.ErrIsNil(ctx, err)
		liberr.ValueIsNil(user, "账号错误")
		//验证密码
		if libUtils.EncryptPassword(req.Password, user.UserSalt) != user.UserPassword {
			liberr.ErrIsNil(ctx, gerror.New("密码错误"))
		}
		//账号状态
		if user.UserStatus == 0 {
			liberr.ErrIsNil(ctx, gerror.New("账号已被冻结"))
		}
	})
	return
}

// 添加用户信息
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
			e = roleService.addRole(ctx, req.RoleIds, userId)
			liberr.ErrIsNil(ctx, e, "设置用户权限失败")
			e = userPostService.addPost(ctx, tx, req.PostIds, userId)
			liberr.ErrIsNil(ctx, e)
		})
		return err
	})
	return
}

// 根据id删除用户
func (s *userImpl) Delete(ctx context.Context, ids []int) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			// 删除账户
			_, err = dao.SysUser.Ctx(ctx).TX(tx).Where(dao.SysUser.Columns().Id+" in(?)", ids).Delete()
			liberr.ErrIsNil(ctx, err, "删除用户失败")
			// 删除账户权限
			enforcer, e := CasbinEnforcer(ctx)
			liberr.ErrIsNil(ctx, e)
			for _, v := range ids {
				enforcer.RemoveFilteredGroupingPolicy(0, gconv.String(v))
			}
			// 删除账户岗位
			_, err = dao.SysUserPost.Ctx(ctx).TX(tx).Delete(dao.SysUserPost.Columns().UserId+" in(?)", ids)
			liberr.ErrIsNil(ctx, err, "删除用户岗位失败")
		})
		return err
	})
	return
}

// 修改用户信息
func (s *userImpl) Update(ctx context.Context, req *apiv1.UserUpdateReq) (err error) {
	err = s.userNameOrMobileExists(ctx, "", req.Mobile, req.UserId)
	if err != nil {
		return
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			_, err = dao.SysUser.Ctx(ctx).TX(tx).WherePri(req.UserId).Update(do.SysUser{
				Mobile:       req.Mobile,
				UserNickname: req.NickName,
				UserStatus:   req.Status,
				UserEmail:    req.Email,
				Sex:          req.Sex,
				DeptId:       req.DeptId,
				Remark:       req.Remark,
				IsAdmin:      req.IsAdmin,
			})
			liberr.ErrIsNil(ctx, err, "修改用户信息失败")
			// 设置用户所属角色信息
			err = roleService.updRole(ctx, req.RoleIds, req.UserId)
			liberr.ErrIsNil(ctx, err, "设置用户权限失败")
			// 设置用户岗位
			err = userPostService.addPost(ctx, tx, req.PostIds, req.UserId)
			liberr.ErrIsNil(ctx, err)
		})
		return err
	})
	return
}

// 获取用户信息
func (s *userImpl) Get(ctx context.Context, id uint64) (res *apiv1.UserGetRes, err error) {
	res = new(apiv1.UserGetRes)
	err = g.Try(func() {
		// 获取用户信息
		res.User, err = s.getUserById(ctx, id)
		liberr.ErrIsNil(ctx, err)
		// 获取用户角色ids
		res.CheckedRoleIds, err = roleService.getRoleIdsByUserId(ctx, id)
		liberr.ErrIsNil(ctx, err)
		// 获取用户部门ids
		res.CheckedPostIds, err = userPostService.getUserPostIdsByUserId(ctx, id)
		liberr.ErrIsNil(ctx, err)
	})
	return
}

// 获取用户分页列表
func (s *userImpl) GetPageList(ctx context.Context, req *apiv1.UserSearchReq) (total int, list []*entity.SysUser, err error) {
	err = g.Try(func() {
		model := dao.SysUser.Ctx(ctx)
		if req.KeyWords != "" {
			keyWords := "%" + req.KeyWords + "%"
			model = model.Where("user_name like ? or  user_nickname like ?", keyWords, keyWords)
		}
		if req.DeptId != "" {
			deptIds, e := deptService.getDeptIdsById(ctx, gconv.Int64(req.DeptId))
			liberr.ErrIsNil(ctx, e)
			model = model.Where("dept_id in (?)", deptIds)
		}
		if req.Status != "" {
			model = model.Where("user_status", gconv.Int(req.Status))
		}
		if req.Mobile != "" {
			model = model.Where("mobile like ?", "%"+req.Mobile+"%")
		}
		if len(req.DateRange) > 0 {
			model = model.Where("created_at >=? AND created_at <=?", req.DateRange[0], req.DateRange[1])
		}
		if req.PageSize == 0 {
			req.PageSize = consts.PageSize
		}
		if req.PageIndex == 0 {
			req.PageIndex = 1
		}
		total, err = model.Count()
		liberr.ErrIsNil(ctx, err, "获取用户数据失败")
		err = model.FieldsEx(dao.SysUser.Columns().UserPassword, dao.SysUser.Columns().UserSalt).
			Page(req.PageIndex, req.PageSize).Order("id asc").Scan(&list)
		liberr.ErrIsNil(ctx, err, "获取用户列表失败")
	})
	return
}

// 获取用户角色 部门信息
func (s *userImpl) GetUserRolesDepts(ctx context.Context, list []*entity.SysUser) (users []*model.SysUserRoleDeptRes, err error) {
	err = g.Try(func() {
		allRoles, e := Role().GetRoleList(ctx)
		liberr.ErrIsNil(ctx, e)
		depts, e := Dept().GetCacheDepts(ctx)
		liberr.ErrIsNil(ctx, e)
		users = make([]*model.SysUserRoleDeptRes, len(list))
		for k, u := range list {
			var dept *entity.SysDept
			users[k] = &model.SysUserRoleDeptRes{
				SysUser: u,
			}
			for _, d := range depts {
				if u.DeptId == uint64(d.DeptId) {
					dept = d
				}
			}
			users[k].Dept = dept
			roles, e := roleService.getRolesByUserId(ctx, u.Id, allRoles)
			liberr.ErrIsNil(ctx, e)
			for _, r := range roles {
				users[k].RoleInfo = append(users[k].RoleInfo, &model.SysUserRoleInfoRes{RoleId: r.Id, Name: r.Name})
			}
		}
	})
	return
}

// 获取用户菜单数据
func (s *userImpl) GetUserMenus(ctx context.Context, userId uint64) (menuList []*model.UserMenus, permissions []string, err error) {
	err = g.Try(func() {
		//是否超管
		isSuperAdmin := false
		//获取无需验证权限的用户id
		s.NotCheckAuthUserIds(ctx).Iterator(func(v interface{}) bool {
			if gconv.Uint64(v) == userId {
				isSuperAdmin = true
				return false
			}
			return true
		})
		//获取用户菜单数据
		allRoles, err := Role().GetRoleList(ctx)
		liberr.ErrIsNil(ctx, err)
		roles, err := roleService.getRolesByUserId(ctx, userId, allRoles)
		liberr.ErrIsNil(ctx, err)
		name := make([]string, len(roles))
		roleIds := make([]uint, len(roles))
		for k, v := range roles {
			name[k] = v.Name
			roleIds[k] = v.Id
		}
		//获取菜单信息
		if isSuperAdmin {
			//超管获取所有菜单
			permissions = []string{"*/*/*"}
			menuList, err = roleService.getRolesMenus(ctx)
			liberr.ErrIsNil(ctx, err)
		} else {
			// menuList, err = s.GetAdminMenusByRoleIds(ctx, roleIds)
			// liberr.ErrIsNil(ctx, err)
			// permissions, err = s.GetPermissions(ctx, roleIds)
			// liberr.ErrIsNil(ctx, err)
		}
	})
	return
}

// 修改用户状态
func (s *userImpl) ChangeStatus(ctx context.Context, req *apiv1.UserStatusReq) (err error) {
	err = g.Try(func() {
		_, err = dao.SysUser.Ctx(ctx).WherePri(req.Id).Update(do.SysUser{UserStatus: req.UserStatus})
		liberr.ErrIsNil(ctx, err, "设置用户状态失败")
	})
	return
}

// 重置密码
func (s *userImpl) ResetPwd(ctx context.Context, req *apiv1.UserResetPwdReq) (err error) {
	salt := grand.S(10)
	password := libUtils.EncryptPassword(req.Password, salt)
	err = g.Try(func() {
		_, err = dao.SysUser.Ctx(ctx).WherePri(req.Id).Update(g.Map{
			dao.SysUser.Columns().UserSalt:     salt,
			dao.SysUser.Columns().UserPassword: password,
		})
		liberr.ErrIsNil(ctx, err, "重置用户密码失败")
	})
	return
}

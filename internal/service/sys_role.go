package service

import (
	"context"

	"github.com/jiayg/liar/apiv1"
	"github.com/jiayg/liar/internal/consts"
	"github.com/jiayg/liar/internal/model"
	"github.com/jiayg/liar/internal/model/entity"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/library/liberr"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type IRole interface {
	GetRoleList(ctx context.Context) (list []*entity.SysRole, err error)
	GetPageList(ctx context.Context, req *apiv1.RoleListReq) (res *apiv1.RoleListRes, err error)
	// AddRole(ctx context.Context, req *apiv1.RoleAddReq) (err error)
	// Get(ctx context.Context, id uint) (res *entity.SysRole, err error)
	// GetFilteredNamedPolicy(ctx context.Context, id uint) (gpSlice []int, err error)
	// EditRole(ctx context.Context, req *apiv1.RoleEditReq) error
	// DeleteByIds(ctx context.Context, ids []int64) (err error)
}

type roleImpl struct{}

var roleService = roleImpl{}

func Role() IRole {
	return IRole(&roleService)
}

// 获取角色ids
func (s *roleImpl) getRoleIdsByUserId(ctx context.Context, userId uint64) (roleIds []uint, err error) {
	enforcer, e := CasbinEnforcer(ctx)
	if e != nil {
		err = e
		return
	}
	//查询关联角色规则
	groupPolicy := enforcer.GetFilteredGroupingPolicy(0, gconv.String(userId))
	if len(groupPolicy) > 0 {
		roleIds = make([]uint, len(groupPolicy))
		//得到角色id的切片
		for k, v := range groupPolicy {
			roleIds[k] = gconv.Uint(v[1])
		}
	}
	return
}

// 获取角色
func (s *roleImpl) getRolesByUserId(ctx context.Context, userId uint64, allRoleList []*entity.SysRole) (roles []*entity.SysRole, err error) {
	var roleIds []uint
	roleIds, err = s.getRoleIdsByUserId(ctx, userId)
	if err != nil {
		return
	}
	roles = make([]*entity.SysRole, 0, len(allRoleList))
	for _, v := range allRoleList {
		for _, id := range roleIds {
			if id == v.Id {
				roles = append(roles, v)
			}
		}
		if len(roles) == len(roleIds) {
			break
		}
	}
	return
}

// 获取树状菜单
func (s *roleImpl) getMenusTree(menus []*model.UserMenus, pid uint) []*model.UserMenus {
	returnList := make([]*model.UserMenus, 0, len(menus))
	for _, menu := range menus {
		if menu.Pid == pid {
			menu.Children = s.getMenusTree(menus, menu.Id)
			returnList = append(returnList, menu)
		}
	}
	return returnList
}

// 获取多角色菜单
func (s *roleImpl) getRolesMenus(ctx context.Context) (menus []*model.UserMenus, err error) {
	// 获取所有开启的菜单
	allMenus, err := authRuleService.GetIsMenuList(ctx)
	if err != nil {
		return
	}
	menus = make([]*model.UserMenus, len(allMenus))
	for k, v := range allMenus {
		var menu *model.UserMenu
		menu = s.setMenuData(menu, v)
		menus[k] = &model.UserMenus{UserMenu: menu}
	}
	menus = s.getMenusTree(menus, 0)
	return
}

// 添加角色
func (s *roleImpl) addRole(ctx context.Context, roleIds []int64, userId int64) (err error) {
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

// 修改角色
func (s *roleImpl) updRole(ctx context.Context, roleIds []int64, userId int64) (err error) {
	err = g.Try(func() {
		enforcer, e := CasbinEnforcer(ctx)
		liberr.ErrIsNil(ctx, e)

		//删除用户旧角色信息
		enforcer.RemoveFilteredGroupingPolicy(0, gconv.String(userId))
		for _, v := range roleIds {
			_, err = enforcer.AddGroupingPolicy(gconv.String(userId), gconv.String(v))
			liberr.ErrIsNil(ctx, err)
		}
	})
	return
}

// 组装菜单数据
func (s *roleImpl) setMenuData(menu *model.UserMenu, entity *model.SysAuthRuleInfoRes) *model.UserMenu {
	menu = &model.UserMenu{
		Id:        entity.Id,
		Pid:       entity.Pid,
		Name:      gstr.CaseCamelLower(gstr.Replace(entity.Name, "/", "_")),
		Component: entity.Component,
		Path:      entity.Path,
		MenuMeta: &model.MenuMeta{
			Icon:        entity.Icon,
			Title:       entity.Title,
			IsLink:      "",
			IsHide:      entity.IsHide == 1,
			IsKeepAlive: entity.IsCached == 1,
			IsAffix:     entity.IsAffix == 1,
			IsIframe:    entity.IsIframe == 1,
		},
	}
	if menu.MenuMeta.IsIframe || entity.IsLink == 1 {
		menu.MenuMeta.IsLink = entity.LinkUrl
	}
	return menu
}

// 从数据库获取所有角色
func (s *roleImpl) getRolesFromDb(ctx context.Context) (value interface{}, err error) {
	err = g.Try(func() {
		var v []*entity.SysRole
		//从数据库获取
		err = dao.SysRole.Ctx(ctx).
			Order(dao.SysRole.Columns().ListOrder + " asc," + dao.SysRole.Columns().Id + " asc").
			Scan(&v)
		liberr.ErrIsNil(ctx, err, "获取角色数据失败")
		value = v
	})
	return
}

// 获取角色分页列表
func (s *roleImpl) GetPageList(ctx context.Context, req *apiv1.RoleListReq) (res *apiv1.RoleListRes, err error) {
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
		if req.PageIndex == 0 {
			req.PageIndex = 1
		}
		res.CurrentPage = req.PageIndex
		if req.PageSize == 0 {
			req.PageSize = consts.PageSize
		}
		err = model.Page(res.CurrentPage, req.PageSize).Order("id asc").Scan(&res.List)
		liberr.ErrIsNil(ctx, err, "获取数据失败")
	})
	return
}

// 获取角色列表
func (s *roleImpl) GetRoleList(ctx context.Context) (list []*entity.SysRole, err error) {
	cache := Cache()
	//从缓存获取
	iList := cache.GetOrSetFuncLock(ctx, consts.CacheSysRole, s.getRolesFromDb, 0, consts.CacheSysAuthTag)
	if iList != nil {
		err = gconv.Struct(iList, &list)
	}
	return
}

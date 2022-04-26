package service

import (
	"context"

	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jiayg/liar/internal/consts"
	"github.com/jiayg/liar/internal/model"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/library/liberr"
)

type IRule interface {
	GetIsMenuList(ctx context.Context) ([]*model.SysAuthRuleInfoRes, error)
	GetMenuList(ctx context.Context) (list []*model.SysAuthRuleInfoRes, err error)
	// GetIsButtonList(ctx context.Context) ([]*model.SysAuthRuleInfoRes, error)
	// Add(ctx context.Context, req *system.RuleAddReq) (err error)
	// Get(ctx context.Context, id uint) (rule *entity.SysAuthRule, err error)
	// GetMenuRoles(ctx context.Context, id uint) (roleIds []uint, err error)
	// Update(ctx context.Context, req *system.RuleUpdateReq) (err error)
	// GetMenuListSearch(ctx context.Context, req *system.RuleSearchReq) (res []*model.SysAuthRuleInfoRes, err error)
	// GetMenuListTree(pid uint, list []*model.SysAuthRuleInfoRes) []*model.SysAuthRuleTreeRes
	// DeleteMenuByIds(ctx context.Context, ids []int) (err error)
}

type authRuleImpl struct {
}

var authRuleService = authRuleImpl{}

func AuthRule() IRule {
	return IRule(&authRuleService)
}

// GetIsMenuList 获取isMenu=0|1
func (s *authRuleImpl) GetIsMenuList(ctx context.Context) ([]*model.SysAuthRuleInfoRes, error) {
	list, err := s.GetMenuList(ctx)
	if err != nil {
		return nil, err
	}
	var gList = make([]*model.SysAuthRuleInfoRes, 0, len(list))
	for _, v := range list {
		if v.MenuType == 0 || v.MenuType == 1 {
			gList = append(gList, v)
		}
	}
	return gList, nil
}

// GetMenuList 获取所有菜单
func (s *authRuleImpl) GetMenuList(ctx context.Context) (list []*model.SysAuthRuleInfoRes, err error) {
	cache := Cache()
	//从缓存获取
	iList := cache.GetOrSetFuncLock(ctx, consts.CacheSysAuthMenu, s.getMenuListFromDb, 0, consts.CacheSysAuthTag)
	if iList != nil {
		err = gconv.Struct(iList, &list)
		liberr.ErrIsNil(ctx, err)
	}
	return
}

// 从数据库获取所有菜单
func (s *authRuleImpl) getMenuListFromDb(ctx context.Context) (value interface{}, err error) {
	err = g.Try(func() {
		var v []*model.SysAuthRuleInfoRes
		//从数据库获取
		err = dao.SysAuthRule.Ctx(ctx).
			Fields(model.SysAuthRuleInfoRes{}).Order("weigh desc,id asc").Scan(&v)
		liberr.ErrIsNil(ctx, err, "获取菜单数据失败")
		value = v
	})
	return
}

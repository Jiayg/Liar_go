package service

import (
	"context"

	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jiayg/liar/internal/consts"
	"github.com/jiayg/liar/internal/model/entity"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/library/liberr"
)

type IDept interface {
	// GetList(ctx context.Context, req *system.DeptSearchReq) (list []*entity.SysDept, err error)
	// Add(ctx context.Context, req *system.DeptAddReq) (err error)
	// Edit(ctx context.Context, req *system.DeptEditReq) (err error)
	GetCacheDepts(ctx context.Context) (list []*entity.SysDept, err error)
	// Delete(ctx context.Context, id int64) (err error)
	// GetListTree(pid int64, list []*entity.SysDept) (deptTree []*model.SysDeptTreeRes)
	FindSonByParentId(deptList []*entity.SysDept, deptId int64) []*entity.SysDept
}

var deptService = deptImpl{}

func Dept() IDept {
	return IDept(&deptService)
}

type deptImpl struct{}

// 获取部门ids
func (s *deptImpl) getDeptIdsById(ctx context.Context, deptId int64) (deptIds []int64, err error) {
	err = g.Try(func() {
		deptAll, e := s.GetCacheDepts(ctx)
		liberr.ErrIsNil(ctx, e)
		deptWithChildren := s.FindSonByParentId(deptAll, gconv.Int64(deptId))
		deptIds = make([]int64, len(deptWithChildren))
		for k, v := range deptWithChildren {
			deptIds[k] = v.DeptId
		}
		deptIds = append(deptIds, deptId)
	})
	return
}

// 从缓存获取部门列表
func (s *deptImpl) GetCacheDepts(ctx context.Context) (list []*entity.SysDept, err error) {
	err = g.Try(func() {
		cache := Cache()
		//从缓存获取
		iList := cache.GetOrSetFuncLock(ctx, consts.CacheSysDept, func(ctx context.Context) (value interface{}, err error) {
			err = dao.SysDept.Ctx(ctx).Scan(&list)
			liberr.ErrIsNil(ctx, err, "获取部门列表失败")
			value = list
			return
		}, 0, consts.CacheSysAuthTag)
		if iList != nil {
			err = gconv.Struct(iList, &list)
			liberr.ErrIsNil(ctx, err)
		}
	})
	return
}

func (s *deptImpl) FindSonByParentId(deptList []*entity.SysDept, deptId int64) []*entity.SysDept {
	children := make([]*entity.SysDept, 0, len(deptList))
	for _, v := range deptList {
		if v.ParentId == deptId {
			children = append(children, v)
			fChildren := s.FindSonByParentId(deptList, v.DeptId)
			children = append(children, fChildren...)
		}
	}
	return children
}

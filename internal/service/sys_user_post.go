package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jiayg/liar/internal/model/entity"
	"github.com/jiayg/liar/internal/service/internal/dao"
	"github.com/jiayg/liar/library/liberr"
)

type IUserPost interface {
}

type userPostImpl struct{}

var userPostService = userPostImpl{}

func UserPost() IUserPost {
	return IUserPost(&roleService)
}

// 获取岗位ids
func (s *userPostImpl) getUserPostIdsByUserId(ctx context.Context, userId uint64) (postIds []int64, err error) {
	err = g.Try(func() {
		var list []*entity.SysUserPost
		err = dao.SysUserPost.Ctx(ctx).Where(dao.SysUserPost.Columns().UserId, userId).Scan(&list)
		liberr.ErrIsNil(ctx, err, "获取用户岗位信息失败")
		postIds = make([]int64, 0)
		for _, entity := range list {
			postIds = append(postIds, entity.PostId)
		}
	})
	return
}

// 添加岗位
func (s *userPostImpl) addPost(ctx context.Context, tx *gdb.TX, postIds []int64, userId int64) (err error) {
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

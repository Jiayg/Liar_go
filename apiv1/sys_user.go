package apiv1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jiayg/liar/internal/model"
	"github.com/jiayg/liar/internal/model/entity"
)

type BaseUserReq struct {
	DeptId   uint64  `p:"deptId" v:"required#用户部门不能为空"` //所属部门
	Email    string  `p:"email" v:"email#邮箱格式错误"`       //邮箱
	NickName string  `p:"nickName" v:"required#用户昵称不能为空"`
	Mobile   string  `p:"mobile" v:"required|phone#手机号不能为空|手机号格式错误"`
	PostIds  []int64 `p:"postIds"`
	Remark   string  `p:"remark"`
	RoleIds  []int64 `p:"roleIds"`
	Sex      int     `p:"sex"`
	Status   uint    `p:"status"`
	IsAdmin  int     `p:"isAdmin"` // 是否后台管理员 1 是  0   否
}

type UserAddReq struct {
	g.Meta `path:"/user" tags:"用户管理" method:"post" summary:"添加用户"`
	*BaseUserReq
	UserName string `p:"userName" v:"required#用户账号不能为空"`
	Password string `p:"password" v:"required|password#密码不能为空|密码以字母开头，只能包含字母、数字和下划线，长度在6~18之间"`
	UserSalt string
}
type UserAddRes struct{}

type UserDeleteReq struct {
	g.Meta `path:"/user" tags:"用户管理" method:"delete" summary:"删除用户"`
	Ids    []int `p:"ids"`
}
type UserDeleteRes struct{}

type UserUpdateReq struct {
	g.Meta `path:"/user" tags:"用户管理" method:"put" summary:"修改用户"`
	*BaseUserReq
	UserId int64 `p:"userId" v:"required#用户id不能为空"`
}
type UserUpdateRes struct{}

type UserGetReq struct {
	g.Meta `path:"/user" tags:"用户管理" method:"get" summary:"获取用户信息"`
	Id     uint64 `p:"id"`
}
type UserGetRes struct {
	g.Meta         `mime:"application/json"`
	User           *entity.SysUser `json:"user"`
	CheckedRoleIds []uint          `json:"checkedRoleIds"`
	CheckedPostIds []int64         `json:"CheckedPostIds"`
}

type UserSearchReq struct {
	g.Meta   `path:"/user" tags:"用户管理" method:"get" summary:"用户分页列表"`
	DeptId   string `p:"deptId"` //部门id
	Mobile   string `p:"mobile"`
	Status   string `p:"status"`
	KeyWords string `p:"keyWords"`
	PageReq
}
type UserSearchRes struct {
	g.Meta   `mime:"application/json"`
	UserList []*model.SysUserRoleDeptRes `json:"userList"`
	ListRes
}

type UserStatusReq struct {
	g.Meta     `path:"/user/changeStatus" tags:"用户管理" method:"put" summary:"修改用户状态"`
	Id         uint64 `p:"userId" v:"required#用户id不能为空"`
	UserStatus uint   `p:"status" v:"required#用户状态不能为空"`
}
type UserStatusRes struct{}

type UserResetPwdReq struct {
	g.Meta   `path:"/user/resetPwd" tags:"用户管理" method:"put" summary:"重置用户密码"`
	Id       uint64 `p:"userId" v:"required#用户id不能为空"`
	Password string `p:"password" v:"required|password#密码不能为空|密码以字母开头，只能包含字母、数字和下划线，长度在6~18之间"`
}
type UserResetPwdRes struct{}

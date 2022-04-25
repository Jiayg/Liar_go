package apiv1

import "github.com/gogf/gf/v2/frame/g"

type SetUserReq struct {
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
	g.Meta `path:"/user/add" tags:"用户管理" method:"post" summary:"添加用户"`
	*SetUserReq
	UserName string `p:"userName" v:"required#用户账号不能为空"`
	Password string `p:"password" v:"required|password#密码不能为空|密码以字母开头，只能包含字母、数字和下划线，长度在6~18之间"`
	UserSalt string
}

type UserAddRes struct {
}

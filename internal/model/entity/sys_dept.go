// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT. Created at 2022-04-24 22:31:11
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDept is the golang structure for table sys_dept.
type SysDept struct {
	DeptId    int64       `json:"deptId"    description:"部门id"`
	ParentId  int64       `json:"parentId"  description:"父部门id"`
	Ancestors string      `json:"ancestors" description:"祖级列表"`
	DeptName  string      `json:"deptName"  description:"部门名称"`
	OrderNum  int         `json:"orderNum"  description:"显示顺序"`
	Leader    string      `json:"leader"    description:"负责人"`
	Phone     string      `json:"phone"     description:"联系电话"`
	Email     string      `json:"email"     description:"邮箱"`
	Status    uint        `json:"status"    description:"部门状态（0正常 1停用）"`
	CreatedBy uint64      `json:"createdBy" description:"创建人"`
	UpdatedBy int64       `json:"updatedBy" description:"修改人"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" description:"修改时间"`
	DeletedAt *gtime.Time `json:"deletedAt" description:"删除时间"`
}
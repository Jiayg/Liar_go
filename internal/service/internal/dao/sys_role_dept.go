// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/jiayg/liar/internal/service/internal/dao/internal"
)

// sysRoleDeptDao is the data access object for table sys_role_dept.
// You can define custom methods on it to extend its functionality as you wish.
type sysRoleDeptDao struct {
	*internal.SysRoleDeptDao
}

var (
	// SysRoleDept is globally public accessible object for table sys_role_dept operations.
	SysRoleDept = sysRoleDeptDao{
		internal.NewSysRoleDeptDao(),
	}
)

// Fill with you ideas below.

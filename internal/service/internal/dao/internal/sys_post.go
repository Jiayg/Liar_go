// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT. Created at 2022-04-24 22:31:11
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysPostDao is the data access object for table sys_post.
type SysPostDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns SysPostColumns // columns contains all the column names of Table for convenient usage.
}

// SysPostColumns defines and stores column names for table sys_post.
type SysPostColumns struct {
	PostId    string // 岗位ID
	PostCode  string // 岗位编码
	PostName  string // 岗位名称
	PostSort  string // 显示顺序
	Status    string // 状态（0正常 1停用）
	Remark    string // 备注
	CreatedBy string // 创建人
	UpdatedBy string // 修改人
	CreatedAt string // 创建时间
	UpdatedAt string // 修改时间
	DeletedAt string // 删除时间
}

//  sysPostColumns holds the columns for table sys_post.
var sysPostColumns = SysPostColumns{
	PostId:    "post_id",
	PostCode:  "post_code",
	PostName:  "post_name",
	PostSort:  "post_sort",
	Status:    "status",
	Remark:    "remark",
	CreatedBy: "created_by",
	UpdatedBy: "updated_by",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewSysPostDao creates and returns a new DAO object for table data access.
func NewSysPostDao() *SysPostDao {
	return &SysPostDao{
		group:   "default",
		table:   "sys_post",
		columns: sysPostColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysPostDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SysPostDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SysPostDao) Columns() SysPostColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SysPostDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysPostDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysPostDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

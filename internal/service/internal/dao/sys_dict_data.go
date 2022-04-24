// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/jiayg/liar/internal/service/internal/dao/internal"
)

// sysDictDataDao is the data access object for table sys_dict_data.
// You can define custom methods on it to extend its functionality as you wish.
type sysDictDataDao struct {
	*internal.SysDictDataDao
}

var (
	// SysDictData is globally public accessible object for table sys_dict_data operations.
	SysDictData = sysDictDataDao{
		internal.NewSysDictDataDao(),
	}
)

// Fill with you ideas below.

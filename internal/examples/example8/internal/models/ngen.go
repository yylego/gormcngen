// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: ngen_test.go:32 -> models.TestGenerate
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yylego/gormcnm"

func (T *Account) TableColumns(decoration gormcnm.ColumnNameDecoration) *AccountColumns {
	return &AccountColumns{
		// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。
		V主键: gormcnm.Cmn(T.ID, "id", decoration),
		V名字: gormcnm.Cmn(T.Name, "name", decoration),
	}
}

type AccountColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	V主键 gormcnm.ColumnName[uint]
	V名字 gormcnm.ColumnName[string]
}

func (T *Purchase) TableColumns(decoration gormcnm.ColumnNameDecoration) *PurchaseColumns {
	return &PurchaseColumns{
		// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。
		V订单主键: gormcnm.Cmn(T.ID, "id", decoration),
		V用户主键: gormcnm.Cmn(T.AccountID, "account_id", decoration),
		V订单金额: gormcnm.Cmn(T.Amount, "amount", decoration),
	}
}

type PurchaseColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	V订单主键 gormcnm.ColumnName[uint]
	V用户主键 gormcnm.ColumnName[uint]
	V订单金额 gormcnm.ColumnName[float64]
}

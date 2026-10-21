// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: ngen_test.go:28 -> models.TestGenerate
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yylego/gormcnm"

func (c *Student) Columns() *StudentColumns {
	return &StudentColumns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
		V班级编码: gormcnm.Cnm(c.ClassCode, "class_code"),
		V学号:   gormcnm.Cnm(c.Num, "num"),
		V名字:   gormcnm.Cnm(c.Name, "name"),
		V科类:   gormcnm.Cnm(c.Track, "track"),
		V生日:   gormcnm.Cnm(c.BornDate, "age"),
	}
}

type StudentColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	V班级编码 gormcnm.ColumnName[string]
	V学号   gormcnm.ColumnName[int]
	V名字   gormcnm.ColumnName[string]
	V科类   gormcnm.ColumnName[TrackType]
	V生日   gormcnm.ColumnName[string]
}

func (c *Class) Columns() *ClassColumns {
	return &ClassColumns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
		V班级编码: gormcnm.Cnm(c.V班级编码, "class_code"),
		V班级名称: gormcnm.Cnm(c.V班级名称, "class_name"),
		V班主任名: gormcnm.Cnm(c.V班主任名, "class_head_name"),
	}
}

type ClassColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	V班级编码 gormcnm.ColumnName[string]
	V班级名称 gormcnm.ColumnName[string]
	V班主任名 gormcnm.ColumnName[string]
}

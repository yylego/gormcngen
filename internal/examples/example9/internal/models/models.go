// Package models demonstrates GORM associations through *T fields
// Auto contains Account and Profile structs with foreignKey relationship
// Used to showcase column generation excluding association fields
//
// models 包演示带有指针字段的 GORM 关联关系
// 自动包含带有 foreignKey 关系的 Account 和 Profile 结构体
// 用于展示排除关联字段的列生成
package models

type Account struct {
	ID      uint
	Name    string
	Profile *Profile `gorm:"foreignKey:AccountID"`
}

type Profile struct {
	ID        uint
	Bio       string
	AccountID uint // 外键
}

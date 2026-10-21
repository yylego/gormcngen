// Package models demonstrates Chinese cnm tag naming with TableName method
// Auto contains Account and Purchase structs with Chinese field aliases
// Used to showcase namespace pattern with custom table names
//
// models 包演示带有 TableName 方法的中文 cnm 标签命名
// 自动包含带有中文字段别名的 Account 和 Purchase 结构体
// 用于展示带有自定义表名的命名空间模式
package models

type Account struct {
	ID   uint   `cnm:"V主键"`
	Name string `cnm:"V名字"`
}

func (*Account) TableName() string {
	return "accounts"
}

type Purchase struct {
	ID        uint    `cnm:"V订单主键"`
	AccountID uint    `cnm:"V用户主键"`
	Amount    float64 `cnm:"V订单金额"`
}

func (*Purchase) TableName() string {
	return "purchases"
}

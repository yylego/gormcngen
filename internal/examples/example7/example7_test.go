// Package example7 demonstrates type-safe join queries with Chinese column names
// Shows comparison between raw SQL and column-based approaches
// Used to showcase table decoration and complex join operations
//
// example7 演示带有中文列名的类型安全联接查询
// 展示原生 SQL 和基于列的查询方法的对比
// 用于展示表装饰和复杂联接操作
package example7

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcngen/internal/examples/example7/internal/models"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormcnm/gormcnmstub"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Shared database instance in the advanced column operation tests
// 用于高级列操作测试的全局数据库实例
var caseDB *gorm.DB

// TestMain sets up the database with accounts and purchases to test join queries
// Creates test data to demonstrate type-safe joins and complex queries
//
// TestMain 为联接查询测试设置用户和订单数据库
// 创建测试数据来演示类型安全的联接和复杂查询
func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Account{}, &models.Purchase{}))

	const accountCount = 10
	accounts := make([]*models.Account, 0, accountCount)
	for idx := 0; idx < accountCount; idx++ {
		accounts = append(accounts, &models.Account{
			ID:   0,
			Name: "name" + strconv.Itoa(idx+1),
		})
	}
	done.Done(db.Create(&accounts).Error)

	const purchaseCount = 20
	purchases := make([]*models.Purchase, 0, purchaseCount)
	for idx := 0; idx < purchaseCount; idx++ {
		accountID := accounts[rand.IntN(len(accounts))].ID

		purchases = append(purchases, &models.Purchase{
			ID:        0,
			AccountID: accountID,
			Amount:    float64(rand.IntN(1000)) + rand.Float64(),
		})
	}
	done.Done(db.Create(&purchases).Error)

	caseDB = db
	m.Run()
}

// TestExample demonstrates comparison between raw SQL and type-safe column queries
// Compares native string-based queries with generated column-based queries
//
// TestExample 演示原生 SQL 和类型安全列查询的对比
// 比较传统基于字符串的查询和生成的基于列的查询
func TestExample(t *testing.T) {
	// Execute native and type-safe queries
	// 执行传统查询和类型安全查询
	expected0Text := neatjsons.S(selectFunc0(t, caseDB))
	expected1Text := neatjsons.S(selectFunc1(t, caseDB))
	// Ensure both results are the same
	// 确保两者结果相同
	require.Equal(t, expected0Text, expected1Text)
}

// selectFunc0 demonstrates native SQL statement approach
// Uses string literals as table names, column names, and joins
//
// selectFunc0 演示传统的原生 SQL 查询方法
// 使用字符串字面量表示表名、列名和联接
func selectFunc0(t *testing.T, db *gorm.DB) []*AccountPurchase {
	var results []*AccountPurchase
	require.NoError(t, db.Table("accounts").
		Select("accounts.id as account_id, accounts.name as account_name, purchases.id as purchase_id, purchases.amount as purchase_amount").
		Joins("left join purchases on purchases.account_id = accounts.id").
		Order("accounts.id asc, purchases.id asc").
		Scan(&results).Error)
	t.Log(neatjsons.S(results))
	return results
}

type AccountPurchase struct {
	AccountID      uint
	AccountName    string
	PurchaseID     uint
	PurchaseAmount float64
}

// selectFunc1 demonstrates type-safe column-based approach
// Uses generated column structs with Chinese column names and type protection
//
// selectFunc1 演示类型安全的基于列的查询方法
// 使用生成的列结构体和中文列名及类型安全
func selectFunc1(t *testing.T, db *gorm.DB) []*AccountPurchase {
	// Create struct instances and get the table-decorated columns
	// 创建模型实例并获取带表装饰的列
	account := &models.Account{}
	accountColumns := account.TableColumns(gormcnm.NewTableDecoration(account.TableName()))
	purchase := &models.Purchase{}
	purchaseColumns := purchase.TableColumns(gormcnm.NewTableDecoration(purchase.TableName()))

	// Create result struct as the column name mapping
	// 为列名映射创建结果结构体
	accountPurchase := &AccountPurchase{}

	// Build the type-safe statement with Chinese column names
	// 使用中文列名构建类型安全查询
	var results []*AccountPurchase
	require.NoError(t, db.Table(account.TableName()).
		Select(gormcnmstub.MergeStmts(
			accountColumns.V主键.AsName(gormcnm.Cnm(accountPurchase.AccountID, "account_id")),
			accountColumns.V名字.AsName(gormcnm.Cnm(accountPurchase.AccountName, "account_name")),
			purchaseColumns.V订单主键.AsName(gormcnm.Cnm(accountPurchase.PurchaseID, "purchase_id")),
			purchaseColumns.V订单金额.AsName(gormcnm.Cnm(accountPurchase.PurchaseAmount, "purchase_amount")),
		)).
		Joins(accountColumns.LEFTJOIN(purchase.TableName()).On(purchaseColumns.V用户主键.OnEq(accountColumns.V主键))).
		Order(accountColumns.V主键.Ob("asc").Ob(purchaseColumns.V订单主键.Ob("asc")).Ox()).
		Scan(&results).Error)
	t.Log(neatjsons.S(results))
	return results
}

// Package example8 demonstrates extended column generation with TableName method
// Shows enhanced matching validation and advanced configuration options
// Used to showcase namespace pattern with custom table names
//
// example8 演示带有 TableName 方法的扩展列生成
// 展示增强查询一致性验证和高级配置选项
// 用于展示带有自定义表名的命名空间模式
package example8

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcngen/internal/examples/example8/internal/models"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormcnm/gormcnmstub"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Shared database instance in the extended column operation tests
// 用于扩展列操作测试的全局数据库实例
var caseDB *gorm.DB

// TestMain sets up database with advanced column generation configuration
// Creates test data to exercise extended column features and search operations
//
// TestMain 为高级列生成配置设置数据库
// 创建测试数据来测试扩展列特性和查询操作
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

// TestExample validates matching between native and enhanced column queries
// Ensures both approaches produce same results with extended features
//
// TestExample 验证传统和增强列查询的一致性
// 确保两种查询方法在扩展功能下产生相同的结果
func TestExample(t *testing.T) {
	// Execute both search methods and compare results
	// 执行两种查询方法并比较结果
	expected0Text := neatjsons.S(selectFunc0(t, caseDB))
	expected1Text := neatjsons.S(selectFunc1(t, caseDB))
	// Confirm results are the same with enhanced configurations
	// 验证在增强配置下结果相同
	require.Equal(t, expected0Text, expected1Text)
}

// selectFunc0 demonstrates native SQL statement method
// Baseline implementation using string literals as the reference
//
// selectFunc0 演示常规的原生 SQL 查询方法
// 使用字符串字面量的基准实现用于比较
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

// selectFunc1 demonstrates the enhanced type-safe column method
// Advanced implementation with extended column generation features
//
// selectFunc1 演示增强的类型安全列查询方法
// 带有扩展列生成特性的高级实现
func selectFunc1(t *testing.T, db *gorm.DB) []*AccountPurchase {
	// Initialize the structs and get the decorated columns
	// 初始化模型并获取其装饰列
	account := &models.Account{}
	accountColumns := account.TableColumns(gormcnm.NewTableDecoration(account.TableName()))
	purchase := &models.Purchase{}
	purchaseColumns := purchase.TableColumns(gormcnm.NewTableDecoration(purchase.TableName()))

	// Define result mapping structure
	// 定义结果映射结构
	accountPurchase := &AccountPurchase{}

	// Build the enhanced statement with extended column features
	// 使用扩展列特性构建增强查询
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

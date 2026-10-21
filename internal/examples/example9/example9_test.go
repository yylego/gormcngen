// Package example9 demonstrates association field handling in column generation
// Shows how association fields are skipped while maintaining relationship queries
// Used to showcase GORM associations through *T fields and foreignKey tags
//
// example9 演示列生成中的关联字段处理
// 展示关联字段如何被跳过同时保持关系查询正常工作
// 用于展示带有指针字段和 foreignKey 标签的 GORM 关联
package example9

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcngen/internal/examples/example9/internal/models"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Shared database instance in the association field tests
// 用于关联字段测试的全局数据库实例
var caseDB *gorm.DB

// TestMain sets up database with association structs
// Creates test data to exercise association field handling and relationship queries
//
// TestMain 为关联模型设置数据库
// 创建测试数据来测试关联字段处理和关系查询
func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	// Auto migrate both structs along with the association
	// 自动迁移两个模型及其关联关系
	done.Done(db.AutoMigrate(&models.Account{}, &models.Profile{}))

	// Create test accounts
	// 创建测试账户
	const accountCount = 5
	accounts := make([]*models.Account, 0, accountCount)
	for idx := 0; idx < accountCount; idx++ {
		accounts = append(accounts, &models.Account{
			Name: "Account" + strconv.Itoa(idx+1),
		})
	}
	done.Done(db.Create(&accounts).Error)

	// Create a profile on each account (has_one relationship)
	// 为每个账户创建配置文件（has_one关系）
	profiles := make([]*models.Profile, 0, accountCount)
	for idx, account := range accounts {
		profiles = append(profiles, &models.Profile{
			Bio:       "Bio of Account" + strconv.Itoa(idx+1),
			AccountID: account.ID,
		})
	}
	done.Done(db.Create(&profiles).Error)

	caseDB = db
	m.Run()
}

// TestExample9 demonstrates association field handling in column generation
// Validates that association fields are skipped as expected while maintaining correct column mappings
//
// TestExample9 演示列生成中的关联字段处理
// 验证关联字段被正确跳过同时保持正确的列映射
func TestExample9(t *testing.T) {
	// Test Account columns - must NOT include the Profile field
	// 测试Account模型列 - 不应包含Profile字段
	account := &models.Account{}
	accountColumns := account.Columns()
	t.Log("Account columns:", neatjsons.S(accountColumns))

	// Confirm Account columns contain the expected fields but not the association field
	// 验证Account列包含预期字段但不包含关联字段
	accountColumnsJSON := neatjsons.S(accountColumns)
	require.Contains(t, accountColumnsJSON, `"ID"`)
	require.Contains(t, accountColumnsJSON, `"id"`)
	require.Contains(t, accountColumnsJSON, `"Name"`)
	require.Contains(t, accountColumnsJSON, `"name"`)
	// Profile field should NOT show in column mappings
	// Profile字段不应出现在列映射中
	require.NotContains(t, accountColumnsJSON, "Profile")

	// Test Profile columns
	// 测试Profile模型列
	profile := &models.Profile{}
	profileColumns := profile.Columns()
	t.Log("Profile columns:", neatjsons.S(profileColumns))

	// Test a basic search on Account
	// 测试Account模型的基础查询
	testBasicAccountSearch(t, caseDB)

	// Test an association search (Account with Profile)
	// 测试关联查询（Account带Profile）
	testAssociationSearch(t, caseDB)
}

// testBasicAccountSearch tests basic Account searches using generated columns
// Demonstrates that association fields don't interfere with standard column operations
//
// testBasicAccountSearch 使用生成的列测试基础Account模型查询
// 演示关联字段不会干扰正常的列操作
func testBasicAccountSearch(t *testing.T, db *gorm.DB) {
	var accounts []*models.Account
	require.NoError(t, db.Find(&accounts).Error)
	require.Greater(t, len(accounts), 0)
	t.Log("Found accounts:", neatjsons.S(accounts))

	// Confirm we can use the generated column name to run queries
	// 验证我们可以使用生成的列名进行查询
	account := &models.Account{}
	accountColumns := account.Columns()

	var foundAccount models.Account
	require.NoError(t, db.Where(accountColumns.Name.Eq("Account1")).First(&foundAccount).Error)
	require.Equal(t, "Account1", foundAccount.Name)
	t.Log("Found account via name:", neatjsons.S(foundAccount))

	// Confirm the column name mapping is correct
	// 验证列名映射是否正确
	require.Equal(t, "name", string(accountColumns.Name))
	require.Equal(t, "id", string(accountColumns.ID))
}

// testAssociationSearch tests GORM association searches
// Demonstrates that association relationships work as expected even though fields are skipped in column generation
//
// testAssociationSearch 测试GORM关联查询
// 演示即使字段在列生成中被跳过，关联关系仍然正常工作
func testAssociationSearch(t *testing.T, db *gorm.DB) {
	// Load Account with Profile preloaded
	// 查询Account并预加载Profile
	var accounts []*models.Account
	require.NoError(t, db.Preload("Profile").Find(&accounts).Error)
	require.Greater(t, len(accounts), 0)

	// Confirm association loading succeeded
	// 验证关联加载正常工作
	for _, account := range accounts {
		if account.Profile != nil {
			require.NotEmpty(t, account.Profile.Bio)
			require.Equal(t, account.ID, account.Profile.AccountID)
			t.Log("Account with profile:", neatjsons.S(map[string]interface{}{
				"AccountName": account.Name,
				"ProfileBio":  account.Profile.Bio,
			}))
		}
	}
}

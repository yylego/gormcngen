// Package gormcngen tests validate position tracking utilities
// Auto verifies GetGenPosFuncMark function with runtime invocation info
//
// gormcngen 测试包验证位置追踪工具
// 自动验证 GetGenPosFuncMark 函数与运行时调用者信息
package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGenPosFuncMark(t *testing.T) {
	result := GetGenPosFuncMark(0)
	require.NotEmpty(t, result)
	t.Log("result:", result)
}

func TestPackageNameOfPath(t *testing.T) {
	for pkgPath, want := range map[string]string{
		"time":                       "time",
		"database/sql":               "sql",
		"gorm.io/datatypes":          "datatypes",
		"gopkg.in/yaml.v3":           "yaml",
		"github.com/a/b/v2":          "b",
		"github.com/yylego/gormcnm":  "gormcnm",
		"github.com/a/version.vnext": "version.vnext",
	} {
		require.Equal(t, want, PackageNameOfPath(pkgPath), pkgPath)
	}
}

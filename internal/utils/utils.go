// Package utils: support functions behind gormcngen code generation
// Provides string manipulation and exportable-status conversion functions
// Supports case conversion and status toggles on generated names
//
// utils: gormcngen 代码生成的内部工具函数
// 提供字符串操作和导出性转换的辅助函数
// 支持生成标识符的大小写转换和可见性切换
package utils

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/yylego/rese"
)

// ConvertToUnexportable converts the first rune of the string to lowercase
// Transforms exported names to unexported ones, used inside the package
// Returns the string with the first rune converted to lowercase
//
// ConvertToUnexportable 将字符串的第一个字符转换为小写
// 将导出标识符转换为未导出标识符供内部使用
// 返回第一个 rune 转换为小写的字符串
func ConvertToUnexportable(s string) string {
	runes := []rune(s)
	if len(runes) > 0 {
		runes[0] = unicode.ToLower(runes[0])
	}
	return string(runes)
}

// IsExportable checks if a name is exportable (starts with uppercase)
// Returns true if the first rune is uppercase, else false
// Used to determine status of generated structs and methods
//
// IsExportable 检查标识符是否可导出（以大写字母开头）
// 如果第一个字符是大写返回 true，否则返回 false
// 用于确定生成的结构体和方法的可见性
func IsExportable(name string) bool {
	return unicode.IsUpper(([]rune(name))[0])
}

// SwitchToggleExportable toggles the exportable status of a name
// Converts an uppercase first rune to lowercase and vice versa
// Used to find alternative struct names during AST analysis
//
// SwitchToggleExportable 切换标识符的可导出性
// 将大写第一个字符转换为小写，反之亦然
// 用于在 AST 分析过程中查找替代结构体名称
func SwitchToggleExportable(name string) string {
	runes := []rune(name)
	if unicode.IsUpper(runes[0]) {
		runes[0] = unicode.ToLower(runes[0])
	} else if unicode.IsLower(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

// GetGenPosFuncMark gets position information used in code generation tracing
// Returns a formatted string showing the source file and function that triggered generation
// Uses runtime.Caller to walk up the stack using the specified frame count
// Supports URL-encoded file paths
//
// GetGenPosFuncMark 获取调用者位置信息，用于代码生成的可追溯性
// 返回格式化字符串，显示触发生成的源文件和函数
// 使用 runtime.Caller 按指定帧数向上遍历调用栈
// 支持URL编码的文件路径
func GetGenPosFuncMark(skip int) string {
	// Get invocation position info using runtime package
	// 使用 runtime 包获取调用者位置信息
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}

	// Extract just the filename, not the complete path, to produce clean output
	// 为了可读性，只提取文件名而不是完整路径
	filename := filepath.Base(file)

	// Decode URL-encoded filename
	// 解码URL编码的文件名
	filename = rese.C1(url.PathUnescape(filename))

	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		return filename
	}

	// Extract just the function name, not the complete package path
	// 只提取函数名，不包含完整的包路径
	funcName := filepath.Base(funcInfo.Name())
	return fmt.Sprintf("%s:%d -> %s", filename, line, funcName)
}

// moduleVersion matches a module version element such as "v2"
// moduleVersion 匹配 "v2" 这类模块版本段
var moduleVersion = regexp.MustCompile(`^v[0-9]+$`)

// PackageNameOfPath guesses the package name of an import path from its last element.
// Both a "/vN" element and a ".vN" suffix mark a module version and get skipped, e.g. "gopkg.in/yaml.v3" -> "yaml", "github.com/a/b/v2" -> "b".
//
// PackageNameOfPath 根据导入路径的最后一段推断包名。
// "/vN" 段或 ".vN" 后缀表示模块版本，会被跳过，如 "gopkg.in/yaml.v3" -> "yaml"、"github.com/a/b/v2" -> "b"。
func PackageNameOfPath(pkgPath string) string {
	elems := strings.Split(pkgPath, "/")
	name := elems[len(elems)-1]
	if len(elems) > 1 && moduleVersion.MatchString(name) {
		name = elems[len(elems)-2]
	}
	if idx := strings.LastIndex(name, ".v"); idx > 0 && moduleVersion.MatchString(name[idx+1:]) {
		name = name[:idx]
	}
	return name
}

// Package gormcngen: Schema-based configuration and single-use code generation
// Handles single instance schema analysis and targeted code generation
// Provides precise management of column struct generation and method creation
//
// gormcngen: 基于 schema 的配置和单次使用代码生成
// 处理单个模式 schema 分析和针对性代码生成
// 提供对列结构体生成和方法创建的精确控制
package gormcngen

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/yylego/erero"
	"github.com/yylego/gormcngen/internal/utils"
	"github.com/yylego/gormcnm"
	"github.com/yylego/must"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"github.com/yylego/tern"
	"github.com/yylego/tern/zerotern"
	"github.com/yylego/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

// SchemaConfig holds configuration used in generating column methods and structures from GORM models
// Contains parsed schema information along with generation options and naming preferences
// Manages the complete lifecycle from schema analysis to code output generation
//
// SchemaConfig 保存从 GORM 模型生成列方法和结构体的配置
// 包含已解析的 schema 信息以及生成选项和命名偏好
// 管理从 schema 分析到代码输出生成的完整生命周期
type SchemaConfig struct {
	sch                    *schema.Schema // Parsed GORM schema from struct // 从模型结构体解析的 GORM schema
	structName             string         // Generated column struct name // 生成的列结构体名称
	methodName             string         // Generated Columns() method name // 生成的 Columns() 方法名称
	methodNameTableColumns string         // Generated TableColumns() method name // 生成的 TableColumns() 方法名称
	options                *Options       // Generation configuration options // 生成行为配置选项
}

// NewSchemaConfig creates a SchemaConfig instance from a GORM structure and generation options
// Parses the structure using GORM schema parsing and applies naming strategies
// Initializes schema analysis, shows debug information, and configures generation parameters
// Returns a configured SchemaConfig prepared to generate code
//
// NewSchemaConfig 从 GORM 模型和生成选项创建 SchemaConfig 实例
// 使用 GORM schema 解析器解析模型结构并应用命名策略
// 初始化 schema 分析，显示调试信息，并配置生成参数
// 返回一个完整配置的 SchemaConfig，准备进行代码生成
func NewSchemaConfig(object interface{}, options *Options) *SchemaConfig {
	sch := rese.P1(schema.Parse(object, &sync.Map{}, &schema.NamingStrategy{
		SingularTable: false, //这是gorm默认的
		NoLowerCase:   false, //这是gorm默认的
	}))

	ShowSchemaEnglish(sch)
	ShowSchemaChinese(sch)
	ShowSchemaRelationshipsEnglish(sch)
	ShowSchemaRelationshipsChinese(sch)

	namingConfig := NewNamingConfig(sch, options)

	return NewConfig(sch, namingConfig, options)
}

// NamingConfig holds naming conventions used in generated code elements
// Contains struct name and method names used in code generation output
//
// NamingConfig 保存生成代码元素的命名约定
// 包含代码生成输出中使用的结构体名称和方法名称
type NamingConfig struct {
	StructName             string // Generated column struct name // 生成的列结构体名称
	MethodNameColumns      string // Columns() method name // Columns() 方法名称
	MethodNameTableColumns string // TableColumns() method name // TableColumns() 方法名称
}

// NewNamingConfig creates naming configuration based on schema and options
// Generates struct name with appropriate exportable status
// Returns configured naming conventions used in code generation
//
// NewNamingConfig 根据 schema 和选项创建命名配置
// 生成具有适当导出性的结构体名称
// 返回代码生成中使用的命名约定
func NewNamingConfig(sch *schema.Schema, options *Options) *NamingConfig {
	const structNameSuffix = "Columns"
	structName := tern.BFF(options.columnClassExportable, func() string {
		return sch.Name + structNameSuffix // 通常定义的结构体名称是导出的
	}, func() string {
		return utils.ConvertToUnexportable(sch.Name) + structNameSuffix
	})

	namingConfig := &NamingConfig{
		StructName:             structName,
		MethodNameColumns:      "Columns",
		MethodNameTableColumns: "TableColumns",
	}
	return namingConfig
}

// Config Configuration of generating column methods and structures.
// Config 根据模型生成列方法和结构的配置。
type Config = SchemaConfig

// NewConfig Creates a new Config instance with the provided schema, struct name, method name, and options.
// NewConfig 创建一个新的 Config 实例，使用提供的 schema、结构体名称、方法名称和选项。
func NewConfig(sch *schema.Schema, namingConfig *NamingConfig, options *Options) *Config {
	return &Config{
		sch:                    sch,
		structName:             namingConfig.StructName,
		methodName:             namingConfig.MethodNameColumns,
		methodNameTableColumns: namingConfig.MethodNameTableColumns,
		options:                options,
	}
}

// GenOutput Structure representing the generated method and struct code with package imports.
// GenOutput 表示生成的方法和结构体代码，以及涉及的包导入信息。
type GenOutput struct {
	methodCode             string          // Generated Columns() method code // 生成的 Columns() 方法代码
	methodTableColumnsCode string          // Generated TableColumns() method code // 生成的 TableColumns() 方法代码
	structCode             string          // Generated structure code // 生成的结构体代码
	pkgImports             map[string]bool // Package imports needed // 生成的代码需要导入的包
}

func (x *GenOutput) GetMethodCode() string {
	return x.methodCode
}

func (x *GenOutput) GetMethodTableColumnsCode() string {
	return x.methodTableColumnsCode
}

func (x *GenOutput) GetStructCode() string {
	return x.structCode
}

func (x *GenOutput) GetPkgImports() map[string]bool {
	return x.pkgImports
}

// Generate Generates the column method and struct based on the configuration.
// Generate 根据配置生成列方法和结构。
func (c *Config) Generate() *GenOutput {
	return c.Gen()
}

// Gen Generates the column method and struct based on the configuration.
// Gen 根据配置生成列方法和结构。
func (c *Config) Gen() *GenOutput {
	structPtx := utils.NewPTX()
	structPtx.Println(fmt.Sprintf("type %s struct{", c.structName))

	methodPtx := utils.NewPTX()
	methodPtx.Println(fmt.Sprintf("func (%s*%s) %s() *%s {", c.options.columnsMethodRecvName, c.sch.Name, c.methodName, c.structName))

	methodTableColumnsPtx := utils.NewPTX()
	methodTableColumnsPtx.Println(fmt.Sprintf("func (%s*%s) %s(decoration gormcnm.ColumnNameDecoration) *%s {", c.options.columnsMethodRecvName, c.sch.Name, c.methodNameTableColumns, c.structName))

	operationClass := reflect.TypeOf(gormcnm.ColumnOperationClass{})
	pkgNameGormCnm := filepath.Base(operationClass.PkgPath())

	var pkgImports = map[string]bool{
		operationClass.PkgPath(): true,
	}

	if c.options.embedColumnOperations {
		structPtx.Println("\t", "// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。")
		structPtx.Println("\t", fmt.Sprintf("%s.%s", pkgNameGormCnm, operationClass.Name()))
	}
	structPtx.Println("\t", "// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。")

	if c.options.isGenFuncTableColumns {
		must.Nice(c.options.columnsMethodRecvName)
		must.True(c.options.columnsCheckFieldType)
		methodPtx.Println(fmt.Sprintf("	return %s.%s(gormcnm.NewPlainDecoration())", c.options.columnsMethodRecvName, c.methodNameTableColumns))
		methodTableColumnsPtx.Println(fmt.Sprintf("	return &%s{", c.structName))
		methodTableColumnsPtx.Println("\t// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。")
	} else {
		methodPtx.Println(fmt.Sprintf("	return &%s{", c.structName))
		methodPtx.Println("\t// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。")
		methodTableColumnsPtx.Println(fmt.Sprintf("	panic(\"METHOD %s.%s IS NOT IMPLEMENTED\")", c.structName, c.methodNameTableColumns))
	}
	for _, field := range c.sch.Fields {
		zaplog.LOG.Debug("processing field column mapping",
			zap.String("name", field.Name),
			zap.String("type", field.FieldType.String()),
			zap.String("dbName", field.DBName),
			zap.String("tag", string(field.Tag)))

		// Skip fields that don't have corresponding database columns (e.g., association fields)
		// 跳过没有对应数据库列的字段（如关联字段等）
		if field.DBName == "" {
			zaplog.LOG.Debug("skipping field without db-name",
				zap.String("name", field.Name),
				zap.String("reason", "without DBName"))
			continue
		}

		// Write the field type as Go source seen from the struct's own package, recording the imports it needs
		// 以模型所在包的视角写出字段类型，并记下它需要的导入
		columnGoTypeName := c.resolveTypeName(field.FieldType, pkgImports)
		structFieldName, ok := c.resolveNewFieldName(field)
		if !ok {
			continue // 某些场景下不需要获得列名
		}
		structPtx.Println("\t", structFieldName, fmt.Sprintf("%s.ColumnName[%s]", pkgNameGormCnm, columnGoTypeName))

		// Fields of a named embedded struct are reached through it, e.g. c.Contact.Name
		// 具名嵌入结构体里的字段要经由它访问，如 c.Contact.Name
		fieldAccess, viaIndirect := c.resolveFieldAccess(field)
		if c.options.isGenFuncTableColumns {
			must.Nice(c.options.columnsMethodRecvName)
			must.True(c.options.columnsCheckFieldType)
			valueColumnName := tern.BFF(viaIndirect, func() string {
				return fmt.Sprintf(`%s.New[%s](decoration.DecorateColumnName("%s"))`, pkgNameGormCnm, columnGoTypeName, field.DBName)
			}, func() string {
				return fmt.Sprintf(`%s.Cmn(%s.%s, "%s", decoration)`, pkgNameGormCnm, c.options.columnsMethodRecvName, fieldAccess, field.DBName)
			})
			methodTableColumnsPtx.Println("\t\t", fmt.Sprintf("%s:%s,", structFieldName, valueColumnName))
		} else {
			valueColumnName := tern.BFF(c.options.columnsMethodRecvName != "" && c.options.columnsCheckFieldType && !viaIndirect, func() string {
				return fmt.Sprintf(`%s.Cnm(%s.%s, "%s")`, pkgNameGormCnm, c.options.columnsMethodRecvName, fieldAccess, field.DBName)
			}, func() string {
				return `"` + field.DBName + `"`
			})
			methodPtx.Println("\t\t", fmt.Sprintf("%s:%s,", structFieldName, valueColumnName))
		}
	}
	structPtx.Println("}")
	if c.options.isGenFuncTableColumns {
		methodTableColumnsPtx.Println("\t}")
	} else {
		methodPtx.Println("\t}")
	}
	methodPtx.Println("}")
	methodTableColumnsPtx.Println("}")

	structCode := strings.TrimSpace(structPtx.String())
	methodCode := strings.TrimSpace(methodPtx.String())
	methodTableColumnsCode := strings.TrimSpace(methodTableColumnsPtx.String())

	zaplog.SUG.Debug("---", "\n", methodCode)
	zaplog.SUG.Debug("---", "\n", methodTableColumnsCode)
	zaplog.SUG.Debug("---", "\n", structCode)
	zaplog.SUG.Debug("---", "\n", neatjsons.S(pkgImports))

	return &GenOutput{
		methodCode:             methodCode,
		methodTableColumnsCode: methodTableColumnsCode,
		structCode:             structCode,
		pkgImports:             pkgImports,
	}
}

// resolveNewFieldName resolves the new field name based on tags and options.
// resolveNewFieldName 根据标签和选项解析新字段名称。
func (c *Config) resolveNewFieldName(field *schema.Field) (string, bool) {
	if c.options.useTagName {
		var tagKeyName = zerotern.VV(c.options.tagKeyName, "cnm")

		name, ok := field.Tag.Lookup(tagKeyName)
		if ok {
			if !utils.IsExportable(name) { // 确保字段名是导出的
				panic(erero.Errorf("name=%v is not exportable", name))
			}
			return name, true
		} else {
			if c.options.excludeUntaggedFields {
				return "", false
			}
			return strings.Join(field.EmbeddedBindNames, ""), true
		}
	}
	// A named embedded struct adds its field name in front, e.g. ContactName; an anonymous one adds nothing
	// 具名嵌入结构体会在前面加上它的字段名，如 ContactName；匿名嵌入不加
	return strings.Join(field.EmbeddedBindNames, ""), true
}

// resolveFieldAccess gives the dotted path that generated code uses to reach a field, e.g. "Contact.Name" in a named embedded struct.
// It also reports a path through a *T hop, since reading such a field panics when that hop holds no value.
//
// resolveFieldAccess 给出生成代码访问字段的路径，如具名嵌入结构体里的 "Contact.Name"。
// 同时报告路径是否经过指针，因为指针为空时读取这样的字段会 panic。
func (c *Config) resolveFieldAccess(field *schema.Field) (string, bool) {
	typ := c.sch.ModelType
	viaIndirect := false
	for _, name := range field.BindNames[:len(field.BindNames)-1] {
		structField, ok := typ.FieldByName(name)
		must.True(ok)
		typ = structField.Type
		if typ.Kind() == reflect.Pointer {
			viaIndirect = true
			typ = typ.Elem()
		}
	}
	return strings.Join(field.EmbeddedBindNames, "."), viaIndirect
}

// typeArgPattern matches a package-qualified name inside generic type args, which reflect writes with the whole import path, e.g. "gorm.io/datatypes.JSONType"
// typeArgPattern 匹配泛型类型参数里带完整导入路径的类型名，reflect 会把它写成 "gorm.io/datatypes.JSONType" 这种形式
var typeArgPattern = regexp.MustCompile(`([A-Za-z0-9_./-]+)\.([A-Za-z_][A-Za-z0-9_]*)`)

// resolveTypeName writes a field type as Go source seen from the struct's own package.
// Types from that package keep bare names; types from elsewhere get the package name and go into pkgImports.
// Pointers, slices, arrays, maps and generic type args are handled at each nesting step.
//
// resolveTypeName 以模型所在包的视角，把字段类型写成 Go 源码。
// 模型包内的类型不加包名；其他包的类型加上包名，并记入 pkgImports。
// 指针、切片、数组、map 和泛型类型参数都逐层处理。
func (c *Config) resolveTypeName(typ reflect.Type, pkgImports map[string]bool) string {
	if typ.Name() == "" {
		switch typ.Kind() {
		case reflect.Pointer:
			return "*" + c.resolveTypeName(typ.Elem(), pkgImports)
		case reflect.Slice:
			return "[]" + c.resolveTypeName(typ.Elem(), pkgImports)
		case reflect.Array:
			return fmt.Sprintf("[%d]%s", typ.Len(), c.resolveTypeName(typ.Elem(), pkgImports))
		case reflect.Map:
			return "map[" + c.resolveTypeName(typ.Key(), pkgImports) + "]" + c.resolveTypeName(typ.Elem(), pkgImports)
		default:
			return typ.String()
		}
	}
	name := typ.Name()
	if idx := strings.IndexByte(name, '['); idx >= 0 {
		name = name[:idx] + c.resolveTypeArgs(name[idx:], pkgImports)
	}
	pkgPath := typ.PkgPath()
	if pkgPath == "" || pkgPath == c.sch.ModelType.PkgPath() { // Predeclared types and types of the struct's own package // 预声明类型和模型包内的类型
		return name
	}
	pkgImports[pkgPath] = true
	pkgName, _, _ := strings.Cut(typ.String(), ".") // reflect puts the package name first // reflect 在开头写的是真实包名
	return pkgName + "." + name
}

// resolveTypeArgs rewrites the import paths inside generic type args into package names, following resolveTypeName.
// Here reflect gives just the import path, so the package name comes from its last element (minus a ".vN" / "/vN" version suffix).
//
// resolveTypeArgs 把泛型类型参数里的导入路径改写成包名，规则同 resolveTypeName。
// 这里 reflect 只给出导入路径，因此包名取路径最后一段（去掉 ".vN" 或 "/vN" 版本后缀）。
func (c *Config) resolveTypeArgs(args string, pkgImports map[string]bool) string {
	return typeArgPattern.ReplaceAllStringFunc(args, func(match string) string {
		groups := typeArgPattern.FindStringSubmatch(match)
		pkgPath, name := groups[1], groups[2]
		if pkgPath == c.sch.ModelType.PkgPath() {
			return name
		}
		pkgImports[pkgPath] = true
		return utils.PackageNameOfPath(pkgPath) + "." + name
	})
}

// ShowSchemaEnglish Displays schema information including struct name, table name, and fields.
// ShowSchemaEnglish 显示模式结构信息，包括结构体名称、表名和字段信息。
func ShowSchemaEnglish(sch *schema.Schema) {
	fmt.Println("---")
	fmt.Println("schema_message", "Struct name:", sch.Name, "Table name:", sch.Table, "Fields: {")
	for _, field := range sch.Fields {
		fmt.Println("   ",
			"Go Field Name:", field.Name,
			" | ",
			"Go Type:", field.FieldType,
			" | ",
			"DB Field Name:", field.DBName,
			" | ",
			"DB Type:", field.DataType,
			" | ",
			"Go Tag:", field.Tag,
		)
	}
	fmt.Println("}")
	fmt.Println("---")
}

// ShowSchemaChinese Displays schema information including struct name, table name, and fields.
// ShowSchemaChinese 显示模式结构信息，包括结构体名称、表名和字段信息。
func ShowSchemaChinese(sch *schema.Schema) {
	fmt.Println("---")
	fmt.Println("schema_message", "结构体名称:", sch.Name, "表名:", sch.Table, "字段信息: {")
	for _, field := range sch.Fields {
		fmt.Println("   ",
			"Go字段名:", field.Name, // Go结构体成员名称
			" | ",
			"Go类型:", field.FieldType, // Go的数据类型
			" | ",
			"DB字段名:", field.DBName, // 数据表中的列名
			" | ",
			"DB类型:", field.DataType, //数据库中的数据类型
			" | ",
			"Go标签:", field.Tag,
		)
	}
	fmt.Println("}")
	fmt.Println("---")
}

// ShowSchemaRelationshipsEnglish displays relationship information to demonstrate schema structure and assist debugging
// ShowSchemaRelationshipsEnglish 显示关系信息，用于教学和调试目的
func ShowSchemaRelationshipsEnglish(sch *schema.Schema) {
	if len(sch.Relationships.Relations) == 0 {
		return
	}

	fmt.Println("---")
	fmt.Println("relationships_message", "Struct name:", sch.Name, "Table name:", sch.Table, "Relationships: {")
	for _, relation := range sch.Relationships.Relations {
		fmt.Println("   ",
			"Association:", relation.Name,
			" | ",
			"Type:", relation.Type,
			" | ",
			"Target:", relation.FieldSchema.Name,
			" | ",
			"FK Columns:", tern.BVF(len(relation.References) == 0, "[]", func() string {
				var columns []string
				for _, ref := range relation.References {
					columns = append(columns, ref.ForeignKey.DBName)
				}
				return "[" + strings.Join(columns, ",") + "]"
			}),
		)
	}
	fmt.Println("}")
	fmt.Println("---")
}

// ShowSchemaRelationshipsChinese displays relationship information in Chinese to demonstrate schema structure and assist debugging
// ShowSchemaRelationshipsChinese 用中文显示关系信息，用于教学和调试目的
func ShowSchemaRelationshipsChinese(sch *schema.Schema) {
	if len(sch.Relationships.Relations) == 0 {
		return
	}

	fmt.Println("---")
	fmt.Println("relationships_message", "结构体名称:", sch.Name, "表名:", sch.Table, "关联关系: {")
	for _, relation := range sch.Relationships.Relations {
		fmt.Println("   ",
			"关联字段:", relation.Name,
			" | ",
			"关系类型:", relation.Type,
			" | ",
			"目标模型:", relation.FieldSchema.Name,
			" | ",
			"外键列名:", tern.BVF(len(relation.References) == 0, "[]", func() string {
				var columns []string
				for _, ref := range relation.References {
					columns = append(columns, ref.ForeignKey.DBName)
				}
				return "[" + strings.Join(columns, ",") + "]"
			}),
		)
	}
	fmt.Println("}")
	fmt.Println("---")
}

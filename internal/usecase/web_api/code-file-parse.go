package web_api

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// StructInfo 保存结构体的字段和方法信息
type StructInfo struct {
	Fields  []string
	Methods []string
}

// ParseResult 解析结果
type ParseResult struct {
	Structs      map[string]*StructInfo
	Interfaces   map[string][]string
	Constants    []string
	ExportedFunc []string
	ExportedVar  []string
}

// PrintResults 打印解析结果
func (p *ParseResult) PrintResults() {
	if len(p.Structs) > 0 {
		fmt.Println("Structs and Methods:")
		for structName, structInfo := range p.Structs {
			fmt.Printf("- %s:\n", structName)
			fmt.Println("  Fields:")
			for _, field := range structInfo.Fields {
				fmt.Printf("    - %s\n", field)
			}
			fmt.Println("  Methods:")
			for _, method := range structInfo.Methods {
				fmt.Printf("    - %s\n", method)
			}
		}
	}

	if len(p.Interfaces) > 0 {
		fmt.Println("\nInterfaces and Methods:")
		for interfaceName, methods := range p.Interfaces {
			fmt.Printf("- %s:\n", interfaceName)
			for _, method := range methods {
				fmt.Printf("  - Method: %s\n", method)
			}
		}
	}

	if len(p.Constants) > 0 {
		fmt.Println("\nConstants:")
		for _, constant := range p.Constants {
			if len(constant) > 64 {
				constant = constant[0:64] + "..."
			}
			fmt.Printf("- %s\n", constant)
		}
	}

	if len(p.ExportedFunc) > 0 {
		fmt.Println("\nExported Functions:")
		for _, fn := range p.ExportedFunc {
			fmt.Printf("- %s\n", fn)
		}
	}

	if len(p.ExportedVar) > 0 {
		fmt.Println("\nExported Variables:")
		for _, v := range p.ExportedVar {
			if len(v) > 64 {
				v = v[0:64] + "..."
			}
			fmt.Printf("- %s\n", v)
		}
	}

}

// Parser 解析器
type Parser struct {
	filePath string
}

// NewParser 创建新的解析器
func NewParser() *Parser {
	return &Parser{}
}

// ParseByFile Parse 解析源代码文件
func (p *Parser) ParseByFile(filePath string) (*ParseResult, error) {
	// 读取文件内容
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 创建文件集
	fset := token.NewFileSet()

	// 解析源代码文件
	f, err := parser.ParseFile(fset, "", string(fileContent), parser.ParseComments)
	if err != nil {
		return nil, err
	}
	result := ParseResult{
		Structs:      make(map[string]*StructInfo),
		Interfaces:   make(map[string][]string),
		Constants:    []string{},
		ExportedFunc: []string{},
		ExportedVar:  []string{},
	}
	// 遍历 AST 树
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.GenDecl:
			// 解析常量和变量声明
			if t.Tok == token.CONST {
				result.Constants = append(result.Constants, parseGenDecl(t)...)
			} else if t.Tok == token.VAR {
				result.ExportedVar = append(result.ExportedVar, parseExportedVars(t)...)
			}

		case *ast.TypeSpec:
			// 解析结构体或接口
			if structType, ok := t.Type.(*ast.StructType); ok {
				parseStruct(t, structType, result.Structs)
			} else if interfaceType, ok := t.Type.(*ast.InterfaceType); ok {
				parseInterface(t, interfaceType, result.Interfaces)
			}

		case *ast.FuncDecl:
			// 解析导出函数或方法
			parseFunc(t, result.Structs, &result.ExportedFunc)
		}
		return true
	})

	return &result, nil
}

// 解析结构体并存储字段和方法
func parseStruct(t *ast.TypeSpec, structType *ast.StructType, structs map[string]*StructInfo) {
	structName := t.Name.Name
	structs[structName] = &StructInfo{
		Fields:  []string{},
		Methods: []string{},
	}

	for _, field := range structType.Fields.List {
		fieldType := exprToString(field.Type)
		for _, name := range field.Names {
			structs[structName].Fields = append(structs[structName].Fields, fmt.Sprintf("%s: %s", name.Name, fieldType))
		}
	}
}

// 解析接口并存储方法
func parseInterface(t *ast.TypeSpec, interfaceType *ast.InterfaceType, interfaces map[string][]string) {
	interfaceName := t.Name.Name
	interfaces[interfaceName] = []string{}

	for _, method := range interfaceType.Methods.List {
		if len(method.Names) > 0 {
			methodName := method.Names[0].Name
			methodSignature := funcTypeToString(method.Type.(*ast.FuncType))
			interfaces[interfaceName] = append(interfaces[interfaceName], fmt.Sprintf("%s %s", methodName, methodSignature))
		}
	}
}

// 解析导出函数和方法
func parseFunc(t *ast.FuncDecl, structs map[string]*StructInfo, exportedFuncs *[]string) {
	if ast.IsExported(t.Name.Name) {
		params := getParamString(t.Type.Params)
		results := getParamString(t.Type.Results)

		if t.Recv != nil {
			// 解析方法的接收者
			receiverType := exprToString(t.Recv.List[0].Type)
			if structInfo, ok := structs[receiverType]; ok {
				structInfo.Methods = append(structInfo.Methods, fmt.Sprintf("%s(%s) (%s)", t.Name.Name, params, results))
			}
		} else {
			// 普通导出函数
			*exportedFuncs = append(*exportedFuncs, fmt.Sprintf("%s(%s) (%s)", t.Name.Name, params, results))
		}
	}
}

// 解析常量或变量声明
func parseGenDecl(genDecl *ast.GenDecl) []string {
	results := []string{}
	for _, spec := range genDecl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for i, name := range valueSpec.Names {
				if i < len(valueSpec.Values) {
					results = append(results, fmt.Sprintf("%s = %s", name.Name, exprToString(valueSpec.Values[i])))
				}
			}
		}
	}
	return results
}

// 解析导出变量
func parseExportedVars(genDecl *ast.GenDecl) []string {
	var exportedVars []string
	for _, spec := range genDecl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
				if ast.IsExported(name.Name) {
					val := ""
					if len(valueSpec.Values) > 0 {
						val = exprToString(valueSpec.Values[0])
					}

					exportedVars = append(exportedVars, fmt.Sprintf("%s = %s", name.Name, val))
				}
			}
		}
	}
	return exportedVars
}

// 获取参数字符串
func getParamString(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	var params []string
	for _, field := range fields.List {
		paramType := exprToString(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
			}
		} else {
			params = append(params, paramType)
		}
	}
	return strings.Join(params, ", ")
}

// 将表达式类型转为字符串
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.BasicLit:
		return t.Value
	case *ast.FuncType:
		return funcTypeToString(t)
	case *ast.MapType:
		// 新增对 map 类型的支持
		return fmt.Sprintf("map[%s]%s", exprToString(t.Key), exprToString(t.Value))
	case *ast.ChanType:
		// 通道类型 chan elemType 或 chan<- elemType
		dir := ""
		if t.Dir == ast.SEND {
			dir = "chan<- "
		} else if t.Dir == ast.RECV {
			dir = "<-chan "
		} else {
			dir = "chan "
		}
		return dir + exprToString(t.Value)
	case *ast.InterfaceType:
		// 接口类型 interface{}
		return "interface{}"
	case *ast.StructType:
		// 结构体类型 struct{}
		return "struct{}"
	default:
		return "unknown"
	}
}

// 将函数类型转为字符串
func funcTypeToString(funcType *ast.FuncType) string {
	params := getParamString(funcType.Params)
	results := getParamString(funcType.Results)
	return fmt.Sprintf("func(%s) (%s)", params, results)
}

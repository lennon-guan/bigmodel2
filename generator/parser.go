package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type (
	PropertyInfo struct {
		Name     string
		Return   string
		Source   string
		GetValue string
	}
	ModelInfo struct {
		Package     string
		Name        string
		Properties  []PropertyInfo
		Sources     []string
		SourceTypes map[string]string
	}
)

const (
	headCommentSource = "bigmodel-source:"
)

func ParseBigmodelInterface(filename string) ([]ModelInfo, error) {
	cf, err := NewCodeFile(filename)
	if err != nil {
		return nil, err
	}
	res := make([]ModelInfo, 0)
	for _, d := range cf.file.Decls {
		ts, inter := filterInterTypeDecl(d)
		if ts == nil {
			continue
		}
		var (
			name        = ts.Name.Name
			properties  = make([]PropertyInfo, 0)
			sourceTypes = make(map[string]string)
		)
		if doc := d.(*ast.GenDecl).Doc; doc != nil && doc.List != nil {
			for _, comment := range doc.List {
				c := getCommentLineContent(comment.Text)
				if strings.HasPrefix(c, headCommentSource) {
					kv := strings.Split(strings.TrimSpace(c[len(headCommentSource):]), " ")
					if len(kv) == 2 {
						sourceTypes[kv[0]] = kv[1]
					}
				}
			}
		}
		for _, method := range inter.Methods.List {
			if p, err := getMethodInfo(cf, name, method); err != nil {
				return nil, err
			} else {
				properties = append(properties, *p)
			}
		}
		allSource, m := []string{}, map[string]int{}
		for _, p := range properties {
			if _, found := m[p.Source]; !found {
				m[p.Source] = 1
				allSource = append(allSource, p.Source)
			}
		}
		res = append(res, ModelInfo{
			Package:     cf.file.Name.Name,
			Name:        name,
			Properties:  properties,
			Sources:     allSource,
			SourceTypes: sourceTypes,
		})
	}
	return res, nil
}

func filterInterTypeDecl(d ast.Decl) (ts *ast.TypeSpec, interf *ast.InterfaceType) {
	decl, isGen := d.(*ast.GenDecl)
	if !isGen || decl.Tok != token.TYPE || len(decl.Specs) < 1 {
		return nil, nil
	}
	ts, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil, nil
	}
	inter, ok := ts.Type.(*ast.InterfaceType)
	if !ok {
		return nil, nil
	}
	if inter.Methods == nil {
		return nil, nil
	}
	return ts, inter
}

func getMethodInfo(cf *CodeFile, typeName string, method *ast.Field) (*PropertyInfo, error) {
	if len(method.Names) != 1 {
		return nil, fmt.Errorf("%s.%s: 方法名数量不是1", typeName, cf.GetText(method))
	}
	methodName := method.Names[0].Name
	ft, ok := method.Type.(*ast.FuncType)
	if !ok {
		return nil, fmt.Errorf("%s.%s: 不是一个方法", typeName, methodName)
	}
	if ft.Params != nil && len(ft.Params.List) > 0 {
		return nil, fmt.Errorf("%s.%s: 参数表必须为空", typeName, methodName)
	}
	if ft.Results == nil || len(ft.Results.List) != 1 || len(ft.Results.List[0].Names) > 1 {
		return nil, fmt.Errorf("%s.%s: 返回值数量必须为1", typeName, methodName)
	}
	pi := &PropertyInfo{
		Name:   methodName,
		Return: cf.GetText(ft.Results.List[0].Type),
	}
	if c := method.Comment; c == nil || len(c.List) != 1 {
		return nil, fmt.Errorf("%s.%s: 解析注释失败", typeName, methodName)
	} else {
		comment := getCommentLineContent(c.List[0].Text)
		if sp := strings.Index(comment, "."); sp != -1 {
			pi.Source = comment[:sp]
			pi.GetValue = comment[sp+1:]
		} else {
			pi.Source = comment
			pi.GetValue = methodName
		}
	}

	return pi, nil
}

func getCommentLineContent(s string) string {
	if strings.HasPrefix(s, "//") {
		s = s[2:]
	} else if strings.HasPrefix(s, "/*") && strings.HasSuffix(s, "*/") {
		s = s[2 : len(s)-2]
	}
	return strings.TrimSpace(s)
}

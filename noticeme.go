// Copyright 2020 Molecula Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package noticeme

import (
	"errors"
	"go/ast"
	"go/types"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check for unused values by type

Invoke with -important <typelist>, with comma-separated types. Names are
checked against the name with no qualifier, just a package name, and the
full import path. Reports any expression statements that have one or more
of the given types.
`

var Analyzer = &analysis.Analyzer{
	Name:     "noticeme",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      noticeme,
}

var importantTypes string

func init() {
	Analyzer.Flags.StringVar(&importantTypes, "types", "", "specify important types")
}

type typeList []string

func typeNameOnly(pkg *types.Package) string {
	return ""
}

func packageNameOnly(pkg *types.Package) string {
	return pkg.Name()
}

// matchImportance returns the first string in the list it found which
// matches the last component of the name of the given type, or of a type
// within it if it's a tuple.
func (tl typeList) matchImportance(t types.Type) (bool, string) {
	if tuple, ok := t.(*types.Tuple); ok {
		for i := 0; i < tuple.Len(); i++ {
			subType := tuple.At(i).Type()
			important, why := tl.matchImportance(subType)
			if important {
				return important, why
			}
		}
	}
	names := []string{
		types.TypeString(t, typeNameOnly),
		types.TypeString(t, packageNameOnly),
		types.TypeString(t, nil),
	}
	for _, w := range tl {
		for _, name := range names {
			if name == w {
				return true, name
			}
		}
	}
	return false, ""
}

var parseImportant sync.Once
var importantList typeList
var relevantTypes = map[types.Type]string{}

func noticeme(pass *analysis.Pass) (_ interface{}, err error) {
	parseImportant.Do(func() {
		if importantTypes == "" {
			err = errors.New("you must specify which types you care about (-important)")
			return
		}
		importantList = strings.Split(importantTypes, ",")
	})
	if err != nil {
		return nil, err
	}
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.ExprStmt)(nil)}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		expr := n.(*ast.ExprStmt).X
		exprType := pass.TypesInfo.Types[expr].Type
		relevant, ok := relevantTypes[exprType]
		if !ok {
			_, relevant = importantList.matchImportance(exprType)
			relevantTypes[exprType] = relevant
		}
		if relevant != "" {
			pass.Reportf(n.Pos(), "unused value of type %s", relevant)
		}
	})
	return nil, nil
}

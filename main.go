// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

var fset = token.NewFileSet()

func main() {
	flag.Parse()

	var conf loader.Config
	conf.Fset = fset
	for _, importPath := range gotool.ImportPaths(flag.Args()) {
		conf.Import(importPath)
	}
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range prog.InitialPackages() {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if s, ok := node.(*ast.RangeStmt); ok {
					danger(s, pkg)
				}
				return true
			})
		}
	}
}

func danger(s *ast.RangeStmt, pkg *loader.PackageInfo) {
	if s.Value == nil {
		return
	}

	tv := pkg.Types[s.X]
	if _, ok := tv.Type.Underlying().(*types.Array); !ok || !tv.Addressable() {
		return
	}

	fmt.Printf("%s: range over addressable array\n", fset.Position(s.Pos()))
}

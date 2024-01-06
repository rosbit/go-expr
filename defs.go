package ex

import (
	"github.com/expr-lang/expr/vm"
)

type fn struct {
	Params []string
	Expr string
	prg *vm.Program
}

type exprScript struct {
	Envs map[string]interface{}
	Funcs map[string]*fn
}

type XExpr struct {
	*exprScript
}


package ex

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"github.com/expr-lang/expr"
)

func bindGoFunc(name string, funcVar interface{}) (goFunc expr.Option, err error) {
	helper, e := elutils.NewGolangFuncHelper(funcVar, name)
	if e != nil {
		err = e
		return
	}

	if len(name) == 0 {
		name = helper.GetRealName()
	}
	goFunc = expr.Function(name, wrapGoFunc(name, helper))
	return
}

func wrapGoFunc(name string, helper *elutils.GolangFuncHelper) func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (val interface{}, err error) {
		getArgs := func(i int) interface{} {
			return args[i]
		}

		val, err = helper.CallGolangFunc(len(args), name, getArgs)
		return
	}
}

package ex

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"github.com/expr-lang/expr"
	"reflect"
	"fmt"
)

func (e *XExpr) bindFunc(funcName string, funcVarPtr interface{}) (err error) {
	helper, er := elutils.NewEmbeddingFuncHelper(funcVarPtr)
	if er != nil {
		err = er
		return
	}

	wf := e.wrapFunc(funcName, helper)
	if wf == nil {
		err = fmt.Errorf("func %s not found")
		return
	}
	helper.BindEmbeddingFunc(wf)
	return
}

func (e *XExpr) wrapFunc(funcName string, helper *elutils.EmbeddingFuncHelper) elutils.FnGoFunc {
	funcMeta, err := e.getFunc(funcName)
	if err != nil {
		return nil
	}

	return func(args []reflect.Value) (results []reflect.Value) {
		exArgs := make(map[string]interface{})

		// copy global env
		for k, v := range e.Envs {
			exArgs[k] = v
		}

		// make expr args
		itArgs := helper.MakeGoFuncArgs(args)
		i := 0
		for arg := range itArgs {
			if i < len(funcMeta.ArgNames) {
				exArgs[funcMeta.ArgNames[i]] = arg
				i += 1
			}
		}

		// call expr function
		res, err := expr.Eval(funcMeta.Expr, exArgs)

		// convert result to golang
		if res != nil {
			resKind := reflect.TypeOf(res).Kind()
			isArray := resKind == reflect.Slice || resKind == reflect.Array
			results = helper.ToGolangResults(res, isArray, err)
		} else {
			results = helper.ToGolangResults(nil, false, err)
		}
		return
	}
}

func (e *XExpr) callFunc(funcMeta *fn, args ...interface{}) (res interface{}, err error) {
	exArgs := make(map[string]interface{})

	// copy global env
	for k, v := range e.Envs {
		exArgs[k] = v
	}

	// set args
	for i, arg := range args {
		if i < len(funcMeta.ArgNames) {
			exArgs[funcMeta.ArgNames[i]] = arg
		}
	}

	res, err = expr.Eval(funcMeta.Expr, exArgs)
	return
}

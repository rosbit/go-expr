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
		itArgs := helper.MakeGoFuncArgs(args)
		res, err := e.evalWithArgs(funcMeta, itArgs)

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
	// copy global env
	exArgs := e.makeEnvs(nil)

	// set args
	for i, arg := range args {
		if i < len(funcMeta.Params) {
			exArgs[funcMeta.Params[i]] = arg
		}
	}

	return funcMeta.eval(exArgs)
}

func (e *XExpr) evalWithArgs(f *fn, args <-chan interface{}) (res interface{}, err error) {
	// copy global env
	exArgs := e.makeEnvs(nil)

	// make expr args
	i := 0
	for arg := range args {
		if i < len(f.Params) {
			exArgs[f.Params[i]] = arg
			i += 1
		}
	}

	return f.eval(exArgs)
}

func (f *fn) eval(envs map[string]interface{}) (res interface{}, err error) {
	// res, err = expr.Eval(f.Expr, envs)
	res, err = expr.Run(f.prg, envs)
	return
}

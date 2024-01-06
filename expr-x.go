package ex

import (
	"github.com/expr-lang/expr"
	"fmt"
	"reflect"
)

func New() *XExpr {
	return &XExpr{}
}

func (e *XExpr) LoadFile(path string, vars map[string]interface{}) (err error) {
	script, er := loadScript(path)
	if er != nil {
		err = er
		return
	}

	e.exprScript = script
	if len(vars) > 0 {
		if err = e.setEnv(vars); err != nil {
			return
		}
	}
	return nil
}

func (e *XExpr) SetEnvs(vars map[string]interface{}) (err error) {
	if e.exprScript == nil {
		e.exprScript = &exprScript{}
	} else {
		e.exprScript.Envs = nil
	}
	if len(vars) > 0 {
		if err = e.setEnv(vars); err != nil {
			return
		}
	}
	return nil
}

func (e *XExpr) Eval(script string, vars map[string]interface{}) (res interface{}, err error) {
	if e.exprScript == nil || e.exprScript.Envs == nil {
		return expr.Eval(script, vars)
	}
	env := e.makeEnvs(vars)
	return expr.Eval(script, env)
}

func (e *XExpr) makeEnvs(vars map[string]interface{}) (map[string]interface{}) {
	envs := make(map[string]interface{})
	for k, v := range e.Envs {
		envs[k] = v
	}
	for k, v := range vars {
		envs[k] = v
	}
	return envs
}

func (e *XExpr) CallFunc(funcName string, args ...interface{}) (res interface{}, err error) {
	funcMeta, er := e.getFunc(funcName)
	if er != nil {
		err = er
		return
	}

	res, err = e.callFunc(funcMeta, args...)
	return
}

// bind a var of golang func with a expr function name, so calling expr function
// is just calling the related golang func.
// @param funcVarPtr  in format `var funcVar func(....) ...; funcVarPtr = &funcVar`
func (e *XExpr) BindFunc(funcName string, funcVarPtr interface{}) (err error) {
	if funcVarPtr == nil {
		err = fmt.Errorf("funcVarPtr must be a non-nil poiter of func")
		return
	}
	t := reflect.TypeOf(funcVarPtr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
		err = fmt.Errorf("funcVarPtr expected to be a pointer of func")
		return
	}

	err = e.bindFunc(funcName, funcVarPtr)
	return
}

func (e *XExpr) BindFuncs(funcName2FuncVarPtr map[string]interface{}) (err error) {
	for funcName, funcVarPtr := range funcName2FuncVarPtr {
		if err = e.BindFunc(funcName, funcVarPtr); err != nil {
			return
		}
	}
	return
}

func (e *XExpr) setEnv(vars map[string]interface{}) (err error) {
	if e.Envs == nil {
		e.Envs = make(map[string]interface{})
	}
	for k, v := range vars {
		if v == nil {
			continue
		}

		e.Envs[k] = v
	}
	return
}

func (e *XExpr) getFunc(name string) (v *fn, err error) {
	if len(e.Funcs) == 0 {
		err = fmt.Errorf("no func named %s found", name)
		return
	}
	r, ok := e.Funcs[name]
	if !ok {
		err = fmt.Errorf("no func named %s found", name)
		return
	}
	v = r
	return
}

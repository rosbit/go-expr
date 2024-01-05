package ex

import (
	"sync"
	"os"
	"time"
)

type exprCtx struct {
	exx *XExpr
	mt   time.Time
}

var (
	exprCtxCache map[string]*exprCtx
	lock *sync.Mutex
)

func InitCache() {
	if lock != nil {
		return
	}
	lock = &sync.Mutex{}
	exprCtxCache = make(map[string]*exprCtx)
}

func LoadFileFromCache(path string, vars map[string]interface{}) (ctx *XExpr, existing bool, err error) {
	lock.Lock()
	defer lock.Unlock()

	exprC, ok := exprCtxCache[path]

	if !ok {
		ctx = New()
		if err = ctx.LoadFile(path, vars); err != nil {
			return
		}
		fi, _ := os.Stat(path)
		exprC = &exprCtx{
			exx: ctx,
			mt: fi.ModTime(),
		}
		exprCtxCache[path] = exprC
		return
	}

	fi, e := os.Stat(path)
	if e != nil {
		err = e
		return
	}
	mt := fi.ModTime()
	if !exprC.mt.Equal(mt) {
		ctx = New()
		if err = ctx.LoadFile(path, vars); err != nil {
			return
		}
		exprC.exx = ctx
		exprC.mt = mt
	} else {
		existing = true
		ctx = exprC.exx
	}
	return
}

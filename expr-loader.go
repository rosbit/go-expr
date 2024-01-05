package ex

import (
	"gopkg.in/yaml.v2"
	"os"
	"fmt"
)

func loadScript(scriptFile string) (script *exprScript, err error) {
	fp, e := os.Open(scriptFile)
	if e != nil {
		err = e
		return
	}
	defer fp.Close()

	var s exprScript
	dec := yaml.NewDecoder(fp)
	if err = dec.Decode(&s); err != nil {
		return
	}

	if err = checkMust(&s); err != nil {
		return
	}
	script = &s
	return
}

func checkMust(script *exprScript) (err error) {
	for fn, f := range script.Funcs {
		if f == nil {
			err = fmt.Errorf("func %s definition expected", fn)
			return
		}
		if len(f.Expr) == 0 {
			err = fmt.Errorf("func %s expr expected", fn)
			return
		}
	}
	return
}

package ex

type fn struct {
	ArgNames []string `yaml:"arg-names"`
	Expr string
}

type exprScript struct {
	Envs map[string]interface{}
	Funcs map[string]*fn
}

type XExpr struct {
	*exprScript
}


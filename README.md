# go-expr, an extending of expr-lange

[Expr](https://github.com/expr-lang/expr) is a Go-centric expression language designed to deliver dynamic
configurations with unparalleled accuracy, safety, and speed

`go-expr` is a package extending the `Expr` and making it a **pragmatic embeddable** language.
With some helper functions provided by `go-expr`, calling Golang functions from `Expr`, 
or calling `Expr` functions from Golang are both very simple. So, with the help of `go-expr`, expr-lang
can be looked as **an embeddable expr-lang**.

### Expr functions

All to be evaluated in `Expr` are expression strings. Function is an extending notion to wrapper `Expr` expressions.
We make use of `YAML` syntax to declare `Expr` functions. e.g., Function `add` is desclared like this:
```yaml
funcs:
  add:
    arg-names: [A, B]
    expr: A + B

  other-funcs:
    ..
```

`arg-names` are argument variable names for function. `expr` is the body of function, is just the `Expr` expressions.

### Usage

The package is fully go-getable, so, just type

  `go get github.com/rosbit/go-expr`

to install.

#### 1. Evaluate expressions

```go
package main

import (
  "github.com/rosbit/go-expr"
  "fmt"
)

func main() {
  ctx := ex.New()

  res, _ := ctx.Eval("1 + 2", nil)
  fmt.Println("result is:", res)
}
```

#### 2. Go calls Expr function:

Suppose there's a Expr file named `a.yaml` like this:

```yaml
funcs:
  add:
    arg-names: [A, B]
    expr: A + B

```

one can call the Expr function `add()` in Go code like the following:

```go
package main

import (
  "github.com/rosbit/go-expr"
  "fmt"
)

var add func(int, int)int

func main() {
  ctx := ex.New()
  if err := ctx.LoadFile("a.yaml", nil); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  if err := ctx.BindFunc("add", &add); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  res := add(1, 2)
  fmt.Println("result is:", res)
}
```

#### 3. Expr calls Go function

Expr calling Go function is also easy. In the Go code, make a Golang function
as Expr built-in function by registering golang functions. There's the example:

```go
package main

import "github.com/rosbit/go-expr"

// function to be called by Expr
func adder(a1 float64, a2 float64) float64 {
    return a1 + a2
}

func main() {
  ctx := ex.New()

  if err := ctx.LoadFile("b.yaml", map[string]interface{}{
     "adder": adder,
  }); err != nil {
     // error handler
     return
  }
  res, _ := ctx.CallFunc("add", 10, 2)
  fmt.Println("result is:", res)
}
```

In Expr code, one can call the registered function directly. There's the example `b.yaml`

```yaml
funcs:
  add:
    arg-names: [A, B]
    expr: adder(A, B)  # the function "adder" is implmented in Go.
```

### Status

The package is not fully tested, so be careful.

### Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.

__Convention:__ fork the repository and make changes on your fork in a feature branch.

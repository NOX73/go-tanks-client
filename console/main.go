// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"go/token"

	"code.google.com/p/go.crypto/ssh/terminal"
	eval "github.com/sbinet/go-eval"

	client "github.com/NOX73/go-tanks-client"
)

var m = client.Message{}
var fset = token.NewFileSet()

type shell struct {
	r io.Reader
	w io.Writer
}

func (sh *shell) Read(data []byte) (n int, err error) {
	return sh.r.Read(data)
}
func (sh *shell) Write(data []byte) (n int, err error) {
	return sh.w.Write(data)
}

type funcV struct {
	target eval.Func
}

func (v *funcV) String() string {
	// TODO(austin) Rob wants to see the definition
	return "func {...}"
}

func (v *funcV) Assign(t *eval.Thread, o eval.Value) { v.target = o.(eval.FuncValue).Get(t) }
func (v *funcV) Get(*eval.Thread) eval.Func          { return v.target }
func (v *funcV) Set(t *eval.Thread, x eval.Func)     { v.target = x }

type testFunc struct{}

var c client.Client

func (*testFunc) NewFrame() *eval.Frame { return &eval.Frame{nil, make([]eval.Value, 2)} }
func (*testFunc) Call(t *eval.Thread) {
	var err error
	c, err = client.ConnectTo("nox73.ru:9000")
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Connected! ")
	}
}

func main() {
	w := eval.NewWorld()
	w.DefineVar("connect", eval.NewFuncType([]eval.Type{}, false, []eval.Type{}), &funcV{&testFunc{}})

	fmt.Println(":: welcome to go-eval...\n(hit ^D to exit)")

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)
	term := terminal.NewTerminal(&shell{r: os.Stdin, w: os.Stdout}, "> ")
	if term == nil {
		panic(errors.New("could not create terminal"))
	}

	for {
		line, err := term.ReadLine()
		if err != nil {
			break
		}
		code, err := w.Compile(fset, line)
		if err != nil {
			term.Write([]byte(err.Error() + "\n"))
			continue
		}
		v, err := code.Run()
		if err != nil {
			term.Write([]byte(err.Error() + "\n"))
			continue
		}
		if v != nil {
			term.Write([]byte(v.String() + "\n"))
		}
	}
}

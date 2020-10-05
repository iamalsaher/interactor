package main

import (
	"fmt"

	"github.com/iamalsaher/interactor/pkg/pty"
)

func main() {
	p, e := pty.NewPTY()
	if e != nil {
		panic(e)
	}
	fmt.Println(p.Master.Fd())
	fmt.Println(p.Slave.Fd())

	if e := p.Close(); e != nil {
		panic(e)
	}
}

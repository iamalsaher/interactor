package main

import (
	"github.com/iamalsaher/interactor/pkg/process"
)

func main() {
	proc := process.NewProcess("binary")
	// proc.SetDirectory("/tmp")
	// proc.SetTimeout(2000)
	proc.ConnectIO(false)
	proc.Start()
}

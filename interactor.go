package main

import (
	"fmt"

	"github.com/iamalsaher/interactor/pkg/process"
)

func main() {
	proc, log := process.NewProcess("binary")
	fmt.Println(log)
	// proc.SetDirectory("/tmp")
	// proc.SetTimeout(2000)
	proc.ConnectIO()
	proc.Start()
	proc.Handle.Wait()
	proc.Show()
}

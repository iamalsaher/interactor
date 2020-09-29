package main

import (
	"fmt"

	"github.com/iamalsaher/interactor/pkg/process"
)

func main() {
	proc, log := process.NewProcess("/bin/sh")
	fmt.Println(log)
	proc.SetDirectory("/tmp")
	// proc.SetIO(os.Stdin, os.Stdout, os.Stderr)
	proc.Start()
	// fmt.Println(proc.Recv(15))
	proc.Handle.Wait()
}

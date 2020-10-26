package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/iamalsaher/interactor/pkg/process"
)

func interactor(input io.Writer, output io.Reader, closer chan bool) {
	var b bytes.Buffer

	b.Grow(8192)
	go io.Copy(&b, output)
	go input.Write([]byte("Swapnil\r\n"))

	for {
		select {
		case <-closer:
			// fmt.Print(string(b.Bytes()))
			return
		default:
			if b.Len() > 0 {
				fmt.Print(string(b.Next(b.Len())))
			}
		}
	}
}

func main() {
	proc := process.NewProcess(`.\binary.exe`, "--tty", "--sleep", "2", "--input")
	// proc.SetEnviron([]string{"SEXYENV=LOL"}, true)
	// proc.SetDirectory("/tmp")
	// proc.SetTimeout(1 * time.Second)
	if e := proc.ConnectIO(true); e != nil {
		panic(e)
	}

	if e := proc.Start(); e != nil {
		panic(e)
	}

	fmt.Printf("Started Process with PID %v\n", proc.PID)
	interactor(proc.Stdin, proc.Stdout, proc.Done)
	fmt.Printf("exit code was: %d\n", proc.State.ExitCode())
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/iamalsaher/interactor/pkg/process"
)

var wg sync.WaitGroup

func interactor(input io.Writer, output io.Reader) {
	var b bytes.Buffer

	b.Grow(8192)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				if b.Len() > 0 {
					fmt.Print(string(b.Next(b.Len())))
				}
			}
		}
	}()

	io.Copy(&b, output)

	wg.Done()
	close(done)
}

func main() {
	proc := process.NewProcess(`./binary`, "--tty", "--sleep", "5")
	// proc.SetEnviron([]string{"SEXYENV=LOL"}, true)
	// proc.SetDirectory("/tmp")
	// proc.SetTimeout(1000)
	if e := proc.ConnectIO(true, false); e != nil {
		panic(e)
	}

	wg.Add(1)
	go interactor(proc.Stdin, proc.Stdout)

	if e := proc.Start(); e != nil {
		panic(e)
	}

	fmt.Printf("Started Process with PID %v\n", proc.PID)
	wg.Wait()
	fmt.Printf("exit code was: %d\n", proc.State.ExitCode())
}

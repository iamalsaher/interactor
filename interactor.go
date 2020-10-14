package main

import (
	"fmt"
	"io"
	"log"

	"github.com/iamalsaher/interactor/pkg/process"
)

func interactor(input io.Writer, output io.Reader) {
	c := make([]byte, 1)
	for {
		if _, e := output.Read(c); e == nil {
			fmt.Print(string(c[0]))
		} else {
			log.Println(e.Error())
			break
		}
	}
}

func main() {
	proc := process.NewProcess("./binary", "--sleep", "2")
	// proc.SetDirectory("/tmp")
	proc.SetTimeout(1000)
	if e := proc.ConnectIO(false, false); e != nil {
		panic(e)
	}

	if e := proc.Start(&process.Interactor{Function: interactor, Input: proc.Pipe.StdinW, Output: proc.Pipe.StdoutR}); e != nil {
		panic(e)
	}

	fmt.Printf("Started Process with PID %v\n", proc.PID)

	ps, err := proc.Wait()
	if err != nil {
		log.Fatalf("Error waiting for process: %v", err)
	}
	fmt.Printf("exit code was: %d\n", ps.ExitCode())
}

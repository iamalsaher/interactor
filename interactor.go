package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/iamalsaher/interactor/pkg/process"
)

func main() {
	fmt.Println("Started")
	proc := process.NewProcess("./binary", "--sleep", "2")
	// proc.SetDirectory("/tmp")
	proc.SetTimeout(3000)
	if e := proc.ConnectIO(false, false); e != nil {
		panic(e)
	}
	// input, output := proc.Pipe.StdinW, proc.Pipe.StdoutR

	if e := proc.Start(); e != nil {
		panic(e)
	}
	fmt.Println(proc.PID)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if proc.IO() {
			c := make([]byte, 1)
			for {
				if _, e := proc.Pipe.StdoutR.Read(c); e == nil {
					fmt.Print(string(c[0]))
				} else {
					log.Println(e.Error())
					break
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

	ps, err := proc.Wait()
	if err != nil {
		log.Fatalf("Error waiting for process: %v", err)
	}
	fmt.Printf("exit code was: %d\n", ps.ExitCode())
}

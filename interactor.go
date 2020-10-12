package main

import (
	"fmt"
	"log"

	"github.com/iamalsaher/interactor/pkg/process"
	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println("Started")
	proc := process.NewProcess("binary.exe")
	// proc.SetDirectory("/tmp")
	// proc.SetTimeout(2000)
	if e := proc.ConnectIO(false); e != nil {
		panic(e)
	}
	master, slave := proc.GetIO()
	buff := make([]byte, 400)

	if e := proc.Start(); e != nil {
		panic(e)
	}
	fmt.Println(proc.PID)

	if n, e := slave.Read(buff); e != nil {
		panic(e)
	} else {
		fmt.Printf("Read %v bytes\nData:%v\n", n, string(buff))
	}
	var n uint32
	windows.WriteFile(windows.Handle(master.Fd()), []byte("lundloda"), &n, nil)

	ps, err := proc.Wait()
	if err != nil {
		log.Fatalf("Error waiting for process: %v", err)
	}
	fmt.Printf("exit code was: %d\n", ps.ExitCode())
}

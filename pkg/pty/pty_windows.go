package pty

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32            = windows.NewLazySystemDLL("kernel32.dll")
	createPseudoConsole = kernel32.NewProc("CreatePseudoConsole")
	closePseudoConsole  = kernel32.NewProc("ClosePseudoConsole")
)

//PTY is the structure which contains the input and output files to pass to processes
type PTY struct {
	Master *os.File
	Slave  *os.File

	size    windows.Coord
	phPC    *windows.Handle
	hInput  windows.Handle
	hOutput windows.Handle
}

//NewPTY returns a master and slave file descriptor for a pty
//https://github.com/ActiveState/termtest/tree/master/conpty
func NewPTY() (*PTY, error) {

	var (
		hPipeIn  windows.Handle
		hPipeOut windows.Handle
	)

	p := PTY{phPC: new(windows.Handle)}

	if err := windows.CreatePipe(&hPipeIn, &p.hInput, nil, 0); err != nil {
		panic(err)
	}

	if err := windows.CreatePipe(&hPipeOut, &p.hOutput, nil, 0); err != nil {
		panic(err)
	}

	r1, _, e := createPseudoConsole.Call(uintptr(unsafe.Pointer(&p.size)), uintptr(hPipeIn), uintptr(hPipeOut), 0, uintptr(unsafe.Pointer(p.phPC)))
	if r1 != 0 {
		fmt.Printf("%x\n", r1)
		return nil, fmt.Errorf("Code: 0x%x Error:%v", r1, e)
	}

	if hPipeIn != windows.InvalidHandle {
		windows.CloseHandle(hPipeIn)
	}

	if hPipeOut != windows.InvalidHandle {
		windows.CloseHandle(hPipeOut)
	}

	p.Master = os.NewFile(uintptr(p.hInput), "|0")
	p.Slave = os.NewFile(uintptr(p.hOutput), "|1")

	return &p, nil
}

//Close is used to release all resources allocated to PTY
func (p *PTY) Close() error {
	r1, _, e := closePseudoConsole.Call(uintptr(*p.phPC))
	if r1 == 0 {
		return fmt.Errorf("Code: 0x%x Error:%v", r1, e)
	}

	p.Master.Close()
	p.Slave.Close()
	return nil
}

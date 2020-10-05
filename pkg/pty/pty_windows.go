package pty

import (
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
}

//NewPTY returns a master and slave file descriptor for a pty
//https://github.com/ActiveState/termtest/tree/master/conpty
func NewPTY() (*PTY, error) {

	var (
		hInput   windows.Handle
		hOutput  windows.Handle
		hPipeIn  windows.Handle
		hPipeOut windows.Handle
	)

	if err := windows.CreatePipe(&hPipeIn, &hInput, nil, 0); err != nil {
		panic(err)
	}

	if err := windows.CreatePipe(&hPipeOut, &hOutput, nil, 0); err != nil {
		panic(err)
	}

	r1, _, e := createPseudoConsole.Call(uintptr(unsafe.Pointer(new(windows.Coord))), uintptr(hPipeIn), uintptr(hPipeOut), 0, uintptr(unsafe.Pointer(new(windows.Handle))))
	if r1 != uintptr(windows.S_OK) {
		return nil, e
	}

	if hPipeIn != windows.InvalidHandle {
		windows.CloseHandle(hPipeIn)
	}

	if hPipeOut != windows.InvalidHandle {
		windows.CloseHandle(hPipeOut)
	}

	master := os.NewFile(uintptr(hInput), "|0")
	slave := os.NewFile(uintptr(hInput), "|1")

	return &PTY{Master: master, Slave: slave}, nil
}

package pty

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	kernel32                          = windows.NewLazySystemDLL("kernel32.dll")
	createPseudoConsole               = kernel32.NewProc("CreatePseudoConsole")
	procResizePseudoConsole           = kernel32.NewProc("ResizePseudoConsole")
	closePseudoConsole                = kernel32.NewProc("ClosePseudoConsole")
	initializeProcThreadAttributeList = kernel32.NewProc("InitializeProcThreadAttributeList")
	updateProcThreadAttribute         = kernel32.NewProc("UpdateProcThreadAttribute")
	deleteProcThreadAttributeList     = kernel32.NewProc("DeleteProcThreadAttributeList")
)

const minConPtySupportBuild int = 18362

//StartupInfoEx exposes the Extended StartupInfo which can be passed to CreateProcess
type StartupInfoEx struct {
	StartupInfo     windows.StartupInfo
	lpAttributeList windows.Handle
}

//PTY is the structure which contains the input and output files to pass to processes
type PTY struct {
	Input  *os.File
	Output *os.File
	SIX    StartupInfoEx

	size    uintptr
	phPC    *windows.Handle
	hInput  windows.Handle
	hOutput windows.Handle

	attrListBuffer []byte
}

func getSize(c, r int16) uintptr {
	return uintptr(c) + (uintptr(r) << 16)
}

func conPTYSupport() bool {

	k, e := registry.OpenKey(registry.LOCAL_MACHINE, `Software`+`\Microsoft`+`\Windows NT`+`\CurrentVersion`, registry.QUERY_VALUE)
	if e != nil {
		panic(e)
	}

	v, _, e := k.GetStringValue("CurrentBuild")
	k.Close()

	s, err := strconv.Atoi(v)
	if err != nil {
		panic(e)
	}

	if s < minConPtySupportBuild {
		return false
	}
	return true
}

//NewPTY returns a master and slave file descriptor for a pty
//https://github.com/ActiveState/termtest/tree/master/conpty
//https://docs.microsoft.com/en-us/windows/console/creating-a-pseudoconsole-session
func NewPTY() (*PTY, error) {

	if !conPTYSupport() {
		return nil, fmt.Errorf("ConPTY support is not present in this version of Windows")
	}

	var (
		hPipeIn  windows.Handle
		hPipeOut windows.Handle
	)
	const procThreadAttributePseudoconsole uintptr = 0x00020016

	p := &PTY{
		phPC: new(windows.Handle),
		size: getSize(80, 32),
	}

	if err := windows.CreatePipe(&hPipeIn, &p.hInput, nil, 0); err != nil {
		panic(err)
	}

	if err := windows.CreatePipe(&p.hOutput, &hPipeOut, nil, 0); err != nil {
		panic(err)
	}

	r, _, e := createPseudoConsole.Call(p.size, uintptr(hPipeIn), uintptr(hPipeOut), 0, uintptr(unsafe.Pointer(p.phPC)))
	if r != 0 {
		return nil, fmt.Errorf("createPseudoConsole Error:%v Code: 0x%x", e, r)
	}

	if hPipeIn != windows.InvalidHandle {
		windows.CloseHandle(hPipeIn)
	}

	if hPipeOut != windows.InvalidHandle {
		windows.CloseHandle(hPipeOut)
	}

	p.Input = os.NewFile(uintptr(p.hInput), "|0")
	p.Output = os.NewFile(uintptr(p.hOutput), "|1")

	//Setting up the process structure to connect conpty to the process
	var attrListSize uint64

	// Prepare Startup Information structure
	p.SIX.StartupInfo.Cb = uint32(unsafe.Sizeof(p.SIX))

	// Discover the size required for the list
	initializeProcThreadAttributeList.Call(0, 1, 0, uintptr(unsafe.Pointer(&attrListSize)))

	// Allocate memory to represent the list
	p.attrListBuffer = make([]byte, attrListSize)

	//Set the location in StartupInfoEx
	p.SIX.lpAttributeList = windows.Handle(unsafe.Pointer(&p.attrListBuffer[0]))

	// Initialize the list memory location
	r, _, e = initializeProcThreadAttributeList.Call(uintptr(p.SIX.lpAttributeList), 1, 0, uintptr(unsafe.Pointer(&attrListSize)))
	if r == 0 {
		return nil, fmt.Errorf("initializeProcThreadAttributeList Error:%v Code: 0x%x", e, r)
	}

	// Set the pseudoconsole information into the list
	r, _, e = updateProcThreadAttribute.Call(uintptr(p.SIX.lpAttributeList), 0, procThreadAttributePseudoconsole, uintptr(*p.phPC), unsafe.Sizeof(*p.phPC), 0, 0)
	if r == 0 {
		return nil, fmt.Errorf("updateProcThreadAttribute Error:%v Code: 0x%x", e, r)
	}
	return p, nil
}

//Resize is used to resize the buffer allocated to PTY
func (p *PTY) Resize(row, column int16) error {
	r, _, e := procResizePseudoConsole.Call(uintptr(*p.phPC), getSize(row, column))
	if r == 0 {
		return fmt.Errorf("procResizePseudoConsole Error:%v Code: 0x%x", e, r)
	}
	return nil
}

//Close is used to release all resources allocated to PTY
func (p *PTY) Close() error {

	var errs []string

	r, _, e := deleteProcThreadAttributeList.Call(uintptr(p.SIX.lpAttributeList))
	if r == 0 {
		errs = append(errs, fmt.Sprintf("closePseudoConsole Error:%v Code: 0x%x", e, r))
	}

	r, _, e = closePseudoConsole.Call(uintptr(*p.phPC))
	if r == 0 {
		errs = append(errs, fmt.Sprintf("closePseudoConsole Error:%v Code: 0x%x", e, r))
	}

	p.Input.Close()
	p.Output.Close()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, " "))
	}
	return nil
}

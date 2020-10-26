package process

import (
	"os"

	"github.com/iamalsaher/interactor/pkg/pty"
)

func (p *Process) osStart() (e error) {

	if p.pty != nil {
		p.proc, e = newConPTYProcess(p.Details.Path, append([]string{p.Details.Path}, p.Details.Args...), p.Details.Dir, p.Details.Env, &p.pty.SIX)
	} else {
		p.proc, e = os.StartProcess(p.Details.Path, append([]string{p.Details.Path}, p.Details.Args...), &os.ProcAttr{
			Dir:   p.Details.Dir,
			Env:   p.Details.Env,
			Sys:   nil,
			Files: []*os.File{p.stdin, p.stdout, p.stderr},
		})
	}

	return e
}

func setPtyIO(p *Process) error {

	// stdinHandle := syscall.Stdin
	// stdoutHandle := syscall.Stdout
	// stderrHandle := syscall.Stderr

	// stdinSafeMode, e := getConsoleMode(stdinHandle)

	// if e != nil {
	// 	return e
	// }

	// if e := setConsoleMode(uintptr(stdoutHandle), uint32(enableProcessedOutput)); e != nil {
	// 	return e
	// }

	// stdinSafeMode, e := getConsoleMode(stdoutHandle)
	// fmt.Println(stdinSafeMode)

	// r, _, e := setConsoleModeProc.Call(uintptr(syscall.Stdin), 0, 0)
	// if r == 0 {
	// 	return fmt.Errorf("setConsoleMode Input handle Error:%v Code: 0x%x", e, r)
	// }

	// r, _, e = setConsoleModeProc.Call(uintptr(syscall.Stdout), 0, 0)
	// if r == 0 {
	// 	return fmt.Errorf("setConsoleMode Output handle Error:%v Code: 0x%x", e, r)
	// }

	if pty, e := pty.NewPTY(); e == nil {
		p.Stdin = pty.Input
		p.Stdout = pty.Output
		p.pty = pty
		p.closers = append(p.closers, pty)
	}
	return nil
}

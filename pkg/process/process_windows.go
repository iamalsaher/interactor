package process

import (
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
	"golang.org/x/sys/windows"
)

//Start is used to finally Start the process
func (p *Process) Start(i *Interactor) (e error) {

	if p.Pipe != nil && i != nil {
		go i.Function(i.Input, i.Output)
	}

	if p.ptys[0] != nil && (windows.StartupInfo{}) != p.ptys[0].SIX.StartupInfo {
		p.proc, e = newConPTYProcess(p.details.path, p.details.args, p.details.rundir, p.details.env, &p.ptys[0].SIX)
	} else {
		p.proc, e = os.StartProcess(p.details.path, append([]string{p.details.path}, p.details.args...), &os.ProcAttr{
			Dir:   p.details.rundir,
			Env:   p.details.env,
			Sys:   nil,
			Files: []*os.File{p.Pipe.stdinR, p.Pipe.stdoutW, p.Pipe.stderrW},
		})
	}

	if e == nil {
		p.PID = p.proc.Pid

		//Setup timeout function
		if p.details.timeout > 0 {
			time.AfterFunc(p.details.timeout, func() {
				p.Kill()
			})
		}
	}
	return e
}

/*
ConnectIO is used to connect the input output handles
It attempts PTY as a primary mechanism else falls back to Pipes
If forcePTY is set then function errors out if pty cannot be aquired
*/
func (p *Process) ConnectIO(forcePTY bool) error {

	p.Pipe = new(Pipes)

	if pty, err := pty.NewPTY(); err == nil {
		p.Pipe.StdinW = pty.Input
		p.Pipe.StdoutR = pty.Output
		p.Pipe.StderrR = pty.Output
		p.ptys = append(p.ptys, pty)
		return nil
	} else if forcePTY {
		p.Pipe = nil
		return err
	}
	return setPipeIO(p)
}

package process

import (
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
	"golang.org/x/sys/windows"
)

//Start is used to finally Start the process
func (p *Process) Start() (e error) {

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

func setPtyIO(p *Process) error {
	pty, e := pty.NewPTY()
	if e == nil {
		p.Pipe = new(Pipes)
		p.Pipe.StdinW = pty.Input
		p.Pipe.StdoutR = pty.Output
		p.Pipe.StderrR = pty.Output
		p.ptys = append(p.ptys, pty)
	}
	return e
}

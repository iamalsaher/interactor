package process

import (
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//Process contains all the info about the Process
type Process struct {
	details *Details
	Pipe    *Pipes
	ptys    []*pty.PTY
	proc    *os.Process
	io      bool
	PID     int
}

//Start is used to finally Start the process
func (p *Process) Start(i *Interactor) (e error) {

	if p.io && i != nil {
		go i.Function(i.Input, i.Output)
	}

	p.proc, e = os.StartProcess(p.details.path, append([]string{p.details.path}, p.details.args...), &os.ProcAttr{
		Dir:   p.details.rundir,
		Env:   p.details.env,
		Sys:   nil,
		Files: []*os.File{p.Pipe.stdinR, p.Pipe.stdoutW, p.Pipe.stderrW},
	})

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
If errPTY is set then function also attempts to assign a pty to stderr
If forcePTY is set then function errors out if pty cannot be aquired
*/
func (p *Process) ConnectIO(errPTY, forcePTY bool) error {
	p.Pipe = new(Pipes)

	in, err := pty.NewPTY()
	if err != nil {
		if forcePTY {
			p.Pipe = nil
			return err
		}
		return setPipeIO(p)
	}

	out, err := pty.NewPTY()
	if err != nil {
		if forcePTY {
			p.Pipe = nil
			return err
		}
		return setPipeIO(p)
	}

	p.Pipe.stdinR = in.Slave
	p.Pipe.StdinW = in.Master
	p.Pipe.StdoutR = out.Master
	p.Pipe.stdoutW = out.Slave
	p.Pipe.StderrR = out.Master
	p.Pipe.stderrW = out.Slave

	p.io = true
	return nil
}

//Kill is a wrapper around os.Process.Kill()
func (p *Process) Kill() error {
	if len(p.ptys) > 0 {
		for _, pty := range p.ptys {
			//This returns error,so have a look at what to do with error
			pty.Close()
		}
	}
	return p.proc.Kill()
}

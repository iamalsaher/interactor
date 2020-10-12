package process

import (
	"os"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//Process contains all the info about the Process
type Process struct {
	details *Details
	pipe    *Pipes
	ptys    []*pty.PTY
	proc    *os.Process
	PID     int
}

//Start is used to finally Start the process
func (p *Process) Start() (e error) {
	p.proc, e = os.StartProcess(p.details.path, p.details.args, &os.ProcAttr{
		Dir:   p.details.rundir,
		Env:   p.details.env,
		Sys:   nil,
		Files: []*os.File{p.pipe.StdinR, p.pipe.StdoutW, p.pipe.StderrW},
	})
	return e
}

/*
ConnectIO is used to connect the input output handles
It attempts PTY as a primary mechanism else falls back to Pipes
If errPTY is set then function also attempts to assign a pty to stderr
If forcePTY is set then function errors out if pty cannot be aquired
*/
func (p *Process) ConnectIO(errPTY, forcePTY bool) error {
	p.pipe = new(Pipes)

	in, err := pty.NewPTY()
	if err != nil {
		if forcePTY {
			p.pipe = nil
			return err
		}
		return setPipeIO(p)
	}

	out, err := pty.NewPTY()
	if err != nil {
		if forcePTY {
			p.pipe = nil
			return err
		}
		return setPipeIO(p)
	}

	p.pipe.StdinR = in.Slave
	p.pipe.StdinW = in.Master
	p.pipe.StdoutR = out.Master
	p.pipe.StdoutR = out.Slave

	if errPTY {
		errout, err := pty.NewPTY()
		if err != nil {
			if forcePTY {
				p.pipe = nil
				return err
			}
			return setPipeIO(p)
		}
		p.pipe.StderrR = errout.Master
		p.pipe.StderrR = errout.Slave

	} else {
		p.pipe.StderrR = out.Master
		p.pipe.StderrR = out.Slave
	}

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

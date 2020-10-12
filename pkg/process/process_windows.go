package process

import (
	"os"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//Process contains all the info about the Process
type Process struct {
	details *Details
	pty     *pty.PTY
	pipe    *Pipes
	proc    *os.Process
	PID     int
}

//Start is used to finally Start the process
func (p *Process) Start() (e error) {
	if p.proc, e = newConPTYProcess(p.details.path, p.details.args, p.details.rundir, p.details.env, &p.pty.SIX); e == nil {
		p.PID = p.proc.Pid
	}
	return e
}

func (p *Process) GetIO() (*os.File, *os.File) {
	return p.pty.Master, p.pty.Slave
}

/*
ConnectIO is used to connect the input output handles
It attempts PTY as a primary mechanism else falls back to Pipes
If forcePTY is set then function errors out if pty cannot be aquired
*/
func (p *Process) ConnectIO(forcePTY bool) error {

	if pty, err := pty.NewPTY(); err == nil {
		p.pty = pty
		return nil
	} else if forcePTY {
		return err
	}
	return setPipeIO(p)
}

//Kill is a wrapper around os.Process.Kill()
func (p *Process) Kill() error {
	if p.pty != nil {
		p.pty.Close()
	}
	return p.proc.Kill()
}

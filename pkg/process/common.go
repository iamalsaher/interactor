package process

import (
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//NewProcess is used to setup details about new process
func NewProcess(path string, args ...string) {

	details := new(Details)
	details.path = path
	details.args = args

	proc := new(Process)
	proc.details = details
}

//SetEnviron is used to add environment variables to the process
//Set def to true to add default environment variables along with the ones specified in env
//If def is False only the ones define in env shall be set
func (p *Process) SetEnviron(env []string, def bool) {
	if def {
		p.details.env = os.Environ()
	}
	p.details.env = append(p.details.env, env...)
}

//SetTimeout is used to add process timeout in milliseconds
func (p *Process) SetTimeout(timeout int64) {
	p.details.timeout = time.Duration(timeout) * time.Millisecond
}

//SetDirectory is used to set directory in which the process should run
func (p *Process) SetDirectory(dir string) {
	p.details.rundir = dir
}

//ConnectIO is used to connect stdin, stdout and stderr
func (p *Process) ConnectIO(forcePTY bool) error {

	if pty, err := pty.NewPTY(); err == nil {
		p.pty = pty
		return nil
	} else if forcePTY {
		return err
	}
	return SetPipeIO(p)
}

//SetPipeIO is used to set OS PIPE based input output
func SetPipeIO(p *Process) error { return nil }

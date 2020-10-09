package process

import (
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//NewProcess is used to setup details about new process
func NewProcess(path string, args ...string) *Process {

	details := &Details{path: path, args: args, env: os.Environ()}
	return &Process{details: details}
}

//SetEnviron is used to add environment variables to the process
//Set keepDefaultEnvironment to false to remove default environment variables
//If keepDefaultEnvironment is False only the ones define in env shall be set
func (p *Process) SetEnviron(env []string, keepDefaultEnvironment bool) {
	if !keepDefaultEnvironment {
		p.details.env = nil
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
//If forcePTY is set then function errors out if pty cannot be aquired
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

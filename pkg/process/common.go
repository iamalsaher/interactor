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

//Release is a wrapper around os.Process.Release()
func (p *Process) Release() error {
	return p.proc.Release()
}

//Signal is a wrapper around os.Process.Signal()
func (p *Process) Signal(sig os.Signal) error {
	return p.proc.Signal(sig)
}

//Wait is a wrapper around os.Process.Wait()
func (p *Process) Wait() (*os.ProcessState, error) {
	return p.proc.Wait()
}

func setPipeIO(p *Process) error { return nil }

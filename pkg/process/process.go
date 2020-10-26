package process

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//Process contains all the info about the Process
type Process struct {
	Details *Details
	PID     int

	Stdin  io.Writer
	Stdout io.Reader
	Stderr io.Reader

	State *os.ProcessState
	Done  chan bool

	proc    *os.Process
	stdin   *os.File
	stdout  *os.File
	stderr  *os.File
	pty     *pty.PTY
	closers []io.Closer
}

//Details is used to store various details about the process
type Details struct {
	Path    string
	Args    []string
	Env     []string
	Dir     string
	Timeout time.Duration
}

//NewProcess is used to setup details about new process
func NewProcess(path string, args ...string) *Process {

	details := &Details{
		Path: path,
		Args: args,
		Env:  os.Environ(),
	}
	return &Process{Details: details}
}

//SetEnviron is used to add environment variables to the process
//Set keepDefaultEnvironment to false to remove default environment variables
//If keepDefaultEnvironment is False only the ones define in env shall be set
func (p *Process) SetEnviron(env []string, keepDefaultEnvironment bool) {
	if !keepDefaultEnvironment {
		p.Details.Env = nil
	}
	p.Details.Env = append(p.Details.Env, env...)
}

//SetTimeout is used to add process timeout in milliseconds
func (p *Process) SetTimeout(t time.Duration) {
	p.Details.Timeout = t
}

//SetDirectory is used to set directory in which the process should run
func (p *Process) SetDirectory(dir string) {
	p.Details.Dir = dir
}

//Start is used to finally Start the process
func (p *Process) Start() (e error) {
	if e = p.osStart(); e == nil {

		p.PID = p.proc.Pid

		//Setup timeout function
		if p.Details.Timeout > 0 {
			time.AfterFunc(p.Details.Timeout, func() {
				p.Kill()
			})
		}

		p.Done = make(chan bool)

		go func() {
			var wErr error
			if p.State, wErr = p.Wait(); wErr != nil {
				panic(fmt.Sprintf("Error waiting for Process: %v", e))
			}
			close(p.Done)
		}()
	}
	return
}

//Release is a wrapper around os.Process.Release()
func (p *Process) Release() error {
	return p.proc.Release()
}

//Signal is a wrapper around os.Process.Signal()
func (p *Process) Signal(sig os.Signal) error {
	return p.proc.Signal(sig)
}

func (p *Process) closeAll() {
	for _, i := range p.closers {
		i.Close()
	}
	p.closers = nil
}

//Wait is a wrapper around os.Process.Wait()
func (p *Process) Wait() (s *os.ProcessState, e error) {
	s, e = p.proc.Wait()
	p.closeAll()
	return
}

//Kill is a wrapper around os.Process.Kill()
func (p *Process) Kill() error {
	return p.proc.Kill()
}

/*
ConnectIO is used to connect the input output handles
It attempts PTY as a primary mechanism else falls back to Pipes
If forcePTY is set then function errors out if pty cannot be aquired
*/
func (p *Process) ConnectIO(forcePTY bool) error {

	if e := setPtyIO(p); e == nil || forcePTY {
		return e
	}

	return setPipeIO(p)
}

func setPipeIO(p *Process) error {

	rin, win, ein := os.Pipe()
	rout, wout, eout := os.Pipe()
	rerr, werr, eerr := os.Pipe()

	if ein != nil || eout != nil || eerr != nil {
		return fmt.Errorf("Cannot aquire pipes for IO, inputErr: %v, outputErr: %v", ein, eout)
	}

	p.closers = append(p.closers, rin, win, rout, wout, rerr, werr)

	p.stdin = rin
	p.stdout = wout
	p.stderr = werr

	p.Stdin = win
	p.Stdout = rout
	p.Stderr = rerr

	return nil
}

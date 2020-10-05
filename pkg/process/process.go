package process

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//NewProcess is used to setup details about new process
func NewProcess(name string, args ...string) (*Process, string) {

	logmsg := "File found in PATH"
	path, lookPathErr := exec.LookPath(name)
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		cwd = fmt.Sprintf("Error while determining working directory: %v", cwdErr.Error())
	}

	if lookPathErr != nil {
		path, lookPathErr = exec.LookPath(filepath.FromSlash("./" + name))
		if lookPathErr != nil {
			return nil, fmt.Sprintf("File not present in PATH and %v", cwd)
		}
		logmsg = fmt.Sprintf("File found in %v", cwd)
	}

	details := new(Details)
	details.path = path
	details.args = args

	proc := new(Process)
	proc.details = details

	return proc, logmsg
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
func (p *Process) ConnectIO() error {

	stdin, err := pty.NewPTY()
	if err != nil {
		return err
	}

	stdout, err := pty.NewPTY()
	if err != nil {
		return err
	}

	stderr, err := pty.NewPTY()
	if err != nil {
		return err
	}

	p.pipes = new(Pipes)
	p.output = new(bytes.Buffer)
	p.errors = new(bytes.Buffer)

	p.details.stdin = stdin.Master
	p.details.stdout = stdout.Slave
	p.details.stderr = stderr.Slave

	p.pipes.stdin = stdin.Slave
	p.pipes.stdout = stdout.Master
	p.pipes.stderr = stderr.Master

	return nil
}

//Start is used to finally Start the process
func (p *Process) Start() error {

	h, e := os.StartProcess(p.details.path, p.details.args, &os.ProcAttr{
		Dir:   p.details.rundir,
		Env:   p.details.env,
		Sys:   nil,
		Files: []*os.File{p.details.stdin, p.details.stdout, p.details.stderr},
	})

	//Setup timeout function
	if e == nil && p.details.timeout > 0 {
		time.AfterFunc(p.details.timeout, func() {
			h.Kill()
		})
	}

	//Start Stdout and Stderr capture
	if p.pipes != nil {
		go copyToBufferFromFile(p.output, p.pipes.stdout)
		go copyToBufferFromFile(p.errors, p.pipes.stderr)
	}

	p.Handle = h
	return e
}

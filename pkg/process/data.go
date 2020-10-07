package process

import (
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//Details is used to store various details about the process
type Details struct {
	path    string
	args    []string
	env     []string
	rundir  string
	timeout time.Duration
	stdin   *os.File
	stdout  *os.File
	stderr  *os.File
}

//Pipes defines the pipes for stdin, stdout and stderr
type Pipes struct {
}

//Process contains all the info about the Process
type Process struct {
	details *Details
	pty     *pty.PTY
	pipe    *Pipes
}

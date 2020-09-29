package process

import (
	"os"
	"time"
)

//IO defines the pipes for stdin, stdout and stderr
type IO struct {
	stdinR  *os.File
	stdinW  *os.File
	stdoutR *os.File
	stdoutW *os.File
	stderrR *os.File
	stderrW *os.File
}

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

//Process contains all the info about the Process
type Process struct {
	details *Details
	io      *IO
	Handle  *os.Process
}

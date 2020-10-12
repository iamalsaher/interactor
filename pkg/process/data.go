package process

import (
	"os"
	"time"
)

//Details is used to store various details about the process
type Details struct {
	path    string
	args    []string
	env     []string
	rundir  string
	timeout time.Duration
}

//Pipes defines the pipes for stdin, stdout and stderr
type Pipes struct {
	stdinR *os.File
	StdinW *os.File

	StdoutR *os.File
	stdoutW *os.File

	StderrR *os.File
	stderrW *os.File
}

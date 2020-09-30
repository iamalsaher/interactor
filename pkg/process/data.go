package process

import (
	"bytes"
	"os"
	"time"
)

//Pipes defines the pipes for stdin, stdout and stderr
type Pipes struct {
	stdin  *os.File
	stdout *os.File
	stderr *os.File
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
	pipes   *Pipes
	Handle  *os.Process
	output  *bytes.Buffer
	errors  *bytes.Buffer
}

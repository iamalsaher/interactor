package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//Details is used to store various details about the process
type Details struct {
	path    string
	args    []string
	env     []string
	rundir  string
	timeout int64
	stdin   *os.File
	stdout  *os.File
	stderr  *os.File
}

//NewProcess is used to setup details about new process
func NewProcess(name string, args ...string) (*Details, string) {

	logmsg := "File found in PATH"
	path, lookPathErr := exec.LookPath(name)

	if lookPathErr != nil {
		path, lookPathErr = exec.LookPath(filepath.FromSlash("./" + name))
		if lookPathErr != nil {
			return nil, fmt.Sprintf("File not present in PATH and %v", os.Getwd())
		}
		logmsg = fmt.Sprintf("File found in %v", os.Getwd())
	}

	details := new(Details)
	details.path = path
	details.args = args

	return *details, logmsg
}

//SetEnviron is used to add environment variables to the process
//Set def to true to add default environment variables along with the ones specified in env
//If def is False only the ones define in env shall be set
func (d *Details) SetEnviron(env []string, def bool) {
	if def {
		d.env = os.Environ()
	}
	d.env = append(d.env, env)
}

//SetTimeout is used to add process timeout
func (d *Details) SetTimeout(timeout int64) {
	d.timeout = timeout
}

//SetDirectory is used to set directory in which the process should run
func (d *Details) SetDirectory(dir string) {
	d.rundir = dir
}

//SetIO is used to set stdin, stdout and stderr
func (d *Details) SetIO(stdin, stdout, stderr *os.File) {
	d.stdin = stdin
	d.stdout = stdout
	d.stderr = stderr
}

//Start is used to start the process
func (d *Details) Start() {
	attr = os.ProcAttr{
		Dir:   d.rundir,
		Env:   d.env,
		Sys:   nil,
		Files: []*os.File{d.stdin, d.stdout, d.stderr},
	}
	os.StartProcess(d.path, d.args)
}

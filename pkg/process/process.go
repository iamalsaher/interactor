package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

//Process is used to setup details about new process
func Process(name string, args ...string) (*Details, string) {

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

	return details, logmsg
}

//SetEnviron is used to add environment variables to the process
//Set def to true to add default environment variables along with the ones specified in env
//If def is False only the ones define in env shall be set
func (d *Details) SetEnviron(env []string, def bool) {
	if def {
		d.env = os.Environ()
	}
	d.env = append(d.env, env...)
}

//SetTimeout is used to add process timeout in milliseconds
func (d *Details) SetTimeout(timeout int64) {
	d.timeout = time.Duration(timeout) * time.Millisecond
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

//Start is used to finally Start the process
func (d *Details) Start() (*os.Process, error) {

	h, e := os.StartProcess(d.path, d.args, &os.ProcAttr{
		Dir:   d.rundir,
		Env:   d.env,
		Sys:   nil,
		Files: []*os.File{d.stdin, d.stdout, d.stderr},
	})

	if e != nil && d.timeout > 0 {
		time.AfterFunc(d.timeout, func() {
			h.Kill()
		})
	}

	return h, e
}

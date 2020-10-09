package process

import "os"

//Start is used to finally Start the process
func (p *Process) Start() error {
	var e error
	p.Proc, e = newConPTYProcess(p.details.path, p.details.args, p.details.rundir, p.details.env, &p.pty.SIX)
	return e
}

func (p *Process) GetIO() (*os.File, *os.File) {
	return p.pty.Master, p.pty.Slave
}

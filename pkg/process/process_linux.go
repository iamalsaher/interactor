package process

import (
	"fmt"
	"os"

	"github.com/iamalsaher/interactor/pkg/pty"
)

func (p *Process) osStart() (e error) {

	p.proc, e = os.StartProcess(p.Details.Path, append([]string{p.Details.Path}, p.Details.Args...), &os.ProcAttr{
		Dir:   p.Details.Dir,
		Env:   p.Details.Env,
		Sys:   nil,
		Files: []*os.File{p.stdin, p.stdout, p.stderr},
	})

	return e
}

func setPtyIO(p *Process, setStderr bool) error {

	ttyIOE, e := pty.NewPTY()
	if e == nil {

		p.stdin = ttyIOE.Slave
		p.stdout = ttyIOE.Slave
		p.stderr = ttyIOE.Slave

		p.Stdin = ttyIOE.Master
		p.Stdout = ttyIOE.Master
		p.Stderr = ttyIOE.Master

		p.closers = append(p.closers, ttyIOE)

		if !setStderr {
			return nil
		}
	}

	ttyE, e := pty.NewPTY()
	if e != nil {
		p.closeAll()
		return fmt.Errorf("Cannot acquire stderr pty: Error: %v", e)
	}

	p.stderr = ttyE.Slave
	p.Stderr = ttyE.Master
	p.closers = append(p.closers, ttyE)

	return nil
}

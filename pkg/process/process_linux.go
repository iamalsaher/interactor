package process

import (
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

func setPtyIO(p *Process) error {

	ttyIE, eIE := pty.NewPTY()
	if eIE != nil {
		return eIE
	}

	ttyO, eO := pty.NewPTY()
	if eO != nil {
		ttyO.Close()
		return eO
	}

	p.stdin = ttyIE.Slave
	p.stdout = ttyO.Slave
	p.stderr = ttyIE.Slave

	p.Stdin = ttyIE.Master
	p.Stdout = ttyO.Master
	p.Stderr = ttyIE.Master

	p.closers = append(p.closers, ttyIE, ttyO)

	return nil
}

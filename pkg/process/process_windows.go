package process

import (
	"os"

	"github.com/iamalsaher/interactor/pkg/pty"
)

func (p *Process) osStart() (e error) {

	if p.pty != nil {
		p.proc, e = newConPTYProcess(p.Details.Path, append([]string{p.Details.Path}, p.Details.Args...), p.Details.Dir, p.Details.Env, &p.pty.SIX)
	} else {
		p.proc, e = os.StartProcess(p.Details.Path, append([]string{p.Details.Path}, p.Details.Args...), &os.ProcAttr{
			Dir:   p.Details.Dir,
			Env:   p.Details.Env,
			Sys:   nil,
			Files: []*os.File{p.stdin, p.stdout, p.stderr},
		})
	}

	return e
}

func setPtyIO(p *Process, setStderr bool) error {
	pty, e := pty.NewPTY()
	if e == nil {
		p.Stdin = pty.Input
		p.Stdout = pty.Output
		p.pty = pty
		p.closers = append(p.closers, pty)
	}
	return e
}

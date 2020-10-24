package process

import (
	"os"
	"time"

	"github.com/iamalsaher/interactor/pkg/pty"
)

//Start is used to finally Start the process
func (p *Process) Start() (e error) {

	p.proc, e = os.StartProcess(p.details.path, append([]string{p.details.path}, p.details.args...), &os.ProcAttr{
		Dir:   p.details.rundir,
		Env:   p.details.env,
		Sys:   nil,
		Files: []*os.File{p.Pipe.stdinR, p.Pipe.stdoutW, p.Pipe.stderrW},
	})

	if e == nil {
		p.PID = p.proc.Pid

		//Setup timeout function
		if p.details.timeout > 0 {
			time.AfterFunc(p.details.timeout, func() {
				p.Kill()
			})
		}
	}
	return e
}

func setPtyIO(p *Process) error {

	pty, e := pty.NewPTY()
	if e == nil {
		p.Pipe = new(Pipes)
		p.Pipe.StdinW = pty.Master
		p.Pipe.StdoutR = pty.Master
		p.Pipe.StderrR = pty.Master
		p.ptys = append(p.ptys, pty)

		p.Pipe.stdinR = pty.Slave
		p.Pipe.stdoutW = pty.Slave
		p.Pipe.stderrW = pty.Slave
	}
	return e
}

// func setPtyIO(p *Process) error {
// 	p.Pipe = new(Pipes)

// 	in, inErr := pty.NewPTY()
// 	out, outErr := pty.NewPTY()

// 	if inErr != nil || outErr != nil {
// 		var errs []string
// 		if inErr != nil {
// 			errs = append(errs, "Input PTY: "+se.Error())
// 		}
// 		if me != nil {
// 			errs = append(errs, "Master: "+me.Error())
// 		}
// 		return errors.New(strings.Join(errs, " "))
// 	}

// 	out, err := pty.NewPTY()
// 	if err != nil {
// 		if forcePTY {
// 			p.Pipe = nil
// 			return err
// 		}
// 		return setPipeIO(p)
// 	}

// 	p.Pipe.stdinR = in.Slave
// 	p.Pipe.StdinW = in.Master
// 	p.Pipe.StdoutR = out.Master
// 	p.Pipe.stdoutW = out.Slave
// 	p.Pipe.StderrR = out.Master
// 	p.Pipe.stderrW = out.Slave
// 	p.ptys = append(p.ptys, in, out)

// 	return nil
// }

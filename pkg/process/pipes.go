package process

import (
	"io"
	"os"
)

func getIO(stdin, stdout, stderr *os.File) (*IO, error) {
	io := new(IO)

	if stdin != nil {
		if r, w, e := os.Pipe(); e == nil {
			io.stdinR = r
			io.stdinW = w
		} else {
			return nil, e
		}
	}

	if stdout != nil {
		if r, w, e := os.Pipe(); e == nil {
			io.stdoutR = r
			io.stdoutW = w
		} else {
			return nil, e
		}
	}

	if stderr != nil {
		if r, w, e := os.Pipe(); e == nil {
			io.stderrR = r
			io.stderrW = w
		} else {
			return nil, e
		}
	}
	return io, nil
}

//Recv is used to receive size bytes from stdout
func (p *Process) Recv(size int) ([]byte, error) {
	b := make([]byte, 0, size)
	_, e := io.ReadFull(p.io.stdoutR, b)
	return b, e
}

//Send is used to send byte array to stdin
func (p *Process) Send(b []byte) (int, error) {
	return p.io.stdinW.Write(b)
}

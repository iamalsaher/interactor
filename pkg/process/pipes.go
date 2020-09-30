package process

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

//Recv is used to receive size bytes from stdout
func (p *Process) Recv(size int) ([]byte, error) {
	b := make([]byte, 0, size)
	_, e := p.pipes.stdout.Read(b)
	return b, e
}

//Send is used to send byte array to stdin
func (p *Process) Send(b []byte) (int, error) {
	return p.pipes.stdin.Write(b)
}

func copyToBufferFromFile(buff *bytes.Buffer, pipe *os.File) {
	c := make([]byte, 1)
	for {
		if _, e := pipe.Read(c); e == nil {
			buff.WriteByte(c[0])
		} else {
			log.Println(e.Error())
			break
		}
	}
}

//Show is used to show the output
func (p *Process) Show() {
	fmt.Println(p.output)
}

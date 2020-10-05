package pty

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	tiocgptn   = 0x80045430 // TIOCGPTN IOCTL used to get the PTY number
	tiocsptlck = 0x40045431 // TIOCSPTLCK IOCT used to lock/unlock PTY
)

//PTY is the structure which contains the input and output files to pass to processes
type PTY struct {
	Master *os.File
	Slave  *os.File
}

func (p *PTY) getPTSNumber() (uint, error) {
	var n uint
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(p.Master.Fd()), uintptr(tiocgptn), uintptr(unsafe.Pointer(&n)))
	if e != 0 {
		return 0, e
	}
	return n, nil
}

func (p *PTY) getPTSName() (string, error) {
	n, err := p.getPTSNumber()
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}

func (p *PTY) Close() error {
	se := fmt.Errorf("Slave FD nil")
	if p.Slave != nil {
		se = p.Slave.Close()
	}

	me := fmt.Errorf("Master FD nil")
	if p.Slave != nil {
		se = p.Slave.Close()
	}

	if se != nil || me != nil {
		var errs []string
		if se != nil {
			errs = append(errs, "Slave: "+se.Error())
		}
		if me != nil {
			errs = append(errs, "Master: "+me.Error())
		}
		return errors.New(strings.Join(errs, " "))
	}
	return nil
}

//NewPTY returns a master and slave file descriptor for a pty
//https://github.com/google/goterm/blob/master/term/termios.go
func NewPTY() (*PTY, error) {
	/*
		Opening /dev/ptmx gives master fd
		https://code.woboq.org/userspace/glibc/login/openpty.c.html#86
	*/
	master, e := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if e != nil {
		return nil, e
	}

	var unlock int
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(master.Fd()), uintptr(tiocsptlck), uintptr(unsafe.Pointer(&unlock))); errno != 0 {
		master.Close()
		return nil, errno
	}

	p := &PTY{Master: master}
	s, err := p.getPTSName()
	if err != nil {
		master.Close()
		return nil, err
	}

	// open pty slave
	p.Slave, err = os.OpenFile(s, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		master.Close()
		return nil, err
	}
	return p, nil
}

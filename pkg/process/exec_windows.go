package process

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/iamalsaher/interactor/pkg/pty"
	"golang.org/x/sys/windows"
)

//Windows Console constants
const (
	enableProcessedInput       = 0x0001
	enableLineInput            = 0x0002
	enableEchoInput            = 0x0004
	enableWindowInput          = 0x0008
	enableMouseInput           = 0x0010
	enableInsertMode           = 0x0020
	enableQuickEditMode        = 0x0040
	enableExtendedFlags        = 0x0080
	enableAutoPosition         = 0x0100
	enableVirtualTerminalInput = 0x0200

	enableProcessedOutput           = 0x0001
	enableWrapAtEolOutput           = 0x0002
	enableVirtualTerminalProcessing = 0x0004
	disableNewlineAutoReturn        = 0x0008
	enableLvbGridWorldwide          = 0x0010
)

const extendedStartupinfoPresent uint32 = 0x00080000

var (
	kernel32DLL = syscall.NewLazyDLL("kernel32.dll")
)

// makeCmdLine builds a command line out of args by escaping "special"
// characters and joining the arguments with spaces.
func makeCmdLine(args []string) string {
	var s string
	for _, v := range args {
		if s != "" {
			s += " "
		}
		s += windows.EscapeArg(v)
	}
	return s
}

func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

func normalizeDir(dir string) (name string, err error) {
	ndir, err := syscall.FullPath(dir)
	if err != nil {
		return "", err
	}
	if len(ndir) > 2 && isSlash(ndir[0]) && isSlash(ndir[1]) {
		// dir cannot have \\server\share\path form
		return "", syscall.EINVAL
	}
	return ndir, nil
}

func volToUpper(ch int) int {
	if 'a' <= ch && ch <= 'z' {
		ch += 'A' - 'a'
	}
	return ch
}

func joinExeDirAndFName(dir, p string) (name string, err error) {
	if len(p) == 0 {
		return "", syscall.EINVAL
	}
	if len(p) > 2 && isSlash(p[0]) && isSlash(p[1]) {
		// \\server\share\path form
		return p, nil
	}
	if len(p) > 1 && p[1] == ':' {
		// has drive letter
		if len(p) == 2 {
			return "", syscall.EINVAL
		}
		if isSlash(p[2]) {
			return p, nil
		}
		d, err := normalizeDir(dir)
		if err != nil {
			return "", err
		}
		if volToUpper(int(p[0])) == volToUpper(int(d[0])) {
			return syscall.FullPath(d + "\\" + p[2:])
		}
		return syscall.FullPath(p)

	}
	// no drive letter
	d, err := normalizeDir(dir)
	if err != nil {
		return "", err
	}
	if isSlash(p[0]) {
		return windows.FullPath(d[:2] + p)
	}
	return windows.FullPath(d + "\\" + p)

}

// createEnvBlock converts an array of environment strings into
// the representation required by CreateProcess: a sequence of NUL
// terminated strings followed by a nil.
// Last bytes are two UCS-2 NULs, or four NUL bytes.
func createEnvBlock(envv []string) *uint16 {
	if len(envv) == 0 {
		return &utf16.Encode([]rune("\x00\x00"))[0]
	}
	length := 0
	for _, s := range envv {
		length += len(s) + 1
	}
	length++

	b := make([]byte, length)
	i := 0
	for _, s := range envv {
		l := len(s)
		copy(b[i:i+l], []byte(s))
		copy(b[i+l:i+l+1], []byte{0})
		i = i + l + 1
	}
	copy(b[i:i+1], []byte{0})

	return &utf16.Encode([]rune(string(b)))[0]
}

// dedupEnvCase is dedupEnv with a case option for testing.
// If caseInsensitive is true, the case of keys is ignored.
func dedupEnvCase(caseInsensitive bool, env []string) []string {
	out := make([]string, 0, len(env))
	saw := make(map[string]int, len(env)) // key => index into out
	for _, kv := range env {
		eq := strings.Index(kv, "=")
		if eq < 0 {
			out = append(out, kv)
			continue
		}
		k := kv[:eq]
		if caseInsensitive {
			k = strings.ToLower(k)
		}
		if dupIdx, isDup := saw[k]; isDup {
			out[dupIdx] = kv
			continue
		}
		saw[k] = len(out)
		out = append(out, kv)
	}
	return out
}

// addCriticalEnv adds any critical environment variables that are required
// (or at least almost always required) on the operating system.
// Currently this is only used for Windows.
func addCriticalEnv(env []string) []string {
	for _, kv := range env {
		eq := strings.Index(kv, "=")
		if eq < 0 {
			continue
		}
		k := kv[:eq]
		if strings.EqualFold(k, "SYSTEMROOT") {
			// We already have it.
			return env
		}
	}
	return append(env, "SYSTEMROOT="+os.Getenv("SYSTEMROOT"))
}

/*
https://golang.org/src/os/exec.go
https://golang.org/src/os/exec_posix.go
https://golang.org/src/syscall/exec_windows.go
https://golang.org/src/os/exec.go
https://golang.org/src/os/exec_windows.go
https://stackoverflow.com/questions/17981651/is-there-any-way-to-access-private-fields-of-a-struct-from-another-package
*/

//NewConPTYProcess is used to start a conpty process
func newConPTYProcess(argv0 string, argv []string, dir string, env []string, six *pty.StartupInfoEx) (o *os.Process, e error) {
	p, h, e := createProcessWithConpty(argv0, argv, dir, env, six)
	if e != nil {
		return nil, e
	}
	o, e = os.FindProcess(p)
	syscall.CloseHandle(syscall.Handle(h))
	return
}

func createProcessWithConpty(argv0 string, argv []string, dir string, env []string, six *pty.StartupInfoEx) (pid int, handle uintptr, err error) {

	if len(argv0) == 0 {
		return 0, 0, syscall.EWINDOWS
	}

	var (
		cmdline = makeCmdLine(argv)
		argvp   *uint16
		dirp    *uint16
		zeroSec = &windows.SecurityAttributes{Length: uint32(unsafe.Sizeof(windows.SecurityAttributes{})), InheritHandle: 1}
	)

	if len(dir) != 0 {
		argv0, err = joinExeDirAndFName(dir, argv0)
		if err != nil {
			return 0, 0, err
		}

		dirp, err = windows.UTF16PtrFromString(dir)
		if err != nil {
			return 0, 0, err
		}
	}

	argv0p, err := windows.UTF16PtrFromString(argv0)
	if err != nil {
		return 0, 0, err
	}

	if len(cmdline) != 0 {
		argvp, err = windows.UTF16PtrFromString(cmdline)
		if err != nil {
			return 0, 0, err
		}
	}

	six.StartupInfo.Flags = windows.STARTF_USESTDHANDLES
	pi := new(windows.ProcessInformation)
	flags := uint32(windows.CREATE_UNICODE_ENVIRONMENT) | extendedStartupinfoPresent

	// fmt.Printf("argv0: %+v\nargv: %+v\npsec: %+v\ntsec: %+v\ninherithandles: %v\nflags: %+v\nenv: %+v\ndir: %v\nSI: %+v\nPI: %+v", argv0, argv, pSec, tSec, false, flags, createEnvBlock(addCriticalEnv(dedupEnvCase(true, env))), dirp, &six.StartupInfo, pi)

	err = windows.CreateProcess(
		argv0p,
		argvp,
		zeroSec, // process handle not inheritable
		zeroSec, // thread handles not inheritable,
		false,
		flags,
		createEnvBlock(addCriticalEnv(dedupEnvCase(true, env))),
		dirp, // use current directory later: dirp,
		&six.StartupInfo,
		pi)

	if err != nil {
		return 0, 0, err
	}

	defer windows.CloseHandle(windows.Handle(pi.Thread))
	return int(pi.ProcessId), uintptr(pi.Process), nil
}

func setConsoleMode(handle uintptr, mode uint32) error {
	r, _, e := kernel32DLL.NewProc("SetConsoleMode").Call(uintptr(syscall.Stdin), uintptr(mode), 0)
	if r == 0 {
		return fmt.Errorf("setConsoleMode Input handle Error:%v Code: 0x%x", e, r)
	}
	return nil
}

func getConsoleMode(handle syscall.Handle) (mode uint32, e error) {
	e = syscall.GetConsoleMode(handle, &mode)
	return
}

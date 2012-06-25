package ProcessManager

import (
	"errors"
	"net"
	"os"

//	"syscall"
//	"unsafe"
)

func Run() (err error) {
	const path = "/system/socket/ProcessSocket"

	// must run as root
	if os.Getuid() != 0 {
		return errors.New("must be run as root")
	}

	//	syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_NAME, uintptr(unsafe.Pointer(syscall.StringBytePtr("ProcessManager"))), 0)

	// check if file exists. if so, delete it.
	if _, err := os.Lstat(path); err == nil {
		os.Remove(path)
	}

	// resolve the address
	addr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return err
	}

	// listen on path
	listener, err := net.ListenUnix("unix", addr)
	if err != nil {
		return err
	}

	// set permission to 777
	os.Chmod(path, 0777)

	// create event handlers
	createEventHandlers()

	// loop for connections
	for {

		conn, err := listener.AcceptUnix()
		if err != nil {
			return err
		}

		go newConnection(conn).readData()

	}
	return
}

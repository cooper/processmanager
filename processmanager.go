package ProcessManager

import (
	"errors"
	"net"
	"os"
	"syscall"
	"time"

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

	// loop for ping checking
	go pingLoop()

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

// check for process ping replies
func pingLoop() {
	for _, conn := range connections {
		if conn.process == nil && time.Since(conn.connected).Seconds() >= 5 {

			// this connection has existed for five seconds and has not registered.
			conn.socket.Close()
			conn.process.Kill(syscall.SIGKILL)

		} else if conn.process != nil && time.Since(conn.lastPong).Seconds() >= 10 {

			// this connection has not responded to pings for a while.
			conn.socket.Close()
			conn.process.Kill(syscall.SIGKILL)

		} else {

			// this connection is doing well. ping it again.
			conn.send("ping", nil)

		}
	}
	time.Sleep(2)
	pingLoop()
}

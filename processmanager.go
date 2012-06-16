package ProcessManager

import (
	"net"
	"os"
)

func Run() (err error) {
	const path = "/system/socket/ProcessSocket"

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

	// create event handlers
	createEventHandlers()

	// loop for connections
	for {

		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go readData(newConnection(conn))

	}
	return
}

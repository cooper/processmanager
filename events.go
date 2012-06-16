package ProcessManager

import "process"

var eventHandlers = make(map[string]func(conn *connection, name string, params map[string]interface{}))

// assign handlers
func createEventHandlers() {
	eventHandlers["register"] = registerHandler
}

// creates a process object for the connected process.
func registerHandler(conn *connection, name string, params map[string]interface{}) {
	pid := params["pid"].(int)

	// this process is already registered...
	if connections[pid] != nil {
		conn.conn.Close()
		return
	}

	conn.process = process.FromPID(pid)

	// store for later
	connections[pid] = conn
}

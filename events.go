package ProcessManager

import "process"

var eventHandlers = make(map[string]func(conn *connection, name string, params map[string]interface{}))

// assign handlers
func createEventHandlers() {
	eventHandlers["register"] = registerHandler
}

// creates a process object for the connected process.
func registerHandler(conn *connection, name string, params map[string]interface{}) {
	pid := int(params["pid"].(float64))
	delete(params, "pid")

	// this process is already registered...
	if connections[pid] != nil {
		conn.socket.Close()
		return
	}

	conn.process = process.SFromPID(pid)

	// assign all of the properties here
	for prop, value := range params {
		conn.process.SetProperty(prop, value.(string))
	}

	// store for later
	connections[pid] = conn
}

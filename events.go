package ProcessManager

import "process"

var eventHandlers = make(map[string]func(conn *connection, name string, params map[string]interface{}))

// assign handlers
func createEventHandlers() {
	eventHandlers["register"] = registerHandler
}

// creates a process object for the connected process.
func registerHandler(conn *connection, name string, params map[string]interface{}) {
	pid := params["pid"].(float64)
	conn.process = process.FromPID(int(pid))
}

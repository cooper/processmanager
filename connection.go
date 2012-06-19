package ProcessManager

import (
	"bufio"
	"encoding/json"
	"net"
	"process"
)

var currentId int = 0
var connections = make(map[int]*connection)

type connection struct {
	socket   net.Conn
	incoming *bufio.Reader
	id       int
	process  *process.SProcess
}

// create a new connection
func newConnection(conn net.Conn) *connection {
	currentId++
	newconn := &connection{
		socket:   conn,
		incoming: bufio.NewReader(conn),
		id:       currentId,
	}
	return newconn
}

// read data from a connection
func (conn *connection) readData() {

	// loop
	for {
		line, _, err := conn.incoming.ReadLine()
		if err != nil {
			break
		}
		handleEvent(conn, line)
	}

	// disconnect
	conn.destroy()
}

// handle a JSON event
func handleEvent(conn *connection, data []byte) bool {
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return false
	}

	// should be an array.
	c := i.([]interface{}) // type assertion -- because an interface is a container

	command := c[0].(string)
	params := c[1].(map[string]interface{})

	// if a handler for this command exists, run it
	if eventHandlers[command] != nil {
		eventHandlers[command](conn, command, params)
	}

	return true
}

// send a JSON event
func (conn *connection) send(command string, params map[string]interface{}) bool {
	b, err := json.Marshal(params)
	if err != nil {
		return false
	}
	_, err = conn.socket.Write(b)
	if err != nil {
		return false
	}
	return true
}

// connection ends
func (conn *connection) destroy() {
	if conn.process != nil {
		delete(connections, conn.process.PID())
		process.Free(conn.process)
	}
}

package go_cluster

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

// Resembles a connection between two nodes.
// The API makes it simpler to customize the connection.
type Connection struct {
	Conn *net.TCPConn
}

// Creates a new connection
func connect(ip string, port int, rip string, rport int) (*Connection, error) {
	if conn, err := net.DialTCP("tcp", &net.TCPAddr{IP: []byte(ip), Port: port}, &net.TCPAddr{IP: []byte(rip), Port: rport}); err != nil {
		return nil, err
	} else {
		return &Connection{Conn: conn}, nil
	}
}

// Handles incoming connections
// This should be ran concurrently in a Go routine
func handleIncoming(host string, port string, node *Node) {
	addr := host + ":" + port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening for incoming connections:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening for incoming connections:", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Couldn't accept the connection from ", conn.RemoteAddr().String(), ":", err.Error())
		}
		go handleMessages(conn, node)
		// TODO introduce new connections to other peers
	}
}

// Handles the messages incoming from the connection
func handleMessages(conn net.Conn, node *Node) {
	for {
		var data Message
		err := gob.NewDecoder(conn).Decode(data)
		if err != nil {
			message := ErrorMessage{Err: err}
			node.Message <- message
		} else {
			node.Message <- data
		}
	}
}

// Write a message to the connection
func (c Connection) Write(msg Message) error {
	return gob.NewEncoder(c.Conn).Encode(msg)
}

// Close the connection
func (c Connection) Close() error {
	return c.Conn.Close()
}

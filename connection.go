package go_cluster

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

// Resembles a connection between two nodes.
// The API makes it simpler to customize the connection.
type Connection struct {
	Conn *net.TCPConn
}

type Transport struct {
	Type    string
	Message Message
}

// Creates a new connection
func connect(address string) (*Connection, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	if conn, err := net.DialTCP("tcp", nil, addr); err != nil {
		return nil, err
	} else {
		return &Connection{Conn: conn}, nil
	}
}

// Handles incoming connections
// This should be ran concurrently in a Go routine
func handleIncoming(address string, node *Node) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println("couldn't resolve address:", err)
		os.Exit(1)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Error listening for incoming connections:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening for incoming connections:", address)
	for node.Ready {
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("Couldn't accept the connection:", err.Error())
			continue
		}
		if conn == nil {
			fmt.Println("conn is nil")
			continue
		}
		connection := &Connection{Conn: conn}
		go handleMessages(connection, node, node.NextId)
		if node.Master {
			if err := connection.Write(ReadyMessage{Id: node.NextId}); err != nil {
				fmt.Println("failed to send ready message to ", conn.RemoteAddr().String(), ":", err.Error())
				continue
			}
			node.Nodes.Store(node.NextId, connection)
			node.NextId++
		}
		// TODO introduce new connections to other peers
	}
}

// Handles the messages incoming from the connection
func handleMessages(connection *Connection, node *Node, remoteid int) {
	conn := connection.Conn
	dec := gob.NewDecoder(conn)
	for {
		data := new(Transport)
		err := dec.Decode(data)
		if err != nil {
			switch err.(type) {
			case *net.OpError:
				// network error connection should be closed
				fmt.Println("Message handling stopped", node.Id, "<->", remoteid)
				node.Nodes.Delete(remoteid)
				return
			default:
				if err == io.EOF { // on EOF reset the decoder
					dec = gob.NewDecoder(conn)
					continue
				}
				fmt.Println("an error has occurred while reading, ", err)
				message := ErrorMessage{Err: err}
				node.Message <- message
			}
		} else {
			if data.Type == "idreq" {
				node.Nodes.Store(data.Message.(IdReqMessage).Id, connection)
				continue
			}
			node.Message <- data.Message
			dec = gob.NewDecoder(conn) // prevent old data from staying in the buffer
		}
	}

}

// Write a message to the connection
func (c Connection) Write(msg Message) error {
	return gob.NewEncoder(c.Conn).Encode(Transport{
		Type:    msg.Type(),
		Message: msg,
	})
}

// Close the connection
func (c Connection) Close() error {
	return c.Conn.Close()
}

// TODO Add heartbeat mechanism
package go_cluster

import (
	"encoding/gob"
	"io"
	"net"
	"os"
)

// Resembles a connection between two nodes.
type Connection struct {
	Conn *net.TCPConn // The connection
	Data interface{}  // The remote node's extra data
}

// The transport protocol
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

// Connects a node to a new node and sets up message handling
func connectNewNode(id int, addr string, data interface{}, node *Node) {
	if conn, err := connect(addr); err != nil {
		node.Log("error while connecting to a new node:", err)
	} else {
		conn.Data = data
		node.Nodes.Store(id, conn)
		msg := GreetingMessage{Id: node.Id, Data: data}
		if err := conn.Write(msg); err != nil {
			node.Log("couldn't send the message to the new node:", err)
		}
		go handleMessages(conn, node, id)
		node.Log("Successfully connected to the new node!")
	}
}

// Handles incoming connections
// This should be ran concurrently in a Go routine
func handleIncoming(address string, node *Node) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		node.Log("couldn't resolve address:", err)
		os.Exit(1)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		node.Log("Error listening for incoming connections:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	node.Log("Listening for incoming connections:", address)
	for node.Ready {
		conn, err := l.AcceptTCP()
		if err != nil {
			node.Log("Couldn't accept the connection:", err.Error())
			continue
		}
		if conn == nil {
			continue
		}
		connection := &Connection{Conn: conn}
		node.Nodes.Store(node.NextId, connection)
		go handleMessages(connection, node, node.NextId)
		if err := connection.Write(ReadyMessage{Id: node.NextId, EntryId: node.Id}); err != nil {
			node.Log("failed to send ready message to ", conn.RemoteAddr().String(), ":", err.Error())
			continue
		}
		node.NextId++
	}
}

// Handles the messages incoming from the connection
func handleMessages(connection *Connection, node *Node, remoteid int) {
	node.Log("Message handling started with", remoteid, "(note that this might not be the ")
	conn := connection.Conn
	dec := gob.NewDecoder(conn)
	for {
		data := new(Transport)
		err := dec.Decode(data)
		if err != nil {
			switch err.(type) {
			case *net.OpError:
				// network error connection should be closed
				node.Log("Message handling stopped with", remoteid)
				node.Nodes.Delete(remoteid)
				return
			default:
				if err == io.EOF { // on EOF reset the decoder
					dec = gob.NewDecoder(conn)
					continue
				}
				node.Log("an error has occurred while reading,", err)
				message := ErrorMessage{Err: err}
				node.Message <- message
			}
		} else {
			if data.Type == "readyreq" {
				if _, ok := node.Nodes.Load(remoteid); !ok {
					readymsg := data.Message.(ReadyMessage)
					node.Id = readymsg.Id
					remoteid = readymsg.EntryId
					node.NextId = node.Id + 1
					if err := connection.Write(IntroduceMessage{Addr: node.Addr, Data: data}); err != nil {
						node.Log("couldn't send an introduce message to the entry point:", err)
						os.Exit(1)
					}
					node.Nodes.Store(remoteid, connection)
					node.Log("Node Intialized!")
				}
			} else if data.Type == "greetreq" {
				remoteid = data.Message.(GreetingMessage).Id
				node.Log("Successfully connected to node", remoteid)
				node.Nodes.Store(remoteid, connection)
			} else if data.Type == "introreq" {
				msg := NewNodeMessage{
					Id:   remoteid,
					Addr: data.Message.(IntroduceMessage).Addr,
					Data: data.Message.(IntroduceMessage).Data,
				}
				if err := node.Broadcast(msg, remoteid); err != nil {
					node.Log("error while broadcasting a new node:", err)
				}
			} else if data.Type == "newnodereq" {
				msg := data.Message.(NewNodeMessage)
				node.NextId = msg.Id + 1
				connectNewNode(msg.Id, msg.Addr, msg.Data, node)
			} else {
				node.Message <- data.Message
			}
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

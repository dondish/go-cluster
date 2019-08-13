// This file contains structs and functions for the nodes.
package go_cluster

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"sync"
)

// The node is the general data type in go-cluster, it resembles a node in the cluster
type Node struct {
	Ready   bool         // Whether is node is ready for connections
	Master  bool         // Whether this node is the master
	Id      int          // This node's ID
	NextId  int          // (Master Only) Next client's id
	Message chan Message // The channel to forward messages to
	Nodes   *sync.Map    // A map that maps other node ids to their connections.
}

func Init() {
	gob.Register(ReadyMessage{})
	gob.Register(ErrorMessage{})
	gob.Register(IdReqMessage{})
	gob.Register(NewNodeMessage{})
	gob.Register(AddressMessage{})
}

// Creates a new master node
// This new node will introduce other nodes to each other.
func CreateMasterNode(address string) *Node {
	node := &Node{
		Ready:   true,
		Master:  true,
		NextId:  1,
		Message: make(chan Message),
		Nodes:   new(sync.Map),
	}
	go handleIncoming(address, node)
	return node
}

// Creates a new node and connects to the master
// TODO add id retrieval via a handshake
// TODO add automatic peer discovery
func CreateNode(address, maddress string) (*Node, error) {
	node := &Node{
		Ready:   true,
		Id:      1,
		Message: make(chan Message),
		Nodes:   new(sync.Map),
	}
	// TODO add id
	if conn, err := connect(maddress); err != nil {
		return nil, err
	} else if conn == nil {
		return nil, errors.New("connection cannot be nil, unexpected error occurred")
	} else {
		node.Nodes.Store(0, conn)
		go handleMessages(conn, node, 0)
		readymsg := <-node.Message
		node.Id = readymsg.(ReadyMessage).Id
		log.Println("Node Intialized! Id:", node.Id, "Address:", address)
		go handleIncoming(address, node)
		if err := node.Send(AddressMessage{Addr: address}, 0); err != nil {
			return nil, err
		}
		return node, nil
	}
}

func (n Node) Log(args ...interface{}) {
	log.Printf("Node %d: %s", n.Id, fmt.Sprintln(args...))
}

// Sends a message to a specific amount of nodes
func (n Node) Send(message Message, ids ...int) error {
	for _, id := range ids {
		if a, ok := n.Nodes.Load(id); !ok {
			return errors.New(fmt.Sprintf("Node with id: %d has not been found.", id))
		} else if err := a.(*Connection).Write(message); err != nil {
			return err
		}

	}
	return nil
}

// Sends a message to all nodes
func (n Node) Broadcast(message Message, except ...int) error {
	var err error = nil
	n.Nodes.Range(func(id, conn interface{}) bool {
		for _, exception := range except {
			if exception == id {
				return true
			}
		}
		if err = conn.(*Connection).Write(message); err != nil {
			return false
		}
		return true
	})
	return err
}

// Shuts down the node, if the node is the master it will broadcast the new master to all nodes before closing.
// It accepts a channel that will receive a boolean when the node is shutdown.
func (n *Node) Close() error {
	// TODO new master broadcast
	n.Log("Shutting down...")
	var err error = nil
	n.Nodes.Range(func(id, conn interface{}) bool {
		if err = conn.(*Connection).Close(); err != nil {
			return false
		}
		return true
	})
	n.Ready = false
	return err
}

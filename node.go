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
	Addr    string       // Incoming connections address
	Ready   bool         // Whether is node is ready for connections
	Id      int          // This node's ID
	NextId  int          // Next client's id
	Message chan Message // The channel to forward messages to
	Nodes   *sync.Map    // A map that maps other node ids to their connections.
	Data    interface{}  // Allows the node to have extra information attached to it
}

func Init() {
	gob.Register(ReadyMessage{})
	gob.Register(ErrorMessage{})
	gob.Register(GreetingMessage{})
	gob.Register(NewNodeMessage{})
	gob.Register(IntroduceMessage{})
	gob.Register(Transport{})
}

// Creates a new master node
// This new node will introduce other nodes to each other.
func CreateCluster(address string, data interface{}) *Node {
	node := &Node{
		Addr:    address,
		Ready:   true,
		NextId:  1,
		Message: make(chan Message),
		Nodes:   new(sync.Map),
		Data:    data,
	}
	go handleIncoming(address, node)
	return node
}

// Creates a new node and connects to the master
func JoinCluster(address, maddress string, data interface{}) (*Node, error) {
	node := &Node{
		Addr:    address,
		Ready:   true,
		Id:      1,
		NextId:  1,
		Message: make(chan Message),
		Nodes:   new(sync.Map),
		Data:    data,
	}
	if conn, err := connect(maddress); err != nil {
		return nil, err
	} else if conn == nil {
		return nil, errors.New("connection cannot be nil, unexpected error occurred")
	} else {
		go handleIncoming(address, node)
		go handleMessages(conn, node, 0)
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
func (n *Node) Close() {
	n.Log("Shutting down...")
	n.Nodes.Range(func(id, conn interface{}) bool {
		if err := conn.(*Connection).Close(); err != nil {
			return true
		}
		return true
	})
	n.Ready = false
}

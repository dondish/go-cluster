// This file contains structs and functions for the nodes.
package go_cluster

import (
	"errors"
)

// The node is the general data type in go-cluster, it resembles a node in the cluster
type Node struct {
	Ready   bool               // Whether is node is ready for connections
	Master  bool               // Whether this node is the master
	Id      int                // This node's ID
	Message chan Message       // The channel to forward messages to
	Nodes   map[int]Connection // A map that maps other node ids to their connections.
}

// Creates a new master node
// This new node will introduce other nodes to each other.
func CreateMasterNode(address string) *Node {
	node := &Node{
		Ready:   true,
		Master:  true,
		Id:      0,
		Message: make(chan Message),
		Nodes:   make(map[int]Connection),
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
		Nodes:   make(map[int]Connection),
	}
	// TODO add id
	if conn, err := connect(address, maddress); err == nil && conn != nil {
		node.Nodes[0] = *conn
	} else if err != nil {
		return nil, err
	} else if conn == nil {
		return nil, errors.New("connection cannot be nil, unexpected error occurred")
	}
	go handleIncoming(address, node)
	return node, nil
}

// Sends a message to a specific amount of nodes
func (n *Node) Send(message Message, ids ...int) error {
	for id := range ids {
		if err := n.Nodes[id].Write(message); err != nil {
			return err
		}
	}
	return nil
}

// Sends a message to all nodes
func (n *Node) Broadcast(message Message) error {
	for _, conn := range n.Nodes {
		if err := conn.Write(message); err != nil {
			return err
		}
	}
	return nil
}

// Shuts down the node, if the node is the master it will broadcast the new master to all nodes before closing.
// It accepts a channel that will receive a boolean when the node is shutdown.
func (n *Node) Close() error {
	n.Ready = false
	close(n.Message)
	// TODO new master broadcast
	for _, conn := range n.Nodes {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

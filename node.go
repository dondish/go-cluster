// This file contains structs and functions for the nodes.
package go_cluster

// The node is the general data type in go-cluster, it resembles a node in the cluster
type Node struct {
	Master  bool               // Whether this node is the master
	Id      int                // This node's ID
	Message chan<- Message     // The channel to forward messages to
	Nodes   map[int]Connection // A map that maps other node ids to their connections.
}

// Creates a new master node
// This new node will introduce other nodes to each other.
func CreateMasterNode(host string, port int) *Node {
	node := &Node{
		Master:  true,
		Id:      0,
		Message: make(chan Message),
		Nodes:   make(map[int]Connection),
	}
	go handleIncoming(host, string(port), node)
	return node
}

// Creates a new node
// Arguments
func CreateNode(ip string, port int, mip string, mport int) (*Node, error) {
	return nil, nil
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
func (*Node) Close() error {
	// TODO finish Close
	return nil
}

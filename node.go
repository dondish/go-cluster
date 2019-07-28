/// This file contains structs and functions for the nodes.
package go_cluster

// The Message interface, this is supposed to be customized (for example Msg is encoded in gob).
type Message interface {
	Msg() string
	Error() error
}

type Node struct {
	master  bool               // Whether this node is the master
	id      int                // This node's id
	message chan<- Message     // The channel to forward messages to
	nodes   map[int]Connection // A map that maps other node ids to their connections.
}

// Starts the node, it will connect it to the master.
// It accepts a channel that will receive a boolean when the node is ready.
func CreateNode(master bool, ip string) (*Node, error) {
	if master {
		return &Node{
			master:  true,
			id:      0,
			message: make(chan Message),
			nodes:   make(map[int]Connection),
		}, nil
	} else {
		// TODO finish CreateNode
		return nil, nil
	}
}

func (*Node) Send(message string, ids ...int) error {
	// TODO finish Send
	return nil
}

func (*Node) Broadcast(message string) error {
	// TODO finish Broadcast
	return nil
}

// Shuts down the node, if the node is the master it will broadcast the new master to all nodes before closing.
// It accepts a channel that will receive a boolean when the node is shutdown.
func (*Node) Close() error {
	// TODO finish Close
	return nil
}

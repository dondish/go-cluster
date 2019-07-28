package go_cluster

import "net"

// Resembles a connection between two nodes.
// In the future, many methods will be added to interact between nodes.
type Connection struct {
	Conn *net.TCPConn
}

// Creates a new connection
func connect(ip string, port int) (*Connection, error) {
	if conn, err := net.DialTCP("tcp", &net.TCPAddr{IP: []byte(ip), Port: port}, nil); err != nil {
		return nil, err
	} else {
		return &Connection{Conn: conn}, nil
	}
}

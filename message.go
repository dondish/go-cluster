package go_cluster

import (
	"encoding/gob"
	"strconv"
)

// The Message interface, this is supposed to be customized (for example Msg is encoded in gob).
type Message interface {
	Type() string
	Msg() interface{}
}

// A message containing only errors
type ErrorMessage struct {
	Err error
}

func (m ErrorMessage) Type() string {
	return "error"
}

func (m ErrorMessage) Msg() interface{} {
	return m.Err.Error()
}

// The message the master sends to a newly connected node with its id
type ReadyMessage struct {
	Id int
}

func (m ReadyMessage) Msg() interface{} {
	return strconv.Itoa(m.Id)
}

func (m ReadyMessage) Type() string {
	return "ready"
}

// The message a node sends to the node it's newly connected to with its id to make authentication easier
type IdReqMessage struct {
	Id int
}

func (m IdReqMessage) Msg() interface{} {
	return strconv.Itoa(m.Id)
}

func (m IdReqMessage) Type() string {
	return "idreq"
}

// The message the master sends when all nodes when a new node joins
type NewNodeMessage struct {
	Id   int    // The Id
	Addr string // The address to connect to
}

func (m NewNodeMessage) Msg() interface{} {
	return m
}

func (m NewNodeMessage) Type() string {
	return "newnode"
}

// The message a new node sends to the master with the incoming connections address
type AddressMessage struct {
	Addr string
}

func (m AddressMessage) Msg() interface{} {
	return m.Addr
}

func (m AddressMessage) Type() string {
	return "addrreq"
}

// Register the message type to gob
func RegisterMessage(msg Message) {
	gob.Register(msg)
}

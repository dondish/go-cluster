package go_cluster

import (
	"encoding/gob"
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
	Id      int
	EntryId int
}

func (m ReadyMessage) Msg() interface{} {
	return m
}

func (m ReadyMessage) Type() string {
	return "readyreq"
}

// The message a node sends to the node it's newly connected to with its id to make authentication easier
type GreetingMessage struct {
	Id   int
	Data interface{}
}

func (m GreetingMessage) Msg() interface{} {
	return m
}

func (m GreetingMessage) Type() string {
	return "greetreq"
}

// The message the master sends when all nodes when a new node joins
type NewNodeMessage struct {
	Id   int         // The Id
	Addr string      // The address to connect to
	Data interface{} // The new node's data
}

func (m NewNodeMessage) Msg() interface{} {
	return m
}

func (m NewNodeMessage) Type() string {
	return "newnodereq"
}

// The message a new node sends to the master with its information
type IntroduceMessage struct {
	Addr string
	Data interface{}
}

func (m IntroduceMessage) Msg() interface{} {
	return m.Addr
}

func (m IntroduceMessage) Type() string {
	return "introreq"
}

// Register the message type to gob
func RegisterMessage(msg Message) {
	gob.Register(msg)
}

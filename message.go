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

func RegisterMessage(msg Message) {
	gob.Register(msg)
}

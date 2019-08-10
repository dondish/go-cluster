package go_cluster

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

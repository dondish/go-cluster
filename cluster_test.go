package go_cluster

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestMessage struct {
	Test string
}

func (t TestMessage) Msg() interface{} {
	return t.Test
}

func (t TestMessage) Type() string {
	return "test"
}

func TestCreateMasterNode(t *testing.T) {
	master := CreateMasterNode("localhost:5555")
	defer master.Close()

	assert.Empty(t, master.Nodes, "the master should not be connected to other nodes implicitly")
	assert.NotNil(t, master.Message, "the message channel should not be nil")
	assert.True(t, master.Id == 0, "the master id should be 0")
	assert.True(t, master.Master, "the Master boolean should be true")
}

func TestCreateNode(t *testing.T) {
	master := CreateMasterNode("localhost:5556")
	defer master.Close()
	node, err := CreateNode("localhost:5557", "localhost:5556")
	defer node.Close()
	if err != nil {
		fmt.Println("couldn't create node:", err)
		t.Fail()
	}

	assert.Contains(t, node.Nodes, 0, "the node should have master in its nodes map")
	assert.True(t, node.Id == 1, "the node's id should be set to 1")
	assert.NotNil(t, master.Message, "the message channel should not be nil")

	go func() {
		if err := node.Send(TestMessage{Test: "testing"}, 0); err != nil {
			t.Fail()
		}
	}()
	c := <-master.Message
	fmt.Println(c.Msg())
	assert.Equal(t, c.Msg(), "teasting", "the message should be equal to testing")
	assert.Equal(t, c.Type(), "test", "the type should be equal to test")
}

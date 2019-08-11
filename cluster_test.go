package go_cluster

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	defer func() {
		err := master.Close()
		assert.Nil(t, err, "There shouldn't be an error while closing the master")
	}()

	assert.Empty(t, master.Nodes, "the master should not be connected to other nodes implicitly")
	assert.NotNil(t, master.Message, "the message channel should not be nil")
	assert.True(t, master.Id == 0, "the master id should be 0")
	assert.True(t, master.Master, "the Master boolean should be true")
}

func TestCreateNode(t *testing.T) {
	master := CreateMasterNode("localhost:5556")
	defer func() {
		err := master.Close()
		assert.Nil(t, err, "There shouldn't be an error while closing the master")
	}()
	node, err := CreateNode("localhost:5557", "localhost:5556")

	if err != nil {
		fmt.Println("couldn't create node:", err)
		assert.NotNil(t, err, "There shouldn't be an error while closing the slave")
	}

	defer func() {
		err := node.Close()
		assert.Nil(t, err, "There shouldn't be an error while closing the slave")
	}()

	assert.Contains(t, node.Nodes, 0, "the node should have master in its nodes map")
	assert.True(t, node.Id == 1, "the node's id should be set to 1")
	assert.NotNil(t, master.Message, "the message channel should not be nil")

	go func() {
		err := node.Send(TestMessage{Test: "testing"}, 0)
		assert.Nil(t, err, "There shouldn't be an error sending the test message")
	}()
	select {
	case c := <-master.Message:
		assert.Equal(t, c.Msg(), "testing", "the message should be equal to testing")
		assert.Equal(t, c.Type(), "test", "the type should be equal to test")
	case <-time.After(time.Second):
		t.Log("The test has timed out")
	}

	go func() {
		err := master.Send(TestMessage{Test: "testing"}, 1)
		assert.Nil(t, err, "There shouldn't be an error sending the test message")
	}()
	select {
	case c := <-node.Message:
		assert.Equal(t, c.Msg(), "testing", "the message should be equal to testing")
		assert.Equal(t, c.Type(), "test", "the type should be equal to test")
	case <-time.After(time.Second):
		t.Log("The test has timed out")
	}

}

func TestMain(m *testing.M) {
	Init()
	RegisterMessage(TestMessage{})
	m.Run()
}

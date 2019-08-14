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

func TestCreateSingleNodeCluster(t *testing.T) {
	master := CreateCluster("localhost:5555", nil)
	defer master.Close()

	assert.Empty(t, master.Nodes, "the node should not be connected to other nodes implicitly")
	assert.NotNil(t, master.Message, "the message channel should not be nil")
	assert.True(t, master.Id == 0, "the node id should be 0")
}

func TestCreateTwoNodeCluster(t *testing.T) {
	master := CreateCluster("localhost:5556", nil)
	defer master.Close()
	time.Sleep(500 * time.Millisecond)
	node, err := JoinCluster("localhost:5557", "localhost:5556", nil)

	if err != nil {
		fmt.Println("couldn't create node:", err)
		assert.NotNil(t, err, "There shouldn't be an error while closing the slave")
	}

	defer node.Close()
	time.Sleep(500 * time.Millisecond)
	_, ok := node.Nodes.Load(0)
	assert.True(t, ok, "the node should have master in its nodes map")
	assert.True(t, node.Id == 1, "the node's id should be set to 1")
	assert.NotNil(t, master.Message, "the message channel should not be nil")

	go func() {
		err := node.Send(TestMessage{Test: "testing"}, 0)
		assert.Nil(t, err, "There shouldn't be an error sending the test message")
	}()
	select {
	case c := <-master.Message:
		assert.Equal(t, "testing", c.Msg(), "the message should be equal to testing")
		assert.Equal(t, "test", c.Type(), "the type should be equal to test")
	case <-time.After(time.Second):
		t.Log("The test has timed out")
	}

	go func() {
		err := master.Send(TestMessage{Test: "testing"}, 1)
		assert.Nil(t, err, "There shouldn't be an error sending the test message")
	}()
	select {
	case c := <-node.Message:
		assert.Equal(t, "testing", c.Msg(), "the message should be equal to testing")
		assert.Equal(t, "test", c.Type(), "the type should be equal to test")
	case <-time.After(time.Second):
		t.Log("The test has timed out")
	}

}

func TestCreateMultiNodeCluster(t *testing.T) {
	master := CreateCluster("localhost:5558", nil)
	defer master.Close()
	time.Sleep(500 * time.Millisecond)
	node1, err := JoinCluster("localhost:5559", "localhost:5558", nil)

	if err != nil {
		fmt.Println("couldn't create node 1:", err)
		assert.NotNil(t, err, "There shouldn't be an error while closing the slave")
	}

	defer node1.Close()
	time.Sleep(500 * time.Millisecond)
	node2, err := JoinCluster("localhost:5560", "localhost:5558", nil)

	if err != nil {
		fmt.Println("couldn't create node 2:", err)
		assert.NotNil(t, err, "There shouldn't be an error while closing the slave")
	}

	defer node2.Close()

	time.Sleep(time.Second * 2) // Wait until the nodes greet each other
	_, ok := node1.Nodes.Load(2)
	assert.True(t, ok, "Node 2 is supposed to be connected to node 1")
	_, ok = node2.Nodes.Load(1)
	assert.True(t, ok, "Node 1 is supposed to be connected to node 2")
}

func TestMain(m *testing.M) {
	Init()
	RegisterMessage(TestMessage{})
	m.Run()
}

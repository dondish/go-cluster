# Go Cluster [![Build Status](https://travis-ci.com/dondish/go-cluster.svg?branch=master)](https://travis-ci.com/dondish/go-cluster)
P2P, Master-Slave model of clustering for Go.

This project aims to be minimal, performance and code size wise.

# Note
The package uses gob for encoding messages, before using the library register all of the message types to gob.

More on that [gob docs Register method](https://golang.org/pkg/encoding/gob/#Register)

# Why?
Go has a great concurrency model that makes the language be frequently used for RPC.

The purpose of this project is to provide a simple and efficient solution to cluster Go microservices easily.

# How?
Peer to peer connectivity, using a custom protocol over TCP. TCP was chosen because of it's reliability.

The model itself is a master-slave model, this lets the connection be distributed faster and easier.
The master introduces new nodes to all of the other nodes. Each node is identified with a custom ID, 
so nodes can communicate with each other as well.

## What happens when the master crashes?
The custom ID is given by the master and it's the number of nodes that have joined before the new node. 
This means the successor of the master is the successor in ID.
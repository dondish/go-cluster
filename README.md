# Go Cluster [![Build Status](https://travis-ci.com/dondish/go-cluster.svg?branch=master)](https://travis-ci.com/dondish/go-cluster)
P2P model of clustering for Go.

This project aims to be minimal, performance and code size wise.

# Note
The package uses gob for encoding messages, before using the library register all of the message types to gob.

More on that [gob docs Register method](https://golang.org/pkg/encoding/gob/#Register)

# Why?
Go has a great concurrency model that makes the language be frequently used for RPC.

The purpose of this project is to provide a simple and efficient solution to cluster Go microservices easily.

# How?
Peer to peer connectivity, using a custom protocol over TCP. TCP was chosen because of it's reliability.

The model itself isn't a master-slave model, this lets the every node to be independent on each other.
Each node can introduce new nodes to all of the other nodes. Each node is identified with a custom ID, 
so nodes can communicate with each other as well.

# Features
* Resilent - Every node is independent of other nodes so when a node crashes nothing happens.
* Fast - The package uses the gob encoding which is very fast and efficient
* Customizable - The Message interface is supposed to be customizable, send every type of message you'd like over the cluster.
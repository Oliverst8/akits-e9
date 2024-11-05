package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "main/grpc"
	"net"
	"os"
	"sync"
)

var lock sync.Mutex

// MutexNode Each node has a server to receive information and a list of clients to send information to all other nodes
type MutexNode struct {
	proto.UnimplementedMutexNodeServer
	port        string
	clients     []proto.MutexNodeClient
	state       string
	myRequests  []proto.RequestMessage
	lamportTime uint64
}

func main() {
	clientPort := os.Args[1] //The port of this node

	node := &MutexNode{}

	node.start_server(clientPort) //Starting a server so that this node listen on the given port

	if len(os.Args) > 2 { //If the client isnt the starting node

		joinOn := os.Args[2] //The port to send the join request to

		node.start_client(joinOn)

		message := proto.JoinMessage{
			Port: clientPort,
		}

		response, err := node.clients[0].Join(context.Background(), &message) //Trying to join using the first node

		if err != nil {
			panic(err)
		}

		if !response.Success {
			panic("Failed to join cluster")
		}

		// get this nodes ID and start clients up for each of the other nodes in the network
		for _, port := range response.Ports {
			node.start_client(port)
		}
		node.state = "RELEASED"
	}
	for {
		if node.state == "RELEASED" {
			panic("Not implemented")
		}

	}
}

func (s MutexNode) start_server(port string) { // start up a new server and listen on the nodes port
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	proto.RegisterMutexNodeServer(grpcServer, s)

	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (s MutexNode) start_client(port string) { // start up a new client for the node to send information through the given port
	conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	lock.Lock()
	s.clients = append(s.clients, proto.NewMutexNodeClient(conn))
	lock.Unlock()
}

// Join Called from client, to make a request to join the network
func (s MutexNode) Join(context context.Context, message *proto.JoinMessage) (*proto.JoinResponse, error) {
	//Sends a
	var ports []string
	lock.Lock()
	for _, client := range s.clients {
		res, err := client.AddNode(context, message)
		if err != nil {
			return nil, err
		}
		ports = append(ports, res.Port)
	}
	lock.Unlock()
	reply := proto.JoinResponse{
		Ports:   ports,
		Success: true,
		Time:    s.lamportTime,
	}
	return &reply, nil
}

// AddNode This function adds a client to the node so it can send information to the newly joined node
func (s MutexNode) AddNode(context context.Context, message *proto.JoinMessage) (*proto.JoinMessage, error) {
	port := message.Port
	s.start_client(port)
	return &proto.JoinMessage{
		Port: s.port,
	}, nil
}

func (s MutexNode) Request(context context.Context, message *proto.RequestMessage) (*proto.Reply, error) {

}

func (s MutexNode) makeRequest(message *proto.RequestMessage) {
	lock.Lock()
	for _, client := range s.clients {
		_, err := client.Request(context.Background(), message)
		if err != nil {
			panic(err)
		}
	}
	lock.Unlock()
}

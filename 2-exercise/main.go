package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "main/grpc"
	"net"
	"os"
	"strconv"
	"sync"
)

// MutexNode Each node has a server to receive information and a list of clients to send information to all other nodes
type MutexNode struct {
	proto.UnimplementedMutexNodeServer
	port        string
	clients     map[string]proto.MutexNodeClient
	state       string
	myRequests  []proto.MutexNodeClient
	lamportTime uint64
	responses   int
}

var responsesLock sync.Mutex
var reqeustLock sync.Mutex

func main() {
	clientPort := os.Args[2] //The port of this node
	desiredNetworkSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	node := &MutexNode{}

	node.start_server(clientPort) //Starting a server so that this node listen on the given port

	if len(os.Args) > 3 { //If the client isnt the starting node

		joinOn := os.Args[3] //The port to send the join request to

		node.start_client(joinOn)

		message := proto.JoinMessage{
			Port: clientPort,
		}

		response, err := node.clients[joinOn].Join(context.Background(), &message) //Trying to join using the first node

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

	}
	for len(node.clients) < desiredNetworkSize {
	}
	node.state = "RELEASED"
	for {
		if node.state == "RELEASED" {
			reqeustLock.Lock()
			for _, client := range node.myRequests {
				reply := proto.Reply{
					Success: true,
					Time:    node.lamportTime,
				}
				_, err = client.RespondToRequest(context.Background(), &reply)
				if err != nil {
					panic(err)
				}
			}
			reqeustLock.Unlock()
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
	s.clients[port] = proto.NewMutexNodeClient(conn)
}

// Join Called from client, to make a request to join the network
func (s MutexNode) Join(context context.Context, message *proto.JoinMessage) (*proto.JoinResponse, error) {
	//Sends a
	var ports []string
	for _, client := range s.clients {
		res, err := client.AddNode(context, message)
		if err != nil {
			return nil, err
		}
		ports = append(ports, res.Port)
	}
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

func (s MutexNode) Request(context context.Context, message *proto.RequestMessage) (*proto.Empty, error) {
	requestingClient := s.clients[message.Port]
	if s.state == "HELD" || (s.state == "WANTED" && s.compare(message)) {
		s.myRequests = append(s.myRequests, requestingClient)
	} else {
		reply := proto.Reply{
			Success: true,
			Time:    s.lamportTime,
		}
		_, err := requestingClient.RespondToRequest(context, &reply)
		if err != nil {
			panic(err)
		}
	}
	return &proto.Empty{}, nil
}

func (s MutexNode) RespondToRequest(context context.Context, reply *proto.Reply) (*proto.Empty, error) {
	responsesLock.Lock()
	s.responses += 1
	responsesLock.Unlock()
	return &proto.Empty{}, nil
}

func (s MutexNode) compare(message *proto.RequestMessage) bool {
	thisPort, err := strconv.Atoi(s.port)
	if err != nil {
		panic(err)
	}
	thatPort, err := strconv.Atoi(message.Port)
	if err != nil {
		panic(err)
	}
	if s.lamportTime < message.Time {
		return true
	}
	if s.lamportTime == message.Time && thisPort < thatPort {
		return true
	}
	return false
}

func (s MutexNode) multicast(message *proto.RequestMessage) {
	s.state = "WANTED"
	s.responses = 0
	for _, client := range s.clients {
		go makeRequest(client, message)
	}
	for {
		if s.responses >= len(s.clients) {
			break
		}
	}
	s.state = "HELD"
	incrementFile()
	s.state = "RELEASED"
}

func incrementFile() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create("log.txt")
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	text := "0"
	for scanner.Scan() {
		text = scanner.Text()
	}

	newNumber, err := strconv.Atoi(text)
	if err != nil {
		panic(err)
	}

	newNumber++

	_, err = fmt.Fprintf(file, "%d\n", newNumber)
	if err != nil {
		panic(err)
	}
}

func makeRequest(client proto.MutexNodeClient, message *proto.RequestMessage) {
	_, err := client.Request(context.Background(), message)
	if err != nil {
		panic(err)
	}
}

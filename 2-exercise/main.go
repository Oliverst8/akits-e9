package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "main/grpc"
	"net"
	"os"
)

var clients []proto.Mutex2Client

type Mutex2 struct {
	proto.UnimplementedMutex2Server
}

func main() {
	port := os.Args[1]
	joinOn := os.Args[2]

	start_client(joinOn)

	server := &Mutex2{}

	server.start_server(port)
	clients[0].Join()
	for {
	}
}

func (s Mutex2) start_server(port string) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	proto.RegisterMutex2Server(grpcServer, s)

	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func start_client(port string) {
	conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	clients = append(clients, proto.NewMutex2Client(conn))
}

func (s Mutex2) Join(message proto.JoinMessage) {

}

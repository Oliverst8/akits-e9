package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	proto "main/grpc"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"time"
)

type MutexNode struct {
	proto.UnimplementedMutexNodeVTwoServer
	port string
}

var hasToken bool
var client proto.MutexNodeVTwoClient

func main() {

	fmt.Println("Starting")

	myPort := os.Args[1]

	node := &MutexNode{
		port: os.Args[2],
	}

	go node.start_server(myPort)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "go" {
			break
		}
	}

	go node.start_client()

	time.Sleep(2 * time.Second)

	if len(os.Args) < 4 {
		hasToken = false
	} else {
		fmt.Println("Starting")
		hasToken = true
	}

	num := rand.Float32()
	wantAccess := false
	i := 0
	go checkforToken()

	for {
		if !wantAccess {
			if num < 0.03 && i < 100 {
				fmt.Println("Requesting access")
				wantAccess = true
			} else {
				num = rand.Float32()
			}
		}

		if hasToken {
			if wantAccess {
				log.Printf("I got access, i: %d\n", i)
				incrementFile()
				wantAccess = false
				i++
			}
			message := &proto.Empty{}
			if client == nil {
				panic("Client is nil")
			}
			hasToken = false
			_, err := client.SendToken(context.Background(), message)
			if err != nil {
				panic(err)
			}
		}

	}

}

func checkforToken() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "?" {
			fmt.Printf("Token status: %s\n", hasToken)
		}
	}
}

func (s MutexNode) SendToken(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	hasToken = true
	fmt.Println("I got access")
	message := proto.Empty{}
	return &message, nil
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

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

func (s *MutexNode) start_client() { // start up a new client for the node to send information through the given port
	conn, err := grpc.NewClient("localhost:"+s.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewMutexNodeVTwoClient(conn)
	fmt.Printf("Client: %s\n", client)
}

func (s *MutexNode) start_server(port string) { // start up a new server and listen on the nodes port
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+port)
	fmt.Println("Created listener")

	if err != nil {
		panic(err)
	}

	proto.RegisterMutexNodeVTwoServer(grpcServer, s)

	err = grpcServer.Serve(listener)
	fmt.Printf("Now listening on port %s\n", port)
	if err != nil {
		panic(err)
	}
}

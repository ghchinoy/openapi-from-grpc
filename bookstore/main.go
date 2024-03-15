package main

import (
	"bookstore/server/bookstore/pb"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedInventoryServer
	pb.UnimplementedEchoServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) GetBooks(ctx context.Context, in *pb.GetBooksRequest) (*pb.GetBooksResponse, error) {
	log.Printf("Received GetBooks request: %v", in.ProtoReflect().Descriptor().FullName())
	return &pb.GetBooksResponse{
		Books: getSampleBooks(),
	}, nil
}

func (s *server) Echo(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	return &pb.EchoMessage{
		Value: in.Value,
	}, nil
}

func getSampleBooks() []*pb.Book {
	return []*pb.Book{
		{
			Pages:  412,
			Title:  "Dune",
			Author: "Herbert, Frank",
		},
		{
			Pages:  256,
			Title:  "Dune Messiah",
			Author: "Herbert, Frank",
		},
		{
			Pages:  232,
			Title:  "Children of Dune",
			Author: "Herbert, Frank",
		},
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// gRPC server instance
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterInventoryServer(s, &server{})
	pb.RegisterEchoServer(s, &server{})

	log.Println("gRPC server started on port 8080")
	go func() {
		log.Fatalln(s.Serve(listener))
	}()

	// proxy
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	gwmux := runtime.NewServeMux()
	err = pb.RegisterEchoHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("gRPC-gateway started on port 8090")
	log.Fatalln(gwServer.ListenAndServe())

}

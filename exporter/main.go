package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tanya.lyubimaya/grpc-exporter/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	s := NewServer()
	grpcServer := grpc.NewServer(opts...)
	server.RegisterExporterServer(grpcServer, s)
	fmt.Println("listening on localhost:8080")
	done := make(chan struct{})
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		done <- struct{}{}
	}()

	<-done
}

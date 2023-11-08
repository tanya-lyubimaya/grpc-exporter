package main

import (
	"log"
	"net"

	"github.com/tanya.lyubimaya/mockConfigStore/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	ms := NewServer()
	grpcServer := grpc.NewServer(opts...)
	server.RegisterMetricsServiceServer(grpcServer, ms)
	done := make(chan struct{})
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		done <- struct{}{}
	}()

	<-done
}

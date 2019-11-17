package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pbsvc "mnovicio.com/nycab/protocol/rpc"
)

// RunServer runs gRPC service to publish NYCab service
func RunServer(ctx context.Context, nyCabService pbsvc.NYCabServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	pbsvc.RegisterNYCabServiceServer(server, nyCabService)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println(fmt.Sprintf("starting gRPC server at port %s...", port))
	return server.Serve(listen)
}

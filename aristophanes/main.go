package main

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const standardPort = ":50052"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	// https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ARISTOPHANES
	logging.System(`
  ____  ____   ____ _____ ______   ___   ____  __ __   ____  ____     ___  _____
 /    ||    \ |    / ___/|      | /   \ |    \|  |  | /    ||    \   /  _]/ ___/
|  o  ||  D  ) |  (   \_ |      ||     ||  o  )  |  ||  o  ||  _  | /  [_(   \_ 
|     ||    /  |  |\__  ||_|  |_||  O  ||   _/|  _  ||     ||  |  ||    _]\__  |
|  _  ||    \  |  |/  \ |  |  |  |     ||  |  |  |  ||  _  ||  |  ||   [_ /  \ |
|  |  ||  .  \ |  |\    |  |  |  |     ||  |  |  |  ||  |  ||  |  ||     |\    |
|__|__||__|\_||____|\___|  |__|   \___/ |__|  |__|__||__|__||__|__||_____| \___|
	`)

	logging.System("βρεκεκεκὲξ κοὰξ κοάξ.τούτῳ γὰρ οὐ νικήσετε.")
	logging.System("Brekekekex koax koax. You never beat me in this play!")
	logging.System("Starting up...")

	// Create the TraceServiceClient using the environment
	traceClient, err := comedy.NewTraceServiceImpl()
	if err != nil {
		log.Fatalf("error creating TraceServiceClient: %v", err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logging.Warn(fmt.Sprintf("Recovered from panic in ManageStartTimeMap: %v", r))
			}
		}()
		traceClient.ManageStartTimeMap()
	}()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server = grpc.NewServer()

	pb.RegisterTraceServiceServer(server, traceClient)

	go func() {
		logging.System(fmt.Sprintf("Server listening on %s", port))
		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-quit
	logging.Warn("Shutting down gRPC server...")

	server.GracefulStop()
	logging.Warn("gRPC server stopped.")
}

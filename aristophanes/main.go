package main

import (
	"github.com/odysseia-greek/attike/aristophanes/app"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const standardPort = ":50052"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	// https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ARISTOPHANES
	log.Print(`
  ____  ____   ____ _____ ______   ___   ____  __ __   ____  ____     ___  _____
 /    ||    \ |    / ___/|      | /   \ |    \|  |  | /    ||    \   /  _]/ ___/
|  o  ||  D  ) |  (   \_ |      ||     ||  o  )  |  ||  o  ||  _  | /  [_(   \_ 
|     ||    /  |  |\__  ||_|  |_||  O  ||   _/|  _  ||     ||  |  ||    _]\__  |
|  _  ||    \  |  |/  \ |  |  |  |     ||  |  |  |  ||  _  ||  |  ||   [_ /  \ |
|  |  ||  .  \ |  |\    |  |  |  |     ||  |  |  |  ||  |  ||  |  ||     |\    |
|__|__||__|\_||____|\___|  |__|   \___/ |__|  |__|__||__|__||__|__||_____| \___|
	`)

	log.Print("βρεκεκεκὲξ κοὰξ κοάξ.τούτῳ γὰρ οὐ νικήσετε.")
	log.Print("Brekekekex koax koax. You never beat me in this play!")
	log.Print("Starting up...")

	env := os.Getenv("ENV")

	// Create the TraceServiceClient using the environment
	traceClient, err := app.NewTraceServiceImpl(env)
	if err != nil {
		log.Fatalf("error creating TraceServiceClient: %v", err)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer()

	pb.RegisterTraceServiceServer(server, traceClient)

	log.Printf("Server listening on %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

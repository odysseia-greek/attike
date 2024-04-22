package main

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/attike/sophokles/proto"
	"github.com/odysseia-greek/attike/sophokles/tragedy"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const standardPort = ":50053"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	// https://patorjk.com/software/taag/#p=display&f=Crawford2&t=SOPHOKLES
	logging.System(`
  _____  ___   ____  __ __   ___   __  _  _        ___  _____
 / ___/ /   \ |    \|  |  | /   \ |  |/ ]| |      /  _]/ ___/
(   \_ |     ||  o  )  |  ||     ||  ' / | |     /  [_(   \_ 
 \__  ||  O  ||   _/|  _  ||  O  ||    \ | |___ |    _]\__  |
 /  \ ||     ||  |  |  |  ||     ||     ||     ||   [_ /  \ |
 \    ||     ||  |  |  |  ||     ||  .  ||     ||     |\    |
  \___| \___/ |__|  |__|__| \___/ |__|\_||_____||_____| \___|

	`)

	logging.System("οὐ γὰρ θανεῖν ἔχθιστον, ἀλλʼ ὅταν θανεῖν. χρῄζων τις εἶτα μηδὲ τοῦτʼ ἔχῃ λαβεῖν.")
	logging.System("For death is not the most odious thing; it is rather craving death, but lacking the means to die.")
	logging.System("Starting up...")

	metricsClient, err := tragedy.NewMetricServiceImpl()
	if err != nil {
		log.Fatalf("error creating MetricsServiceClient: %v", err)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer()

	pb.RegisterMetricsServiceServer(server, metricsClient)

	logging.System(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

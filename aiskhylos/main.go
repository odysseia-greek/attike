package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/aiskhylos/tragedy"
)

func main() {
	// https://patorjk.com/software/taag/#p=display&f=Crawford2&t=Aiskhylos
	logging.System(`
  ____  ____ _____ __  _  __ __  __ __  _       ___   _____
 /    ||    / ___/|  |/ ]|  |  ||  |  || |     /   \ / ___/
|  o  | |  (   \_ |  ' / |  |  ||  |  || |    |     (   \_ 
|     | |  |\__  ||    \ |  _  ||  ~  || |___ |  O  |\__  |
|  _  | |  |/  \ ||     \|  |  ||___, ||     ||     |/  \ |
|  |  | |  |\    ||  .  ||  |  ||     ||     ||     |\    |
|__|__||____|\___||__|\_||__|__||____/ |_____| \___/  \___|

	`)

	logging.System("πάθει μάθος ὁ δ' ἄνθρωπος.")
	logging.System("He who learns must suffer.")
	logging.System("Starting up...")

	metricsClient, err := tragedy.NewGatherer()
	if err != nil {
		log.Fatalf("error creating Gatherer: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	metricsClient.Start(ctx)
}

package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/aiskhylos/tragedy"
	"log"
	"os"
	"time"
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

	env := os.Getenv("ENV")

	metricsClient, err := tragedy.NewMetricGathererImpl(env)
	if err != nil {
		log.Fatalf("error creating MetricsServiceClient: %v", err)
	}

	ticker := time.NewTicker(metricsClient.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := metricsClient.GatherMetricsOnTimerFull()
			if err != nil {
				logging.Error(err.Error())
			}
		}
	}
}

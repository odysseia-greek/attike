package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/sophokles/tragedy"
)

func main() {
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
	logging.System("Starting Sophokles collector...")

	collector, err := tragedy.NewCollector()
	if err != nil {
		logging.Error(err.Error())
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := collector.Run(ctx); err != nil {
		logging.Error(err.Error())
		panic(err)
	}
}

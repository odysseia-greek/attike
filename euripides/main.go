package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/euripides/gateway"
	"github.com/odysseia-greek/attike/euripides/routing"
)

const standardPort = ":8080"

func main() {

	logging.System(`
   ___  __ __  ____   ____  ____ ____  ___      ___  _____
  /  _]|  |  ||    \ |    ||    \    ||   \    /  _]/ ___/
 /  [_ |  |  ||  D  ) |  | |  o  )  | |    \  /  [_(   \_ 
|    _]|  |  ||    /  |  | |   _/|  | |  D  ||    _]\__  |
|   [_ |  :  ||    \  |  | |  |  |  | |     ||   [_ /  \ |
|     ||     ||  .  \ |  | |  |  |  | |     ||     |\    |
|_____| \__,_||__|\_||____||__| |____||_____||_____| \___|
                                                          
`)
	logging.System(`κακῶς φρονοῦντες· ὡς τρὶς ἂν παρ’ ἀσπίδα στῆναι θέλοιμ’ ἂν μᾶλλον ἢ τεκεῖν ἅπαξ.`)
	logging.System("How wrong they are! I would rather stand three times with a shield in battle than give birth once.\n")
	logging.System("Starting up...")

	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort

	}

	handler, err := gateway.CreateNewConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	graphqlServer := routing.InitRoutes(handler)

	logging.System(fmt.Sprintf("Server running on port %s", port))
	err = http.ListenAndServe(port, graphqlServer)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

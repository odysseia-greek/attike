package main

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/euripides/schemas"
	"github.com/odysseia-greek/attike/euripides/tragedy"
	"net/http"
	"os"
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
	logging.System(`	κακῶς φρονοῦντες· ὡς τρὶς ἂν παρ’ ἀσπίδα στῆναι θέλοιμ’ ἂν μᾶλλον ἢ τεκεῖν ἅπαξ.`)
	logging.System("How wrong they are! I would rather stand three times with a shield in battle than give birth once.\n")
	logging.System("Starting up...")

	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort

	}

	schemas.InitEuripidesHandler()
	srv := tragedy.InitRoutes()

	logging.System(fmt.Sprintf("%s : %s", "running on port", port))
	err := http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}

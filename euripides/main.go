package main

import (
	"github.com/odysseia-greek/attike/euripides/app"
	"github.com/odysseia-greek/attike/euripides/schemas"
	"log"
	"net/http"
	"os"
)

const standardPort = ":8080"

func main() {

	log.Print(`
   ___  __ __  ____   ____  ____ ____  ___      ___  _____
  /  _]|  |  ||    \ |    ||    \    ||   \    /  _]/ ___/
 /  [_ |  |  ||  D  ) |  | |  o  )  | |    \  /  [_(   \_ 
|    _]|  |  ||    /  |  | |   _/|  | |  D  ||    _]\__  |
|   [_ |  :  ||    \  |  | |  |  |  | |     ||   [_ /  \ |
|     ||     ||  .  \ |  | |  |  |  | |     ||     |\    |
|_____| \__,_||__|\_||____||__| |____||_____||_____| \___|
                                                          
`)
	log.Print(`	κακῶς φρονοῦντες· ὡς τρὶς ἂν παρ’ ἀσπίδα στῆναι θέλοιμ’ ἂν μᾶλλον ἢ τεκεῖν ἅπαξ.`)
	log.Print("How wrong they are! I would rather stand three times with a shield in battle than give birth once.\n")
	log.Print("Starting up...")

	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort

	}

	schemas.InitEuripidesHandler()
	srv := app.InitRoutes()

	log.Printf("%s : %s", "running on port", port)
	err := http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}

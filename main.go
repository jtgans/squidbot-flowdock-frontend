package main;

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/jtgans/squidbot-flowdock-frontend/frontend"
)

var httpPort = flag.String("port", "", "the host:port where varz and healthz should be served from. Required.")
var brainHostPort = flag.String("brain-hostport", "", "the host:port where the brain is running. Required.")
var authToken = flag.String("auth-token", "", "the authentication token to use for connecting to Flowdock. Required.")

var fe *frontend.Frontend

func healthHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(fe.IsOk(), "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func main() {
	flag.Parse()

	if flag.Lookup("port").DefValue == *httpPort {
		log.Fatalf("No http port specified to listen on.")
	}

	if flag.Lookup("brain-hostport").DefValue == *brainHostPort {
		log.Fatalf("No brain host:port specified to connect to.")
	}

	if flag.Lookup("auth-token").DefValue == *authToken {
		log.Fatalf("No auth token specified.")
	}

	fe = frontend.NewFrontend(*brainHostPort, *authToken)
	fe.Start()

	http.HandleFunc("/debug/health", healthHandler)

	log.Printf("Serving /debug on %v", *httpPort)
	log.Fatalf("Error during ListenAndServe: %v", http.ListenAndServe(*httpPort, nil))
}

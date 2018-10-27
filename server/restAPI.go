package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	)

type RestAPI struct {
}

type Address struct {
	Address string
}

type VoteReg struct {
	S_timestamp int64
	E_timestamp int64
	Name string
	Meta string
	Candidate []string
}

type Voting struct {
	Account string
	Voting  string
}

func nodeDiscovery(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var address Address

	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		log.Fatal(err)
	}

	Connect(address.Address)
}

func nodeList(w http.ResponseWriter, r *http.Request) {
	for _, node := range NodeList.NodeList {
		log.Printf("[list] %s:%s \n", node.IP, node.Port)
	}
	defer r.Body.Close()
}

func voteReg(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var voteReg VoteReg

	if err := json.NewDecoder(r.Body).Decode(&voteReg); err != nil {
		log.Fatal(err)
	}
}

func voting(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var voting Voting

	if err := json.NewDecoder(r.Body).Decode(&voting); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "voting : " + voting.Voting)
	fmt.Fprintf(w, "account : %x", []byte(voting.Account)) // byte로 어떻게 받아오는지
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func (rest RestAPI) handleRequest(restPort string) {
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/peers", nodeDiscovery).Methods("POST")
	r.HandleFunc("/list", nodeList).Methods("GET")
	r.HandleFunc("/voting", voting).Methods("POST")
	r.HandleFunc("/voteReg", voteReg).Methods("POST")
	semiTcpPort := ":" + restPort
	if err := http.ListenAndServe(semiTcpPort, r); err != nil {
		log.Fatal(err)
	}
}

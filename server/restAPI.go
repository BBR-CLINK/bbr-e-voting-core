package server

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
		)

type RestAPI struct{

}

type Address struct{
	Address string
}


func nodeDiscovery(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	var address Address

	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		log.Fatal(err)
	}

	Connect(address.Address)
}

func nodeList(w http.ResponseWriter, r *http.Request){
	for _, node := range NodeList.NodeList {
		log.Printf("[list] %s:%s \n", node.IP, node.Port)
	}
	defer r.Body.Close()
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func(rest RestAPI) handleRequest(restPort string){
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/peers", nodeDiscovery).Methods("POST")
	r.HandleFunc("/list", nodeList).Methods("GET")
	semiTcpPort := ":" + restPort
	if err := http.ListenAndServe(semiTcpPort,r); err != nil {
		log.Fatal(err)
	}
}

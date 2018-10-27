package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"bbr-e-voting-core/blockchain"
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
	Account []byte
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
	candidate := [][]byte{}

	if err := json.NewDecoder(r.Body).Decode(&voteReg); err != nil {
		log.Fatal(err)
	}

	for _, value := range voteReg.Candidate{
		candidate = append(candidate, []byte(value))
	}

	voteType := blockchain.VoteType{
		S_timestamp: voteReg.S_timestamp,
		E_timestamp: voteReg.E_timestamp,
		Name: []byte(voteReg.Name),
		Meta: []byte(voteReg.Meta),
		Candidate: candidate,
	}
	vote := blockchain.NewVote(blockchain.Account{}, []byte{}, &voteType)
	vote.SetID()

	lastBlock := Bc.GetLastBlock()
	currentBlock := blockchain.NewBlock([]*blockchain.Vote{vote}, lastBlock.Hash, lastBlock.Index)
	blockchain.BlockPool.Block = append(blockchain.BlockPool.Block, currentBlock)

	w.Write([]byte("VoteReg Success"))
}

func voting(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var voting Voting

	if err := json.NewDecoder(r.Body).Decode(&voting); err != nil {
		log.Fatal(err)
	}
	/*
	1. token 확인
	2. VoteType 확인
	3. Vote 설정
	 */

	//for _, node := range NodeList.NodeList {
	//	randNum := rand.Intn(len(NodeList.NodeList))
	//
	//}

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

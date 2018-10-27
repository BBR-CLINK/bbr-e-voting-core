package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"bbr-e-voting-core/blockchain"
	"strconv"
		"os"
	"github.com/gorilla/handlers"
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
	Meta string
	Voting  string
}

type VotingType struct {
	Meta string
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
	enableCors(&w)

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
	currentBlock := blockchain.NewBlock([]*blockchain.Vote{vote}, lastBlock.Hash, lastBlock.Index+1)
	blockchain.BlockPool.Block = append(blockchain.BlockPool.Block, currentBlock)
	Bc.AddBlock(currentBlock)
	w.Write([]byte("VoteReg Success"))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func voting(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var voting Voting

	if err := json.NewDecoder(r.Body).Decode(&voting); err != nil {
		log.Fatal(err)
	}

	//voteType, err := Bc.FindVoteReg([]byte(voting.Meta))
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//
	//
	//vote := blockchain.Vote{
	//	N_timestamp: time.Now().Unix(),
	//	Account: blockchain.Account{
	//		PublicKey: voting.Account,
	//		Token1: 0,
	//		Token2:
	//	}
	//}

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

func voteType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var votingType VotingType

	if err := json.NewDecoder(r.Body).Decode(&votingType); err != nil {
		log.Fatal(err)
	}
	enableCors(&w)

	setupResponse(&w, r)
	voteType, err := Bc.FindVoteReg([]byte(votingType.Meta))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(voteType)
	respondWithJSON(w, http.StatusOK, voteType)
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	parameter, found := mux.Vars(r)["index"]
	fmt.Println("parameter : ", parameter)
	fmt.Println("found : ", found)
	query, _ := r.URL.Query()["index"]
	index, _ := strconv.Atoi(query[0])
	fmt.Println(index)
	block := Bc.FindBlockByIndex(index)

	respondWithJSON(w, http.StatusOK, block)
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
	r.HandleFunc("/voteType", voteType).Methods("POST")
	r.HandleFunc("/getBlock", getBlock).Methods("GET")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	semiTcpPort := ":" + restPort
	//if err := http.ListenAndServe(semiTcpPort, r); err != nil {
	//	log.Fatal(err)
	//}
	log.Fatal(http.ListenAndServe(semiTcpPort, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
package blockchain

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"os"
	"log"
	"time"
	)

const dbFile = "blockchain_%s.db"
const genesisData = "C-LINK E-VOTING SYSTEM"

type Blockchain struct {
	tip []byte
	db *leveldb.DB
}

func CreateBlockchain(nodeAddress string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeAddress)
	if dbExists(dbFile) {
		log.Printf("[Blockchain] blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	voteType := &VoteType{0, 0, []byte(genesisData), []byte{}, [][]byte{}}
	vote := &Vote{time.Now().Unix(), Account{}, []byte{}, nil, voteType}
	genesis := CreateGenesisBlock(vote)

	db, err := leveldb.OpenFile(dbFile, nil)
	if err != nil {
		log.Printf("[Blockchain] Invalid dbFile")
		log.Panic(err)
	}

	err = db.Put(genesis.Hash, genesis.Serialize(), nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Put([]byte("1"), genesis.Hash, nil)
	if err != nil {
		log.Panic(err)
	}
	tip = genesis.Hash

	bc := Blockchain{tip, db}

	return &bc
}

// Load Blockchain. if dosen't exist, Create Blockchain
func LoadBlockchain(nodeAddress string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeAddress)
	if dbExists(dbFile) == false {
		log.Printf("[Blockchain] No existing blockchain found. Create one first.")
		os.Exit(1)
		// 여기에 blockchain 새로 생성하는거 해줘야 할까?
	}

	var tip []byte
	db, err := leveldb.OpenFile(dbFile, nil)
	if err != nil {
		log.Panic(err)
	}

	tip, err = db.Get([]byte("1"), nil)
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}

	return &bc
}

func (bc *Blockchain) AddBlock(block *Block){
	// 블록체인과 비교해서 최신것이 아닌지 판단하는거 넣어야됨
	err := bc.db.Put([]byte(block.Hash), block.Serialize(), nil)
	if err != nil {
		log.Panic(err)
	}

	err = bc.db.Put([]byte("1"), []byte(block.Hash), nil)
	if err != nil {
		log.Panic(err)
	}
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
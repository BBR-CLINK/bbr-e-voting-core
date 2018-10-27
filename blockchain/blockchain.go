package blockchain

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"os"
	"log"
	"time"
	)

var dbFile = "%s_blockData"
const genesisData = "C-LINK E-VOTING SYSTEM"

type Blockchain struct {
	tip []byte
	db *leveldb.DB
}

func CreateBlockchain(port string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, port)
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
func LoadBlockchain(port string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, port)
	if dbExists(dbFile) == false {
		log.Printf("[Blockchain] No existing blockchain found. Create one first.")
		return CreateBlockchain(port)
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

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetLastBlock() *Block {
	var lastBlock Block

	// db, err := leveldb.OpenFile(dbFile, nil)
	// defer db.Close()

	lastHash, err := bc.db.Get([]byte("1"), nil)
	if err != nil {
		log.Panic(err)
	}
	blockData, err := bc.db.Get(lastHash, nil)
	if err != nil {
		log.Panic(err)
	}

	lastBlock = *DeserializeBlock(blockData)

	return &lastBlock
}

func (bc *Blockchain) GetLastHash() []byte {
	lastHash, err := bc.db.Get([]byte("1"), nil)
	if err != nil {
		log.Panic(err)
	}

	return lastHash
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
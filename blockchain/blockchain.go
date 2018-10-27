package blockchain

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"os"
	"log"
	"time"
	"bbr-e-voting-core/util"
	"errors"
)

var dbFile = "%s_blockData"
const genesisData = "C-LINK E-VOTING SYSTEM"

type Blockchain struct {
	tip []byte
	db *leveldb.DB
}

type BlockchainIterator struct {
	currentHash []byte
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

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	lastHash := bc.GetLastHash()
	blockData, err := bc.db.Get(lastHash, nil)
	if err != nil {
		log.Panic(err)
	}
	lastBlock = *DeserializeBlock(blockData)

	return lastBlock.Index
}

// 블록체인 상의 블록 해쉬들을 반환
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	return blocks
}

// Iterator 체인의 끝을 가리킨다.
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// Next ...
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	encodedBlock, err := i.db.Get([]byte(i.currentHash), nil)
	if err != nil {
		log.Panic(err)
	}

	block = DeserializeBlock(encodedBlock)
	i.currentHash = block.PreviousHash

	return block
}

func (bc *Blockchain) FindVoteReg(meta []byte) (*VoteType, error) {
	bci := bc.Iterator()
	var voteType *VoteType
	for {
		block := bci.Next()

		if util.Equal(block.Votes[0].VoteType.Meta, meta) && time.Now().Unix() > block.Votes[0].VoteType.S_timestamp && time.Now().Unix() < block.Votes[0].VoteType.E_timestamp {
			voteType = block.Votes[0].VoteType
			return voteType, nil
		}

		if len(block.PreviousHash) == 0 {
			break;
		}
	}

	return nil, errors.New("VoteReg nil")
}
//
//func (bc *Blockchain) FindVoteToken(publicKey []byte, meta []byte) bool {
//	bci := bc.Iterator()
//	metaStr := string(meta[:])
//
//	for {
//		block := bci.Next()
//
//		if block.Votes[0].Voting == nil && util.Equal(block.Votes[0].Account.PublicKey, publicKey) && util.Equal(block.Votes[0].VoteType.Meta, meta) {
//			return true
//		}
//
//		if len(block.PreviousHash) == 0 {
//			return false
//		}
//	}
//
//}

func (bc *Blockchain) FindBlockByIndex(index int) *Block {
	bci := bc.Iterator()
	for {
		block := bci.Next()

		if block.Index == index{
			return block
		}
	}
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
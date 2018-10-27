package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"strconv"
	"fmt"
	"crypto/sha256"
)

type Block struct {
	Timestamp    int64
	PreviousHash []byte
	Hash         []byte
	Index        int
	Votes        []*Vote
}

func (b *Block) HashVotes() []byte {
	var votes [][]byte

	for _, vote := range b.Votes {
		votes = append(votes, vote.Serialize())
	}

	//merkle tree root 반환
	return []byte{} // 루트 반환
}

func NewBlock(votes []*Vote, previousHash []byte, index int) *Block {
	block := &Block{time.Now().Unix(), previousHash, []byte{}, index, votes }
	// poa := NewPoA()
	// 합의 알고리즘 : hash 반환
	return block
}

// Create genesis block
func CreateGenesisBlock(vote *Vote) *Block{
	return NewBlock([]*Vote{vote}, []byte{}, 0)
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10)) // b.Timestamp(type : Time)를 10진수로 변형
	headers := bytes.Join([][]byte{b.PreviousHash, b.HashVotes(), timestamp}, []byte{})
	fmt.Printf("hear %x\n", headers)
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func (b *Block) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

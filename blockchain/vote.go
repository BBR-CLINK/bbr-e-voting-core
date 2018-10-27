package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"bbrHack/util"
	"time"
)
type VoteType struct {
	S_timestamp int64 // 투표 시작
	E_timestamp int64 // 투표 끝
	Name        []byte // 투표이름
	Meta		[]byte // 해당 단과대에 대한 meta data, 총학(즉, 모두가 포함되는 투표)은 nil
	Candidate	[][]byte // 후보자
}

type Vote struct {
	N_timestamp int64 // 투표 시간
	Account     Account //누가
	Voting      []byte // 어디에 투표
	Hash		[]byte // Vote ID
	VoteType	*VoteType // VoteType
}

func NewVote(account Account, voting []byte, voteType *VoteType) *Vote{
	vote := &Vote{time.Now().Unix(), account, voting, nil, voteType}
	vote.SetID()
	return vote
}

func(v *Vote) Verify() error {
	S_timestamp := v.VoteType.S_timestamp
	E_timestamp := v.VoteType.E_timestamp
	//Name := v.VoteType.Name
	Candidate := v.VoteType.Candidate

	if !(v.N_timestamp > S_timestamp && v.N_timestamp < E_timestamp) {
		return ErrVoteTime
	}

	for idx, a  := range v.GetID() {
		if v.Hash[idx] != a {
			return ErrVoteHash
		}
	}

	for _, candidate := range Candidate {
		if !util.Equal(candidate, v.Voting) {
			return ErrVoteCandidate
		}
	}

	return nil
}

func(v *Vote) GetID() []byte{
	var hash [32]byte

	vote := *v
	vote.Hash = []byte{}

	hash = sha256.Sum256(vote.Serialize())

	return hash[:]
}

func(v Vote) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(v)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	v.Hash = hash[:]
}

func(v Vote) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(v)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
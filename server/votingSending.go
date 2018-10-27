package server

import (
	"bbr-e-voting-core/blockchain"
	"bbr-e-voting-core/node"
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"net"
)

type VotingSending struct {
}

func (be BlockExchange) VotingSending(ln net.Listener, connNeighbor node.Node, vote *blockchain.Vote) {
	go func() {
		connMe, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		//sendVote(block, connNeighbor)
		go handleVote(connMe, connNeighbor)
	}()
}

func sendVote(block *blockchain.Block, connNeighbor node.Node) {
	request := append(CommandToBytes("vote"), block.Serialize()...)
	log.Printf("[Block] Send Vote to %s", connNeighbor.IP)
	//spew.Dump(block)
	_, err := io.Copy(connNeighbor.Conn, bytes.NewReader(request))
	if err != nil {
		log.Panic(err)
	}
}

func handleVote(connMe net.Conn, connNeighbor node.Node) {
	request := make([]byte, 4096)

	n, err := connMe.Read(request)
	if err != nil {
		log.Panic(err)
	}

	request = request[:n]

	command := BytesToCommand(request[:CommandLength])

	switch command {
	case "vote":
		receiveVote(request, connNeighbor)
	}

}
func receiveVote(request []byte, connNeighbor node.Node) {
	var buff bytes.Buffer
	var payload blockchain.Block

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("[Block] Receive Block from %s", connNeighbor.IP)
	if !(payload.Timestamp == 0 && payload.PreviousHash == nil && payload.Hash == nil && payload.Index == 0 && payload.Votes == nil) {
		Bc.AddBlock(&payload)
	}
	//spew.Dump(payload)
}
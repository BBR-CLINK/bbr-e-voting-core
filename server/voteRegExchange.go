package server

import (
	"net"
	"bbr-e-voting-core/node"
	"log"
	"io"
	"bytes"
	"encoding/gob"
	"bbr-e-voting-core/blockchain"
)

type voteRegExchange struct {

}

func (vre voteRegExchange) VoteRegExchange(ln net.Listener, connNeighbor node.Node){

}

func sendVoteReg(block blockchain.Block, connNeighbor node.Node) {
	request := append(CommandToBytes("block"), block.Serialize()...)
	log.Printf("[Block] Send Block to %s", connNeighbor.IP)
	//spew.Dump(block)
	_, err := io.Copy(connNeighbor.Conn, bytes.NewReader(request))
	if err != nil {
		log.Panic(err)
	}
}

func handleVoteReg(connMe net.Conn, connNeighbor node.Node) {
	request := make([]byte, 4096)

	n, err := connMe.Read(request)
	if err != nil {
		log.Panic(err)
	}
	request = request[:n]

	command := BytesToCommand(request[:CommandLength])

	switch command {
	case "block":
		receiveBlock(request, connNeighbor)
	}

}
func receiveVoteReg(request []byte, connNeighbor node.Node) {
	var buff bytes.Buffer
	var payload blockchain.Block

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("[Block] Receive Block from %s", connNeighbor.IP)
	//spew.Dump(payload)
}
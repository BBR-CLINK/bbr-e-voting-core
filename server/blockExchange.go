package server

import (
	"bbrHack/blockchain"
	"bbrHack/node"
	"bytes"
	"encoding/gob"
		"io"
	"log"
	"net"
		"github.com/davecgh/go-spew/spew"
)

type BlockExchange struct {
}

func (be BlockExchange) BlockExchange(ln net.Listener, connNeighbor node.Node, block *blockchain.Block) {
	go func() {
		connMe, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		sendBlock(block, connNeighbor)
		go handleBlock(connMe, connNeighbor)
	}()
}

func sendBlock(block *blockchain.Block, connNeighbor node.Node) {
	request := append(CommandToBytes("block"), block.Serialize()...)
	log.Printf("[Block] Send Block to %s", connNeighbor.IP)
	spew.Dump(block)
	_, err := io.Copy(connNeighbor.Conn, bytes.NewReader(request))
	if err != nil {
		log.Panic(err)
	}
}

func handleBlock(connMe net.Conn, connNeighbor node.Node) {
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
func receiveBlock(request []byte, connNeighbor node.Node) {
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
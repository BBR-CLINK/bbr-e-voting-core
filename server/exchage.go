package server

import (
	"bbrHack/blockchain"
	"bbrHack/node"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

const commandLength = 12

var mutex = &sync.Mutex{}

type BlockExchange struct {
}

func (be BlockExchange) DataExchange(ln net.Listener, connNeighbor node.Node) {
	go func() {

		connMe, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		//sendBlock(block, connNeighbor)
		go handleBlock(connMe, connNeighbor)

	}()
}

func sendBlock(block blockchain.Block, connNeighbor node.Node) {
	request := append(commandToBytes("block"), block.Serialize()...)
	log.Printf("[Block] Send Block to %s", connNeighbor.IP)
	//spew.Dump(block)
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

	command := bytesToCommand(request[:commandLength])

	switch command {
	case "block":
		receiveBlock(request, connNeighbor)
	}

}
func receiveBlock(request []byte, connNeighbor node.Node) {
	var buff bytes.Buffer
	var payload blockchain.Block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("[Block] Receive Block from %s", connNeighbor.IP)
	//spew.Dump(payload)
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

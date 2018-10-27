package node

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type Node struct {
	IP   string
	Port string
	Conn net.Conn
}

var Seed = "172.16.5.29:3000"

func (node Node) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(node)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

package server

import (
	"net"
	"log"
	"strings"
	"bbrHack/node"
	"bbrHack/util"
)

func Connect(nodeAddress string) {
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		log.Printf("[connect] %s is not available\n", nodeAddress)
	} else {

		node := node.Node{
			IP:   strings.Split(nodeAddress, ":")[0],
			Port: strings.Split(nodeAddress, ":")[1],
			Conn: conn,
		}

		if util.NodeExists(NodeList, node) {
			log.Printf("[connect] %s is already connected", nodeAddress)
		} else {
			log.Printf("[connect] Success to handshake : %s", nodeAddress)
			NodeList.NodeList = append(NodeList.NodeList, node)
		}
	}
}
package server

import (
	"bbr-e-voting-core/node"
	"log"
	"net"
	"sync"

		"bbr-e-voting-core/blockchain"
	"time"
	"fmt"
)

var mutex = &sync.Mutex{}
var NodeList = node.NodeList{} // 전역 변수 어떻게 없애지
var Bc = &blockchain.Blockchain{}
var Ln = struct{
	ln net.Listener
}{}

func StartServer(tcpPort string, restPort string) {
	nodeIP := GetOutboundIP()
	log.Printf("Start with : %s:%s ", nodeIP, tcpPort)
	//nodeAddress := fmt.Sprintf("%s:%s", nodeIP, tcpPort)
	semiTcpPort := ":" + tcpPort
	ln, err := net.Listen("tcp", semiTcpPort)
	if err != nil {
		log.Panic(err)
	}

	Ln.ln = ln // 글로벌 극혐 코드 망함

	defer ln.Close()

	Bc = blockchain.LoadBlockchain(tcpPort)

	time.Sleep(5 * time.Second)

// 이걸 어떻게 넣어서 연결시킬지. 그냥 글로벌하게 할까...
	go func() {
			for {
				time.Sleep(5 * time.Second)
				mutex.Lock()
				fmt.Println("size : %d	",len(blockchain.BlockPool.Block))
				if len(blockchain.BlockPool.Block) == 0  {
					exchange(ln, NodeList, &blockchain.Block{})
				} else {
					exchange(ln, NodeList, blockchain.BlockPool.Block[0])
					blockchain.BlockPool.Block = append(blockchain.BlockPool.Block[:0], blockchain.BlockPool.Block[1:]...)
					fmt.Println("size : %d	",len(blockchain.BlockPool.Block))
				}
				mutex.Unlock()
			}
	}()

	restAPI := RestAPI{}
	restAPI.handleRequest(restPort)
}

func exchange(ln net.Listener, nodeList node.NodeList, block *blockchain.Block) {
	dataExchange := BlockExchange{}
	// 교환 주기 생각할 것
	for i := 0; i < len(nodeList.NodeList); i++ {
		address := nodeList.NodeList[i].IP + ":" +nodeList.NodeList[i].Port
		nodeList.NodeList[i].Conn = Connect(address)
		dataExchange.BlockExchange(ln, nodeList.NodeList[i], block)
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

package server

import (
	"bbrHack/node"
	"log"
	"net"
	"sync"
	"time"
		)

var mutx = &sync.Mutex{}
var NodeList = node.NodeList{} // 전역 변수 어떻게 없애지

func StartServer(tcpPort string, restPort string) {
	nodeIP := GetOutboundIP()
	log.Printf("Start with : %s:%s ", nodeIP, tcpPort)
	//nodeAddress := fmt.Sprintf("%s:%s", nodeIP, tcpPort)
	semiTcpPort := ":" + tcpPort
	ln, err := net.Listen("tcp", semiTcpPort)
	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()

	restAPI := RestAPI{}
	restAPI.handleRequest(restPort)
//	bc := blockchain.CreateBlockchain(nodeAddress)
// 이걸 어떻게 넣어서 연결시킬지. 그냥 글로벌하게 할까...
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			exchange(ln, NodeList)
			mutex.Unlock()
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}

}

func exchange(ln net.Listener, nodeList node.NodeList) {

	dataExchange := BlockExchange{}
	// 교환 주기 생각할 것
	for i := 0; i < len(nodeList.NodeList); i++ {
		dataExchange.DataExchange(ln, nodeList.NodeList[i])
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

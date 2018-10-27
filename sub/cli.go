package sub

import (
	"github.com/urfave/cli"
			"bbrHack/server"
)

//var StartCmd = cli.Command{
//	Name: "start",
//	Usage: "CLINK start [port]",
//	Action: func(c *cli.Context) error {
//		nodeID := GetOutboundIP()
//		port := c.Args().Get(0)
//		server.StartServer(nodeID, port)
//		return nil
//	},
//}

var ConnectCmd = cli.Command{
	Name: "connect",
	Usage: "CLINK connect [nodeAddress]",
	Action: func(c *cli.Context) error {
		nodeAddress := c.Args().Get(0)
		server.Connect(nodeAddress)
		return nil
	},
}
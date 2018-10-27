package main

import (
	"bbrHack/server"
	"bbrHack/sub"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "CLINK"
	app.Usage = "BBR Hackathon"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "tcpPort, tp",
			Value: "",
			Usage: "set node tcpPort",
		},
		cli.StringFlag{
			Name:  "restPort, rp",
			Value: "",
			Usage: "set node restPort",
		},
	}
	app.Commands = []cli.Command{}
	//app.Commands = append(app.Commands, sub.StartCmd)
	app.Commands = append(app.Commands, sub.ConnectCmd)
	/*
		connect, disconnect, 정도만
	*/
	app.Before = func(c *cli.Context) error {
		tcpPort := c.String("tcpPort")
		restPort := c.String("restPort")

		if tcpPort != "" && restPort != "" {
			server.StartServer(tcpPort, restPort)
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"os"

	"github.com/MohamedBassem/DSSS/agent"
	"github.com/MohamedBassem/DSSS/client"
	"github.com/MohamedBassem/DSSS/server"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "DSSS"
	app.Usage = "To start agent/server/download/upload"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Starts the server",
			Action: func(c *cli.Context) {
				logger := log.New(os.Stdout, "Server ", log.LstdFlags)
				server.Main(8082, 8081, logger)
			},
		},

		{
			Name:  "agent",
			Usage: "Starts the Agent",
			Action: func(c *cli.Context) {
				logger := log.New(os.Stdout, "Agent ", log.LstdFlags)
				agent.Main(logger)
			},
		},
		{
			Name:  "upload",
			Usage: "<file_name> <output_mainifist_name>",
			Action: func(c *cli.Context) {
				logger := log.New(os.Stdout, "Client ", log.LstdFlags)
				client.Upload(c.Args().First(), c.Args().Get(1), logger)
			},
		},
		{
			Name:  "download",
			Usage: "Downloads a file",
			Action: func(c *cli.Context) {
				logger := log.New(os.Stdout, "Client ", log.LstdFlags)
				client.Download(c.Args().First(), c.Args().Get(1), logger)
			},
		},
	}
	app.HideVersion = true

	app.RunAndExitOnError()
}

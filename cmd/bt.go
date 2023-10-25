package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	commands := []*cli.Command{
		{
			Name:   "info",
			Usage:  "Parse and print out information about a torrent file",
			Action: info,
		},
		{
			Name:   "peers",
			Usage:  "Connect to peers and exit",
			Action: peers,
		},
	}

	app := cli.App{
		Name:     "bt",
		Usage:    "A BitTorrent client",
		Commands: commands,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/yinfredyue/bittorrent-go/client"
	"github.com/yinfredyue/bittorrent-go/torrent"
)

func peers(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("expect [Torrent_filepath] argument")
	}

	torrentFilepath := ctx.Args().Get(0)
	torrent, err := torrent.OfFile(torrentFilepath)
	if err != nil {
		return err
	}

	client := client.NewClient(torrent)
	return client.ConnectToPeers()
}

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/yinfredyue/bittorrent-go/torrent"
)

func info(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("expect [Torrent_filepath] argument")
	}

	torrentFilepath := ctx.Args().Get(0)
	torrent, err := torrent.OfFile(torrentFilepath)
	if err != nil {
		return err
	}

	fmt.Printf("Tracker: %v\n", torrent.Tracker)
	fmt.Printf("Length: %v\n", torrent.Info.Length)
	fmt.Printf("Piece length: %v\n", torrent.Info.PieceLength)
	fmt.Printf("Piece hashes:\n")
	for _, pieceHash := range torrent.Info.PieceHashes {
		fmt.Printf("%v\n", hex.EncodeToString(pieceHash[:]))
	}

	return nil
}

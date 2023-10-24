package torrent

import (
	"bytes"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
	"github.com/yinfredyue/bittorrent-go/util"
)

const (
	pieceHashLength = 20
)

type Info struct {
	Name        string
	PieceLength int
	PieceHashes []([pieceHashLength]byte)
	Length      int
}

type Torrent struct {
	Tracker string
	Info    Info
}

func OfBytes(b []byte) (Torrent, error) {
	decodeErr := fmt.Errorf("fail to decode torrent file")

	data, err := bencode.Decode(bytes.NewReader(b))
	if err != nil {
		return Torrent{}, err
	}

	torrentDict, ok := data.(map[string]interface{})
	if !ok {
		return Torrent{}, decodeErr
	}

	tracker, ok := torrentDict["announce"].(string)
	if !ok {
		return Torrent{}, decodeErr
	}
	util.DPrintf("Tracker: %v\n", tracker)

	infoDict, ok := torrentDict["info"].(map[string]interface{})
	if !ok {
		return Torrent{}, decodeErr
	}

	name, ok := infoDict["name"].(string)
	if !ok {
		return Torrent{}, decodeErr
	}
	util.DPrintf("Name: %v\n", name)

	pieceLength, ok := infoDict["piece length"].(int64)
	if !ok {
		return Torrent{}, decodeErr
	}
	util.DPrintf("Piece length: %v\n", pieceLength)

	pieceHashesBytes, ok := infoDict["pieces"].(string)
	if !ok {
		return Torrent{}, decodeErr
	}
	util.DPrintf("len(Pieces): %v\n", len(pieceHashesBytes))

	numPieces := len(pieceHashesBytes) / pieceHashLength
	pieceHashes := make([]([pieceHashLength]byte), numPieces)
	for i := 0; i < numPieces; i++ {
		bytes := pieceHashesBytes[i*pieceHashLength : (i+1)*pieceHashLength]
		util.AssertEqual(len(bytes), pieceHashLength)
		copy(pieceHashes[i][:], bytes)
	}

	length, ok := infoDict["length"].(int64)
	if !ok {
		return Torrent{}, decodeErr
	}
	util.DPrintf("Length: %v\n", length)

	torrent := Torrent{
		Tracker: tracker,
		Info: Info{
			Name:        name,
			PieceLength: int(pieceLength),
			PieceHashes: pieceHashes,
			Length:      int(length),
		},
	}

	return torrent, nil
}

func OfFile(filename string) (Torrent, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Torrent{}, err
	}

	return OfBytes(content)
}

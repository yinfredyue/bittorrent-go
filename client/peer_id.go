package client

import "fmt"

var nextId = 0

func NewPeerId() []byte {
	nextId++
	return []byte(fmt.Sprintf("%02dDEADBEEF_ERROR_404", nextId))
}

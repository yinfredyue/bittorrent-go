package util

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
)

// Debugging
const Debug = 1

func DPrintf(format string, a ...interface{}) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
}

func AssertEqual(left interface{}, right interface{}) {
	if !reflect.DeepEqual(left, right) {
		errMsg := fmt.Sprintf("unequal: left=%v, right=%v", left, right)
		panic(errMsg)
	}
}

func EncodeUint32(value uint32, length int) ([]byte, error) {
	if length < 4 {
		return nil, fmt.Errorf("buffer length < 4")
	}

	buffer := make([]byte, length)
	binary.BigEndian.PutUint32(buffer[length-4:], value)

	return buffer, nil
}

func ConcatBytes(bs []([]byte)) []byte {
	res := make([]byte, 0)
	for _, b := range bs {
		res = append(res, b...)
	}
	return res
}

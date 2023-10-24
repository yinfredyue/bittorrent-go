package util

import (
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

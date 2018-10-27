package server

import (
	"sync"
	"fmt"
)

const CommandLength = 12
var flag bool
var Mutex = &sync.Mutex{}


func CommandToBytes(command string) []byte {
	var bytes [CommandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func BytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

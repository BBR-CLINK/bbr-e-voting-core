package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {

}

func (b Block) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

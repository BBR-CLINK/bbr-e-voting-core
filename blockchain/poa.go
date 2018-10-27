package blockchain

import "bbrHack/util"

type PoA struct {
	block *Block
}

func NewPoA() *PoA{
	poa := &PoA{}
	return poa
}

func (poa PoA) Validate(currentHash []byte) error {
	if !util.Equal(currentHash, poa.block.PreviousHash){
		return ErrBlockPreviousHash
	}
	return nil
}
package blockchain

type PoA struct {
	block *Block
}

func NewPoA(block *Block) *PoA{
	poa := &PoA{
		block: block,
	}
	return poa
}

func (poa PoA) Validate() error {
	//if !util.Equal(currentHash, poa.block.PreviousHash){
	//	return ErrBlockPreviousHash
	//}
	return nil
}
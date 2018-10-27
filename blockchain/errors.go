package blockchain

import "errors"

var ErrVoteTime = errors.New("Invalid Vote Time")
var ErrVoteHash = errors.New("Invalid Vote Hash")
var ErrVoteCandidate = errors.New("Invalid Vote Candidate")
var ErrBlockPreviousHash = errors.New("Invalid Previous Hash")
package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
)

type snowflake struct {
	epoch     int64
	nodeID    uint64
	lastTsSeq atomic.Uint64
}

const (
	nodeBits     = 10
	sequenceBits = 12
	maxSequence  = (1 << sequenceBits) - 1

	timeShift = nodeBits + sequenceBits
	nodeShift = sequenceBits
)

var (
	sfInstance *snowflake
	once       sync.Once
)

func NewSnowFlake(epoch int64, nodeID uint64) *snowflake {
	once.Do(func() {
		sfInstance = &snowflake{
			epoch:  epoch,
			nodeID: nodeID,
		}
	})
	return sfInstance
}

func GenerateID() uint64 {
	if sfInstance == nil {
		NewSnowFlake(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(), 1)
	}
	return sfInstance.Generate()
}

func (s *snowflake) Generate() uint64 {
	for {
		now := uint64(time.Now().UnixMilli() - s.epoch)
		packed := s.lastTsSeq.Load()
		lastTs := packed >> sequenceBits
		seq := packed & maxSequence

		if now == lastTs {
			if seq >= maxSequence {
				for now <= lastTs {
					now = uint64(time.Now().UnixMilli() - s.epoch)
				}
				seq = 0
			} else {
				seq++
			}
		} else {
			seq = 0
		}

		newPacked := (now << sequenceBits) | seq
		if s.lastTsSeq.CompareAndSwap(packed, newPacked) {
			id := (now << timeShift) | (s.nodeID << nodeShift) | seq
			return id
		}
	}
}

func Generate6DigitCode() (int, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number")
	}

	return int(num.Int64()) + 100000, nil
}

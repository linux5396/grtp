package util

import (
	"sync"
	"time"
)

//base on twitter 's snowflake algorithm to dispatch the global-id
type IdWorker struct {
	//with lock
	sync.Mutex
	workerId     int64 //need to be init
	dataCenterId int64 //need to be init
	sequence     int64 //need to be init
	//below fields value should set private to keep right
	workerIdBits     int64
	dataCenterIdBits int64
	maxWorkerId      int64
	maxCenterId      int64
	sequenceBits     int64
	epoch            int64
	//next is some fields composite
	workerIdShift      int64
	dataCenterIdShift  int64
	timestampLeftShift int64
	sequenceMask       int64
	//record last Timestamp
	lastTimestamp int64
}

func NewIdWorker(workerId, dataCenterId, sequence int64) *IdWorker {
	//FIXME lack of verifying the input args
	worker := new(IdWorker)
	worker.workerId = workerId
	worker.dataCenterId = dataCenterId
	worker.sequence = sequence
	//default set
	worker.lastTimestamp = -1
	worker.epoch = 1288834974657
	worker.workerIdBits = 5
	worker.dataCenterIdBits = 5
	worker.maxWorkerId = (-1) ^ ((-1) << worker.workerIdBits)
	worker.maxCenterId = (-1) ^ ((-1) << worker.dataCenterIdBits)
	worker.sequenceBits = 12
	//
	worker.workerIdShift = worker.sequenceBits
	worker.dataCenterIdShift = worker.sequenceBits + worker.workerIdBits
	worker.timestampLeftShift = worker.sequenceBits + worker.workerIdBits + worker.dataCenterIdBits
	worker.sequenceMask = (-1) ^ ((-1) << worker.sequenceBits)

	return worker
}
func (idWorker *IdWorker) NextId() int64 {
	idWorker.Lock()
	defer idWorker.Unlock()
	currentTimestamp := timeGen()
	if currentTimestamp < idWorker.lastTimestamp {
		panic("occurred error")
	}
	if idWorker.lastTimestamp == currentTimestamp {
		idWorker.sequence = (idWorker.sequence + 1) & idWorker.sequenceMask
		if idWorker.sequence == 0 {
			currentTimestamp = tilNextMillis(idWorker.lastTimestamp)
		}
	} else {
		idWorker.sequence = 0
	}
	idWorker.lastTimestamp = currentTimestamp
	return ((currentTimestamp - idWorker.epoch) << idWorker.timestampLeftShift) |
		(idWorker.dataCenterId << idWorker.dataCenterIdShift) |
		(idWorker.workerId << idWorker.workerIdShift) |
		idWorker.sequence
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

//get timestamp
func timeGen() int64 {
	return time.Now().Unix()
}

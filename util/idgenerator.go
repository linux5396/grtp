package util

import (
	"sync"
	"time"
)

/**
*《SNOWFLAKE》
*1、生成的分布式ID是一个64bits的整数，结构如下：
*|不用|                     41bit的时间戳               |10bit的机器号|12bit序列号|
*| 0  |-00000000-00000000-00000000-00000000-00000000-0-|00000000-00|00000000-0000|
*2、description:
*@ 0bit:由于二进制中最高位是1则代表负数，但是ID一般都是用整数，因此，最高位不用
*@ 41bit:时间戳（毫秒），41位可以表示2^41 - 1个数字，可以用69年 :(2^41-1)/(1000*60*60*24*365)=69 years
*@ 10bit机器：机器由dataCenterId和workerId确认，因此，最多支持1024个节点
*@ 12bit序列号：用来记录同毫秒内产生的不同id。12位支持4095，因此，同机器同一时间戳（MS）产生4095个ID
*3、snowflake promises that all ID is increased by time,and there is not any duplicated ID in the distributed system,because the dataCenterID & workerId diff.
 */

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
	epoch            int64 //begin offset timestamp
	//next is some fields composite
	workerIdShift      int64
	dataCenterIdShift  int64
	timestampLeftShift int64
	sequenceMask       int64
	//record last Timestamp
	lastTimestamp int64
}

//the args I didn't verified,if you ....
//the bit range can modify by your specified need.
func NewIdWorker(workerId, dataCenterId, sequence int64) *IdWorker {
	//FIXME lack of verifying the input args
	worker := new(IdWorker)
	worker.workerId = workerId
	worker.dataCenterId = dataCenterId
	worker.sequence = sequence
	//default set
	worker.lastTimestamp = -1
	worker.epoch = 0
	worker.workerIdBits = 5
	worker.dataCenterIdBits = 5
	worker.maxWorkerId = (-1) ^ ((-1) << worker.workerIdBits)     //通过补码（-1）异或（-1）左移5位，等低5位，即31
	worker.maxCenterId = (-1) ^ ((-1) << worker.dataCenterIdBits) //通过补码（-1）异或（-1）左移5位，等低5位，即31
	worker.sequenceBits = 12
	//
	worker.workerIdShift = worker.sequenceBits                                                      //12
	worker.dataCenterIdShift = worker.sequenceBits + worker.workerIdBits                            //17
	worker.timestampLeftShift = worker.sequenceBits + worker.workerIdBits + worker.dataCenterIdBits //22
	worker.sequenceMask = (-1) ^ ((-1) << worker.sequenceBits)                                      //低12位掩码
	return worker
}
func (idWorker *IdWorker) NextId() int64 {
	idWorker.Lock()
	defer idWorker.Unlock()
	currentTimestamp := timeGen()
	if currentTimestamp < idWorker.lastTimestamp {
		panic("occurred error")
	}
	//if last timestamp equals cur timestamp,mean it should increased by seq
	if idWorker.lastTimestamp == currentTimestamp {
		//mask is 4095 with 12bit is 1.
		//this (idWorker.sequence + 1) & idWorker.sequenceMask can promises the range is from 0	to 4095
		idWorker.sequence = (idWorker.sequence + 1) & idWorker.sequenceMask
		if idWorker.sequence == 0 {
			currentTimestamp = tilNextMillis(idWorker.lastTimestamp)
		}
	} else {
		idWorker.sequence = 0
	}
	idWorker.lastTimestamp = currentTimestamp
	//currentTimestamp - epoch as init Id range
	//below code is core of snow flake.
	return ((currentTimestamp - idWorker.epoch) << idWorker.timestampLeftShift) | //move shift to promise it in the true hold.
		(idWorker.dataCenterId << idWorker.dataCenterIdShift) | //centerId move shift to promise it in the true hold.
		(idWorker.workerId << idWorker.workerIdShift) | //workerId move shift to promise it in the true hold.
		idWorker.sequence //because the sequence id is the rightest bits .
	//By the end,all bits collected by the OR operation to produce the final ID
}

//solve the
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

package buffer

import (
	"errors"
	"grtp/util"
)

var ErrorNull = errors.New("ring buffer is null,can't use")

//these thresholds are based on the way:
//machine memory is 16 GB
//numbers of buffer is C10K
//at least cost 19.53125 GB
//so I design these thresholds, if 0.75 is free,compact it quickly.
//TODO make these thresholds can modify by caller.And with a global starting button that decide these use or not.
const (
	GcThreshold           = 1 << 21 //2MB
	BusyPercentThreshold  = 0.25
	LowestNarrowThreshold = 64
)

//define a ring buffer
type RingBuffer struct {
	buf     []byte //storage
	size    int    //buf size
	mask    int    //mask
	r       int    //next r position
	w       int    //next w position
	isEmpty bool   //now is empty or not
}

func New(size int) *RingBuffer {
	if size == 0 {
		return &RingBuffer{isEmpty: true}
	}
	if size < 0 {
		panic("size can not be less than 0")
	}
	size = util.Ceil(size)
	return &RingBuffer{
		buf:     make([]byte, size),
		size:    size,
		mask:    size - 1, //because the size is power of 2,so size -1 mean low 1 bit with all bit'1'
		isEmpty: true,
	}
}

//Caller should promise that the length of bytes can read the buffer bytes completely.
//And then ,this ring buffer is reusable,we read it by moving the rb.r position.
//The data  is available can read only one times.
//So if you take the readable bytes from it and you will use it once and once again,you may cache the dst bytes.
//this method returns length of read,and err maybe occurred.
func (rb *RingBuffer) Read(bytes []byte) (n int, err error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	if rb.isEmpty {
		return 0, ErrorNull
	}
	//inRange read or Named flat read
	if rb.w > rb.r {
		n = rb.w - rb.r //cnt readable length
		if n > len(bytes) {
			n = len(bytes) //only fill bytes full.
		}
		copy(bytes, rb.buf[rb.r:rb.r+n]) //copy
		rb.r += n
		if rb.r == rb.w {
			rb.isEmpty = true
		}
		//check narrow.
		if rb.size >= GcThreshold && (float64(rb.Available())/float64(rb.size)) < BusyPercentThreshold {
			rb.narrow()
		}
		return
	}
	//ring read
	//example: size=8 r=5 w=1 so readable bytes is from 5,6,7,0   equals 8-5+1=4
	n = rb.size - rb.r + rb.w //cnt readable length
	if n > len(bytes) {
		n = len(bytes)
	}
	if rb.r+n <= rb.size { //although w<r,means ring,but the readSize satisfied less than size.
		copy(bytes, rb.buf[rb.r:rb.r+n]) //only copy
	} else {
		//read ring. may regards the buffer as 2 partitions
		r1 := rb.size - rb.r
		copy(bytes, rb.buf[rb.r:]) //first part
		r2 := n - r1
		copy(bytes[r1:], rb.buf[:r2]) //second part
	}
	//it is equaled with rb.r=rb.r+n-rb.size
	rb.r = (rb.r + n) & rb.mask //by & mask which is size-1,so if rb.r +n  was  overflow,can reset to the begin of ring
	if rb.r == rb.w {
		rb.isEmpty = true
	}
	//check narrow.
	if rb.size >= GcThreshold && (float64(rb.Available())/float64(rb.size)) < BusyPercentThreshold {
		rb.narrow()
	}
	return n, err
}

//read in the view of buf.never move the read offset.
//supporting is by copying.
//so the return bytes is the replication.
//modify the bytes will not modify the RingBuffer.buf
func (rb *RingBuffer) ReadView(size int) (bytes []byte, err error) {
	if rb.isEmpty || size <= 0 {
		return
	}
	bytes = make([]byte, size) //make
	//flat read
	if rb.r < rb.w {
		n := rb.w - rb.r //len
		if n > size {
			n = size
		}
		copy(bytes, rb.buf[rb.r:rb.r+n])
		return
	}
	n := rb.size - rb.r + rb.w //ring read
	if n > size {
		n = size
	}
	//still flat
	if rb.r+n <= rb.size {
		copy(bytes, rb.buf[rb.r:rb.r+n])
	} else {
		//ring read
		first := rb.size - rb.r
		copy(bytes, rb.buf[rb.r:])
		next := n - first
		copy(bytes[first:], rb.buf[:next])
	}
	return
}

func (rb *RingBuffer) Write(source []byte) (n int, err error) {
	n = len(source)
	if n == 0 {
		return 0, nil
	}
	//get the remaining of buffer
	if free := rb.Remain(); free < n {
		//not enough .so call resize()
		rb.enlarge(n - free)
	}
	//w >=r
	if rb.w >= rb.r {
		fromWtoEnd := rb.size - rb.w
		if fromWtoEnd >= n {
			copy(rb.buf[rb.w:], source)
			rb.w += n
		} else {
			copy(rb.buf[rb.w:], source[:fromWtoEnd])
			rare := n - fromWtoEnd
			copy(rb.buf, source[fromWtoEnd:])
			rb.w = rare
		}
	} else {
		//enough and w<r, always enough!!
		//because if not enough,the Remain had called resize().
		copy(rb.buf[rb.w:], source)
		rb.w += n
	}
	if rb.w == rb.size {
		rb.w = 0
	}
	rb.isEmpty = false
	return n, err
}

//remain of buffer
//if r==w   mean empty or 0
//if r<w    like ring read:rb.size - rb.w + rb.r
//if w<r    only can write the range from w to r-1,the length is r-w
func (rb *RingBuffer) Remain() int {
	if rb.r == rb.w {
		if rb.isEmpty {
			return rb.size
		}
		return 0
	}
	//remain
	if rb.w < rb.r {
		return rb.r - rb.w
	}
	return rb.size - rb.w + rb.r //rb.w>rb.r
}

// Available() returns the length of available read bytes.
func (rb *RingBuffer) Available() int {
	if rb.r == rb.w {
		if rb.isEmpty {
			return 0
		}
		return rb.size
	}
	if rb.w > rb.r {
		return rb.w - rb.r
	}
	return rb.size - rb.r + rb.w
}

//return buffer.buf len
func (rb *RingBuffer) Len() int {
	return len(rb.buf)
}

//return buffer capacity
func (rb *RingBuffer) Cap() int {
	return rb.size
}

//enlarge and move the old bytes to the new array.
func (rb *RingBuffer) enlarge(cap int) {
	newCapacity := util.Ceil(rb.size + cap)
	newBuf := make([]byte, newCapacity)
	oldLength := rb.Available() //the available bytes which never read.
	_, _ = rb.Read(newBuf)
	rb.r = 0
	rb.w = oldLength
	rb.size = newCapacity
	rb.mask = newCapacity - 1
	rb.buf = newBuf
}

//if buf remains above 75% space and Cap() is more than 2MB
func (rb *RingBuffer) narrow() {
	newCapacity := util.Ceil(rb.Available())
	if newCapacity < LowestNarrowThreshold {
		newCapacity = LowestNarrowThreshold
	}
	newBuf := make([]byte, newCapacity)
	oldLength := rb.Available()
	_, _ = rb.Read(newBuf)
	rb.r = 0
	rb.w = oldLength
	rb.size = newCapacity
	rb.mask = newCapacity - 1
	rb.buf = newBuf
}

package tunnel

import "sync/atomic"

func Nextp(ports []string, next *uint32) *string {

	n := atomic.AddUint32(next, 1)
	return &ports[(int(n)-1)%len(ports)]
}

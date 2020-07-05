package tracker

import (
	"strconv"
	_ "unsafe"
)

func IPPort(ip string, port uint16) string {
	return ip + ":" + strconv.Itoa(int(port))
}

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

func Now() uint64 {
	return uint64(nanotime())
}

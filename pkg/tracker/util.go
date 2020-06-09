package tracker

import (
	"strconv"
)

func FormatCID(cid ConnectionID) string {
	return cid.SAddr + ":" + strconv.Itoa(int(cid.SPort)) + "-" + cid.DAddr + ":" + strconv.Itoa(int(cid.DPort))
}

func IPPort(ip string, port uint16) string {
	return ip + ":" + strconv.Itoa(int(port))
}

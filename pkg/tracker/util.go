package tracker

import(
	"strconv"
)

func FormatCID(cid ConnectionID) string {
	return cid.SAddr + ":" + strconv.Itoa(int(cid.SPort)) + "-" + cid.DAddr + ":" + strconv.Itoa(int(cid.DPort))
}

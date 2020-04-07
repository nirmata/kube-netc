package tracker

import(
	"strconv"
)

func FormatCID(cid ConnectionID) string {
	return cid.DAddr + ":" + strconv.Itoa(int(cid.DPort))
}

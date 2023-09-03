package Util

import "strconv"

func TimeStampTransfer(timeStamp string) int64 {
	var numTime int64
	if timeStamp != "0" {
		timeStamp = timeStamp[0:10]
	}
	numTime, _ = strconv.ParseInt(timeStamp, 10, 64)
	return numTime
}

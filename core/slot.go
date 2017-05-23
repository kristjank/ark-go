package core

import "time"

//GetTime return time slot difference in secods. This timestamp is
//added to the transaction.
func GetTime() int32 {
	mainNetStart := time.Date(2017, 3, 21, 13, 00, 0, 0, time.UTC)
	now := time.Now()

	diff := now.Sub(mainNetStart)
	return int32(diff.Seconds())
}

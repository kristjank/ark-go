package core

import "time"

var mainNetStart = time.Date(2017, 3, 21, 13, 00, 0, 0, time.UTC)

//GetTime return time slot difference in secods. This timestamp is
//added to the transaction.
func GetTime() int32 {
	now := time.Now()
	diff := now.Sub(mainNetStart)

	return int32(diff.Seconds())
}

//Calculates duration between now and provided timestamp
func GetDurationTime(timestamp int32) int32 {
	var durationSeconds time.Duration = time.Duration(timestamp) * time.Second
	timeCalculcated := mainNetStart.Add(durationSeconds)

	now := time.Now()
	diff := now.Sub(timeCalculcated)

	return int32(diff.Hours())
}

package service

import "time"

func TimeTamperProofService(startTime string, deadline string) uint64 {
	startTimeTime, _ := time.Parse("2006-01-02 15:04:05.000", startTime)
	deadlineTime, _ := time.Parse("2006-01-02 15:04:05.000", deadline)
	return uint64(startTimeTime.Unix() + deadlineTime.Unix())
}

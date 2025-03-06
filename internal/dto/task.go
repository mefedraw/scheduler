package dto

import "time"

type Task struct {
	Id                uint
	Title             string
	URL               string
	CreationTimeStamp time.Time
}

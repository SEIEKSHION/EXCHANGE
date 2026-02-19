package models

import "time"

type Measurement struct {
	User    string    `json:"user"`
	Muscle  string    `json:"muscle"`
	Measure float32   `json:"measure"`
	Date    time.Time `json:"date"`
}

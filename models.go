package main

import "time"

type Client struct {
	name      string
	startTime time.Time
}
type Res struct {
	spentTime float64
	profit    float64
}

type InputData struct {
	clientTime time.Time
	eventID    string
	clientName string
	tableID    int
}

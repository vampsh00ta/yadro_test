package main

import (
	"bufio"
	"time"
)

func strToTime(str string) time.Time {
	layout := "15:04"
	t, _ := time.Parse(layout, str)
	return t
}

func readLine(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

package main

import (
	"fmt"
	"time"
)

func textOk(t time.Time, ID, clientName string) string {
	return fmt.Sprintf(timeIDName, t.Hour(), t.Minute(), ID, clientName)
}

func textOKTableID(t time.Time, ID, clientName string, tableID int) string {
	return fmt.Sprintf(timeIDName+" %d ", t.Hour(), t.Minute(), ID, clientName, tableID)
}

func textError(t time.Time, ID, err string) string {
	return fmt.Sprintf(timeIDName, t.Hour(), t.Minute(), "13", err)
}

package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

const formatRegex = `^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9] (1|2|3|4) [a-zA-Z0-9]+\s?[1-9]*$`

func getInputData(str string, lastTime *string) (InputData, bool) {
	re, _ := regexp.Compile(formatRegex)

	if !re.Match([]byte(str)) {
		fmt.Println(str)
		return InputData{}, false
	}

	var inputData InputData
	splitStr := strings.Split(str, " ")
	clientTimeStr := splitStr[0]
	if strToTime(clientTimeStr).Before(strToTime(*lastTime)) {
		fmt.Println(str)
		return InputData{}, false
	}
	*lastTime = clientTimeStr
	inputData.clientTime = strToTime(clientTimeStr)

	inputData.eventID = splitStr[1]
	inputData.clientName = splitStr[2]
	if inputData.eventID == "2" {
		tableID, _ := strconv.Atoi(splitStr[3])
		inputData.tableID = tableID
	}
	return inputData, true
}

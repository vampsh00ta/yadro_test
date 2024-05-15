package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	name      string
	startTime time.Time
}

const (
	NotOpenYet       = "NotOpenYet"
	ICanWaitNoLonger = "ICanWaitNoLonger!"
	YouShallNotPass  = "YouShallNotPass"
	PlaceIsBusy      = "PlaceIsBusy"
	ClientUnknown    = "ClientUnknown"
)

const (
	timeIDName = "%02d:%02d %s %s"
)

func main() {
	inFile, err := os.Open("file.txt")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	computerCount, _ := strconv.Atoi(readLine(scanner))

	workTimeStr := readLine(scanner)
	workTimeStrSplited := strings.Split(workTimeStr, " ")
	workStartTime := strToTime(workTimeStrSplited[0])
	workEndTime := strToTime(workTimeStrSplited[1])

	cost, _ := strconv.Atoi(readLine(scanner))
	// учет столов tableId : client {name, startTime}
	tables := make(map[int]Client)
	// учет клинетов name : tableID
	clients := make(map[string]int)
	// очередь ожидания
	q := NewQueue()
	// итоговые данные по столам

	res := make([]float64, computerCount)

	// итоговые ушедших
	goneClients := make([]string, 0)
	fmt.Printf("%02d:%02d\n", workStartTime.Hour(), workStartTime.Minute())

	for scanner.Scan() {
		str := scanner.Text()
		splitedStr := strings.Split(str, " ")
		clientTime := strToTime(splitedStr[0])
		ID := splitedStr[1]
		clientName := splitedStr[2]
		switch ID {
		case "1":

			fmt.Println(textOk(clientTime, ID, clientName))
			if clientTime.Before(workStartTime) || clientTime.After(workEndTime) {
				fmt.Println(textError(clientTime, "13", NotOpenYet))
				continue
			}
			if _, ok := clients[clientName]; ok {
				fmt.Println(textError(clientTime, "13", YouShallNotPass))

				continue
			}
			clients[clientName] = -1
		case "2":
			tableID, _ := strconv.Atoi(splitedStr[3])
			fmt.Println(textOKTableID(clientTime, ID, clientName, tableID))

			if _, ok := tables[tableID]; ok {
				fmt.Println(textError(clientTime, "13", PlaceIsBusy))

				continue
			}
			if _, ok := clients[clientName]; !ok {

				fmt.Println(textError(clientTime, "13", ClientUnknown))

				continue
			}
			delete(clients, clientName)
			clients[clientName] = tableID

			client := Client{name: clientName, startTime: clientTime}
			tables[tableID] = client

		case "3":
			fmt.Println(textOk(clientTime, ID, clientName))

			if len(tables) < computerCount {
				fmt.Println(textError(clientTime, "13", ICanWaitNoLonger))
				continue
			}
			if q.Size() >= computerCount {
				goneClients = append(goneClients, clientName)

				continue
			}
			q.Push(clientName)

		case "4":
			fmt.Println(textOk(clientTime, ID, clientName))

			if _, ok := clients[clientName]; !ok {
				fmt.Println(textError(clientTime, "13", ClientUnknown))

				continue
			}
			// считаем время за компом
			tableID := clients[clientName]
			leftClient := tables[tableID]

			spentTime := clientTime.Sub(leftClient.startTime)

			res[tableID-1] += spentTime.Minutes()
			delete(clients, clientName)

			// берем из очереди клиента
			qClient := q.Pop()
			if qClient != "" {
				clients[qClient] = tableID
				tables[tableID] = Client{name: qClient, startTime: clientTime}
				fmt.Println(textOKTableID(clientTime, "12", qClient, tableID))

			}

		}

	}
	for _, tableID := range clients {
		spentTime := workEndTime.Sub(tables[tableID].startTime)
		res[tableID-1] += spentTime.Minutes()

		goneClients = append(goneClients, tables[tableID].name)
	}
	for _, clientName := range goneClients {
		fmt.Println(textOk(workEndTime, "11", clientName))
	}
	fmt.Printf("%02d:%02d \n", workEndTime.Hour(), workEndTime.Minute())

	for i, t := range res {
		cashRes := math.Ceil(t/60) * float64(cost)
		fmt.Printf("%d %d %02d:%02d\n", i+1, int(cashRes), int(t/60), int(t)%60)

	}
}

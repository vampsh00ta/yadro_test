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
type Res struct {
	spentTime float64
	profit    float64
}

const (
	NotOpenYet       = "NotOpenYet"
	ICanWaitNoLonger = "ICanWaitNoLonger!"
	YouShallNotPass  = "YouShallNotPass"
	PlaceIsBusy      = "PlaceIsBusy"
	ClientUnknown    = "ClientUnknown"
)

func countProfit(t float64, cost int) float64 {
	return math.Ceil(t/60) * float64(cost)
}
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
	res := make([]Res, computerCount)

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
			tableID := clients[clientName]
			leftClient := tables[tableID]

			// считаем время за компом в минутах
			spentTime := clientTime.Sub(leftClient.startTime)
			//добавляет в итоговые значения
			res[tableID-1].spentTime += spentTime.Minutes()
			res[tableID-1].profit += countProfit(spentTime.Minutes(), cost)

			delete(clients, clientName)

			// берем из очереди клиента
			qClientName := q.Pop()
			if qClientName != "" {
				clients[qClientName] = tableID
				tables[tableID] = Client{name: qClientName, startTime: clientTime}
				fmt.Println(textOKTableID(clientTime, "12", qClientName, tableID))

			}

		}

	}
	for _, tableID := range clients {
		//добавляем в результат оставшихся в компьютерном клубе
		spentTime := workEndTime.Sub(tables[tableID].startTime)
		res[tableID-1].spentTime += spentTime.Minutes()
		res[tableID-1].profit += countProfit(spentTime.Minutes(), cost)

		goneClients = append(goneClients, tables[tableID].name)
	}
	for _, clientName := range goneClients {
		fmt.Println(textOk(workEndTime, "11", clientName))
	}
	fmt.Printf("%02d:%02d \n", workEndTime.Hour(), workEndTime.Minute())

	for i, value := range res {

		fmt.Printf("%d %d %02d:%02d\n", i+1, int(value.profit), int(value.spentTime/60), int(value.spentTime)%60)

	}
}

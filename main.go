package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"

	"strconv"
	"strings"
)

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

var lastTime = "00:00"

func main() {
	//принимает имя файла
	args := os.Args[1:]
	if len(args) == 0 {
		panic("no input file ")

	}
	fileName := args[0]
	inFile, err := os.Open("./" + fileName)
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

	// очередь ожидания клиентов
	q := NewQueue()

	// итоговые данные по столам
	res := make([]Res, computerCount)

	// итоговые ушедшие/оставшийся до закрытия
	goneClients := make([]string, 0)

	fmt.Printf("%02d:%02d\n", workStartTime.Hour(), workStartTime.Minute())

	for scanner.Scan() {
		str := scanner.Text()
		inputData, ok := getInputData(str, &lastTime)
		if !ok {
			return
		}
		clientTime := inputData.clientTime
		eventID := inputData.eventID
		clientName := inputData.clientName

		switch eventID {
		case "1":

			fmt.Println(textOk(clientTime, eventID, clientName))
			if clientTime.Before(workStartTime) || clientTime.After(workEndTime) {
				fmt.Println(textError(clientTime, "13", NotOpenYet))
				continue
			}
			if _, ok := clients[clientName]; ok {
				fmt.Println(textError(clientTime, "13", YouShallNotPass))

				continue
			}
			//добавляет клиента в мапу со  значеним -1, чтобы  можно было проверить его присутствие
			clients[clientName] = -1
		case "2":
			tableID := inputData.tableID
			fmt.Println(textOKTableID(clientTime, eventID, clientName, tableID))

			if _, ok := tables[tableID]; ok {
				fmt.Println(textError(clientTime, "13", PlaceIsBusy))

				continue
			}
			if _, ok := clients[clientName]; !ok {

				fmt.Println(textError(clientTime, "13", ClientUnknown))

				continue
			}
			//удаляет клиента на случай, если он решил поменять стол
			delete(clients, clientName)
			clients[clientName] = tableID

			client := Client{name: clientName, startTime: clientTime}
			tables[tableID] = client

		case "3":
			fmt.Println(textOk(clientTime, eventID, clientName))

			if len(tables) < computerCount {
				fmt.Println(textError(clientTime, "13", ICanWaitNoLonger))
				continue
			}
			//добавляет клиента в слайс недождавшихся/ушедших, если длина очереди больше кол-ва компов

			if q.Size() >= computerCount {
				goneClients = append(goneClients, clientName)

				continue
			}
			//добавляет клиента в очередь ожидания
			q.Push(clientName)

		case "4":
			fmt.Println(textOk(clientTime, eventID, clientName))

			if _, ok := clients[clientName]; !ok {
				fmt.Println(textError(clientTime, "13", ClientUnknown))

				continue
			}
			//берем ушедшего клиента из мапы
			tableID := clients[clientName]
			leftClient := tables[tableID]

			// считаем время за компом в минутах
			spentTime := clientTime.Sub(leftClient.startTime)

			//добавляет в итоговые значения
			res[tableID-1].spentTime += spentTime.Minutes()
			res[tableID-1].profit += countProfit(spentTime.Minutes(), cost)

			delete(clients, clientName)

			// берем из очереди клиента(если он есть)
			qClientName := q.Pop()
			if qClientName != "" {
				clients[qClientName] = tableID
				tables[tableID] = Client{name: qClientName, startTime: clientTime}
				fmt.Println(textOKTableID(clientTime, "12", qClientName, tableID))

			}

		}

	}

	for _, tableID := range clients {
		//добавляем в результат оставшихся клиентов в компьютерном клубе
		if tableID == -1 {
			continue
		}
		spentTime := workEndTime.Sub(tables[tableID].startTime)
		res[tableID-1].spentTime += spentTime.Minutes()
		res[tableID-1].profit += countProfit(spentTime.Minutes(), cost)

		goneClients = append(goneClients, tables[tableID].name)
	}
	//сортировка недождавшися/оставшихся клиентов
	slices.Sort(goneClients)
	for _, clientName := range goneClients {
		fmt.Println(textOk(workEndTime, "11", clientName))
	}
	fmt.Printf("%02d:%02d \n", workEndTime.Hour(), workEndTime.Minute())

	for i, value := range res {

		fmt.Printf("%d %d %02d:%02d\n", i+1, int(value.profit), int(value.spentTime/60), int(value.spentTime)%60)

	}
}

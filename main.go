package main

import (
	"bufio"
	"fmt"
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

func strToTime(str string) time.Time {
	layout := "15:04"
	t, _ := time.Parse(layout, str)
	return t
}
func readLine(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

func main() {
	inFile, err := os.Open("file.txt")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	computerCount, _ := strconv.Atoi(readLine(scanner))
	fmt.Println(computerCount)

	workTimeStr := readLine(scanner)
	workTimeStrSplited := strings.Split(workTimeStr, " ")
	fmt.Println(workTimeStrSplited)
	workStartTime := strToTime(workTimeStrSplited[0])
	workEndTime := strToTime(workTimeStrSplited[1])
	fmt.Println(workStartTime, workEndTime)

	cost := readLine(scanner)
	fmt.Println(cost)
	//учет столов tableId : client {name, startTime}
	tables := make(map[int]Client)
	// учет клинетов name : tableID
	clients := make(map[string]int)
	//очередь ожидания
	q := NewQueue()
	//итоговые данные по столам

	res := make([]time.Duration, computerCount)

	//итоговые ушедших
	goneClients := make([]string, 0)

	for scanner.Scan() {
		str := scanner.Text()
		splitedStr := strings.Split(str, " ")
		clientTime := strToTime(splitedStr[0])
		ID := splitedStr[1]
		clientName := splitedStr[2]
		switch ID {
		case "1":
			//fmt.Printf("%s %s", ID, clientName)
			fmt.Println(clientTime.Hour(), ":", clientTime.Minute(), ID, " ", clientName)

			if clientTime.Before(workStartTime) || clientTime.After(workEndTime) {
				fmt.Println("13" + NotOpenYet)
				continue
			}
			if _, ok := clients[clientName]; ok {
				fmt.Println("13" + YouShallNotPass)
				continue
			}
			clients[clientName] = -1
		case "2":
			tableID, _ := strconv.Atoi(splitedStr[3])
			//fmt.Printf("%s %s %s", ID, clientName, tableID)
			fmt.Println(clientTime.Hour(), ":", clientTime.Minute(), ID, " ", clientName, tableID)

			if _, ok := tables[tableID]; ok {
				fmt.Println("13" + PlaceIsBusy)
				continue
			}
			if _, ok := clients[clientName]; !ok {
				fmt.Println("13" + ClientUnknown)
				continue
			}
			delete(clients, clientName)
			clients[clientName] = tableID

			client := Client{name: clientName, startTime: clientTime}
			tables[tableID] = client

		case "3":
			fmt.Println(clientTime.Hour(), ":", clientTime.Minute(), ID, " ", clientName)

			if len(tables) < computerCount {
				fmt.Println("13" + ICanWaitNoLonger)
				continue
			}
			if q.Size() >= computerCount {
				goneClients = append(goneClients, clientName)

				continue
			}
			q.Push(clientName)

		case "4":
			fmt.Println(clientTime.Hour(), ":", clientTime.Minute(), ID, " ", clientName)

			if _, ok := clients[clientName]; !ok {
				fmt.Println("13" + ClientUnknown)
				continue
			}
			//считаем время за компом
			tableID := clients[clientName]
			leftClient := tables[tableID]

			spentTime := clientTime.Sub(leftClient.startTime)
			res[tableID] = res[tableID].Truncate(spentTime)
			delete(clients, clientName)
			//берем из очереди клиента
			qClient := q.Pop()
			if qClient != "" {
				clients[qClient] = tableID
				tables[tableID] = Client{name: qClient, startTime: clientTime}
				fmt.Println(clientTime.Hour(), ":", clientTime.Minute(), 12, " ", clientName, tableID)
			}
			//берем из очереди клиента

		}

	}
	for _, id := range clients {
		goneClients = append(goneClients, tables[id].name)
	}
	for _, clientName := range goneClients {
		fmt.Println(11, clientName)
	}
	for i, t := range res {
		fmt.Println(i, t)
	}

}

//
//
//09:41 1 client1
//09:48 1 client2
//09:52 3 client1
//09:52 13 ICanWaitNoLonger!
//09:54 2 client1 1
//10:25 2 client2 2
//10:58 1 client3
//10:59 2 client3 3
//11:30 1 client4
//11:35 2 client4 2
//11:35 13 PlaceIsBusy
//11:45 3 client4
//12:33 4 client1
//12:33 12 client4 1
//12:43 4 client2
//15:52 4 client4

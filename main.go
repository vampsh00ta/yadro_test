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
	timeFormat     = "%02d:%02d "
	idCLientFormat = "\t%s %s "
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

	workTimeStr := readLine(scanner)
	workTimeStrSplited := strings.Split(workTimeStr, " ")
	workStartTime := strToTime(workTimeStrSplited[0])
	workEndTime := strToTime(workTimeStrSplited[1])

	cost, _ := strconv.Atoi(readLine(scanner))
	//учет столов tableId : client {name, startTime}
	tables := make(map[int]Client)
	// учет клинетов name : tableID
	clients := make(map[string]int)
	//очередь ожидания
	q := NewQueue()
	//итоговые данные по столам

	res := make([]float64, computerCount)

	//итоговые ушедших
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

			fmt.Printf("%02d:%02d %s %s\n", clientTime.Hour(), clientTime.Minute(), ID, clientName)

			if clientTime.Before(workStartTime) || clientTime.After(workEndTime) {
				fmt.Printf(timeFormat+idCLientFormat+"\n", clientTime.Hour(), clientTime.Minute(), "13", NotOpenYet)
				continue
			}
			if _, ok := clients[clientName]; ok {
				fmt.Printf(timeFormat+idCLientFormat+"\n", clientTime.Hour(), clientTime.Minute(), "13", YouShallNotPass)

				continue
			}
			clients[clientName] = -1
		case "2":
			tableID, _ := strconv.Atoi(splitedStr[3])
			fmt.Printf("%02d:%02d  %s %s %d\n", clientTime.Hour(), clientTime.Minute(), ID, clientName, tableID)

			if _, ok := tables[tableID]; ok {
				fmt.Printf(timeFormat+idCLientFormat+"\n", clientTime.Hour(), clientTime.Minute(), "13", PlaceIsBusy)

				continue
			}
			if _, ok := clients[clientName]; !ok {
				fmt.Printf(timeFormat+idCLientFormat+"\n", clientTime.Hour(), clientTime.Minute(), "13", ClientUnknown)

				continue
			}
			delete(clients, clientName)
			clients[clientName] = tableID

			client := Client{name: clientName, startTime: clientTime}
			tables[tableID] = client

		case "3":
			fmt.Printf("%02d:%02d %s %s\n", clientTime.Hour(), clientTime.Minute(), ID, clientName)

			if len(tables) < computerCount {
				fmt.Printf(timeFormat+idCLientFormat+"\n", clientTime.Hour(), clientTime.Minute(), "13", ICanWaitNoLonger)

				continue
			}
			if q.Size() >= computerCount {
				goneClients = append(goneClients, clientName)

				continue
			}
			q.Push(clientName)

		case "4":
			fmt.Printf("%02d:%02d %s %s\n", clientTime.Hour(), clientTime.Minute(), ID, clientName)

			if _, ok := clients[clientName]; !ok {
				fmt.Printf(timeFormat+idCLientFormat+"\n", clientTime.Hour(), clientTime.Minute(), "13", ClientUnknown)

				continue
			}
			//считаем время за компом
			tableID := clients[clientName]
			leftClient := tables[tableID]

			spentTime := clientTime.Sub(leftClient.startTime)

			res[tableID-1] += spentTime.Minutes()
			delete(clients, clientName)
			//берем из очереди клиента
			qClient := q.Pop()
			if qClient != "" {
				clients[qClient] = tableID
				tables[tableID] = Client{name: qClient, startTime: clientTime}
				fmt.Printf("%02d:%02d %d %s %d\n", clientTime.Hour(), clientTime.Minute(), 12, clientName, tableID)
			}
			//берем из очереди клиента

		}

	}
	for _, tableID := range clients {
		spentTime := workEndTime.Sub(tables[tableID].startTime)
		res[tableID-1] += spentTime.Minutes()

		goneClients = append(goneClients, tables[tableID].name)
	}
	for _, clientName := range goneClients {
		fmt.Printf("%02d:%02d %d %s \n", workEndTime.Hour(), workEndTime.Minute(), 11, clientName)
	}
	for i, t := range res {
		cashRes := math.Ceil(t/60) * float64(cost)
		fmt.Println(i+1, cashRes)
	}
	fmt.Printf("%02d:%02d\n", workEndTime.Hour(), workEndTime.Minute())

}

//09:00
//08:48 1 client1
//13 NotOpenYet
//09:41 1 client1
//09:48 1 client2
//09:52 3 client1
//13 ICanWaitNoLonger!
//09:54  2 client1 1
//10:25  2 client2 2
//10:58 1 client3
//10:59  2 client3 3
//11:30 1 client4
//11:35  2 client4 2
//13 PlaceIsBusy
//11:45 3 client4
//12:33 4 client1
//12:33 12 client1 1
//12:43 4 client2
//15:52 4 client4
//19:00 11 client3
//1 60
//2 30
//3 90
//19:00
//
//
//09:00
//08:48 1 client1
//08:48 13 NotOpenYet
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
//19:00 11 client3
//19:00
//1 70 05:58
//2 30 02:18
//3 90 08:01

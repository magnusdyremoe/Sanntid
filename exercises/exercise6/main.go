package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"./Network/network/bcast"
)

var message msg

type msg struct {
	id     int
	number int
}

var status alive

type alive struct {
	id int
}

var count int

func main() {
	localid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Print("panicyo")
	}
	message.id = localid
	status.id = localid

	tx := make(chan msg)
	rx := make(chan msg)
	acktx := make(chan alive)
	ackrx := make(chan alive)
	alivemsgtimer := time.NewTimer(0)
	fmt.Print("id: ", localid, "\n")

	timeOut := time.NewTimer(0)
	go bcast.Receiver(15648, rx)
	go bcast.Transmitter(15648, tx)
	go bcast.Receiver(15649, ackrx)
	go bcast.Transmitter(15649, acktx)
	alivemsgtimer.Reset(200 * time.Millisecond)
	for {
		select {
		case receivedmsg := <-rx:

			msgid := receivedmsg.id

			fmt.Println("idmsg: ", msgid)
			num := receivedmsg.number
			if msgid == 1 && localid == 2 {
					message.number = num
			}

		case alivemsg := <-ackrx:
			fmt.Print(alivemsg.id, "is id\n")
			if alivemsg.id != localid {
				timeOut.Reset(1 * time.Second)
			}

		case <-timeOut.C:
			localid = 1

		case <-alivemsgtimer.C:
			msg := alive{status.id}
			acktx <- msg

		}
		if localid == 1 {
			count = count + 1
			fmt.Print("mastercount", count, "\n")
			msg := msg{message.id, message.number}
			tx <- msg
		}
		if	localid == 2 {
			msg := msg{message.id, message.number}
			tx <- msg
		}

		time.Sleep(1 * time.Second)

	}
}

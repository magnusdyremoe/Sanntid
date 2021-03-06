package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"./elevator"
	"./elevio"
	"./network/bcast"
	"./variables"
)

//sudo iptables -A INPUT -p udp --dport 15648 -m statistic --mode random --probability 0.2 -j DROP

//go run main.go portnr id
func main() {

	cmd := os.Args[1]
	ElevatorID, err := strconv.Atoi(os.Args[2])
	fmt.Println(ElevatorID)
	time.Sleep(1 * time.Second)

	if err != nil {
		panic(err)
	}

	elevio.Init("localhost:"+cmd, variables.N_FLOORS)
	elevator.ElevatorInit(ElevatorID)
	elevator.LocalQueueInit()
	fmt.Println("Initialized")

	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvStop := make(chan bool)
	elevTx := make(chan variables.ElevatorMessage)
	elevRx := make(chan variables.ElevatorMessage)
	timeOut := time.NewTimer(0)
	DoorTimer := elevator.Elev.DoorTimer
	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollStopButton(drvStop)
	go bcast.Receiver(15648, elevRx)
	go bcast.Transmitter(15648, elevTx)

	for {
		select {
		case atFloor := <-drvFloors:
			elevator.ElevatorFloorUpdate(ElevatorID, atFloor)
			elevator.FsmFloor(atFloor, elevator.Elev.Dir)
			msg := variables.ElevatorMessage{ElevatorID, "FLOOR", -1, atFloor, int(elevator.Elev.Dir), elevator.Elev.ElevState}
			elevTx <- msg
			elevTx <- msg
			elevTx <- msg
			elevTx <- msg
		case stop := <-drvStop:
			elevator.FsmStop(stop)
		case elevatorMessageReceived := <-elevRx:
			elevator.FsmMessageReceivedHandler(elevatorMessageReceived, ElevatorID)
			if !elevator.QueueCheckEmpty(variables.LOCAL) || !elevator.QueueCheckEmpty(variables.REMOTE) {
				timeOut.Reset(variables.FAULT_TIME * time.Second)
			} else {
				timeOut.Stop()
			}
		case buttonCall := <-drvButtons:
			if buttonCall.Button == elevator.Cab {
				elevator.FsmOnButtonRequest(buttonCall, true)
			} else {
				msg := variables.ElevatorMessage{ElevatorID, "ORDER", int(buttonCall.Button), buttonCall.Floor, int(elevator.Elev.Dir), elevator.Elev.ElevState}
				elevTx <- msg
				elevTx <- msg
				elevTx <- msg
				elevTx <- msg
			}

		case <-timeOut.C:
			elevator.ElevatorSetConnectionStatus(ElevatorID, variables.ELEV_OFFLINE)
			msg := variables.ElevatorMessage{ElevatorID, "NOT_RESPONDING", -1, -1, int(elevator.Elev.Dir), elevator.Elev.ElevState}
			elevTx <- msg
			elevTx <- msg
			elevTx <- msg
			elevTx <- msg

		case <-DoorTimer.C:
			elevator.FsmExitDoorState(elevator.Elev.DoorTimer)
		}
	}

}

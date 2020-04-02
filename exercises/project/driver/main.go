package main

import (
	"flag"
	"fmt"
	"os"

	"./elevator"
	"./elevio"
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"./variables"
)

func main() {

	cmd := os.Args[1]
	fmt.Println(cmd)
	elevio.Init("localhost:"+cmd, variables.N_FLOORS)
	//go run main.go portnr

	elevator.ElevatorInit()
	elevator.QueueInit()
	//elevator.backupInit()
	fmt.Println("Initialized")

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// Channels
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	//drvObstr := make(chan bool)
	drvStop := make(chan bool)

	elevTx := make(chan elevator.ElevatorMessage)
	elevRx := make(chan elevator.ElevatorMessage)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	//go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)
	go elevator.FsmPollButtonRequest(drvButtons)
	go bcast.Receiver(15648, elevRx)
	go bcast.Transmitter(15648, elevTx)

	for {
		select {
		case a := <-drvFloors:
			elevator.FsmFloor(a)
			msg := elevator.ElevatorMessage{"FLOOR", a, a}
			elevTx <- msg
			elevTx <- msg
			elevTx <- msg
			fmt.Printf("New Floor Sent\n")
		case a := <-drvStop:
			elevator.FsmStop(a)
		case p := <-elevRx:
			elevator.FsmMessageReceived(p)
			fmt.Printf("New ButtonPress Sent\n")
		case s := <-drvButtons:
			msg := elevator.ElevatorMessage{"ORDER", int(s.Button), s.Floor}
			elevTx <- msg
			fmt.Printf("New message sent\n")

		//case s:= <- drvFloors:
		//msg := elevator.ElevatorMessage{"FINISHED", s,0}

		case f := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", f.Peers)
			fmt.Printf("  New:      %q\n", f.New)
			fmt.Printf("  Lost:     %q\n", f.Lost)
		}
	}

}

// chmod +x ElevatorServer

// cant just run main. correct command:
// go run *.go

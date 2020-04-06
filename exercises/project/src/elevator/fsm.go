package elevator

import (
	"fmt"
	"time"

	"../elevio"
	"../variables"
)

//seems like there is a bug related to cab calls. The elevator sometimes go
//out of bounds

func FsmFloor(newFloor int) {

	elevatorSetNewFloor(newFloor)
	if queueCheckCurrentFloorSameDir(newFloor, elev.dir) {
		elevatorSetMotorDir(Stop)
		fsmDoorState()
		//elevatorExitDoorState()
		queueRemoveOrder(newFloor, elev.dir)
		elevatorLightsMatchQueue()

	}
	elevatorSetDir(queueReturnElevDir(newFloor, elev.dir))

}

/*func FsmPollButtonRequest(drvButtons chan elevio.ButtonEvent) {
	for {
		fsmOnButtonRequest(<-drvButtons)
	}
}
*/

func fsmOnButtonRequest(a elevio.ButtonEvent) {
	fmt.Print("New order recieved")
	fmt.Printf("%+v\n", a)
	queueRecieveOrder(a)
	elevatorLightsMatchQueue()
	//backupSync()

	if elevatorGetDir() == Stop {
		if a.Floor == elevatorGetFloor() && elevatorGetDir() == Stop {
			fsmDoorState()
			FsmFloor(elevatorGetFloor())
		}
		if ElevatorGetDoorState() == false{ 
			elevatorSetDir(queueReturnElevDir(elevatorGetFloor(), elevatorGetDir()))
		}
	}
}

func FsmMessageReceivedHandler(msg ElevatorMessage) {
	//sync the new message with queue
	fmt.Println("received a message")
	msgType := msg.MessageType
	button := msg.Button
	floor := msg.Floor
	if msgType == "ORDER" {
		event := elevio.ButtonEvent{floor, elevio.ButtonType(button)}
		fsmOnButtonRequest(event)
	} else if msgType == "FLOOR" {
		FsmFloor(floor)
	} else {
		fmt.Print("invalid message")
	}
	elevatorLightsMatchQueue()

}


func fsmDoorState() {
	fmt.Print("Door state")
	ElevatorSetDoorState(true)
	elevio.SetDoorOpenLamp(true)
	elev.doorTimer.Reset(variables.DOOROPENTIME * time.Second)
	<-elev.doorTimer.C
	ElevatorSetDoorState(false)

	elevio.SetDoorOpenLamp(false)

}

//From project destription in the course embedded systems
func FsmStop(a bool) {
	fmt.Print("Stop state")
	fmt.Printf("%+v\n", a)
	ElevatorInit()
	QueueInit()
	elevatorLightsMatchQueue()
}

package elevator

import (
	"fmt"
	"time"

	"../elevio"
	"../variables"
)

func FsmFloor(newFloor int) {

	//decisionAlgorithm(newFloor, elev.dir)
	elevatorSetNewFloor(newFloor)
	if localQueueCheckCurrentFloorSameDir(newFloor, elev.dir) == true {
		fsmDoorState()
	}
	elevatorSetDir(localQueueReturnElevDir(newFloor, elev.dir))

}

func fsmOnButtonRequest(a elevio.ButtonEvent) {
	fmt.Print("New order recieved")
	fmt.Printf("%+v\n", a)
	//remoteQueueRecieveOrder(a)
	//decisionAlgorithm(elev.currentFloor, elev.dir)
	localQueueRecieveOrder(a)
	elevatorLightsMatchQueue()
	elev = ElevatorGetElev()
	if elev.dir == Stop && !ElevatorGetDoorOpenState() {
		if a.Floor == elev.currentFloor && elev.dir == Stop {
			FsmFloor(elev.currentFloor)
		}
		if ElevatorGetDoorOpenState() == false {
			elevatorSetDir(localQueueReturnElevDir(elev.currentFloor, elev.dir))
		}
	}
}

func FsmMessageReceivedHandler(msg variables.ElevatorMessage, ID string) {
	//sync the new message with queue
	fmt.Println("received a message")
	msgType := msg.MessageType
	msgID := msg.ElevID
	floor := msg.Floor
	button := msg.Button
	event := elevio.ButtonEvent{floor, elevio.ButtonType(button)}
	switch msgType {
	case "ORDER":
		fmt.Println(msgID+"+"+ID, "+", button)
		if button == 2 { //cabcall
			if msgID == ID {
				fsmOnButtonRequest(event)
			} else {
				fmt.Println("cabcall other elev")
			}
		} else {
			fsmOnButtonRequest(event)
		}
	case "FLOOR":
		if msgID == ID{
			fmt.Print("Floor")
			FsmFloor(floor)
		}
	case "ALIVE":
		fmt.Println("Alive from", msgID)
	default:
		fmt.Print("invalid message")
	}
	elevatorLightsMatchQueue()

}


func fsmDoorState(){
	elev = ElevatorGetElev()
	elevatorSetMotorDir(Stop)
	ElevatorSetDoorOpenState(true)
	elevio.SetDoorOpenLamp(true)
	
	elev.doorTimer.Reset(variables.DOOROPENTIME * time.Second)
	select{
		case <- elev.doorTimer.C:
		default:
	}
	fsmRemoveOrderHandler(elev)
	fsmExitDoorState()

	
}

func fsmExitDoorState(){
	<-elev.doorTimer.C
	elevio.SetDoorOpenLamp(false)
	ElevatorSetDoorOpenState(false)
	fmt.Print("yo")
}
//From project destription in the course embedded systems
func FsmStop(a bool) {
	elev = ElevatorGetElev()
	fmt.Print("Stop state")
	fmt.Printf("%+v\n", a)
	ElevatorInit(elev.ElevID)
	LocalQueueInit()
	elevatorLightsMatchQueue()
}

func fsmRemoveOrderHandler(elev Elevator){
	localQueueRemoveOrder(elev.currentFloor,elev.dir)
	elevatorLightsMatchQueue()
}
package elevator

import (
	"fmt"

	"../elevio"
	"../variables"
)

//decisionAlgorithm making a choice every time a new order is stashed into queue remote (on button push).

//Calculates cost of new order for N_ELEVATORS, finds best elevator (lowest cost). Returns best elevator
func DecisionAlgorithm(buttonPush elevio.ButtonEvent) int {
	var CostArray [variables.N_ELEVATORS + 1]int

	//Init cost array
	for elev := 1; elev < variables.N_ELEVATORS+1; elev++ {
		CostArray[elev] = 0
	}

	//Calculate cost
	for elevator := 1; elevator < variables.N_ELEVATORS+1; elevator++ {

		cost := buttonPush.Floor - Elev.ElevState[elevator][0]
		if cost < 0 {
			cost = -cost
		}

		if Elev.ElevState[elevator][2] > 0 {
			cost = cost + variables.ELEV_OFFLINE
		}

		CostArray[elevator] = cost
		fmt.Println("Elevator #: ", elevator, "%n Cost: ", cost)
	}

	//Find best elevator
	var bestElev int
	bestElev = 1
	for elevator := 2; elevator < variables.N_ELEVATORS+1; elevator++ {
		if CostArray[elevator] < CostArray[bestElev] {
			bestElev = elevator
		}

	}
	/*
		fmt.Println("Elevator 1 cost :", CostArray[1])
		fmt.Println("Elevator 2 cost :", CostArray[2])
		fmt.Println("Best elevator : ", bestElev)
		fmt.Println("Elevator ID : ", Elev.ElevID)
		fmt.Println("*------------------------_*")
		fmt.Println("Elev 1 at floor: ", Elev.ElevState[1][0])
		fmt.Println("Elev 2 at floor: ", Elev.ElevState[2][0])
	*/
	return bestElev

}

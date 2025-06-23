package main

func main() {
	floors := 10
	elevatorCount := 2
	totalPeople := 40 // 題目要求模擬 40 人

	building := NewBuilding(floors, elevatorCount, totalPeople)
	dispatcher := NewDispatcher(building)
	dispatcher.Run()
}

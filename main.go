package main

func main() {
	// 題目要求模擬 10 層樓、2 部電梯、40 人次
	floors := 10
	elevatorCount := 2
	totalPeople := 40

	building := NewBuilding(floors, elevatorCount, totalPeople)
	dispatcher := NewDispatcher(building)
	dispatcher.Run()
}

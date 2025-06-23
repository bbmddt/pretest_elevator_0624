package main

import (
	"fmt"
	"sync"
)

type Elevator struct {
	ID           int
	CurrentFloor int
	Direction    string    // 移動方向："up", "down", "idle"
	Passengers   []*Person // 電梯內的乘客
	Mutex        sync.Mutex
}

func NewElevator(id int) *Elevator {
	e := &Elevator{
		ID:           id,
		CurrentFloor: 1, // 電梯初始在 1 樓
		Direction:    "idle",
		Passengers:   []*Person{},
	}
	return e
}

func (e *Elevator) Step(building *Building) {
	e.Mutex.Lock()
	defer e.Mutex.Unlock()

	// 檢查是否無任務且應停留
	if e.Direction == "idle" && len(e.Passengers) == 0 {
		building.Mutex.Lock()
		hasWaiting := false
		for _, q := range building.WaitQueue {
			if len(q.UpQueue) > 0 || len(q.DownQueue) > 0 {
				hasWaiting = true
				break
			}
		}
		building.Mutex.Unlock()
		if !hasWaiting {
			return
		}
	}

	// 處理當前樓層的乘客上下車
	building.Mutex.Lock()
	hasOperation := false

	// 1. 先處理下車的乘客
	remaining := []*Person{}
	for _, p := range e.Passengers {
		if p.Destination == e.CurrentFloor {
			hasOperation = true
			fmt.Printf("[Log]🛗 %d 秒 電梯 %d 乘客 %d 在 %d 樓 出電梯\n",
				building.TimeElapsed, e.ID, p.ID, e.CurrentFloor)
			building.PeopleDone++
		} else {
			remaining = append(remaining, p)
		}
	}
	e.Passengers = remaining

	// 2. 再處理上車的乘客
	queueStruct := building.WaitQueue[e.CurrentFloor]
	var queue []*Person
	switch e.Direction {
	case "up":
		queue = queueStruct.UpQueue
		queueStruct.UpQueue = []*Person{}
	case "down":
		queue = queueStruct.DownQueue
		queueStruct.DownQueue = []*Person{}
	default:
		queue = append(queueStruct.UpQueue, queueStruct.DownQueue...)
		queueStruct.UpQueue = []*Person{}
		queueStruct.DownQueue = []*Person{}
	}

	// 處理上車
	remain := []*Person{}
	for _, p := range queue {
		if len(e.Passengers) < 5 {
			e.Passengers = append(e.Passengers, p)
			hasOperation = true
			fmt.Printf("[Log]🛗 %d 秒 電梯 %d 接乘客 %d 在 %d 樓 (目標: %d 樓)\n",
				building.TimeElapsed, e.ID, p.ID, e.CurrentFloor, p.Destination)
		} else {
			remain = append(remain, p)
		}
	}

	// 將未能上車的乘客放回等待隊列
	switch e.Direction {
	case "up":
		queueStruct.UpQueue = remain
	case "down":
		queueStruct.DownQueue = remain
	default:
		upLen := len(queueStruct.UpQueue)
		if upLen > len(remain) {
			upLen = len(remain)
		}
		queueStruct.UpQueue = remain[:upLen]
		queueStruct.DownQueue = remain[upLen:]
	}
	building.WaitQueue[e.CurrentFloor] = queueStruct
	building.Mutex.Unlock()

	// 判斷乘客進出停留
	if hasOperation {
		fmt.Printf("[Log]🛗 %d 秒 電梯 %d 停留在 %d 樓 有乘客進出中...\n",
			building.TimeElapsed, e.ID, e.CurrentFloor)
		return
	}

	// 決定移動方向
	e.updateDirection(building)

	// 移動電梯
	if e.Direction == "up" && e.CurrentFloor < building.Floors {
		e.CurrentFloor++
	} else if e.Direction == "down" && e.CurrentFloor > 1 {
		e.CurrentFloor--
	}

	fmt.Printf("[Log]🛗 %d 秒 電梯 %d 移動到 %d 樓，方向: %s，乘客數: %d\n",
		building.TimeElapsed, e.ID, e.CurrentFloor, e.Direction, len(e.Passengers))
}

func (e *Elevator) updateDirection(building *Building) {
	if len(e.Passengers) > 0 {
		// 優先服務車內乘客
		upNeeded := false
		downNeeded := false
		for _, p := range e.Passengers {
			if p.Destination > e.CurrentFloor {
				upNeeded = true
			} else if p.Destination < e.CurrentFloor {
				downNeeded = true
			}
		}
		if upNeeded && !downNeeded {
			e.Direction = "up"
		} else if downNeeded && !upNeeded {
			e.Direction = "down"
		} else if upNeeded && downNeeded {
			if e.Direction == "idle" {
				e.Direction = "up"
			}
		} else {
			e.Direction = "idle"
		}
	} else {
		// 無車內乘客，尋找最近的等待樓層
		minDist := building.Floors + 1
		targetFloor := -1
		building.Mutex.Lock()
		for floor, q := range building.WaitQueue {
			if len(q.UpQueue) > 0 || len(q.DownQueue) > 0 {
				dist := abs(floor - e.CurrentFloor)
				if dist < minDist {
					minDist = dist
					targetFloor = floor
				}
			}
		}
		building.Mutex.Unlock()

		if targetFloor == -1 {
			e.Direction = "idle"
		} else if targetFloor > e.CurrentFloor {
			e.Direction = "up"
		} else if targetFloor < e.CurrentFloor {
			e.Direction = "down"
		} else {
			e.Direction = "idle"
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getPersonIDs(people []*Person) []int {
	ids := make([]int, len(people))
	for i, p := range people {
		ids[i] = p.ID
	}
	return ids
}

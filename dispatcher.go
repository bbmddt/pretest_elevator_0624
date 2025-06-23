package main

import (
	"fmt"
	"sync"
)

type Dispatcher struct {
	Building *Building
}

func NewDispatcher(b *Building) *Dispatcher {
	return &Dispatcher{
		Building: b,
	}
}

func (d *Dispatcher) Run() {
	personID := 1
	for {
		if personID <= d.Building.TotalPeople {
			d.Building.GeneratePerson(personID)
			personID++
		}

		// 獲取等待樓層
		waitingFloors := d.getWaitingFloors()

		// 分配電梯
		assignedFloors := make(map[int]bool)
		for _, e := range d.Building.Elevators {
			e.Mutex.Lock()
			if len(waitingFloors) == 0 {
				e.Direction = "idle"
			} else {
				// 尋找最適合的等待樓層
				minDist := d.Building.Floors + 1
				targetFloor := -1
				for _, wf := range waitingFloors {
					if assignedFloors[wf] {
						continue
					}
					dist := abs(wf - e.CurrentFloor)
					canAssign := (e.Direction == "up" && e.CurrentFloor <= wf) ||
						(e.Direction == "down" && e.CurrentFloor >= wf) ||
						e.Direction == "idle"
					if canAssign && dist < minDist {
						minDist = dist
						targetFloor = wf
					}
				}
				if targetFloor != -1 {
					assignedFloors[targetFloor] = true
					if targetFloor > e.CurrentFloor {
						e.Direction = "up"
					} else if targetFloor < e.CurrentFloor {
						e.Direction = "down"
					} else {
						e.Direction = "idle"
					}
				} else {
					e.Direction = "idle"
				}
			}
			e.Mutex.Unlock()
		}

		var wg sync.WaitGroup
		for _, e := range d.Building.Elevators {
			wg.Add(1)
			go func(e *Elevator) {
				defer wg.Done()
				e.Step(d.Building)
			}(e)
		}
		wg.Wait()

		d.Building.Mutex.Lock()
		d.Building.TimeElapsed++
		if d.Building.PeopleDone > d.Building.TotalPeople {
			fmt.Printf("[Error] PeopleDone %d 超過 TotalPeople %d，強制終止\n", d.Building.PeopleDone, d.Building.TotalPeople)
			d.Building.Mutex.Unlock()
			break
		}
		d.Building.Mutex.Unlock()

		if d.isSimulationDone(personID) {
			break
		}
	}

	fmt.Printf("\n✅ 所有乘客運送完成，總耗時: %d 秒\n", d.Building.TimeElapsed)
}

func (d *Dispatcher) getWaitingFloors() []int {
	d.Building.Mutex.Lock()
	defer d.Building.Mutex.Unlock()

	floors := []int{}
	for f, q := range d.Building.WaitQueue {
		if len(q.UpQueue) > 0 || len(q.DownQueue) > 0 {
			floors = append(floors, f)
		}
	}
	return floors
}

func (d *Dispatcher) isSimulationDone(personID int) bool {
	d.Building.Mutex.Lock()
	defer d.Building.Mutex.Unlock()

	// 檢查是否所有乘客都已生成
	if personID <= d.Building.TotalPeople {
		// fmt.Printf("[Debug] 未結束: 仍有乘客未生成 (personID: %d, TotalPeople: %d)\n", personID, d.Building.TotalPeople)
		return false
	}

	// // 檢查等待隊列
	// for floor, queue := range d.Building.WaitQueue {
	// 	if len(queue.UpQueue) > 0 || len(queue.DownQueue) > 0 {
	// 		fmt.Printf("[Debug] 未結束: 樓層 %d 有等待乘客 (UpQueue: %v, DownQueue: %v)\n", floor, getPersonIDs(queue.UpQueue), getPersonIDs(queue.DownQueue))
	// 		return false
	// 	}
	// }

	// 檢查電梯狀態
	for _, e := range d.Building.Elevators {
		e.Mutex.Lock()
		hasPassenger := len(e.Passengers) > 0
		isIdle := e.Direction == "idle"
		e.Mutex.Unlock()
		if hasPassenger || !isIdle {
			// fmt.Printf("[Debug] 未結束: 電梯 %d 有乘客 (%d) 或非閒置 (方向: %s)\n", e.ID, len(e.Passengers), e.Direction)
			return false
		}
	}

	// 檢查是否所有乘客都已送達
	if d.Building.PeopleDone < d.Building.TotalPeople {
		// fmt.Printf("[Debug] 未結束: 已送達乘客 %d 小於總人數 %d\n", d.Building.PeopleDone, d.Building.TotalPeople)
		return false
	}

	// fmt.Printf("[Debug] 模擬結束: 所有乘客已送達，耗時 %d 秒\n", d.Building.TimeElapsed)
	return true
}

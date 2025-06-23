package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Person struct {
	ID          int
	Origin      int
	Destination int
}

type Building struct {
	Floors      int
	TotalPeople int
	WaitQueue   map[int]struct {
		UpQueue   []*Person
		DownQueue []*Person
	}
	Elevators   []*Elevator
	PeopleDone  int
	TimeElapsed int
	Mutex       sync.Mutex
}

func NewBuilding(floors, elevatorCount, totalPeople int) *Building {
	b := &Building{
		Floors:      floors,
		TotalPeople: totalPeople,
		WaitQueue:   make(map[int]struct{ UpQueue, DownQueue []*Person }),
		Elevators:   make([]*Elevator, elevatorCount),
		PeopleDone:  0,
		TimeElapsed: 0,
	}
	for i := 0; i < elevatorCount; i++ {
		b.Elevators[i] = NewElevator(i + 1)
	}
	for i := 1; i <= floors; i++ {
		b.WaitQueue[i] = struct{ UpQueue, DownQueue []*Person }{UpQueue: []*Person{}, DownQueue: []*Person{}}
	}
	// rand.Seed(time.Now().UnixNano())
	fmt.Printf("[Init] 大樓初始化完成，樓層數: %d，電梯數: %d，總乘客數: %d\n", floors, elevatorCount, totalPeople)
	return b
}

func (b *Building) GeneratePerson(id int) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	origin := rand.Intn(b.Floors) + 1
	var dest int
	for {
		dest = rand.Intn(b.Floors) + 1
		if dest != origin {
			break
		}
	}

	p := &Person{
		ID:          id,
		Origin:      origin,
		Destination: dest,
	}
	// 取出 WaitQueue[origin] 的結構體
	queue := b.WaitQueue[origin]
	if dest > origin {
		queue.UpQueue = append(queue.UpQueue, p)
		fmt.Printf("[Log]👤 %d 秒 乘客 %d 於 %d 樓按 🔼 要去 %d 樓\n", b.TimeElapsed, p.ID, p.Origin, p.Destination)
	} else {
		queue.DownQueue = append(queue.DownQueue, p)
		fmt.Printf("[Log]👤 %d 秒 乘客 %d 於 %d 樓按 🔽 要去 %d 樓\n", b.TimeElapsed, p.ID, p.Origin, p.Destination)
	}
	// 將修改後的結構體寫回 map
	b.WaitQueue[origin] = queue
}

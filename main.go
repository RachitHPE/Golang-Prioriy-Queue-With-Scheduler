package main

import (
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"time"
)

type TaskScheduler struct {
	ItemName string
}

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	originTime time.Time
	currentTime time.Time
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func priorityQueueTickChecker(pq *PriorityQueue) {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				countOfItems := len(*pq)
				if countOfItems == 0 {
					fmt.Println("No element present in the priority queue")

					os.Exit(0)
				}

				item := heap.Pop(pq).(*Item)
				fmt.Println("Element that gets expired first in given set of data is ", item.value)
				fmt.Println("CurrentTime is ", item.currentTime)
				fmt.Println("OriginTime is ", item.originTime)
				fmt.Println("Item priority in seconds is ", item.priority)
				fmt.Println()
			}
		}
	}()
	time.Sleep(50 * time.Second)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

func createTaskInLoop(pq *PriorityQueue) {
	originTime := time.Now()

	for i := 0; i < 7; i++ {
		time.Sleep(2 * time.Second)
		currentTime := time.Now()
		taskScheduler := TaskScheduler{
			ItemName: "taskItem_" + strconv.Itoa(i),
		}

		item := &Item{
			value:    taskScheduler.ItemName,
			originTime: originTime,
			currentTime: currentTime,
			priority: int(currentTime.Sub(originTime).Seconds()),
		}
		heap.Push(pq, item)

		fmt.Println("created task scheduler for ", taskScheduler.ItemName)
	}
}

func main() {
	fmt.Println("Entering the main functionnn!!!")
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	createTaskInLoop(&pq)
	priorityQueueTickChecker(&pq)
}

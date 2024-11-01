package main

import (
	"fmt"
	"sync"
)

func setLeader(leader int) {
	//log.Printf("The leader is now: %d at time: %d	\n", leader, time.Now().Nanosecond())
	counter++
}

var lock sync.Mutex
var counter int

func thread(wg *sync.WaitGroup, id int) {
	lock.Lock()
	setLeader(id)
	lock.Unlock()
	wg.Done()
}

func main() {

	// Dette er en kommentar
	counter = 0
	var wg sync.WaitGroup
	wg.Add(100000)
	for i := 0; i < 100000; i++ {
		go thread(&wg, i)
	}
	wg.Wait()
	fmt.Println(counter)
}

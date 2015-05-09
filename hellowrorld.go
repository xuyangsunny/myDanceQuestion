package main

import (
	"fmt"
	"runtime"
	"sync"
)

var counter int = 0

// func Add(x, y int) {
// 	z := x + y
// 	fmt.Println(z)
// }
func Count(lock *sync.Mutex) {
	lock.Lock()
	counter++
	fmt.Println(counter)
	lock.Unlock()

}
func main() {
	//fmt.Println("helloworld")
	lock := &sync.Mutex{}
	for i := 0; i < 10; i++ {
		//go Add(1, 1)
		go Count(lock)
	}

	for {
		lock.Lock()
		c := counter
		lock.Unlock()
		runtime.Gosched()
		if c >= 10 {
			break
		}
	}

}

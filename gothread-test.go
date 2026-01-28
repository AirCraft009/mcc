package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	n := 10
	syncArr := make([]chan int, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		syncArr[i] = make(chan int)

		go func() {
			defer wg.Done()
			dosmthBlocking(syncArr[i])
		}()
	}

	for i := 0; i < n; i++ {
		syncArr[i] <- i
		<-syncArr[i]
	}
	wg.Wait()
}

func dosmthBlocking(channel chan int) {
	x := <-channel
	fmt.Println("sent unblocking req:", x)
	channel <- x
}

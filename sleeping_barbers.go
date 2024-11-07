// Sleeping barbers

package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numSeats                = 3                      // Number of seats in waiting room
	haircutDuration         = time.Second            // Duration for each haircut
	customerArrivalInterval = 500 * time.Millisecond // Time between customer arrivals
)

var (
	customerQueue  = make(chan int, numSeats) // Buffered channel as a waiting room queue
	seatsAvailable = numSeats
	seatMutex      sync.Mutex
	barberCond     = sync.NewCond(&sync.Mutex{}) // Condition for barber to wait on
)

func barber() {
	for {
		barberCond.L.Lock() // Lock to protect condition check
		for len(customerQueue) == 0 {
			fmt.Println("Barber is sleeping...")
			barberCond.Wait() // Wait until there is a customer in the queue
		}
		barberCond.L.Unlock() // Unlock after waking up to serve a customer

		// Serve the next customer in the queue
		customer := <-customerQueue
		fmt.Printf("Barber starts cutting hair of customer %d.\n", customer)
		time.Sleep(haircutDuration) // Simulate time taken to cut hair
		fmt.Printf("Barber finishes cutting hair of customer %d.\n", customer)
	}
}

func customer(id int) {
	seatMutex.Lock() // Protect seat count updates
	if seatsAvailable > 0 {
		seatsAvailable--
		fmt.Printf("Customer %d takes a seat. Seats available: %d\n", id, seatsAvailable)
		seatMutex.Unlock() // Unlock after modifying seats

		// Add the customer to the queue
		customerQueue <- id

		// Notify barber that a customer is waiting
		barberCond.L.Lock()
		barberCond.Signal()
		barberCond.L.Unlock()

		// Simulate customer leaving after haircut
		seatMutex.Lock()
		seatsAvailable++
		seatMutex.Unlock()
	} else {
		fmt.Printf("Customer %d leaves due to no available seats.\n", id)
		seatMutex.Unlock() // Unlock if no seats are available
	}
}

func main() {
	go barber() // Start the barber goroutine

	customerID := 1
	for {
		go customer(customerID) // Each customer runs as a Goroutine
		customerID++
		time.Sleep(customerArrivalInterval) // Delay before next customer arrives
	}
}

// Team: AJ Valentino and Lauren Ferlito
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numSeats        = 3               // Number of seats in the waiting room
	haircutDuration = time.Second     // Duration for each haircut
)

var (
	seatsAvailable  = numSeats
	seatMutex       sync.Mutex                  // Mutex to protect seat availability
	barberSleeping  = make(chan bool, 1)        // Channel to signal if the barber is asleep
	customerChan    = make(chan int, numSeats)  // Channel to represent the waiting room with limited seats
)

func barber() {
	for {
		// Barber waits to be woken up by a customer if there are no customers waiting
		select {
		case customer := <-customerChan:
			fmt.Printf("Barber starts cutting hair of customer %d.\n", customer)
			time.Sleep(haircutDuration)
			fmt.Printf("Barber finishes cutting hair of customer %d.\n", customer)
		default:
			// If no customers, barber goes to sleep
			fmt.Println("Barber is sleeping...")
			time.Sleep(500 * time.Millisecond) // Avoid busy waiting
		}
	}
}

func customer(id int) {
	select {
	case customerChan <- id:
		// Customer takes a seat in the channel if available
		fmt.Printf("Customer %d takes a seat. Seats available: %d\n", id, cap(customerChan)-len(customerChan))
	default:
		// No seats are available, so the customer leaves
		fmt.Printf("Customer %d leaves because there are no available seats.\n", id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Start the barber goroutine
	go barber()

	// Infinite loop to simulate continuous customer arrivals with adjusted rates
	customerID := 1
	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)+500)) // Random arrival time with minimum wait
		go customer(customerID) // Each customer arrives as a goroutine
		customerID++
	}
}

// Team: AJ Valentino and Lauren Ferlito
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numSeats                = 3             // Number of seats in waiting room
	haircutDuration         = time.Second   // Duration for each haircut
	customerArrivalInterval = 500 * time.Millisecond // Time between customer arrivals
)

var (
	customerQueue  = make([]int, 0, numSeats) // Queue of waiting customers
	seatsAvailable = numSeats
	seatMutex      sync.Mutex                    // Mutex to protect seat availability
	barberMutex    sync.Mutex                    // Mutex to protect the barber's work
	conditionMutex sync.Mutex                    // Mutex to protect condition checking
	wg              sync.WaitGroup                // WaitGroup to wait for customer goroutines to complete
)

// Function for the barber to serve customers
func barber() {
	for {
		// Lock the barber mutex to avoid race condition while serving customers
		barberMutex.Lock()
		conditionMutex.Lock()

		if len(customerQueue) == 0 {
			// No customers in the queue, barber sleeps
			conditionMutex.Unlock()
			barberMutex.Unlock()
			fmt.Println("Barber is sleeping...")
			time.Sleep(time.Second) // Simulate barber sleeping
			continue
		}

		// There is a customer to serve
		customer := customerQueue[0] // Get the first customer in the queue
		customerQueue = customerQueue[1:] // Remove the customer from the queue

		conditionMutex.Unlock() // Unlock the condition mutex since we're going to work on a customer
		barberMutex.Unlock() // Unlock the barber mutex to allow other customers to arrive

		// Serve the customer
		fmt.Printf("Barber starts cutting hair of customer %d.\n", customer)
		time.Sleep(haircutDuration) // Simulate haircut time
		fmt.Printf("Barber finishes cutting hair of customer %d.\n", customer)

		// Lock the seatMutex again to modify the available seats
		seatMutex.Lock()
		seatsAvailable++
		seatMutex.Unlock()

		wg.Done() // Indicate that this customer has finished
	}
}

// Function for the customers
func customer(id int) {
	// Lock the seat mutex to modify the seat availability
	seatMutex.Lock()
	if seatsAvailable > 0 {
		// If there is space, customer takes a seat
		seatsAvailable--
		fmt.Printf("Customer %d takes a seat. Seats available: %d\n", id, seatsAvailable)
		customerQueue = append(customerQueue, id) // Add the customer to the queue
		seatMutex.Unlock() // Release the seat mutex

		// Signal the barber (implicitly handled by the barber function running concurrently)
		wg.Add(1) // Add to the wait group for this customer

		// The barber will start cutting hair once a customer is added to the queue
	} else {
		// If no seats are available, customer leaves
		fmt.Printf("Customer %d leaves because there are no available seats.\n", id)
		seatMutex.Unlock()
	}
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Start the barber goroutine
	go barber()

	// Simulate customers arriving at random intervals
	customerID := 1
	for {
		// Simulate customer arrival at random intervals
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+200)) // Random customer arrival

		go customer(customerID) // Launch each customer as a goroutine
		customerID++
	}
}
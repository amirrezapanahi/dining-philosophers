/*
The classical Dining philosophers problem.

Implemented with forks (aka chopsticks) as atomics.
*/

/*
Implement the dining philosopher’s problem with the following constraints/modifications.

There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.
Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).
In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
The host allows no more than 2 philosophers to eat concurrently.
Each philosopher is numbered, 1 through 5.

When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>”
on a line by itself, where <number> is the number of the philosopher.

When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>”
on a line by itself, where <number> is the number of the philosopher.
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var eatWgroup sync.WaitGroup

type fork struct {
	id        int
	picked_up int32
}

type philosopher struct {
	id                  int
	leftFork, rightFork *fork
}

// Goes from thinking to hungry to eating and done eating then starts over.
// Adapt the pause values to increase or decrease contentions
// around the forks.
func (p philosopher) eat() {
	defer eatWgroup.Done()
	for {
		p.pickup_fork(p.leftFork)
		p.pickup_fork(p.rightFork)

		fmt.Printf("Philosopher #%d picked up left fork %d\n", p.id+1, p.leftFork.id)
		fmt.Printf("Philosopher #%d picked up right fork %d\n\n", p.id+1, p.rightFork.id)

		say("eating\n", p.id)

		rand.Seed(time.Now().UnixNano())
		randomSec := rand.Intn(5-1) + 1
		time.Sleep(time.Duration(randomSec) * time.Second)

		p.drop_fork(p.rightFork)
		p.drop_fork(p.leftFork)

		fmt.Printf("Philosopher #%d dropped left fork %d\n", p.id+1, p.leftFork.id)
		fmt.Printf("Philosopher #%d dropped right fork %d\n\n", p.id+1, p.rightFork.id)

		say("finished eating\n", p.id)
		time.Sleep(time.Second)
	}

}

func (p philosopher) pickup_fork(fork *fork) {
	if fork == p.leftFork {
		for atomic.LoadInt32(&fork.picked_up) == 1 {
		}
		atomic.StoreInt32(&p.leftFork.picked_up, 1)
	} else if fork == p.rightFork {
		for atomic.LoadInt32(&fork.picked_up) == 1 {
		}
		atomic.StoreInt32(&p.rightFork.picked_up, 1)
	}
}

func (p philosopher) drop_fork(fork *fork) {
	if fork == p.leftFork {
		atomic.StoreInt32(&p.leftFork.picked_up, 0)
	} else if fork == p.rightFork {
		atomic.StoreInt32(&p.rightFork.picked_up, 0)
	}
}

func say(action string, id int) {
	fmt.Printf("Philosopher #%d is %s\n", id+1, action)
}

func main() {
	// How many philosophers and forks

	count := 5

	// Create forks
	forks := make([]*fork, count)
	for i := 0; i < count; i++ {
		forks[i] = &fork{
			id: i + 1, picked_up: 0,
		}
	}

	// Create philospoher, assign them 2 forks and send them to the dining table
	philosophers := make([]*philosopher, count)
	for i := 0; i < count; i++ {
		philosophers[i] = &philosopher{
			id: i, leftFork: forks[i], rightFork: forks[(i+1)%count]}
		eatWgroup.Add(1)
		go philosophers[i].eat()

	}

	eatWgroup.Wait()

}
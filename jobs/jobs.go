package main

/*
# basic example of having
*/

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const exit = ":exit:"

func main() {
	// create threads to consume jobs, total cpus -1, so pc have resource to keep feeding jobs and doing other tasks
	cpus := runtime.NumCPU() - 1
	ch := make(chan string, cpus)
	var wg sync.WaitGroup
	wg.Add(cpus)
	for i := 0; i < cpus; i++ {
		go startThread(i, ch, &wg)
	}

	alpha := "abcdefghijklmn"
	// alpha := "ab"
	jobs := []string{}
	for _, a := range alpha {
		jobs = append(jobs, string(a))
	}
	fmt.Println("%+v", jobs)

	// fake loading jobs later, to simulate slow jobs creation
	time.Sleep(time.Second * 3)
	for _, job := range jobs {
		ch <- job
	}
	fmt.Println("jobs loaded")
	// fill the exit, using this method to help coping unknown job length at runtime
	for i := 0; i < cpus; i++ {
		ch <- exit
	}

	// wait until all threads are done
	wg.Wait()
	fmt.Println("Finished")
}

func startThread(cpu int, c <-chan string, wgg *sync.WaitGroup) {
Loop:
	for {
		select {
		case x := <-c:
			if x == exit {
				fmt.Println("exit proc", cpu)
				break Loop
			}

			time.Sleep(time.Second * time.Duration(rand.Intn(10)))
			fmt.Printf("proc %d, done %s\n", cpu, x)
		}
	}
	wgg.Done()

}

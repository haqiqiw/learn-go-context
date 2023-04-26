package learn_go_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	contextG := context.WithValue(contextF, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f")) // yes
	fmt.Println(contextF.Value("c")) // yes, get from parent
	fmt.Println(contextF.Value("b")) // no, different parent
	fmt.Println(contextA.Value("b")) // no, can't get child data
}

// goroutine leak
// func CreateCounter() chan int {
// 	destination := make(chan int)

// 	go func() {
// 		defer close(destination)

// 		counter := 1
// 		for {
// 			destination <- counter
// 			counter++
// 		}
// 	}()

// 	return destination
// }

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)

		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // simulate slow process
			}
		}
	}()

	return destination
}

func TestContextCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	cancel()
	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func addTen(ptr *int) {
	*ptr += 10
}

func mul2(ptr []int) {
	for i := 0; i < len(ptr); i++ {
		ptr[i] *= 2
	}
}

func twoRoutine() {
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 1 {
				fmt.Println("打印奇数的协程，结果是：", i)
			}
		}
	}()
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Println("打印偶数的协程，结果是：", i)
			}
		}
	}()
}

type task func()

func taskScheduler(tasks []task) []time.Duration {

	result := make([]time.Duration, len(tasks))

	for idx, tas := range tasks {

		go func(i int, t task) {
			start := time.Now()
			t()
			result[i] = time.Since(start)
		}(idx, tas)
	}

	return result
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() string {
	return fmt.Sprintf("Name: %s, Age: %d, EmployeeID: %s", e.Name, e.Age, e.EmployeeID)
}

func routineCommunication() {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for v := range ch { // 消费者
			fmt.Println("消费者打印：", v)
		}
	}()
}

func bufferedRoutineCommunication() {
	ch := make(chan int, 20)
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
		}
	}()
	go func() {
		for v := range ch { // 消费者
			fmt.Println("消费者打印：", v)
		}
	}()
}

var counter int = 0
var mutex sync.Mutex

func safeIncrement() {
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}

		}()
	}
	time.Sleep(time.Second)
	fmt.Println("Counter:", counter)
}

var counter64 int64 = 0

func atomicIncrement() {
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter64, 1)
			}

		}()
	}
	time.Sleep(time.Second)
	fmt.Println("Counter:", counter64)
}

func main() {
	//var x int = 10
	//addTen(&x)
	//fmt.Println(x)

	//arr := []int{1, 2, 3, 4, 5}
	//mul2(arr)
	//fmt.Println(arr)
	//twoRoutine()
	//time.Sleep(time.Second)
	//tasks := []task{
	//	func() {
	//		fmt.Println("Task 1")
	//		time.Sleep(200 * time.Millisecond)
	//	},
	//	func() {
	//		fmt.Println("Task 2")
	//		time.Sleep(210 * time.Millisecond)
	//	},
	//	func() {
	//		time.Sleep(220 * time.Millisecond)
	//		fmt.Println("Task 3")
	//	},
	//}
	//result := taskScheduler(tasks)
	//time.Sleep(time.Second)
	//fmt.Println(result)
	// rectangle := Rectangle{Width: 10, Height: 20}
	// circle := Circle{Radius: 10}
	// fmt.Println(rectangle.Area())
	// fmt.Println(rectangle.Perimeter())
	// fmt.Println(circle.Area())
	// fmt.Println(circle.Perimeter())
	// employee := Employee{Person: Person{Name: "John", Age: 30}, EmployeeID: "123456"}
	// fmt.Println(employee.PrintInfo())

	//routineCommunication()
	// bufferedRoutineCommunication()
	// time.Sleep(time.Second)
	//safeIncrement()
	atomicIncrement()

}

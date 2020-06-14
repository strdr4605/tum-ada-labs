package main

import "time"

func fibonacci(number int) int64 {
	if number <= 0 {
		println("Error! Passed number is negative %i. Expected only positive number as input.", number)
		return 0
	}

	if number == 1 { return 0 }
	if number == 2 { return 1 }

	var number1 int64 = 0
	var number2 int64 = 1
	var fib int64

	for i := 3; i <= number; i++ {
		fib = number1 + number2
		number1 = number2
		number2 = fib
	}

	return fib
}

func sleepyFibonacci(number int) int64 {
	time.Sleep(250 * time.Millisecond)

	return fibonacci(number)
}

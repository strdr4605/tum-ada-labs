package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readInput(fileName string) ([]int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var numbers []int

	for scanner.Scan() {
		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}

	return numbers, nil
}

func writeOutput(fileName string, numbers []int64) error {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, value := range numbers {
		_, err := fmt.Fprintln(writer, strconv.FormatInt(value, 10))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

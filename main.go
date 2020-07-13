package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1) // nonzero means error exit
}

// accepts a slice of slices that contain strings
func parseLines(lines [][]string) []problem {
	arr := make([]problem, len(lines)) // returns the value, not a pointer with new

	// populate the array with structs --> problems
	for index, line := range lines {
		arr[index] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return arr
}

func main() {
	// if you wanted your program to accept flags
	csvFileName := flag.String("csv", "problems.csv", "a csv file in format of question,answer") // problems is default on -csv flag
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open csv file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file) // creates a reader for the file
	lines, err := r.ReadAll()

	if err != nil {
		exit("failed to parse the csv file")
	}

	fmt.Println(lines) // slices are printed out without comma delimiter

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	numCorrect := 0

problemloop:
	for index, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", index+1, problem.question)
		answerChannel := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d!\n", numCorrect, len(problems))
			break problemloop
		case answer := <-answerChannel:
			if answer == problem.answer {
				numCorrect++
			}
		}

	}

}

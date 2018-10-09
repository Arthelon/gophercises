package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Question struct {
	text   string
	answer int
}

func handleError(err error) {
	fmt.Printf("Error: %s", err.Error())
	os.Exit(1)
}

func initQuestions(filePath string) *[]Question {
	output := []Question{}

	file, err := os.Open(filePath)
	if err != nil {
		handleError(err)
	}
	reader := csv.NewReader(bufio.NewReader(file))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			handleError(err)
		}
		if len(record) != 2 {
			handleError(errors.New("Invalid CSV input, each record must have 2 fields"))
		}
		answer, err := strconv.Atoi(record[1])
		if err != nil {
			handleError(err)
		}
		output = append(output, Question{
			text:   record[0],
			answer: answer,
		})
	}

	return &output
}

func main() {
	var timeLimit int
	var csvPath string
	flag.IntVar(&timeLimit, "limit", 10, "Time limit for each question in seconds")
	flag.StringVar(&csvPath, "path", "problems.csv", "Path to problems source file")
	flag.Parse()

	questions := initQuestions(csvPath)

	score := 0
	for idx, question := range *questions {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("#%d: %s\n", idx+1, question.text)
		scanner.Scan()
		resp := scanner.Text()
		respConv, err := strconv.Atoi(resp)
		for err != nil {
			fmt.Printf("Invalid input: %s. Please try again\n", resp)
			scanner.Scan()
			resp = scanner.Text()
			respConv, err = strconv.Atoi(resp)
		}
		if respConv == question.answer {
			score++
		}
	}
	fmt.Printf("Total Score: %d out of %d", score, len(*questions))
}

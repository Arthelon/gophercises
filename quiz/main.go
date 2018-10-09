package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	text   string
	answer int
}

func handleError(err error) {
	fmt.Printf("Error: %s", err.Error())
	os.Exit(1)
}

func initQuestions(filePath string) []Question {
	output := []Question{}

	file, err := os.Open(filePath)
	if err != nil {
		handleError(err)
	}
	defer file.Close()
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

	return output
}

func main() {
	var timeLimit int
	var csvPath string
	var shuffle bool
	flag.IntVar(&timeLimit, "limit", 10, "Time limit for each question in seconds")
	flag.StringVar(&csvPath, "path", "problems.csv", "Path to problems source file")
	flag.BoolVar(&shuffle, "shuffle", false, "Will shuffle questions if true")
	flag.Parse()

	score := 0
	done := make(chan bool)
	questions := initQuestions(csvPath)
	if shuffle {
		// Seed with current time
		rand.Seed(time.Now().UTC().UnixNano())
		shuffler := rand.Perm(len(questions))
		shuffled := make([]Question, len(questions))
		for i, v := range shuffler {
			shuffled[i] = questions[v]
		}
		questions = shuffled
	}

	go func() {
		for idx, question := range questions {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("#%d: %s\n", idx+1, question.text)
			timer := time.NewTimer(time.Second * time.Duration(timeLimit))
			go func() {
				<-timer.C
				done <- true
			}()
			scanner.Scan()
			resp := strings.Trim(scanner.Text(), " ")
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
			timer.Stop()
		}
		done <- true
	}()
	<-done
	fmt.Printf("You scored %d out of %d", score, len(questions))
}

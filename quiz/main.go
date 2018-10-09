package main

import (
	"flag"
	"fmt"
)

func main() {
	var timeLimit int
	var csvPath string
	flag.IntVar(&timeLimit, "limit", 10, "Time limit for each question in seconds")
	flag.StringVar(&csvPath, "path", "problems.csv", "Path to problems source file")
	flag.Parse()

	fmt.Println(csvPath)

}

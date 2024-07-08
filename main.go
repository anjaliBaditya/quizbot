package quiz

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Question struct that stores question with answer
type Question struct {
	question string
	answer   string
}

func main() {
	filename, timeLimit, shuffleQuestions, outputScore, retry := readArguments()
	f, err := openFile(filename)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer f.Close()

	questions, err := readCSV(f)
	if err != nil {
		log.Fatalf("Error reading questions: %v", err)
	}

	if shuffleQuestions {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}

	for {
		score, err := askQuestions(questions, timeLimit)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Your Score: %d/%d\n", score, len(questions))

		if outputScore != "" {
			saveScore(outputScore, score, len(questions))
		}

		if !retry {
			break
		}

		fmt.Println("Do you want to retry? (yes/no)")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(answer)) != "yes" {
			break
		}
	}
}

func readArguments() (string, int, bool, string, bool) {
	filename := flag.String("filename", "problems.csv", "CSV File containing quiz questions")
	timeLimit := flag.Int("limit", 30, "Time Limit for each question in seconds")
	shuffleQuestions := flag.Bool("shuffle", false, "Shuffle the questions")
	outputScore := flag.String("output", "", "File to save the score")
	retry := flag.Bool("retry", false, "Retry the quiz")
	flag.Parse()
	return *filename, *timeLimit, *shuffleQuestions, *outputScore, *retry
}

func readCSV(f io.Reader) ([]Question, error) {
	reader := csv.NewReader(f)
	allQuestions, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var questions []Question
	for _, line := range allQuestions {
		if len(line) != 2 {
			return nil, fmt.Errorf("Invalid format in CSV")
		}
		questions = append(questions, Question{question: line[0], answer: line[1]})
	}

	return questions, nil
}

func openFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

func getInput(input chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		result, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input <- result
	}
}

func askQuestions(questions []Question, timeLimit int) (int, error) {
	score := 0
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	done := make(chan string)
	go getInput(done)

	for _, question := range questions {
		fmt.Printf("%s: ", question.question)
		select {
		case <-timer.C:
			return score, fmt.Errorf("Time out")
		case answer := <-done:
			if strings.TrimSpace(strings.ToLower(answer)) == strings.TrimSpace(strings.ToLower(question.answer)) {
				score++
			}
		}
	}
	return score, nil
}

func saveScore(filename string, score, total int) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Could not open file to save score: %v", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "Score: %d/%d\n", score, total)
	if err != nil {
		log.Fatalf("Could not write score to file: %v", err)
	}
}

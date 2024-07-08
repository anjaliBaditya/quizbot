
# Go Quiz Program

This project is a simple quiz application written in Go. The quiz reads questions and answers from a CSV file, asks the user the questions, and calculates the score based on the user's answers.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Code Explanation](#code-explanation)
  - [Main Function](#main-function)
  - [Reading Arguments](#reading-arguments)
  - [Reading CSV File](#reading-csv-file)
  - [Opening File](#opening-file)
  - [Getting Input](#getting-input)
  - [Asking Questions](#asking-questions)
  - [Each Question](#each-question)
- [Conclusion](#conclusion)

## Features

- Reads quiz questions from a CSV file.
- Allows setting a time limit for answering each question.
- Calculates and displays the user's score.

## Prerequisites

- Go installed on your system.

## Installation

1. Clone this repository:
   \`\`\`bash
   git clone <repository-url>
   cd <repository-directory>
   \`\`\`

2. Build the quiz program:
   \`\`\`bash
   go build -o quiz
   \`\`\`

## Usage

Run the quiz program with the following command:

\`\`\`bash
./quiz -filename=problem.csv -limit=30
\`\`\`

- \`-filename\`: The CSV file containing quiz questions and answers (default: \`problem.csv\`).
- \`-limit\`: Time limit for each question in seconds (default: 30 seconds).

## Code Explanation

### Main Function

The `main` function orchestrates the quiz by reading the arguments, opening the CSV file, reading the questions, and asking the questions.

```go
func main() {
    filename, timeLimit := readArguments()
    f, err := openFile(filename)
    if err != nil {
        return
    }
    questions, err := readCSV(f)

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if questions == nil {
        return
    }
    score, err := askQuestion(questions, timeLimit)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("Your Score %d/%d
", score, totalQuestions)
}
```

### Reading Arguments

The `readArguments` function reads the command-line arguments for the CSV filename and time limit.

```go
func readArguments() (string, int) {
    filename := flag.String("filename", "problem.csv", "CSV File that conatins quiz questions")
    timeLimit := flag.Int("limit", 30, "Time Limit for each question")
    flag.Parse()
    return *filename, *timeLimit
}
```

### Reading CSV File

The `readCSV` function reads the questions and answers from the provided CSV file.

```go
func readCSV(f io.Reader) ([]Question, error) {
    allQuestions, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return nil, err
    }

    numOfQues := len(allQuestions)
    if numOfQues == 0 {
        return nil, fmt.Errorf("No Question in file")
    }

    var data []Question
    for _, line := range allQuestions {
        ques := Question{}
        ques.question = line[0]
        ques.answer = line[1]
        data = append(data, ques)
    }

    return data, nil
}
```

### Opening File

The `openFile` function opens the specified file.

```go
func openFile(filename string) (io.Reader, error) {
    return os.Open(filename)
}
```

### Getting Input

The `getInput` function reads user input from the standard input.

```go
func getInput(input chan string) {
    for {
        in := bufio.NewReader(os.Stdin)
        result, err := in.ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }

        input <- result
    }
}
```

### Asking Questions

The `askQuestion` function asks the questions to the user and calculates the score.

```go
func askQuestion(questions []Question, timeLimit int) (int, error) {
    totalScore := 0
    timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
    done := make(chan string)

    go getInput(done)

    for i := range [totalQuestions]int{} {
        ans, err := eachQuestion(questions[i].question, questions[i].answer, timer.C, done)
        if err != nil && ans == -1 {
            return totalScore, nil
        }
        totalScore += ans
    }
    return totalScore, nil
}
```

### Each Question

The `eachQuestion` function asks a single question and checks the user's answer.

```go
func eachQuestion(Quest string, answer string, timer <-chan time.Time, done <-chan string) (int, error) {
    fmt.Printf("%s: ", Quest)

    for {
        select {
        case <-timer:
            return -1, fmt.Errorf("Time out")
        case ans := <-done:
            score := 0
            if strings.Compare(strings.Trim(strings.ToLower(ans), "\n"), answer) == 0 {
                score = 1
            } else {
                return 0, fmt.Errorf("Wrong Answer")
            }

            return score, nil
        }
    }
}
```

## Conclusion

You now have a functional quiz application in Go that reads questions from a CSV file, allows setting a time limit for each question, and calculates the user's score. This setup provides a solid foundation for further enhancements and customization based on your specific requirements.
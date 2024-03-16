package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func problemPuller(filename string) ([]problem, error) {
	fObj, err := os.Open(filename)
	if err == nil {
		csvR := csv.NewReader(fObj)
		clines, err := csvR.ReadAll()
		if err == nil {
			return parseProblem(clines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in scv"+"format from %s file; %s", filename, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", filename, err.Error())
	}
}

func parseProblem(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

func main() {
	fName := flag.String("f", "quiz.csv", "path of the file")
	timer := flag.Int("t", 30, "time for the quiz")
	flag.Parse()
	problems, err := problemPuller(*fName)
	if err != nil {
		exit(fmt.Sprintf("somthig went wrong: %s", err.Error()))
	}
	correctAns := 0
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

	problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.q)
		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()
		select{
		case <-tObj.C:
			fmt.Println();
			break problemLoop;
		case iAns := <-ansC:
			if iAns == p.a {
			    correctAns++;
			}
			if i == len(problems) - 1{
				close(ansC);
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correctAns, len(problems));
	fmt.Printf("Press any key to exit.");

}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

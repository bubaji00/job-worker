package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (jd *JobDispatcher) startTimer(timeValue int, timeUnit string) {
	if jd.timerLimit() {
		return
	}
	for timeUnit == "" {
		input := getInput("Enter a time followed by valid unit. (sec, min, hr) or exit to go back: ")
		words := strings.Fields(input)

		if words[0] == EXIT {
			return
		}
		if len(words) != 2 {
			fmt.Println("Wrong format entered, please try again. e.g. 5 sec, 6 hr, 20 min")
			continue
		}
		//parse time and check for validation
		parsedTime, err := strconv.Atoi(words[0])
		if err != nil || parsedTime < 0 {
			fmt.Printf("The first word '%s' is not a valid positive integer.\n", words[0])
			continue
		}

		if !isValidUnit(words[1]) {
			fmt.Println("Error: The 'unit' must be one of 'sec', 'min', or 'hr'.")
			continue
		}

		timeValue = parsedTime
		timeUnit = words[1]
		break
	}

	duration := convertTime(timeValue, timeUnit)

	// Check if the duration exceeds the limit
	if duration > TimeLimit {
		fmt.Println(TimeLimitPrompt)
		return
	}

	newTimer := jd.newJob(duration, timeUnit)

	fmt.Printf("Timer %s started for %d %s.\n", newTimer.id, timeValue, newTimer.unit)

	go func(t *job, timeValue int) {
		select {
		case <-time.After(t.duration):
			fmt.Printf("Timer %d %s for %s is completed!\n", timeValue, t.unit, t.id)
			fmt.Printf("Enter a new command: ")
			t.changeState(COMPLETED)
		case <-t.stopChan:
			t.changeState(STOPPED)
		}
	}(&newTimer, timeValue)
}

func (jd *JobDispatcher) stopTimer() {
	id := getInput("Enter the user ID to stop its timer: ")
	timer, exists := jd.findUser(id)

	if !exists {
		return
	} else if timer.state != STARTED {
		fmt.Printf("Timer has already been %s\n", timer.state)
	} else {
		fmt.Printf("Timer %s has been stopped!\n", timer.id)
		timer.stopChan <- true
	}

}

func (jd *JobDispatcher) queryTimer() {
	id := getInput("Enter the user ID to check the status of the timer: ")
	timer, exists := jd.findUser(id)
	if !exists {
		fmt.Printf("No timer with ID %s found.\n", id)
	} else {
		fmt.Printf("Timer %s status: %s\n", id, timer.state)
	}
}

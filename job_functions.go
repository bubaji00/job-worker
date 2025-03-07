package worker

import (
	"fmt"
	"strconv"
	"strings"
)

func (jd *JobDispatcher) StartTimerCLI(timeValue int, timeUnit string) {
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

		if !IsValidUnit(words[1]) {
			fmt.Println("Error: The 'unit' must be one of 'sec', 'min', or 'hr'.")
			continue
		}

		timeValue = parsedTime
		timeUnit = words[1]
		break
	}

	duration := ConvertTime(timeValue, timeUnit)

	// Check if the duration exceeds the limit
	if duration > TimeLimit {
		fmt.Println(TimeLimitPrompt)
		return
	}

	newTimer := jd.newJob(duration, timeUnit)

	fmt.Printf("Timer %s started for %d %s.\n", newTimer.id, timeValue, newTimer.unit)

	jobID, err := jd.StartTimerCore(duration, timeUnit)
	if err != nil {
		fmt.Printf("Error starting timer: %v\n", err)
		return
	}

	fmt.Printf("Timer %s started for %d %s.\n", jobID, timeValue, timeUnit)
}

func (jd *JobDispatcher) stopTimerCLI() {
	id := getInput("Enter the user ID to stop its timer: ")

	if err := jd.StopTimerCore(id); err != nil {
		fmt.Printf("Error stopping timer: %v\n", err)
		return
	}
	fmt.Printf("Timer %s has been stopped!\n", id)

}

func (jd *JobDispatcher) queryTimerCLI() {
	id := getInput("Enter the user ID to check the status: ")
	state, err := jd.QueryTimerCore(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Timer %s status: %s\n", id, state)
}

package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutex sync.RWMutex
var reader = bufio.NewReader(os.Stdin)

type JobDispatcher struct {
	jobs map[string]*job // map of job id to job initialized as empty
}

type job struct {
	duration time.Duration
	start    time.Time
	state    string
	id       string
	unit     string
	stopChan chan bool
	wg       *sync.WaitGroup // Pass the WaitGroup to the job
}

func (jd *JobDispatcher) timerLimit() bool {
	if len(jd.jobs) < 5 {
		return false
	}

	fmt.Println("Can not start a new timer, 1000 timer limit reached")
	// offer to delete timer that are stopped or completed
	var jobsToDelete []string
	for id, job := range jd.jobs {
		if job.state == "completed" || job.state == "stopped" {
			jobsToDelete = append(jobsToDelete, id)
		}
	}

	// If there are no jobs to delete, exit the function
	if len(jobsToDelete) == 0 {
		fmt.Println("No completed or stopped timers to delete. Stop a timer before starting a new timer")
		return true
	}

	// List jobs to delete and ask for confirmation
	fmt.Printf("There are %d completed or stopped timers ready for deletion:\n", len(jobsToDelete))

	// Ask for user confirmation
	input := validInput("Do you want to delete these jobs from the map? (y/n): ")
	if input != "yes" && input != "y" {
		fmt.Println("No jobs were deleted.")
		return true
	}

	// Delete the selected jobs
	for _, id := range jobsToDelete {
		delete(jd.jobs, id)
	}
	fmt.Printf("Deleted %d jobs from the map.\n", len(jobsToDelete))
	return false
}

func (jd *JobDispatcher) getInput() {
	fmt.Print("\nEnter a new command\n" + "start: start a timer		|stop: stop a timer\n" +
		"query: return status of a timer	|end: terminate the program\n")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	words := strings.Fields(input)

	if len(words) != 1 {
		fmt.Println("No command or more than 1 command entered, please try again.")
		return
	}
	command := words[0]
	switch command {
	case "start":
		jd.startTimer(-1, "")
	case "stop":
		jd.stopTimer()
	case "query":
		jd.queryTimer()
	case "end":
		fmt.Println("Exiting program. Bye!")
		os.Exit(0)
	default:
		fmt.Println("Please only enter valid command: start/stop/query/end.")
	}
}

//starting a timer in go routine, and store it in map for future reference
func (jd *JobDispatcher) startTimer(timeValue int, timeUnit string) {
	if jd.timerLimit() {
		return
	}
	for timeUnit == "" {
		fmt.Println("Enter a time followed by valid unit. (sec, min, hr) or exit to go back: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		words := strings.Fields(input)

		if words[0] == "exit" {
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

	// Convert the time to time.Duration
	var duration time.Duration
	switch timeUnit {
	case "sec":
		duration = time.Duration(timeValue) * time.Second
	case "min":
		duration = time.Duration(timeValue) * time.Minute
	case "hr":
		duration = time.Duration(timeValue) * time.Hour
	}

	// Check if the duration exceeds 100 hours
	if duration > 100*time.Hour {
		fmt.Println("Error: The entered time exceeds the 100-hour limit.")
		return
	}

	timerId := uuid.NewString()[:6]
	var wg sync.WaitGroup
	wg.Add(1) // Increment WaitGroup counter
	newTimer := job{
		id:       timerId,
		duration: duration,
		start:    time.Now(),
		state:    "started",
		unit:     timeUnit,
		stopChan: make(chan bool),
		wg:       &wg, // Pass the WaitGroup to the job
	} // construct a new job struct

	mutex.Lock()                 // write lock
	jd.jobs[timerId] = &newTimer // store job in map
	mutex.Unlock()
	fmt.Printf("Timer %s started for %d %s.\n", newTimer.id, timeValue, newTimer.unit)

	go func(t *job, timeValue int) {
		defer t.wg.Done() // Signal that goroutine is done when exiting
		tempState := ""
		select {
		case <-time.After(t.duration):
			fmt.Printf("Timer %d %s for %s is completed!\n", timeValue, t.unit, t.id)
			tempState = "completed"
		case <-t.stopChan:
			fmt.Printf("Timer %d %s for %s has been stopped!\n", timeValue, t.unit, t.id)
			tempState = "stopped"
		}
		mutex.Lock()
		t.state = tempState
		mutex.Unlock()
		// Once complete, change the state
	}(&newTimer, timeValue)
}

func (jd *JobDispatcher) stopTimer() {
	userId := validInput("Enter the user ID to stop its timer")
	mutex.RLock()
	timer, exists := jd.jobs[userId]
	mutex.RUnlock() // Unlock before returning

	if !exists {
		fmt.Printf("No timer with ID %s found.\n", userId)
	} else if timer.state != "started" {
		fmt.Printf("Timer has already been %s\n", timer.state)
	} else {
		// Send stop signal before waiting
		timer.stopChan <- true
		// Wait for the goroutine to finish
		timer.wg.Wait()
		fmt.Printf("Timer %s has been stopped.\n", timer.id)
	}
}

func (jd *JobDispatcher) queryTimer() {
	userId := validInput("Enter the user ID to check the status of the timer")
	if userId == "" {
		return
	}

	mutex.RLock()
	timer, exists := jd.jobs[userId] // Check if the timer exists in the map
	mutex.RUnlock()
	if !exists {
		fmt.Printf("No timer with ID %s found.\n", userId)
		return
	} else if timer.state != "started" {
		fmt.Printf("Timer %s has already been %s\n", timer.id, timer.state)
	} else {
		fmt.Printf("Timer %s has been %s\n", timer.id, timer.state)
	}
}

func isValidUnit(unit string) bool {
	validUnits := map[string]bool{"sec": true, "min": true, "hr": true}
	return validUnits[unit]
}
func validInput(prompt string) string {
	fmt.Println(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	words := strings.Fields(input)
	if len(words) != 1 {
		fmt.Println("No argument or more than 1 argument received.")
		return ""
	}
	return words[0]
}

package worker

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"os"
	"strings"
	"sync"
	"time"
)

var Mutex sync.RWMutex
var reader = bufio.NewReader(os.Stdin)

type JobDispatcher struct {
	Jobs map[string]*Job // map of job id to job initialized as empty
}

type Job struct {
	duration time.Duration
	start    time.Time
	State    string
	id       string
	unit     string
	stopChan chan bool
}

func (jd *JobDispatcher) timerLimit() bool {
	if len(jd.Jobs) < JobLimit {
		return false
	}

	fmt.Println(JobLimitPrompt)
	// offer to delete timer that are stopped or completed
	var jobsToDelete []string
	for id, job := range jd.Jobs {
		if job.State == COMPLETED || job.State == STOPPED {
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
	input := getSingleInput("Do you want to delete these jobs from the map? (y/n): ")
	if input != YES && input != Y {
		fmt.Println("No jobs were deleted.")
		return true
	}

	// Delete the selected jobs
	for _, id := range jobsToDelete {
		delete(jd.Jobs, id)
	}
	fmt.Printf("Deleted %d jobs from the map.\n", len(jobsToDelete))
	return false
}

func (jd *JobDispatcher) findUser(id string) (*Job, bool) {
	Mutex.RLock()
	timer, exists := jd.Jobs[id]
	Mutex.RUnlock() // Unlock before returning
	return timer, exists
}

func (jb *Job) changeState(changedState string) {
	Mutex.Lock()
	jb.State = changedState
	Mutex.Unlock()
}

func (jd *JobDispatcher) newJob(timeDuration time.Duration, timeUnit string) Job {
	timerId := uuid.NewString()[:6]
	newTimer := Job{
		id:       timerId,
		duration: timeDuration,
		start:    time.Now(),
		State:    STARTED,
		unit:     timeUnit,
		stopChan: make(chan bool),
	} // construct a new job struct

	Mutex.Lock()                 // write lock
	jd.Jobs[timerId] = &newTimer // store job in map
	Mutex.Unlock()
	return newTimer
}

func (jd *JobDispatcher) Start() {
	c := cron.New()
	go c.Start()
	//defer c.Stop()

	prompt := "Available Commands:\n" +
		"• start: Start a timer\n" +
		"• stop: Stop a timer\n" +
		"• query: Return the status of a timer\n" +
		"• cron: Start a cron job\n" +
		"• end: Terminate the program\n" +
		"Enter a command: "
	for {
		command := getSingleInput(prompt)
		switch command {
		case START:
			jd.StartTimerCLI(-1, EMPTY)
		case STOP:
			jd.stopTimerCLI()
		case QUERY:
			jd.queryTimerCLI()
		case END:
			fmt.Println("Exiting program. Bye!")
			c.Stop()
			os.Exit(0)
		case CRON:
			cronJob(c)
		default:
			fmt.Printf("Please only enter valid command: start/stop/query/cron/end\n")
		}
		fmt.Println()
	}
}

// starting a timer in go routine, and store it in map for future reference
func IsValidUnit(unit string) bool {
	validUnits := map[string]bool{SEC: true, MIN: true, HR: true}
	return validUnits[unit]
}

func getSingleInput(prompt string) string {
	input := getInput(prompt)
	words := strings.Fields(input)
	if len(words) != 1 {
		fmt.Println("No argument or more than 1 argument received.")
		return EMPTY
	}
	return strings.ToLower(words[0])
}

func getInput(prompt string) string {
	fmt.Printf(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func ConvertTime(val int, unit string) time.Duration {
	switch unit { // Convert the time to time.Duration
	case SEC:
		return time.Duration(val) * time.Second
	case MIN:
		return time.Duration(val) * time.Minute
	case HR:
		return time.Duration(val) * time.Hour
	default:
		return -1
	}
}

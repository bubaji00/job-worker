package worker

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

var a = 0

func getCurrentTime() string {
	return time.Now().In(time.Local).Format(time.RFC3339)
}

func cronJobsExecutes(name string) {
	fmt.Printf("\nTask '%s' executed at %s\n",
		name, getCurrentTime())
	fmt.Printf("hello world %d\n", a)
	fmt.Printf("Enter a new command: ")
	a = a + 1
}

func cronJob(c *cron.Cron) {
	taskName := getInput("Enter a task name or type 'exit' to terminate: ")

	if strings.ToLower(taskName) == "exit" {
		fmt.Println("Exiting cron job...")
		return
	}

	// asking user input
	cronExpr := getInput("Enter a cron expression (e.g., '*/1 * * * *'): ")
	// executes

	// TODO this is for debugging convenience
	if cronExpr == "." {
		cronExpr = "*/1 * * * *"
	}

	_, err := c.AddFunc(cronExpr, func() {
		cronJobsExecutes(taskName)
	})
	if err != nil {
		fmt.Printf("invalid cron expression: %s\n", err)
		return
	}

	fmt.Printf("Task '%s' successfully scheduled with cron: %s at %s\n",
		taskName, cronExpr, getCurrentTime())
	//c.Start()
}

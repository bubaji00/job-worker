// go build -o timer.exe
// .\timer.exe start --time=[int] --unit=[sec/min/hr]

//since last meeting: adjusted mutex lock and unlock, refactoring into different file, OOD, prevent large time input (<100hr)

package main

import (
	"fmt"
)

func main() {
	var jobDispatch = JobDispatcher{
		jobs: make(map[string]*job),
	}
	rootCmd := setupCommands(&jobDispatch)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

	for {
		jobDispatch.getInput()
	}
}

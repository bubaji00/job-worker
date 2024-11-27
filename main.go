// go build -o timer.exe
// .\timer.exe start --time=[int] --unit=[sec/min/hr]

//since last meeting: adjusted mutex lock and unlock, refactoring into different file, OOD, prevent large time input (<100hr)

// 10/8 suggestion from teacher: have another file for global number, const, prompt, "compeled", "stopped"
// format input string prompt
/* pull out
mutex.Lock()                 // write lock
jd.jobs[timerId] = &newTimer // store job in map
mutex.Unlock()
*/
// no big deal for circular buffer
//timer wheel is better
// at the end of the project, user can start a timer and when the time is reached it will execute some cmd command on a remote machine

//parse cron express
//understand grpc, no need to mix first
//do a hello world grpc, say hello,
//add authentication, as simple as password, and code

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
	jobDispatch.start()
	//or {
	//	jobDispatch.start()
	//}
}

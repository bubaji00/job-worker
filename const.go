package main

import (
	"strconv"
	"time"
)

const JobLimit int = 5

var JobLimitPrompt = "Can not start a new timer, " + strconv.Itoa(JobLimit) + " timer limit reached"

const TimeLimit = 100 * time.Hour
const TimeLimitPrompt = "Error: The entered time exceeds the 100-hour limit."

const CRON string = "cron"
const START string = "start"
const STOP string = "stop"
const QUERY string = "query"
const END string = "end"
const EXIT string = "exit"

const STARTED string = "started"
const COMPLETED string = "completed"
const STOPPED string = "stopped"

const SEC string = "sec"
const MIN string = "min"
const HR string = "hr"

const Root string = "job"
const TIME string = "time"
const UNIT string = "unit"

const YES string = "yes"
const Y string = "y"
const EMPTY string = ""

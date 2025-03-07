package worker

import (
	"fmt"
	"time"
)

func (jd *JobDispatcher) StartTimerCore(duration time.Duration, timeUnit string) (string, error) {
	newTimer := jd.newJob(duration, timeUnit)
	go func(t *Job) {
		select {
		case <-time.After(t.duration):
			t.changeState(COMPLETED)
		case <-t.stopChan:
			t.changeState(STOPPED)
		}
	}(&newTimer)
	return newTimer.id, nil
}

func (jd *JobDispatcher) StopTimerCore(jobID string) error {
	timer, exists := jd.findUser(jobID)
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}
	if timer.State == STOPPED || timer.State == COMPLETED {
		return fmt.Errorf("job %s is already %s", jobID, timer.State)
	}
	timer.stopChan <- true
	return nil
}

func (jd *JobDispatcher) QueryTimerCore(jobID string) (string, error) {
	job, exists := jd.findUser(jobID)
	if !exists {
		return "", fmt.Errorf("job %s not found", jobID)
	}
	return job.State, nil
}

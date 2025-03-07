package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"jobWorker/worker"
)

var rootCmd = &cobra.Command{
	Use:   worker.Root,
	Short: "A simple CLI countdown timer",
	Long:  "This CLI supports starting, stopping, and checking the status of the countdown timer.",
}

func SetupCommands(dispatcher *worker.JobDispatcher) *cobra.Command {
	var startCmd = &cobra.Command{
		Use:   worker.START,
		Short: "starting to count down",
		Long:  "starting a countdown timer with given time and return user ID",

		RunE: func(cmd *cobra.Command, args []string) error {
			time1, err := cmd.Flags().GetInt(worker.TIME)
			if err != nil || time1 <= 0 {
				return fmt.Errorf("error: -time must be a positive integer")
			}

			unit, err := cmd.Flags().GetString(worker.UNIT)
			if err != nil || !worker.IsValidUnit(unit) {
				return fmt.Errorf("error: invalid or missing unit (valid options: sec, min, hr)")
			}

			dispatcher.StartTimerCLI(time1, unit)
			return nil
		},
	}
	//job1 := newJob("max", time.Now(), "as")
	startCmd.Flags().Int(worker.TIME, 0, "the time to count down")
	startCmd.Flags().String(worker.UNIT, worker.EMPTY, "the unit to countdown (sec, min, hr)")
	rootCmd.AddCommand(startCmd)
	return rootCmd
}

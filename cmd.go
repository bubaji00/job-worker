package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "job",
	Short: "A simple CLI countdown timer",
	Long:  "This CLI supports starting, stopping, and checking the status of the countdown timer.",
}

func setupCommands(dispatcher *JobDispatcher) *cobra.Command {
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "starting to count down",
		Long:  "starting a countdown timer with given time and return user ID",

		RunE: func(cmd *cobra.Command, args []string) error {
			time1, err := cmd.Flags().GetInt("time")
			if err != nil || time1 <= 0 {
				return fmt.Errorf("error: -time must be a positive integer")
			}

			unit, err := cmd.Flags().GetString("unit")
			if err != nil || !isValidUnit(unit) {
				return fmt.Errorf("error: invalid or missing unit (valid options: sec, min, hr)")
			}

			dispatcher.startTimer(time1, unit)
			return nil
		},
	}
	//job1 := newJob("max", time.Now(), "as")
	startCmd.Flags().Int("time", 0, "the time to count down")
	startCmd.Flags().String("unit", "", "the unit to countdown (sec, min, hr)")
	rootCmd.AddCommand(startCmd)
	return rootCmd
}

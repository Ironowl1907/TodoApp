package cmd

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var CMDCreate = &cobra.Command{
	Use:   "create [task name]",
	Short: "Creates a task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		DB, err := checkOrCreateDB(DefaultDBName)
		if err != nil {
			panic(err)
		}
		defer DB.Close()

		priorLevel, _ := cmd.Flags().GetInt("priority")

		// Update IDCounter based on existing records
		if err := updateIDCounter(DB); err != nil {
			fmt.Printf("[Error] Couldn't update ID counter: %s\n", err)
			return
		}

		writer := csv.NewWriter(DB)
		record := []string{
			strconv.Itoa(IDCounter),
			strings.Join(args, " "),
			strconv.Itoa(priorLevel),
			"0",
		}
		err = writer.Write(record)
		if err != nil {
			fmt.Printf("[Error] Could not write to file: %s\n", err)
		} else {
			IDCounter++
			fmt.Printf("Task created: %s\n", record[1])
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			fmt.Printf("[Error] Could not flush to file: %s\n", err)
		}
	},
}

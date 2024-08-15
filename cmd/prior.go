package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CMDPrior = &cobra.Command{
	Use:   "prioritize [task ID] [new taks priority]",
	Short: "Edits a task priority",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var edited bool = false
		DB, err := checkOrCreateDB(DefaultDBName)
		if err != nil {
			panic(err)
		}
		defer DB.Close()

		var filteredDBBuffer [][]string
		reader := csv.NewReader(DB)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, line := range records {
			if len(line) < 3 {
				fmt.Println("[Error] Invalid record format:", line)
				continue
			}

			if line[0] != args[0] {
				filteredDBBuffer = append(filteredDBBuffer, line)
			} else {
				line = []string{line[0], line[1], args[1], line[3]}
				filteredDBBuffer = append(filteredDBBuffer, line)
				edited = true
			}
		}

		if !edited {
			println("No matching ID in the database")
			return
		}

		// Recreate the DB with the edited one
		DB, err = os.Create(DefaultDBName)
		if err != nil {
			fmt.Println("[Error] Error while writing to new file:", err)
			return
		}
		defer DB.Close()

		csvWriter := csv.NewWriter(DB)
		err = csvWriter.WriteAll(filteredDBBuffer)
		if err != nil {
			fmt.Println("[Error] Couldn't write data:", err)
			return
		}

		csvWriter.Flush()
		fmt.Printf("Task with ID %s prioritized to level %s\n", args[0], args[1])
	},
}

package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CMDDone = &cobra.Command{
	Use:   "done [task ID]",
	Short: "Marks a task as done",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		DB, err := checkOrCreateDB(DefaultDBName)
		if err != nil {
			panic(err)
		}
		defer DB.Close()

		var found bool
		var filteredDBBuffer [][]string
		reader := csv.NewReader(DB)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("[Error] Couldn't read data:", err)
			return
		}

		for _, line := range records {
			if len(line) < 3 {
				fmt.Println("[Error] Invalid record format:", line)
				continue
			}

			if line[0] == args[0] {
				line[3] = "1" // Mark as done
				found = true
			}

			filteredDBBuffer = append(filteredDBBuffer, line)
		}

		// Recreate the DB file with the updated task status
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
		if found {
			fmt.Println("Task marked as done successfully.")
		} else {
			fmt.Println("Couldn't find record")
		}
	},
}

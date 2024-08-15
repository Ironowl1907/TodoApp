package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CMDDelete = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Deletes a task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var deleted bool = false
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
				deleted = true
			}
		}

		if !deleted {
			println("No matching ID in the database")
			return
		}

		// Recreate the DB file without the deleted task
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
	},
}

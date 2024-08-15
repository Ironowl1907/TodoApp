package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var CMDList = &cobra.Command{
	Use:   "list",
	Short: "Lists tasks",
	Run: func(cmd *cobra.Command, args []string) {
		DB, err := checkOrCreateDB(DefaultDBName)
		if err != nil {
			fmt.Printf("[Error] Couldn't access the database: %v\n", err)
			return
		}
		defer DB.Close()

		unDone, _ := cmd.Flags().GetBool("undone")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		reader := csv.NewReader(DB)

		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("[Error] Couldn't read data from the CSV file.")
			return
		}

		fmt.Fprintln(w, "ID\tTask\tPriority\tDone")
		for _, line := range records[1:] {
			if len(line) < 3 {
				fmt.Println("[Warning] Skipping incomplete record:", line)
				continue
			}
			if line[3] == "0" || !unDone {
				fmt.Fprintf(w, "%s\t%s\t%s ", line[0], line[1], line[2])
				if line[3] == "0" {
					fmt.Fprint(w, "\t[ ]\n")
				} else {
					fmt.Fprint(w, "\t[X]\n")
				}
			}
		}

		if err := w.Flush(); err != nil {
			fmt.Printf("[Error] Couldn't flush writer: %v\n", err)
		}
	},
}

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

const DefaultDBName = "DB.csv"

var IDCounter int = 0

// checkOrCreateDB checks if the database file exists and has the correct format; if not, it creates a new one.
func checkOrCreateDB(name string) (*os.File, error) {
	// Try to open the database file
	DBFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// Check if the file is empty
	fileInfo, err := DBFile.Stat()
	if err != nil {
		return nil, err
	}

	// If the file is empty, write the header row
	if fileInfo.Size() == 0 {
		writer := csv.NewWriter(DBFile)
		err = writer.Write([]string{"ID", "Name", "Done"})
		if err != nil {
			return nil, err
		}
		writer.Flush()
	}

	// Check if the header row has the correct format
	reader := csv.NewReader(DBFile)
	firstLine, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if firstLine[0] != "ID" || firstLine[1] != "Name" || firstLine[2] != "Done" {
		return nil, errors.New("format error in the DB file")
	}

	// Reset the file pointer to the beginning
	_, err = DBFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return DBFile, nil
}

// updateIDCounter updates the global IDCounter based on existing records in the database.
func updateIDCounter(DB *os.File) error {
	reader := csv.NewReader(DB)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Find the maximum ID in the existing records
	for _, record := range records[1:] { // Skip header
		if len(record) > 0 {
			id, err := strconv.Atoi(record[0])
			if err == nil && id >= IDCounter {
				IDCounter = id + 1
			}
		}
	}
	return nil
}

func main() {
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

			// Update IDCounter based on existing records
			if err := updateIDCounter(DB); err != nil {
				fmt.Printf("[Error] Couldn't update ID counter: %s\n", err)
				return
			}

			writer := csv.NewWriter(DB)
			record := []string{
				strconv.Itoa(IDCounter),
				strings.Join(args, " "),
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
				}	else {
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

	var CMDShow = &cobra.Command{
		Use:   "show",
		Short: "Show tasks",
		Run: func(cmd *cobra.Command, args []string) {
			DB, err := checkOrCreateDB(DefaultDBName)
			if err != nil {
				fmt.Printf("[Error] Couldn't access the database: %v\n", err)
				return
			}
			defer DB.Close()

			allTask, _ := cmd.Flags().GetBool("all")

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			reader := csv.NewReader(DB)

			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println("[Error] Couldn't read data from the CSV file.")
				return
			}

			fmt.Fprintln(w, "ID\tTask\tDone")
			for _, line := range records[1:] {
				if len(line) < 3 {
					fmt.Println("[Warning] Skipping incomplete record:", line)
					continue
				}
				if line[2] == "0" || allTask {
					fmt.Fprintf(w, "%s\t%s ", line[0], line[1])
					if line[2] == "0" {
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
					line[2] = "1" // Mark as done
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

	var rootCmd = &cobra.Command{
		Use:   "toDo",
		Short: "Basic to-do app",
	}
	rootCmd.AddCommand(CMDCreate)
	rootCmd.AddCommand(CMDDelete)
	rootCmd.AddCommand(CMDShow)
	CMDShow.Flags().BoolP("all", "a", false, "Show all the created (non deleted) tasks")
	rootCmd.AddCommand(CMDDone)
	rootCmd.Execute()
}

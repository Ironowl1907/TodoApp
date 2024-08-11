package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const DefaultDBName = "DB.csv"

var IDCounter int = 0

// Function to check if a database exists; if not, prompt to create it
func checkOrCreateDB(name string) (*os.File, error) {
	// Try to open the database file
	DBFile, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		// If the file does not exist, prompt to create it
		fmt.Println("Database file does not exist.")
		var response string
		fmt.Print("Create new DB? (y/N): ")
		fmt.Scanf("%s", &response)
		if response == "y" || response == "Y" {
			// Create a new DB file with the specified name
			DBFile, err = os.Create(name)
			if err != nil {
				return nil, errors.New("error creating the database")
			}
			fmt.Println("Database created:", name)
		} else {
			return nil, errors.New("database not created")
		}
	} else {
	}
	return DBFile, nil
}

func main() {
	var CMDCreate = &cobra.Command{
		Use:   "create [task name]",
		Short: "Creates a task",
		Long:  `later`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			DB, err := checkOrCreateDB(DefaultDBName)
			if err != nil {
				panic(err)
			}
			defer DB.Close()
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
		Long:  `later`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

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
				panic("[Error] Couldn't read data")
			}

			for _, line := range records {
				if len(line) < 3 {
					fmt.Println("[Error] Invalid record format:", line)
					continue
				}

				if line[0] != args[0] {
					filteredDBBuffer = append(filteredDBBuffer, line)
				}
			}

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
		Short: "Creates a task",
		Long:  `later`,
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			DB, err := checkOrCreateDB(DefaultDBName)
			if err != nil {
				panic(err)
			}
			defer DB.Close()

			reader := csv.NewReader(DB)
			records, err := reader.ReadAll()
			if err != nil {
				println("[Error] Couldn't read data")
			}
			fmt.Println(records)
		},
	}

	var CMDDone = &cobra.Command{
		Use:   "done [task name]",
		Short: "Creates a task",
		Long:  `later`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			DB, err := checkOrCreateDB(DefaultDBName)
			if err != nil {
				panic(err)
			}
			defer DB.Close()

			var found bool = false

			var filteredDBBuffer [][]string
			reader := csv.NewReader(DB)
			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println("[Error] Couldn't read data:", err)
				panic(err)
			}

			for _, line := range records {
				if len(line) < 3 {
					fmt.Println("[Error] Invalid record format:", line)
					continue
				}

				if line[0] == args[0] {
					line[2] = "1"
					found = true
				}
				filteredDBBuffer = append(filteredDBBuffer, line)
			}

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
	rootCmd.AddCommand(CMDDone)
	rootCmd.Execute()

}

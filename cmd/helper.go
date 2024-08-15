package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

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
		fmt.Println("[Warning] Creating and formating new database file")
		writer := csv.NewWriter(DBFile)
		err = writer.Write([]string{"ID", "Name", "Priority", "Done"})
		if err != nil {
			return nil, err
		}
		writer.Flush()
	}

	// Reset the file pointer to the beginning
	_, err = DBFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Check if the header row has the correct format
	reader := csv.NewReader(DBFile)
	firstLine, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if firstLine[0] != "ID" || firstLine[1] != "Name" || firstLine[2] != "Priority" || firstLine[3] != "Done" {
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

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Basic to-do app",
}

func Execute() {
	Init()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Init() {
	rootCmd.AddCommand(CMDCreate)
	CMDCreate.Flags().IntP("priority", "p", 1, "Set a priority for the task")
	rootCmd.AddCommand(CMDDelete)
	rootCmd.AddCommand(CMDList)
	rootCmd.AddCommand(CMDEdit)
	CMDList.Flags().BoolP("undone", "u", false, "Lists all the undone (non deleted) tasks")
	rootCmd.AddCommand(CMDDone)
	rootCmd.AddCommand(CMDPrior)
}

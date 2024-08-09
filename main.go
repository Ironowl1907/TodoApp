package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const DBfileName = "DB.csv"
	var DBFile *os.File
	var err error

	fmt.Println("[To-Do app] by ironowl")
	DBFile, err = os.OpenFile(DBfileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("[Warning] no file named '%s' was found\n", DBfileName)
		var response string
		fmt.Print("Create new DB? (y/N): ")
		fmt.Scanf("%s", &response)

		if response == "y" || response == "Y" {
			// Create a new DB file
			DBFile, err = os.Create(DBfileName)
			if err != nil {
				fmt.Printf("[Error] Could not create file: %s\n", err)
				return
			}
			fmt.Println("New database created successfully.")
		} else {
			fmt.Println("No database created. Exiting.")
			return
		}
	}
	defer DBFile.Close()

	fmt.Printf("Opening database '%s'\n", DBfileName)

	writer := csv.NewWriter(DBFile)

	var IDCounter int = 0

	var CLInput string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		CLInput = scanner.Text()
		var InputSplited []string = strings.Split(CLInput, " ")

		switch InputSplited[0] {
		case "create":
			if len(InputSplited) < 2 {
				fmt.Println("[Error] Task description required.")
				continue
			}
			record := []string{
				strconv.Itoa(IDCounter),
				strings.Join(InputSplited[1:], " "),
				"0",
			}
			err := writer.Write(record)
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
		case "show":
			fmt.Println("Show functionality not implemented yet.")
		case "remove":
			fmt.Println("Remove functionality not implemented yet.")
		case "done":
			fmt.Println("Done functionality not implemented yet.")
		case "exit":
			return
		default:
			fmt.Printf("[Error] undefined command -> '%s'\n", InputSplited[0])
		}
		fmt.Print("> ")
	}
}

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	const DBfileName = "DB.csv"
	var DBFile *os.File
	fmt.Println("[To-Do app] by ironowl")
	DBFile, err := os.Open(DBfileName)
	if err != nil {
		fmt.Printf("[Warning] no file named '%s' was found\n", DBfileName)
		var response string
		fmt.Print("Create new DB? (y/N): ")
		fmt.Scanf("%s", &response)

		if response == "y" || response == "Y" {
			// Create a new DB file
			_, err := os.Create(DBfileName)
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

	var CLInput string
	for {
		fmt.Print("> ")
		fmt.Scanln(&CLInput)
		var InputSplited []string = strings.Split(CLInput, " ")

		switch InputSplited[0] {
		case "create":
			fmt.Println("create")
		case "show":
			fmt.Println("create")
		case "remove":
			fmt.Println("create")
		case "done":
			fmt.Println("create")
		default:
			fmt.Printf("[Error] undefined command -> '%s'\n", InputSplited[0])
		}
	}
}

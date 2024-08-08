package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	type TODOs struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		Status      bool   `json:"status"`
	}

	counter := 0
	data := make([]TODOs, 0)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("CMDs: show create remove done\n> ")

	for scanner.Scan() {
		line := scanner.Text()
		var split []string = strings.Split(line, " ")
		switch split[0] {
		case "create":
			data = append(data, TODOs{
				ID:          counter,
				Description: strings.Join(split[1:], " "),
				Status:      false,
			})
			counter++
			fmt.Println("Created TODO item")
		case "remove":
			index, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			found := false
			for i, todo := range data {
				if todo.ID == index {
					data = append(data[:i], data[i+1:]...)
					fmt.Println(data)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("TODO item with ID %d not found\n", index)
			}

		case "show":
			for _, todo := range data {
				fmt.Println(todo)
			}
		case "done":
			index, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			for i, todo := range data {
				if todo.ID == index {
					data[i].Status = true
					break
				}
			}
		}
		fmt.Print("> ")
	}

}

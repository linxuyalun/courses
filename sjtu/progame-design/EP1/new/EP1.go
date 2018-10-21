package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	// Reading files
	file, err := os.Open("./shakespeare.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// input the word
	fmt.Println("Please enter the word:")
	var input string
	fmt.Scan(&input)

	// `\w` matches a letter or a number
	split := regexp.MustCompile(`[\w]+`)
	lineNum := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		text := split.FindAllString(scanner.Text(), -1)
		if text == nil {
			continue
		}
		if text[0] == input {
			if len(text) == 1 {
				//The line only has one word
				fmt.Printf("%d: %s, %s, %s\n", lineNum, "NULL", text[0], "NULL")
				continue
			}
			fmt.Printf("%d: %s, %s, %s\n", lineNum, "NULL", text[0], text[1])
		}
		if text[len(text)-1] == input {
			fmt.Printf("%d: %s, %s, %s\n", lineNum, text[len(text)-2], text[len(text)-1], "NULL")
		}
		for i := 1; i < len(text)-1; i++ {
			if text[i] == input {
				fmt.Printf("%d: %s, %s, %s\n", lineNum, text[i-1], text[i], text[i+1])
			}
		}
	}

}

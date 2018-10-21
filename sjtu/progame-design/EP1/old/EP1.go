package main

import (
	"fmt"
	"io/ioutil"
)

// Split the string.
func split(letter uint8) bool {
	dim := []uint8{'\t', ',', '.', ':', '#', '*', '(', ')', '[', ']', ' '}
	for i := range dim {
		if letter == dim[i] {
			return true
		}
	}
	return false
}

func main() {

	// Reading files.
	dat, err := ioutil.ReadFile("./shakespeare.txt")
	if err != nil {
		panic(err)
	}
	text := string(dat)

	fmt.Println("Please enter the word:")
	var input string
	fmt.Scan(&input)

	// flag is the line number
	// k is a string to store a word
	// content is a slice to store words
	// j is a number to mark the location of content
	flag := 1
	k := ""
	var content = [3]string{}
	j := 0

	// Iterate over the text
	for i := 0; i < len(text); i++ {
		if text[i] == '\r' {
			// '\r\n' always appears together, "i++" skips \n
			i++
			flag++
			if k != "" {
				content[j] = k
				// Compare two words
				if content[j] == input {
					fmt.Printf("LINE NUMBER: %d,  THE PREVIOUS WORD: %s, THE WORD: %s, THE FOLLOWING WORD: %s\n",
						flag, content[(j+2)%3], content[j], "no word!")
				}
				j = (j + 1) % 3
				k = ""
			}
			// Avoid repeated "no word!"
			if content[(j+2)%3] != "no word!" {
				content[j] = "no word!"
				j = (j + 1) % 3
			}
			continue
		}

		if split(text[i]) {
			if k == "" {
				continue
			}
			content[j] = k
			// Compare two words
			if content[(j+2)%3] == input {
				fmt.Printf("LINE NUMBER: %d,  THE PREVIOUS WORD: %s, THE WORD: %s, THE FOLLOWING WORD: %s\n",
					flag, content[(j+1)%3], content[(j+2)%3], content[j%3])
			}
			k = ""
			j = (j + 1) % 3
			continue
		}
		k += string(text[i])
	}
}

package main

import (
	"fmt"
	"strings"
)

// https://www.codewars.com/kata/5727bb0fe81185ae62000ae3/train/go

// 6kyu - Backspaces in string

// Assume "#" is like a backspace in string. This means that string "a#bc#d" actually is "bd"
// Your task is to process a string with "#" symbols.
// Examples
// "abc#d##c"      ==>  "ac"
// "abc##d######"  ==>  ""
// "#######"       ==>  ""
// ""              ==>  ""

func CleanString(s string) string {
	var output []string

	for _, rune := range s {
		char := string(rune)

		if char == "#" {
			if len(output) > 0 {
				output = output[:len(output)-1]
			}
		} else {
			output = append(output, char)
		}
	}

	fmt.Println(output)
	return strings.Join(output, "")
}

func main() {
	CleanString("abc#d##c")
}

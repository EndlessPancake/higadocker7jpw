package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// match, _ := regexp.MatchString("もなか", "")
	//  fmt.Println(match)

	open := func(path string) []string {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", path, err)
		}
		defer f.Close()
		lines := []string{}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if serr := scanner.Err(); serr != nil {
			fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", path, err)
		}
		return lines
	}

	// file := open("./item.txt")
	// fmt.Println(len(file))
	for _, v := range open("./item.txt") {
		r, _ := regexp.Compile("もなか")
		if r.MatchString(v) {
			fmt.Println("もなかやないか！")
		} else {
			fmt.Println("ほな、もなかと違うかあ...")
		}
	}
}

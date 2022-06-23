package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
        args := strings.Fields("-T -P 443 -r 1.1.1.1") // => ["-T", "-P"]
	cmd := exec.Command("mtr", args...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		// fmt.Println("\t---")
	}

	cmd.Wait()
}

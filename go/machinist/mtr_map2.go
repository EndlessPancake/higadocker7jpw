package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"regexp"
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

	lines := []string{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
	    lines = append(lines, scanner.Text())
	    // fmt.Println(scanner.Text())
   	    // fmt.Println("\t---")
	}

	for _, v := range lines {
		regEx := `^[0-9]`
		reg, _ := regexp.Compile(regEx)
		tmp := strings.Fields(v)

		if reg.MatchString(tmp[0]) {
			hop := strings.Replace(tmp[0], ".|--", "", 1) // get HOP Number from result
			avg, _ := strconv.ParseFloat(tmp[5], 64)
			best, _ := strconv.ParseFloat(tmp[6], 64)
			worst, _ := strconv.ParseFloat(tmp[7], 64)
			stdev, _ := strconv.ParseFloat(tmp[8], 64)

			fmt.Println(hop, avg, best, worst, stdev)
		}
	}
	cmd.Wait()
}

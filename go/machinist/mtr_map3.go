package main

import (
	"bufio"
	"fmt"
	"time"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"regexp"
	"encoding/json"
        "github.com/garyburd/redigo/redis"
)

type Rttdata struct {
    DateAtStore int64
    Hops        string
    Avg         float64
    Best        float64
    Worst       float64
    Stdev       float64
}

func redis_connect() redis.Conn {
	const IP_PORT = "localhost:6379"

	c, err := redis.Dial("tcp", IP_PORT)
	if err != nil {
	  panic(err)
	}
	return c
}

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

	// t := time.Now()           // get localtime, after mtr commands
	t := time.Now().Unix()    // get localtime as UNIXTIME, after mtr commands

	for _, v := range lines {
	    regEx := `^[0-9]`
	    reg, _ := regexp.Compile(regEx)
	    tmp := strings.Fields(v)

	    if reg.MatchString(tmp[0]) {
		hop := strings.Replace(tmp[0], ".|--", "", 1)   // How many HOPs ? from result
		avg, _ := strconv.ParseFloat(tmp[5], 64)
	 	best, _ := strconv.ParseFloat(tmp[6], 64)
		worst, _ := strconv.ParseFloat(tmp[7], 64)
		stdev, _ := strconv.ParseFloat(tmp[8], 64)
		// fmt.Println(hop, avg, best, worst, stdev)

		// - Struct to Json
		rttdata := &Rttdata{DateAtStore: t, Hops: hop, Avg: avg, Best: best, Worst: worst, Stdev: stdev}
		serialized, _ := json.Marshal(rttdata)

		// store data to redis
		c := redis_connect()
		defer c.Close()

		// set JSON
		c.Do("SET", "json_test", serialized)

		// get JSON
		data, _ := redis.Bytes(c.Do("GET", "json_test"))

		// JSON to struct
		if data != nil {
		   deserialized := new(Rttdata)
		   json.Unmarshal(serialized, deserialized)
        	   fmt.Println("deserialized : ", deserialized)
		}
	    }
	}
	cmd.Wait()
}

package main

import (
        "fmt"
        "os"
        "os/exec"
)

func main() {
        out, err := exec.Command("ls", "-la").Output()

        if err != nil {
           fmt.Println(err)
           os.Exit(1)
        }

        fmt.Println(string(out))
}

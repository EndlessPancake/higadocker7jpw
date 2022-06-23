package main
import "fmt"
import "regexp"
func main() {
 match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
 fmt.Println(match)

 // others
 r, _ := regexp.Compile("p([a-z]+)ch")

 fmt.Println(r.MatchString("peach"))

}

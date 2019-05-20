package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	for i := 92000; i < 150000; i += 1000 {
		log.Printf("Running %s;", strconv.Itoa(i))
		c := exec.Command("python", "C:/dev/projects/EliCDavis/smartest-reddits/backfiller/gsen.py", strconv.Itoa(i))
		out, err := c.Output()

		if err != nil {
			fmt.Fprintln(c.Stderr, err)
			fmt.Fprintln(os.Stdout)
			log.Println(string(out))
			log.Printf("Error Running Script: %s", err.Error())
			return
		}
	}
}

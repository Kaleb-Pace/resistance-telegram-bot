package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readFile(number int, out *os.File) {
	file, err := os.Open(fmt.Sprintf("gsen-result%d.txt", number))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	for scanner.Scan() {
		_, err := out.WriteString(scanner.Text() + "\n")
		if err != nil {
			panic(err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {

	f, err := os.Create("all-g-results.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	f.WriteString("messageID, polarity, magnitude\n")

	for i := 0; i < 144; i++ {
		readFile(i*1000, f)
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}

}

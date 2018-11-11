package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func usersSpiteLeft(chatId string, poster string) (int, error) {
	f, err := os.OpenFile("yeets/"+chatId, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		if poster == entry[0] {
			return strconv.Atoi(entry[1])
		}
	}

	return 3, nil
}

type spiteEntry struct {
	user  string
	spite int
}

func tradeSpite(chatID string, from string, to string) error {
	f, err := os.OpenFile("yeets/"+chatID, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	yeeters := make([]spiteEntry, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		yeets, err := strconv.Atoi(entry[1])
		if err != nil {
			return err
		}
		yeeters = append(yeeters, spiteEntry{
			user:  entry[0],
			spite: yeets,
		})
		y := yeeters[len(yeeters)-1]
		log.Printf("[%s]:%d\n", y.user, y.spite)
	}

	yeeterFound := false
	yeetedFound := false
	for i, yeeter := range yeeters {
		if yeeter.user == from {
			yeeters[i].spite--
			yeeterFound = true
			log.Printf("found yeeter")
		} else if yeeter.user == to {
			yeeters[i].spite++
			yeetedFound = true
			log.Printf("found yeeted")
		}
	}

	if yeeterFound == false {
		yeeters = append(yeeters, spiteEntry{
			user:  from,
			spite: 2,
		})
	}

	if yeetedFound == false {
		yeeters = append(yeeters, spiteEntry{
			user:  to,
			spite: 4,
		})
	}

	f.Truncate(0)
	f.Seek(0, 0)

	for _, yeeter := range yeeters {
		if _, err = f.WriteString(fmt.Sprintf("%s %d\n", yeeter.user, yeeter.spite)); err != nil {
			return err
		}
	}

	return nil
}

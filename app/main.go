// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RedditResponse A typical resonse when searching reddit
type RedditResponse struct {
	Data struct {
		Children []struct {
			Data *Submission
		}
	}
}

// Resposible for sending a message to the appropriate group chat
func sendMessage(message string, url string, update Update) {

	// Send Message to telegram's api
	resp, err := http.Post(url+"sendMessage", "application/json", bytes.NewBuffer([]byte(`{
		"chat_id": `+strconv.FormatInt(update.Message.Chat.ID, 10)+`,
		"text": "`+message+`",
		"parse_mode": "HTML"
	}`)))

	// Catch errors
	if err != nil {
		log.Println("Error sending message:")
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	// Read and print message
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("\nTelegram Said: ")
	log.Println(string(body))
}

func rule34Search(term string, url string, update Update, errorLogger func(string)) {
	log.Println("searching rule 34: " + term)
	searchURL := "https://www.reddit.com/r/rule34/search.json?q=" + term + "&restrict_sr=on&sort=relevance&t=all"
	resp, err := http.Get(searchURL)

	if err != nil {
		errorLogger("Error Searching Reddit: " + err.Error())
	}

	defer resp.Body.Close()

	r := RedditResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf(string(body))
	json.Unmarshal([]byte(body), &r)
	if err != nil {
		errorLogger("Error Parsing Reddit Response: " + err.Error())
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	log.Println("Succesful")
	log.Println(submissions)

	if len(submissions) > 0 {
		log.Println("How's this? : " + submissions[0].URL)
		// sendMessage("How's this? : "+submissions[0].URL, url, update)
	} else {
		log.Println("Couldn't find any porn for: " + term)
		// sendMessage("Couldn't find any porn for: "+term, url, update)
	}
}

// Builds and returns commands with url.
func getCommands(url string, errorLogger func(string)) []func(Update) {

	return []func(update Update){

		// Kill command
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "kill")
			if len(commands) > 1 {
				n := rand.Int() % len(killStatements)
				go sendMessage(strings.TrimSpace(commands[1])+killStatements[n], url, update)
			}
		},

		// Traps command
		func(update Update) {
			if strings.Contains(update.Message.Text, "traps") {
				go sendMessage("https://www.youtube.com/watch?v=9E1YYSZ9qrk", url, update)
			}
		},

		// God command
		func(update Update) {
			if strings.Contains(update.Message.Text, "gg") {
				go sendMessage("GOD IS GREAT", url, update)
			}
		},

		// Rule34 command
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "rule34")
			if len(commands) > 1 {
				go rule34Search(strings.TrimSpace(commands[1]), url, update, errorLogger)
			}
		},

		//pokedexSerach
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "pokedex")
			if len(commands) > 1 {
				go pokedexSerach(strings.TrimSpace(commands[1]), url, update, errorLogger)
			}
		},
	}
}

// Create our routes
func initRoutes(router *gin.Engine, teleurl string, errors *[]string) {

	router.LoadHTMLGlob("templates/*")

	//errors := [...]string{"help", "me"}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":  "Telegram Bot",
			"errors": errors,
		})
	})

}

func listenForUpdates(teleurl string, errorLogger func(string)) {
	var lastUpdate = -1

	commands := getCommands(teleurl, errorLogger)

	for {
		time.Sleep(time.Second)

		resp, err := http.Get(teleurl + "getUpdates?offset=" + strconv.Itoa(lastUpdate))
		if err != nil {
			errorLogger("Error Obtaining Updates: " + err.Error())
			return
		}

		defer resp.Body.Close()

		var updates BatchUpdates
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorLogger("Error Reading Body: " + err.Error())
			return
		}

		err = json.Unmarshal([]byte(body), &updates)
		if err != nil {
			errorLogger("Error Parsing Telegram getUpdates Response: " + err.Error())
			return
		}

		// Dispatch incoming messages to appropriate functions
		for _, update := range updates.Result {
			log.Println("Msg: " + update.Message.Text)
			lastUpdate = update.UpdateID + 1
			for _, command := range commands {
				command(update)
			}
		}

	}
}

func main() {

	// Can't run a server without a port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable was not set")
		return
	}

	log.Println("Starting bot!")

	errorMessages := []string{}
	var errorLogger = func(msg string) {
		t := time.Now()
		log.Println(msg)
		newMsg := [...]string{t.Format("Mon Jan _2 15:04:05 2006") + ": " + msg}

		errorMessages = append(newMsg[:], errorMessages...)
	}

	teleurl := "https://api.telegram.org/bot" + os.Getenv("TELE_KEY") + "/"

	go listenForUpdates(teleurl, errorLogger)

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	initRoutes(r, teleurl, &errorMessages)
	r.Run(":" + port)

}
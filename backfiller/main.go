package main

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

func findChildNode(parent *html.Node, elementType, class string) *html.Node {
	if parent == nil {
		return nil
	}

	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == elementType {
			for _, a := range child.Attr {
				if a.Key == "class" && a.Val == class {
					return child
				}
			}
		}
	}

	return nil
}

func getTextFromNode(messageNode *html.Node) (string, error) {
	if messageNode == nil {
		return "", errors.New("Require a node to operate on")
	}

	textNode := findChildNode(messageNode, "div", "text")
	if textNode != nil {
		return strings.TrimSpace(textNode.FirstChild.Data), nil
	}
	return "", nil
}

func getFromUserFromNode(messageNode *html.Node) (string, error) {
	if messageNode == nil {
		return "", errors.New("Require a node to operate on")
	}

	textNode := findChildNode(messageNode, "div", "from_name")
	if textNode != nil {
		return strings.TrimSpace(textNode.FirstChild.Data), nil
	}
	return "", nil
}

func processHistoryNode(historyNode *html.Node) {
	for child := historyNode.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "div" {
			for _, a := range child.Attr {
				if a.Key == "class" && strings.Contains(a.Val, "message default clearfix") {
					body := findChildNode(child, "div", "body")
					message, err := getTextFromNode(body)
					if err != nil {
						panic(err)
					}

					fromUser, err := getFromUserFromNode(body)
					if err != nil {
						panic(err)
					}

					log.Printf("%s: %s", fromUser, message)
				}
			}
		}
	}
}

func processFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err.Error())
	}
	doc, err := html.Parse(file)
	if err != nil {
		panic(err.Error())
	}
	var f func(*html.Node)
	found := 0
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "history" {
					processHistoryNode(n)
					found++
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if found != 1 {
		panic("found history: " + string(found))
	}
}

func main() {

	folderToExplore := "ChatExport_08_03_2019"

	processFile(path.Join(folderToExplore, "messages.html"))

	// files, err := ioutil.ReadDir(folderToExplore)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, file := range files {
	// 	if file.IsDir() == false {
	// 		processFile(path.Join(folderToExplore, file.Name()))
	// 	}
	// }

}

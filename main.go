package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ilkyun/PubScrapper/scrapper"
	"github.com/labstack/echo"
)

const fileName string = "results.csv"

// to add open browser with localhost

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello. This is our fist Go web app")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	impact := c.FormValue("impact")
	imp, _ := strconv.ParseFloat(impact, 32)
	reldate := c.FormValue("reldate")
	retmax := c.FormValue("retmax")
	scrapper.Scrape(term, retmax, imp, reldate)

	return c.Attachment(fileName, fileName)
}

//GetPort get the Port from the environment
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == " " {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected")
	}
	return ":" + port
}

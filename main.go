package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/ilkyun/PubScrapper/scrapper"
	"github.com/labstack/echo"
)

const fileName string = "results.csv"

// to add open browser with localhost

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Start(port)

}

func handleHome(c echo.Context) error {
	return c.File("index.html")
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

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
		fmt.Println("darwin!")
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

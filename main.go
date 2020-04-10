package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/ilkyun/PubScrapper/scrapper"
	"github.com/labstack/echo"
)

const fileName string = "results.csv"

// to add open browser with localhost

func main() {
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	log.Fatal("$PORT must be set")
	// }
	openbrowser("http://localhost:1323/")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Static("/", dir+"/"+"index.html")
	fmt.Println(dir + "index.html")
	e.POST("/scrape", handleScrape)
	e.Start(":1323")

}

// func handleHome(c echo.Context) error {
// 	return c.File("index.html")
// }

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
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

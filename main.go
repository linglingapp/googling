package main

import (
	"hello/googlescraper"
	"os"

	"github.com/labstack/echo"
)

const fileName string = "results.csv"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	searchTerm := c.FormValue("searchTerm")
	countryCode := c.FormValue("countryCode")
	languageCode := c.FormValue("languageCode")
	googlescraper.GoogleScrape(searchTerm, countryCode, languageCode)
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}

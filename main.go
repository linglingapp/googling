package main

import (
	"bingscraper"
	"googlescraper"
	"os"

	"github.com/labstack/echo"
)

const fileName string = "results.csv"

func handleHome(c echo.Context) error {
	return c.File("index.html")
}

func handleGoogleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	searchTerm := c.FormValue("searchTerm")
	countryCode := c.FormValue("countryCode")
	googlescraper.GoogleScrape(searchTerm, countryCode)
	return c.Attachment(fileName, fileName)
}

func handleBingScrape(c echo.Context) error {
	defer os.Remove(fileName)
	searchTerm := c.FormValue("searchTerm")
	countryCode := c.FormValue("countryCode")
	bingscraper.BingScrape(searchTerm, countryCode)
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/googlescrape", handleGoogleScrape)
	e.POST("/bingscrape", handleBingScrape)
	e.Logger.Fatal(e.Start(":1323"))
}

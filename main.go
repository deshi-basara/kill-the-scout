package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asdine/storm/v3"
	"github.com/deshi-basara/kill-the-scout/database"
	"github.com/deshi-basara/kill-the-scout/notifier"
	"github.com/deshi-basara/kill-the-scout/scraper"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	url := os.Getenv("IMMO_SEARCH_URL")
	mailTo := os.Getenv("MAIL_TO")

	// open embedded database
	db := database.Setup()

	//  get all latest (at max 20) Exposes from the result-list and check if they were already scraped
	results := scraper.ScrapeResultList(url)
	for _, exposeID := range results {
		_, err := database.GetExpose(db, exposeID)
		if err == storm.ErrNotFound {
			fmt.Printf("Expose %s not found, should scrape \n", exposeID)

			// scarpe and save expose
			expose := scraper.ScrapeExpose(exposeID)
			err := database.SaveExpose(db, expose)
			if err != nil {
				log.Fatalf("SaveExpose.Fatal: %s \n", err)
			}

			// notify
			notifier.SendMail(mailTo, expose)
		}
	}

	defer db.Close()
}
